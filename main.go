package main

import (
	"context"
	"log"
	"main/src"

	"github.com/jackc/pgx/v4"
)

func main() {
	secret := src.Getopts()
	connection, err := pgx.Connect(context.Background(), src.ConnectionString())
	if err != nil {
		log.Fatal("pgx.Connect", err)
	}
	secret.Conn = connection
	secret.Conf = src.ReadConfig()
	switch {
	case secret.Scope == src.ScopeCreate:
		//log.Println(src.ScopeCreate)
		secret.Save()
	case secret.Scope == src.ScopeCreateFile:
		//log.Println(src.ScopeCreateFile)
        secret.SaveFile()
	case secret.Scope == src.ScopeSelect:
		secret.Select()
	case secret.Scope == src.ScopeGet:
		//log.Println(src.ScopeGet)
		secret.Get()
	case secret.Scope == src.ScopePut:
		//log.Println(src.ScopePut)
        secret.Update()
	case secret.Scope == src.ScopeDelete:
		//log.Println(src.ScopeDelete)
	default:
		log.Println("Scope is not defined")
	}
}
