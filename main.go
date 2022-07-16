package main

import (
	"context"
	"log"
	"main/src"

	"github.com/jackc/pgx/v4"
)

func seed(secret src.Secret) {
	payload := src.ReadJson()
	for _, v := range payload {
		secret.Key = v.Name
		secret.Value = v.Login.Password
		secret.Username = v.Login.Username
        secret.Notes = src.EncryptString(v.Notes)
		//fmt.Printf("%s %s %s ", v.Name, v.Login.Username, v.Login.Password)
		if len(v.Login.Uris) == 1 {
			for _, u := range v.Login.Uris {
				//log.Printf("%s", u.Uri)
				secret.Uri = u.Uri
			}
		}
        //secret.SaveValue()
	}
}

func main() {
	secret := src.Getopts()
	connection, err := pgx.Connect(context.Background(), src.ConnectionString())
	if err != nil {
		log.Fatal("pgx.Connect", err)
	}
	secret.Conn = connection
	secret.Conf = src.ReadConfig()
	secret.Migrate()
	switch {
	case secret.Scope == src.ScopeCreateValue:
		secret.SaveValue()
	case secret.Scope == src.ScopeCreateFile:
		secret.SaveFile()
	case secret.Scope == src.ScopeSelect:
		secret.Select()
	case secret.Scope == src.ScopeGet:
		log.Println(secret.Scope, secret.Key)
		secret.Get()
	case secret.Scope == src.ScopePutValue:
		secret.UpdateValue()
	case secret.Scope == src.ScopePutFile:
		secret.UpdateFile()
	case secret.Scope == src.ScopeDelete:
		secret.Remove()
	case secret.Scope == src.ScopeDrop:
		secret.Drop()
	default:
		log.Println("Scope is not defined")
	}
}
