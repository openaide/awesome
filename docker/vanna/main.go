package main

import (
	"flag"
	"os"

	"github.com/openaide/docker/vanna/proxy"
)

func main() {
	// Define command-line flags
	dbHost := flag.String("host", "localhost", "Database host")
	dbPort := flag.String("port", "5432", "Database port")
	dbUsername := flag.String("user", "postgres", "Database user")
	dbPassword := flag.String("pass", "", "Database password")
	dbName := flag.String("dbname", "postgres", "Database name")
	serverAddress := flag.String("address", ":58080", "Server address")

	// New flags for train and store paths
	trainPath := flag.String("train", "", "Path to training data")
	storePath := flag.String("store", "", "Path to store data")

	venvPath := flag.String("venv", "", "Specify the Python venv path")
	appScript := flag.String("script", "app.py", "Python app script to execute")

	// Parse the flags
	flag.Parse()

	// Set dbUsername from environment if not provided via flags
	if *dbUsername == "" {
		*dbUsername = os.Getenv("POSTGRES_USER")
		if *dbUsername == "" {
			*dbUsername = "postgres"
		}
	}

	// Set dbPassword from environment if not provided via flags
	if *dbPassword == "" {
		*dbPassword = os.Getenv("POSTGRES_PASSWORD")
	}

	// Set trainPath from environment if not provided via flags
	if *trainPath == "" {
		*trainPath = os.Getenv("TRAIN_BASE")
		if *trainPath == "" {
			*trainPath = "local/train" // Default to "local/train" if not set
		}
	}

	// Set storePath from environment if not provided via flags
	if *storePath == "" {
		*storePath = os.Getenv("STORE_BASE")
		if *storePath == "" {
			*storePath = "local/store" // Default to "local/store" if not set
		}
	}

	if *venvPath == "" {
		*venvPath = os.Getenv("PYTHON_VENV")
		if *venvPath == "" {
			*venvPath = "./local/venv/"
		}
	}

	// Use venvPath for spawning the Python app
	proxy.StartServer(&proxy.WebConfig{
		App: &proxy.AppConfig{
			Name:    "SQL Copilot",
			Version: "0.0.1",
			DBInfo: &proxy.DBInfo{
				Host:     *dbHost,
				Port:     *dbPort,
				Username: *dbUsername,
				Password: *dbPassword,
				DBName:   *dbName,
			},
			TrainPath: *trainPath,
			StorePath: *storePath,
			VenvPath:  *venvPath,
			AppScript: *appScript,
		},
		Address: *serverAddress,
	})
}
