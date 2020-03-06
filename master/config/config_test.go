package config

import "testing"

func TestParseConfig(t *testing.T) {
	fileName := "../master.json"
	err := ParseConfig(fileName)
	if err != nil {
		t.Errorf("test config.ParseConfig err: %v", err)
	}
	if G_Config == nil {
		t.Errorf("test config.ParseConfig err: empty G_config")
	}
	if G_Config.ApiServerPort != 8080 {
		t.Errorf("test config.ParseConfig err, wrong server apiport, expected 8080,"+
			" actual: %d", G_Config.ApiServerPort)
	}
	if G_Config.ApiServerReadTimeOut != 5000 {
		t.Errorf("test config.ParseConfig err, wrong server apireadTimeout, expected 5000,"+
			" actual: %d", G_Config.ApiServerReadTimeOut)
	}
	if G_Config.ApiServerWriteTimeOut != 5000 {
		t.Errorf("test config.ParseConfig err, wrong server apiWriteTimeout, expected 5000,"+
			" actual: %d", G_Config.ApiServerWriteTimeOut)
	}
}
