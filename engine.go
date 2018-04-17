package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/oxtoacart/bpool"
	"golang.org/x/crypto/acme/autocert"
)

type Data map[string]interface{}

var templates map[string]*template.Template
var bufpool *bpool.BufferPool
var mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`
var templateConfig TemplateConfig

type TemplateConfig struct {
	TemplateLayoutPath  string
	TemplateIncludePath string
}

func loadLetsencript() *http.Server {
	if os.Getenv("ENV") == "prod" {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("secret-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("suncork.net"),
		}
		go http.ListenAndServe(":http", m.HTTPHandler(nil))
		s := &http.Server{
			Addr:      ":https",
			TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		}

		return s
	}
	s := &http.Server{
		Addr: "127.0.0.1:1234",
	}

	return s
}

func loadConfiguration() {
	templateConfig.TemplateLayoutPath = "views/layout/"
	templateConfig.TemplateIncludePath = "views/partials/"
}

func loadTemplates() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	layoutFiles, e := filepath.Glob(templateConfig.TemplateLayoutPath + "*.html")
	crash(e)

	includeFiles, e := filepath.Glob(templateConfig.TemplateIncludePath + "*.html")
	crash(e)

	mainTemplate := template.New("main")

	mainTemplate, e = mainTemplate.Parse(mainTmpl)
	crash(e)
	for _, file := range includeFiles {
		fileName := filepath.Base(file)
		files := append(layoutFiles, file)
		templates[fileName], e = mainTemplate.Clone()
		crash(e)
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}

	print("templates loading successful")

	bufpool = bpool.NewBufferPool(64)
	print("buffer allocation successful")
}

func renderTemplate(w http.ResponseWriter, name string, data Data, s *Session) {
	print("Rendering process start -> Populating environment")
	translations := getTrans(s.lang)
	data["t"] = translations
	data["lang"] = s.lang
	data["env"] = os.Getenv("ENV")
	data["admin"] = s.admin
	data["config"] = getConfig()
	print("All environment vars assigned -> Finding template")
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("Server error (templating %s)", name),
			http.StatusInternalServerError)
	}
	print("Got template -> Executing template")

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	print("Template executed -> Serving template")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
