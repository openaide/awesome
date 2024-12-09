package main

import (
	"flag"
	"os"

	"github.com/openaide/docker/vanna/proxy"
)

func main() {
	// Set default values using environment variables
	defaultHost := os.Getenv("POSTGRES_HOST")
	if defaultHost == "" {
		defaultHost = "localhost"
	}

	defaultPort := os.Getenv("POSTGRES_PORT")
	if defaultPort == "" {
		defaultPort = "5432"
	}

	defaultUser := os.Getenv("POSTGRES_USER")
	if defaultUser == "" {
		defaultUser = "postgres"
	}

	defaultDbName := os.Getenv("POSTGRES_DBNAME")
	if defaultDbName == "" {
		defaultDbName = "postgres"
	}

	//
	dbHost := flag.String("host", defaultHost, "Postgres host")
	dbPort := flag.String("port", defaultPort, "Postgres port")
	dbUser := flag.String("user", defaultUser, "Postgres user")
	dbName := flag.String("dbname", defaultDbName, "Postgres database name")

	dbPassword := flag.String("pass", "", "Database password")

	//
	serverAddress := flag.String("address", ":58080", "Server address")

	//
	trainPath := flag.String("train", "", "Path to training data")
	storePath := flag.String("store", "", "Path to store data")

	venvPath := flag.String("venv", "", "Specify the Python venv path")
	appScript := flag.String("script", "app.py", "Python app script to execute")

	// Parse the flags
	flag.Parse()

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
				Username: *dbUser,
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
