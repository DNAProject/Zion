package maas_config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	InitMaasConfig()
	os.Exit(m.Run())
}

