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
	default:
		log.Println("Scope is not defined")
	}
}
