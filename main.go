package main

import (
	"context"
	"log"
	"main/src"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx := src.Getopts()
	connection, err := pgxpool.Connect(context.Background(), src.ConnectionString())
	if err != nil {
		log.Fatal("pgx.Connect", err)
	}
	ctx.Pool = connection
	ctx.Conf = src.ReadConfig()
	ctx.Migrate()
	switch {
    case ctx.Scope == src.ScopeImportBitwarden:
        ctx.ImportBitwarden()
	case ctx.Scope == src.ScopeCreateValue:
		ctx.SaveValue()
	case ctx.Scope == src.ScopeCreateFile:
		ctx.SaveFile()
	case ctx.Scope == src.ScopeSelect:
		ctx.Select()
	case ctx.Scope == src.ScopeGet:
		ctx.Get()
	case ctx.Scope == src.ScopePutValue:
		ctx.UpdateValue()
	case ctx.Scope == src.ScopePutFile:
		ctx.UpdateFile()
	case ctx.Scope == src.ScopeDelete:
		ctx.Remove()
	case ctx.Scope == src.ScopeDrop:
		ctx.Drop()
	default:
		log.Println("Scope is not defined")
	}
}
