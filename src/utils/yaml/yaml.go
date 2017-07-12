/**************************************************************************

Copyright:YOYA

Author: shaozhenyu

Date:2017-06-22

Description:  functions to handle yaml file

**************************************************************************/

package yaml

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Read(configFile string, conf interface{}) error {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		return err
	}

	return nil
}
