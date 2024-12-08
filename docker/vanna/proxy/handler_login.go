package proxy

import (
	"fmt"
	"log"
	"net/http"
)

const loginTemplate = `
<div class="container">
	<h1>Welcome to {{ .App.Name }}</h1>
	<p>Version: {{ .App.Version }}</p>
	<p>Postgres SQL Copilot is a web-based tool for querying Postgres databases.</p>
	<div>
	  <form action="/login" method="post">
	  	<label for="host">Host:</label><br>
		<input type="text" id="host" name="host" value="{{ .App.DBInfo.Host }}"><br>
		<label for="port">Port:</label><br>
		<input type="text" id="port" name="port" value="{{ .App.DBInfo.Port }}"><br>
	    <label for="username">Username:</label><br>
	    <input type="text" id="username" name="username" value="{{ .App.DBInfo.Username }}"><br>
		<label for="password">Password:</label><br>
		<input type="password" id="password" name="password"><br>
		<label for="dbname">Database Name:</label><br>
		<input type="text" id="dbname" name="dbname" value="{{ .App.DBInfo.DBName }}"><br>
		<input type="submit" value="Submit" />
	  </form>
	</div>
</div>
`

type LoginPage struct {
	app   *AppConfig
	title string
}

func NewLoginPage(c *WebConfig) *LoginPage {
	return &LoginPage{
		app:   c.App,
		title: "Login",
	}
}

func (p LoginPage) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("LoginPage.Handler: Path=%s, Referrer=%s\n", r.URL.Path, r.Header.Get("Referer"))
		if r.Method == "GET" {
			renderTemplate(w, p.app, loginTemplate, p.title, nil)
			return
		}

		if host := r.FormValue("host"); host != "" {
			p.app.DBInfo.Host = host
		}

		if port := r.FormValue("port"); port != "" {
			p.app.DBInfo.Port = port
		}

		if username := r.FormValue("username"); username != "" {
			p.app.DBInfo.Username = username
		}

		if dbname := r.FormValue("dbname"); dbname != "" {
			p.app.DBInfo.DBName = dbname
		}

		if r.FormValue("password") != "" {
			p.app.DBInfo.Password = r.FormValue("password")
		}

		if !p.app.DBInfo.IsValid() {
			renderError(w, p.app, http.StatusBadRequest, fmt.Errorf("missing required fields"))
			return
		}

		pg, err := Connect(p.app.DBInfo)
		if err != nil {
			renderError(w, p.app, http.StatusInternalServerError, err)
			return
		}
		defer pg.Close()

		v, err := RetrievePGVersion(r.Context(), pg)
		if err != nil {
			renderError(w, p.app, http.StatusInternalServerError, err)
			return
		}
		p.app.DBInfo.Version = v.Version

		http.Redirect(w, r, "/db", http.StatusSeeOther)
	}
}
