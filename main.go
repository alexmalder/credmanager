package main

import (
	"log"
	"main/src"
)

func main() {
	opts := src.Getopts()
	log.Println(opts)
	switch {
	case opts.Scope == src.ScopeCreate:
		log.Println(src.ScopeCreate)
	case opts.Scope == src.ScopeCreateFile:
		log.Println(src.ScopeCreateFile)
	case opts.Scope == src.ScopeGet:
		log.Println(src.ScopeGet)
	case opts.Scope == src.ScopePut:
		log.Println(src.ScopePut)
	case opts.Scope == src.ScopeDelete:
		log.Println(src.ScopeDelete)
	default:
		log.Println("Scope is not defined")
	}
}
