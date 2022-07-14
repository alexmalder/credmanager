package main

import (
	"log"
	"main/src"
	"os"
)

func test() {
	rawData := src.FileAsString("test.yml")
	encStr, err := src.EncTest(rawData)
	if err != nil {
		log.Fatal(err)
	}
	decStr, err := src.DecTest(encStr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Decrypted Secret: \n%s", decStr)
	if rawData == decStr {
		log.Println("Test successfull")
	}
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
		//test()
		src.PgTest()
	}
}
