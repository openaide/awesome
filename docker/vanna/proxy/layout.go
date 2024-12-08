package proxy

import (
	"fmt"
	"html/template"
	"net/http"
)

type Crumb struct {
	Title string
	Path  string
}

type PageData struct {
	App *AppConfig

	Title  string
	Crumbs []*Crumb
	Data   interface{}
}

const layout = `
<!DOCTYPE html>
<html>
<head>
<style>
:root {
}

html {
    overflow-y: scroll;
}

body {
    margin: 0px auto 200px;
}

header {
}

footer {
	text-align: center;
}

a {
	color: #0275d8;
	text-decoration: none;
}
a:hover {
	color: #01447e;
	text-decoration: underline;
}

div.container {
	margin-left: 1em;
	margin-right: 1em;
    height: auto; 
}

/* breadcrumb */
ul.breadcrumb {
	padding: 10px 16px;
	list-style: none;
	background-color: #eee;
}
ul.breadcrumb li {
	display: inline;
	font-size: 18px;
}
ul.breadcrumb li+li:before {
	padding: 8px;
	color: black;
	content: "/";
}

/* table */
table {
	font-family: arial, sans-serif;
	border-collapse: collapse;
	width: 100%%;
}
td, th {
	border: 1px solid #dddddd;
	text-align: left;
	padding: 8px;
}
tr:nth-child(even) {
	background-color: #dddddd;
}

dl {
	display: grid;
	grid-template-columns: max-content auto;
}
dt {
	grid-column-start: 1;
}
dd {
	grid-column-start: 2;
}
</style>
<meta charSet="UTF-8"/>
<link href="/favicon.ico" rel="icon"/>
<title>{{ .Title }}</title>
</head>
<body>
	<header>
		<nav>
			<ul class="breadcrumb">
				<li><a href="/home">{{ .App.Name }} {{ .App.Version }}</a></li>
				{{ range .Crumbs }}
					<li><a href="{{ .Path }}">{{ .Title }}</a></li>
				{{ end }}
				{{ if ne .Title "Home" }}
					<li><a href="#">{{ .Title }}</a></li>
				{{ end }}
			</ul>
		</nav>
	</header>
	<div class="container">
		<h3>{{ .Title }}</h3>
		<br />
		<div>%s</div>
	</div>
	<footer>
		<p>© 2024. All Rights Reserved.</p>
	</footer>
</body>
</html>
`

func renderTemplate(w http.ResponseWriter, app *AppConfig, tpl string, title string, data interface{}, crumbs ...*Crumb) {
	t, err := template.New("page").Parse(fmt.Sprintf(layout, tpl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = t.Execute(w, &PageData{
		App:    app,
		Title:  title,
		Data:   data,
		Crumbs: crumbs,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ErrorData struct {
	Error error
}

const errorContent = `
<p>An error occurred:</p>
<pre>{{ .Data.Error }}</pre>
`

func renderError(w http.ResponseWriter, app *AppConfig, code int, err error) {
	w.WriteHeader(code)
	renderTemplate(w, app, errorContent, fmt.Sprintf("Error %d", code), &ErrorData{
		Error: err,
	})
}
