package config

import (
	"encoding/json"
	"fmt"
	"github.com/FirewineXie/govm/inner/arch"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func SaveSettings() {

	marshal, _ := json.Marshal(env)
	err := ioutil.WriteFile(env.Settings, []byte(marshal), 0644)
	if err != nil {
		log.Fatalf("save setting error + %v", err.Error())
		return
	}
}

func ReadSettings() {

	file, err := ioutil.ReadFile(env.Settings)
	if err != nil {
		log.Fatalf("read setting error + %v", err.Error())
		return
	}
	m := make(map[string]string)
	if err = json.Unmarshal(file, m); err != nil {
		log.Fatalf("read setting error + %v", err.Error())
		return
	}

	if val, ok := m["root"]; ok {
		env.Root = filepath.Clean(val)
	}

	if val, ok := m["arch"]; ok {
		env.Arch = val
	}

	env.Arch = arch.Validate(env.Arch)

	// Make sure the directories exist
	_, e := os.Stat(env.Root)
	if e != nil {
		fmt.Println(env.Root + " could not be found or does not exist. Exiting.")
		return
	}
}
