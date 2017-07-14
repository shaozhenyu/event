package fluent

import (
	"yoya/utils/log/logrus"
)

const (
	TagName      = "fluent"
	TagField     = "tag"
	MessageField = "message"
)

// FluentHook is logrus hook for fluentd.
type FluentHook struct {
	Fluent *Fluent

	host   string
	port   int
	levels []logrus.Level
	tag    *string

	ignoreFields map[string]struct{}
	filters      map[string]func(interface{}) interface{}
}

// NewHook returns initialized logrus hook for fluentd.
func NewHook(host string, port int) *FluentHook {
	fluent, err := NewClient(Config{
		FluentHost: host,
		FluentPort: port,
	})
	if err != nil {
		return nil
	}
	return &FluentHook{
		Fluent:       fluent,
		host:         host,
		port:         port,
		levels:       logrus.AllLevels,
		tag:          nil,
		ignoreFields: make(map[string]struct{}),
		filters:      make(map[string]func(interface{}) interface{}),
	}
}

// Levels returns logging level to fire this hook.
func (hook *FluentHook) Levels() []logrus.Level {
	return hook.levels
}

// SetLevels sets logging level to fire this hook.
func (hook *FluentHook) SetLevels(levels []logrus.Level) {
	hook.levels = levels
}

// Tag returns custom static tag.
func (hook *FluentHook) Tag() string {
	if hook.tag == nil {
		return ""
	}

	return *hook.tag
}

// SetTag sets custom static tag to override tag in the message fields.
func (hook *FluentHook) SetTag(tag string) {
	hook.tag = &tag
}

// AddIgnore adds field name to ignore.
func (hook *FluentHook) AddIgnore(name string) {
	hook.ignoreFields[name] = struct{}{}
}

// AddFilter adds a custom filter function.
func (hook *FluentHook) AddFilter(name string, fn func(interface{}) interface{}) {
	hook.filters[name] = fn
}

// Fire is invoked by logrus and sends log to fluentd logger.
func (hook *FluentHook) Fire(entry *logrus.Entry) error {
	var err error

	// Create a map for passing to FluentD
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		if _, ok := hook.ignoreFields[k]; ok {
			continue
		}
		if fn, ok := hook.filters[k]; ok {
			v = fn(v)
		}
		data[k] = v
	}
	setLevelString(entry, data)
	tag := hook.getTagAndDel(entry, data)
	if tag != entry.Message {
		setMessage(entry, data)
	}

	fluentData := ConvertToValue(data, TagName)
	err = hook.Fluent.PostWithTime(tag, entry.Time, fluentData)
	return err
}

// getTagAndDel extracts tag data from log entry and custom log fields.
// 1. if tag is set in the hook, use it.
// 2. if tag is set in custom fields, use it.
// 3. if cannot find tag data, use entry.Message as tag.
func (hook *FluentHook) getTagAndDel(entry *logrus.Entry, data logrus.Fields) string {
	// use static tag from
	if hook.tag != nil {
		return *hook.tag
	}

	tagField, ok := data[TagField]
	if !ok {
		return entry.Message
	}

	tag, ok := tagField.(string)
	if !ok {
		return entry.Message
	}

	// remove tag from data fields
	delete(data, TagField)
	return tag
}

func setLevelString(entry *logrus.Entry, data logrus.Fields) {
	data["level"] = entry.Level.String()
}

func setMessage(entry *logrus.Entry, data logrus.Fields) {
	if _, ok := data[MessageField]; !ok {
		data[MessageField] = entry.Message
	}
}
