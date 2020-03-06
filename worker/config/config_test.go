package config

import "testing"

func TestParseConfig(t *testing.T) {
	fileName := "../worker.json"
	err := ParseConfig(fileName)
	if err != nil {
		t.Errorf("test config.ParseConfig err: %v", err)
	}
	if G_Config == nil {
		t.Errorf("test config.ParseConfig err: empty G_config")
	}
	if G_Config.EtcdClientTimeOut != 3000 {
		t.Errorf("test config.ParseConfig err, wrong server EtcdClientTimeOut, expected 3000,"+
			" actual: %d", G_Config.EtcdClientTimeOut)
	}
}
