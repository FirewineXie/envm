package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func SaveSettings() {
	marshal, _ := json.Marshal(env)
	err := ioutil.WriteFile(env.SettingPath, []byte(marshal), 0644)
	if err != nil {
		log.Fatalf("save setting error + %v", err.Error())
		return
	}
}

func ReadSettings() {

	if _, err := os.Stat(env.SettingPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("settings does not exist")
		} else {
			fmt.Println("please check settings file")
		}
	} else {
		file, err := ioutil.ReadFile(env.SettingPath)
		if err != nil {
			log.Fatalf("read setting error + %v", err.Error())
			return
		}
		m := make(map[string]string)
		if err = json.Unmarshal(file, &m); err != nil {
			log.Fatalf("read setting error + %v", err.Error())
			return
		}

		if val, ok := m["download"]; ok {
			env.Downloads = filepath.Clean(val)
		}
		// Make sure the directories exist
		_, e := os.Stat(env.Root)
		if e != nil {
			fmt.Println(env.Root + " could not be found or does not exist. Exiting.")
			return
		}
	}

}
