package proxy

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/koding/websocketproxy"
)

const UserIDHeader = "X-User-Id"

var proxyClient = &http.Client{
	Timeout: time.Second * 10,
}

type ProxyServer struct {
	app *AppConfig

	port   int
	user   string
	dbname string

	started int
	ready   bool
	proxy   *httputil.ReverseProxy

	sync.RWMutex
}

func (p *ProxyServer) isReady() bool {
	p.RLock()
	defer p.RUnlock()

	return p.ready
}

func (p *ProxyServer) setReady(b bool) {
	p.Lock()
	defer p.Unlock()

	p.ready = b
}

func (p *ProxyServer) start() {
	p.Lock()
	defer p.Unlock()

	//only start once
	if p.started > 0 {
		return
	}

	p.started++

	go p.spawn()

	go p.check()
}

func (p *ProxyServer) stop() {
	p.Lock()
	defer p.Unlock()
	p.started--
}

func (p *ProxyServer) spawn() {
	log.Printf("Proxy port: %v user: %s", p.port, p.user)

	getEnv := func(key, defVal string) string {
		if v, exist := os.LookupEnv(key); exist {
			return v
		}
		return defVal
	}

	llmAPIKey := getEnv("LLM_API_KEY", "sk-1234")
	llmBaseUrl := getEnv("LLM_BASE_URL", "http://host.docker.internal:4000")
	llmModel := getEnv("LLM_MODEL", "gpt-4o")

	db := p.app.DBInfo

	storeBase := getEnv("STORE_BASE", p.app.StorePath)
	trainBase := getEnv("TRAIN_BASE", p.app.TrainPath)

	envVars := []string{
		fmt.Sprintf("HOST=%v", "0.0.0.0"),
		fmt.Sprintf("PORT=%v", p.port),
		//
		fmt.Sprintf("LLM_API_KEY=%v", llmAPIKey),
		fmt.Sprintf("LLM_BASE_URL=%v", llmBaseUrl),
		fmt.Sprintf("LLM_MODEL=%v", llmModel),
		//
		fmt.Sprintf("POSTGRES_HOST=%v", db.Host),
		fmt.Sprintf("POSTGRES_PORT=%v", db.Port),
		fmt.Sprintf("POSTGRES_USER=%v", db.Username),
		fmt.Sprintf("POSTGRES_PASSWORD=%v", db.Password),

		fmt.Sprintf("POSTGRES_DBNAME=%v", p.dbname),
		//
		fmt.Sprintf("STORE_PATH=%v", filepath.Join(storeBase, p.user, p.dbname)),
		fmt.Sprintf("TRAIN_PATH=%v", filepath.Join(trainBase, p.dbname)),
	}

	python := filepath.Join(p.app.VenvPath, "bin/python")
	args := []string{p.app.AppScript, "serve"}

	err := Retry(func() error {
		//localhost 127.0.0.1 0.0.0.0
		e := p.runScript(python, args, envVars)
		log.Printf("Proxy %v user: %v  %v", p.app, p.user, e)
		return e
	})

	//
	p.stop()

	log.Printf("Proxy %v failed or exited. user: %v error: %v", p.started, p.user, err)
}

func (p *ProxyServer) check() {
	// check
	uri := fmt.Sprintf("http://127.0.0.1:%v/", p.port)
	err := Retry(func() error {
		b, e := isServerReady(uri)

		if e == nil {
			p.setReady(b)
		}

		log.Printf("Proxy isServerReady %v port: %v ready: %v error: %v", p.user, p.port, b, e)
		return e
	}, NewBackOff(12, 1*time.Second))

	log.Printf("Proxy started: %v %v port: %v  error: %v", p.started, p.user, p.port, err)
}

func NewProxyServer(app *AppConfig, port int, user, dbname string) *ProxyServer {
	s := ProxyServer{
		app:     app,
		port:    port,
		user:    user,
		dbname:  dbname,
		started: 0,
	}
	return &s
}

type ProxyPage struct {
	app *AppConfig

	servers map[string]*ProxyServer

	sync.RWMutex
}

func (p *ProxyPage) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("ProxyPage.Handler: %s\n", r.URL.Path)

		vars := mux.Vars(r)
		dbname := vars["dbname"]

		user := p.app.DBInfo.Username
		if uid, ok := r.Header[UserIDHeader]; ok {
			user = uid[0]
		}

		prefix := fmt.Sprintf("/app/%s/%s/", user, dbname)
		http.Redirect(w, r, prefix, http.StatusFound)
	}
}

func (p *ProxyPage) Proxy() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxy  %s\n", r.URL.Path)

		//
		if p.app.DBInfo == nil || !p.app.DBInfo.IsValid() {
			http.Error(w, "Invalid DBInfo", http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		user := vars["user"]
		dbname := vars["dbname"]

		prefix := fmt.Sprintf("/app/%s/%s", user, dbname)
		subpath := fmt.Sprintf("/%s", vars["subpath"])

		//
		getServer := func() (*ProxyServer, bool) {
			p.Lock()
			defer p.Unlock()

			server, ok := p.servers[prefix]
			if ok {
				log.Printf("Proxy  already exists. proxy: %s port: %v\n", prefix, server.port)
				return server, true
			}

			port := FreePort()
			server = NewProxyServer(p.app, port, user, dbname)
			server.proxy = httputil.NewSingleHostReverseProxy(&url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("localhost:%v", port),
				Path:   "/",
			})
			server.start()
			p.servers[prefix] = server

			log.Printf("Proxy  created. prefix: %s port: %v\n", prefix, port)
			return server, false
		}

		server, exist := getServer()
		if !exist {
			// serve loading page
			loadingHandler(w, r)
			return
		}

		if !server.isReady() {
			serviceUnavailableHandler(w, r)
			return
		}

		// proxy request
		// remove prefix
		r.URL.Scheme = "http"
		r.URL.Host = fmt.Sprintf("localhost:%v", server.port)
		r.URL.Path = subpath

		ws := isWebSocketRequest(r)
		log.Printf("Proxy  %s ws: %v subpath: %s\n", r.URL.Path, ws, subpath)
		if ws {
			wsProxy := websocketproxy.NewProxy(&url.URL{
				Scheme: "ws",
				Host:   r.URL.Host,
				Path:   r.URL.Path,
			})

			wsProxy.ServeHTTP(w, r)
			log.Printf("Proxy  ws proxy: %s started: %v", r.URL, server.started)
			return
		}

		server.proxy.ServeHTTP(w, r)

		log.Printf("Proxy  http proxy: %s started: %v", r.URL, server.started)
	}
}

func isWebSocketRequest(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Connection"), "Upgrade") &&
		strings.EqualFold(r.Header.Get("Upgrade"), "websocket")
}

func NewProxyPage(c *WebConfig) *ProxyPage {
	return &ProxyPage{
		app:     c.App,
		servers: make(map[string]*ProxyServer),
	}
}

func isServerReady(uri string) (bool, error) {
	res, err := proxyClient.Get(uri)

	if err != nil {
		log.Printf("Proxy isServerReady: %v", err)
		return false, err
	}
	defer res.Body.Close()

	log.Printf("Proxy isServerReady: %v", res)

	return (res.StatusCode == 200), nil
}

func (p *ProxyServer) runScript(bin string, args []string, env []string) error {
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), env...)

	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("### Proxy runScript StdoutPipe %v", err)
		return err
	}
	stdout := bufio.NewScanner(outPipe)
	go func() {
		for stdout.Scan() {
			fmt.Printf("OUT> %s\n", stdout.Text())
		}
	}()

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("### Proxy runScript StderrPipe %v", err)
		return err
	}
	stderr := bufio.NewScanner(errPipe)
	go func() {
		for stderr.Scan() {
			fmt.Printf("STDERR> %s\n", stderr.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Printf("### Proxy runScript Start %v", err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		log.Printf("### Proxy runScript Wait %v", err)
	}

	log.Printf("### Proxy runScript exited due to error or inactivity, user: %s", p.user)
	return nil
}
