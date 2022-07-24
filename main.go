package main

import (
	"bytes"
	//"context"
	"fmt"
	"log"
	"main/src"
	"strings"

	"html/template"
	"net/http"

	//"github.com/jackc/pgx/v4/pgxpool"
)

func serve() {
	ctx := src.Getopts()
    ctx.Init()
    data := ctx.Select()
	t, err := template.ParseFiles("templates/hello.html")
	if err != nil {
		log.Fatal(err)
	}
	var b bytes.Buffer
	err = t.Execute(&b, data)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/provision", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, b.String())
	})
	http.Handle("/static/", http.StripPrefix(strings.TrimRight("/static/", "/"), http.FileServer(http.Dir("static"))))

	// listen to port
	http.ListenAndServe(":5050", nil)

}

func main() {
	ctx := src.Getopts()
    ctx.Init()
	switch {
	case ctx.Scope == src.ScopeMigrate:
		ctx.Migrate()
	case ctx.Scope == src.ScopeImportBitwarden:
		ctx.ImportBitwarden()
	case ctx.Scope == src.ScopeCreateValue:
		ctx.Save(src.TypeValue)
	case ctx.Scope == src.ScopeCreateFile:
		ctx.CliSecret.Value = src.FileAsString(ctx.Filepath)
		ctx.Save(src.TypeFile)
	case ctx.Scope == src.ScopeSelect:
		ctx.Select()
	case ctx.Scope == src.ScopeGet:
		ctx.Get()
	case ctx.Scope == src.ScopePutValue:
		ctx.Update()
	case ctx.Scope == src.ScopePutFile:
		ctx.CliSecret.Value = src.FileAsString(ctx.Filepath)
		ctx.Update()
	case ctx.Scope == src.ScopeDrop:
		ctx.Drop()
	case ctx.Scope == "serve":
		serve()
	default:
		log.Println("Scope is not defined")
	}
}
