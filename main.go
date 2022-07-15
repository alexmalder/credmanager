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
    value, err := src.EncTest(secret.Value)
    if err != nil {
        log.Fatal(err)
    }
    secret.Value = value
	log.Println(secret)
    secret.Migrate()
	switch {
	case secret.Scope == src.ScopeCreate:
		log.Println(src.ScopeCreate)
        secret.Save()
	case secret.Scope == src.ScopeCreateFile:
		log.Println(src.ScopeCreateFile)
	case secret.Scope == src.ScopeGet:
		log.Println(src.ScopeGet)
        secret.List()
	case secret.Scope == src.ScopePut:
		log.Println(src.ScopePut)
	case secret.Scope == src.ScopeDelete:
		log.Println(src.ScopeDelete)
	default:
		log.Println("Scope is not defined")
	}
}
