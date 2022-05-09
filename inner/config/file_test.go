package config

import (
	"testing"
)

func TestSaveSettings(t *testing.T) {
	env.Settings = "testdata/settings"
	env.Root = "D:\\ProgramData\\Go"
	SaveSettings()
}
