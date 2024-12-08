package proxy

import (
	"os"
	"testing"
)

// Start web server for local dev
func TestStartServer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test server ...")
	}

	os.Setenv("STORE_BASE", "/tmp/store")
	os.Setenv("TRAIN_BASE", "/tmp/train")

	StartServer(&WebConfig{
		App: &AppConfig{
			Name:    "SQL Copilot",
			Version: "0.0.1",
			DBInfo: &DBInfo{
				Host:     "localhost",
				Port:     "5432",
				Username: "postgres",
				Password: "",
				DBName:   "postgres",
			},
		},
		Address: ":58080",
	})
}
