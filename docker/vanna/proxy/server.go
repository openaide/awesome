package proxy

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type DBInfo struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string

	// Postgres version string
	Version string
}

// DSN returns the data source name for connecting to the database.
func (d *DBInfo) DSN() string {
	host := d.Host
	if host == "" {
		host = "localhost"
	}
	port := d.Port
	if port == "" {
		port = "5432"
	}
	dbname := d.DBName
	if dbname == "" {
		dbname = "postgres"
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, d.Username, d.Password, dbname)
}

func (d *DBInfo) IsValid() bool {
	return d.Username != "" && d.Password != ""
}

type AppConfig struct {
	Name    string
	Version string

	DBInfo *DBInfo

	TrainPath string
	StorePath string

	VenvPath  string
	AppScript string
}

type WebConfig struct {
	App     *AppConfig
	Address string
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func StartServer(config *WebConfig) error {
	loginPage := NewLoginPage(config)
	dbPage := NewDBPage(config)
	proxyPage := NewProxyPage(config)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if config.App.DBInfo.IsValid() {
			http.Redirect(w, r, "/db", http.StatusFound)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	}).Methods("GET")

	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/login", loginPage.Handler()).Methods("GET", "POST")
	r.HandleFunc("/db", dbPage.Handler()).Methods("GET")
	r.HandleFunc("/db/{dbname}", proxyPage.Handler()).Methods("GET")
	r.HandleFunc("/app/{user}/{dbname}/{subpath:.*}", proxyPage.Proxy())

	addr := config.Address

	log.Printf("SQL Copilot running at: %s\n", addr)

	return http.ListenAndServe(addr, r)
}
