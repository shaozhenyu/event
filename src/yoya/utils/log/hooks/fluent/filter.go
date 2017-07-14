package fluent

// function to convert error type to string type.
func FilterError(v interface{}) interface{} {
	if err, ok := v.(error); ok {
		return err.Error()
	}
	return v
}
