package proxy

import (
	"fmt"
	"log"
	"net/http"
	"sort"
)

type DBNameData struct {
	DBInfo *DBInfo

	Names []string
}

const dbNamesTemplate = `
<div>
	<p>{{ .Data.DBInfo.Version }}
</div>
<div>
	<p>Host: {{ .Data.DBInfo.Host }}</p>
	<p>Port: {{ .Data.DBInfo.Port }}</p>
	<p>Username: {{ .Data.DBInfo.Username }}</p>
	<p>Database: {{ .Data.DBInfo.DBName }}</p>
</div>
<div>
	<table>
		<tr>
			<th>Database</th>
		</tr>
		{{ range $index, $name := .Data.Names }}
		<tr>
			<td><a href="/db/{{ $name }}" target="_blank">{{ $name }}</a></td>
		</tr>
		{{ end }}
	</table>
</div>
`

type DBPage struct {
	app *AppConfig
}

func NewDBPage(c *WebConfig) *DBPage {
	return &DBPage{
		app: c.App,
	}
}

func (p DBPage) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DBPage.Handler: %s\n", r.URL.Path)

		pg, err := Connect(p.app.DBInfo)
		if err != nil {
			renderError(w, p.app, http.StatusInternalServerError, err)
			return
		}
		defer pg.Close()

		dbs, err := RetrieveDatabases(r.Context(), pg)
		if err != nil {
			renderError(w, p.app, http.StatusInternalServerError, err)
			return
		}

		names := []string{}
		for _, db := range dbs {
			names = append(names, db.Datname)
		}

		if len(names) == 0 {
			renderError(w, p.app, http.StatusNotFound, fmt.Errorf("no database: %q", p.app.DBInfo.Username))
			return
		}

		sort.Strings(names)

		data := &DBNameData{
			DBInfo: p.app.DBInfo,
			Names:  names,
		}
		renderTemplate(w, p.app, dbNamesTemplate, p.app.DBInfo.Username, data)
	}
}
