package main

import (
	"log"
	"main/src"
	"os"
)

func test() {
	encStr, err := src.EncTest(src.FileAsString("config.yml"))
	if err != nil {
		log.Fatal(err)
	}
	decStr, err := src.DecTest(encStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Decrypted Secret: \n%s", decStr)
}

func main() {
	target := os.Args[1]
	if target == "server" {
		err := src.ZMQServer()
		if err != nil {
			log.Fatal(err)
		}
	} else if target == "client" {
		err := src.ZMQClient("ok")
		if err != nil {
			log.Fatal(err)
		}
	} else if target == "load" {
		//log.Println(src.FileAsString(os.Args[2]))
		err := src.ZMQClient(src.FileAsString(os.Args[2]))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		test()
	}
}
