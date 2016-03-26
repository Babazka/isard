package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

func HttpStartServer(addr string) {
	http.HandleFunc("/", HttpView(HttpIndex))
	http.HandleFunc("/plot.png", HttpView(HttpPlot))
	log.Fatal(http.ListenAndServe(addr, nil))
}

func HttpView(handler HttpHandler) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 500)
			}
		}()
		handler(w, r)
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	IndexTemplate = template.Must(template.New("http_index").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
			<title>Monitor</title>
		</head>
		<body>
		<style>
		div.item {
			width: auto;
			float: left;
			margin: 8px;
			padding: 8px;
			background-color: #eee;
		}
		</style>
		<div>
			<a href="/?hours=1">1h</a>
			<a href="/?hours=4">4h</a>
			<a href="/?hours=24">24h</a>
			<a href="/?hours=72">72h</a>
			<a href="/?hours=168">7d</a>
			<a href="/?hours=744">31d</a>
			<a href="/?hours=2232">3m</a>
		</div>
		<hr/>
		{{ range .Items }}
		<div class="item">
			<p><strong>{{ .Name }}</strong></p>
			<img src="/plot.png?item={{ .Name }}&hours={{ $.Hours }}&width=490&height=300" width="490" height="300"/>
		</div>
		{{ end }}
		</body>
	</html>
	`))
)

func HttpIndex(w http.ResponseWriter, r *http.Request) {
	items, err := FindItems()
	panicOnError(err)

	hours := IntOrDefault(r.FormValue("hours"), 24)

	ctx := &struct {
		Hours int
		Items []Item
	}{hours, items}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = IndexTemplate.Execute(w, ctx)
	panicOnError(err)
}

func HttpPlot(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Printf("%s", r.Form)
	itemName := r.Form.Get("item")
	if itemName == "" {
		http.Error(w, "`item` required", 400)
		return
	}

	opts := PlotOptions{}
	opts.Width = IntOrDefault(r.Form.Get("width"), 1400)
	opts.Height = IntOrDefault(r.Form.Get("height"), 600)
	opts.LastNHours = IntOrDefault(r.Form.Get("hours"), 24)
	opts.LastNDays = IntOrDefault(r.Form.Get("days"), 0)

	png, err := Plot(itemName, opts)
	panicOnError(err)

	w.Header().Set("Content-Type", "image/png")
	w.Write(png.Data)
}