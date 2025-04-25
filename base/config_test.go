package base

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const cliConfigJSON = `[
	{"project_id":"org-bdks4e","region":"cn-bj2","zone":"cn-bj2-04","base_url":"%s","timeout_sec":15,"profile":"uweb","active":true},
	{"project_id":"org-oxjwoi","region":"hk","zone":"hk-02","base_url":"%s","timeout_sec":15,"profile":"test","active":false}
]`

const credentialJSON = `[
	{"public_key":"4E9UU*****3ZAPWQ==","private_key":"6945*****a0d45","profile":"uweb"},
	{"public_key":"YSQG*****zgnCRQ=","private_key":"jtma*****Avms","profile":"test"}
]`

func TestAggConfigManager(t *testing.T) {
	os.MkdirAll(ConfigPath, 0700)
	baseURL := fmt.Sprintf("https://api.%s/", BrandURL)
	configJSON := fmt.Sprintf(cliConfigJSON, baseURL, baseURL)
	err := ioutil.WriteFile(ConfigPath+"/config.json", []byte(configJSON), LocalFileMode)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile(ConfigPath+"/credential.json", []byte(credentialJSON), LocalFileMode)
	if err != nil {
		t.Error(err)
	}
	defer func() {
		err := os.RemoveAll(ConfigPath)
		if err != nil {
			t.Error(err)
		}
	}()

	configFile, err := os.OpenFile(ConfigPath+"/config.json", os.O_CREATE|os.O_RDONLY, LocalFileMode)
	if err != nil {
		t.Error(err)
	}

	credFile, err := os.OpenFile(ConfigPath+"/credential.json", os.O_CREATE|os.O_RDONLY, LocalFileMode)
	if err != nil {
		t.Error(err)
	}

	acManager, err := NewAggConfigManager(configFile, credFile)
	if err != nil {
		t.Error(err)
	}

	if len(acManager.configs) != 2 {
		t.Errorf("expect length of configs is 2, accpet %d", len(acManager.configs))
	}

}

func TestEmptyAggConfigManager(t *testing.T) {
	os.MkdirAll(ConfigPath, 0700)
	defer func() {
		err := os.RemoveAll(ConfigPath)
		if err != nil {
			t.Error(err)
		}
	}()

	configFile, err := os.OpenFile(ConfigPath+"/config.json", os.O_CREATE|os.O_RDONLY, LocalFileMode)
	if err != nil {
		t.Error(err)
	}

	credFile, err := os.OpenFile(ConfigPath+"/credential.json", os.O_CREATE|os.O_RDONLY, LocalFileMode)
	if err != nil {
		t.Error(err)
	}

	acManager, err := NewAggConfigManager(configFile, credFile)
	if err != nil {
		t.Error(err)
	}

	err = acManager.Load()
	if err != nil {
		t.Fatal(err)
	}

	if len(acManager.configs) != 0 {
		t.Errorf("expect length of configs is 2, accpet %d", len(acManager.configs))
	}
}
