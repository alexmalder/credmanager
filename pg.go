package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v4"
)

type Student struct {
	Name string
	ID   int
	conn *pgx.Conn
}

func (s Student) Print() {
	fmt.Println(s)
}

func (s Student) Save() {
	_, err := s.conn.Exec(context.Background(), "insert into students (name) values($1)", s.Name)
	if err != nil {
		log.Fatal(err)
	}
}

func (s Student) Remove() {
	_, err := s.conn.Exec(context.Background(), "delete from students where id=$1", s.ID)
	if err != nil {
		log.Fatal(err)
	}
}
func connString() string {
	connection := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	return connection
}

func main_() {
	conn, err := pgx.Connect(context.Background(), connString())
	if err != nil {
		log.Fatal("pgx.Connect", err)
	}
	student := Student{"Jack", 1, conn}
	student.Print()
	//student.Save()
	//student.Remove()
}

func FingerPrint() {
	ifaces, err := net.Interfaces()
	if err != nil {
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			hostname, err := os.Hostname()
			log.Println(ip, hostname)
			if err != nil {
			}
		}
	}
}
