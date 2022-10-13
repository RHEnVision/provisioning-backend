package services

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/version"
)

const welcomeTmpl = `<!DOCTYPE html>
<html>
	<head>
		<title>Provisioning backend</title>
		<meta charset="utf-8"/>
	</head>
	<body>
		<h1>Provisioning backend API {{ .APIVersion }} {{ .BuildCommit }}</h1>
		<ul>
			<li><a href="/docs">OpenAPI documentation</a></li>
			<li><a href="/ping">Ping service</a> (identity not needed)</li>
			<li><a href="/api/provisioning/{{ .APIVersion }}/openapi.json">OpenAPI JSON</a></li>
			<li><a href="/api/provisioning/{{ .APIVersion }}/ready">Ready service</a> (identity needed)</li>
		</ul>
		<p>Built at {{ .BuildTime }}</p>
	</body>
</html>
`

type Vars struct {
	APIVersion  string
	BuildCommit string
	BuildTime   string
}

func WelcomeService(w http.ResponseWriter, r *http.Request) {
	vars := Vars{
		APIVersion:  version.APIPathVersion,
		BuildCommit: version.BuildCommit,
		BuildTime:   version.BuildTime,
	}

	tmpl := template.Must(template.New("welcome").Parse(welcomeTmpl))
	buf := bytes.NewBuffer(nil)
	_ = tmpl.Execute(buf, vars)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
