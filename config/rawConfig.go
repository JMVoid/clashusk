package config

import (
	"gopkg.in/yaml.v2"
)

var (
	GlobalObject   interface{}
	CurrentCfgBuff []byte
	InitOverwrite  bool
	BackupCfg      bool
)

func jsonConverter(input interface{}) interface{} {
	switch it := input.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range it {
			m2[k.(string)] = jsonConverter(v)
		}
		return m2
	case []interface{}:
		for mk, mv := range it {
			it[mk] = jsonConverter(mv)
		}
	}
	return input
}

func JsonRawConfig(ymlBuf []byte) error {
	CurrentCfgBuff = ymlBuf
	var um interface{}
	err := yaml.Unmarshal(ymlBuf, &um)
	if err != nil {
		return err
	}
	GlobalObject = jsonConverter(um)
	return nil
}

func GetRawObject() interface{} {
	return GlobalObject
}
