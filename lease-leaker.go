package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	leases "github.com/npotts/go-dhcpd-leases"
)

type rleases struct {
	Error  error          `json:"error"`
	Leases []leases.Lease `json:"leases"`
}

type app struct {
	port      *int
	leasefile *string
}

func (a *app) leases() rleases {
	f, e := os.Open(*a.leasefile)
	if e != nil {
		return rleases{Error: e, Leases: nil}
	}
	defer f.Close()
	return rleases{Error: nil, Leases: leases.Parse(f)}
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.RequestURI == "/json" {
		en := json.NewEncoder(w)
		en.Encode(a.leases())
		return
	}
	e := tmpl.Execute(w, a.leases())
	if e != nil {
		log.Print(e)
	}
}

func (a *app) Serve() error {
	http.Handle("/", a)
	http.Handle("/json", a)
	return http.ListenAndServe(fmt.Sprintf(":%d", *a.port), a)
}

var (
	daApp = app{
		port:      flag.Int("port", 3000, "HTTP Listening port"),
		leasefile: flag.String("lease-file", "/var/dhcpd/var/db/dhcpd.leases", "Path to dhcpd's  dhcpd.leases file"),
	}

	tmpl = template.Must(template.New("tabular").Parse(`
<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="dhcpd leases">
	<title>dhcpd's leases</title>
	<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/pure-min.css" integrity="sha384-nn4HPE8lTHyVtfCBi5yW9d20FjT8BJwUXyWZT9InLYax14RDjBj46LmSztkmNP9w" crossorigin="anonymous">
</head>
<body>


<div id="layout" class="pure-g">
	<div class="pure-u-1-24"></div>
	<div class="pure-u-22-24">{{if .Error}}<p>Read error: {{.Error}}</p>{{else}}
		<span>{{len .Leases}} client lease(s) issued</span>
		<hr />
		<table class="pure-table pure-table-bordered" style="width:100%;">
			<thead>
				<tr>
					<th>#</th>
					<th>Client</th>
					<th>IP Address</th>
					<th>MAC Addr</th>
					<th>Lease Start</th>
				</tr>
			</thead>
			<tbody>{{range $i, $elem :=  .Leases}}
				<tr>
					<td>{{$i}}</td>
					<td>{{$elem.ClientHostname}}</td>
					<td>{{$elem.IP}}</td>
					<td>{{$elem.Hardware.MAC}}</td>
					<td>{{$elem.Starts}}</td>
				</tr>
			{{end}}
			</tbody>
		</table>
	</div>
	<div class="pure-u-1-24"></div>
	{{end}}
</div>
</body>
</html>
`))
)

func main() {
	flag.Parse()
	daApp.Serve()
}
