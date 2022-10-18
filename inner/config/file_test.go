package config

import (
	"testing"
)

func TestSaveSettings(t *testing.T) {
	env.SettingPath = "testdata/settings"
	env.Root = "D:\\ProgramData\\Go"
	SaveSettings()
}
