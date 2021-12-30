package config

import (
	"gopkg.in/yaml.v2"
)

type HuskType struct {
	GlobalObject   interface{}
	CurrentCfgBuff []byte
	AllowOverWrite bool
	IsBackupCfg    bool
}

func (ht *HuskType) SetGlobalObject(obj interface{}) {
	ht.GlobalObject = obj
}

func (ht *HuskType) SetCfgBuff(buff []byte) {
	ht.CurrentCfgBuff = buff
}

func (ht *HuskType) SetAllowOverWrite(allow bool) {
	ht.AllowOverWrite = allow
}

func (ht *HuskType) SetBackupCfg(isBackup bool) {
	ht.IsBackupCfg = isBackup
}

var GlobalHusk = &HuskType{
	AllowOverWrite: false,
	IsBackupCfg:    false,
}

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
	GlobalHusk.SetCfgBuff(ymlBuf)
	var um interface{}
	err := yaml.Unmarshal(ymlBuf, &um)
	if err != nil {
		return err
	}
	GlobalHusk.SetGlobalObject(jsonConverter(um))
	//GlobalObject = jsonConverter(um)
	return nil
}

func GetRawObject() interface{} {
	return GlobalHusk.GlobalObject
}
