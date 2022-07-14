package src

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v4"
)

type Secret struct {
	ID     int
	Key    string
	Value  string
	config Config
	conn   *pgx.Conn
}

func (s Secret) Print() {
	fmt.Println(s)
}

func (s Secret) Migrate() {
	for _, v := range s.config.Queries {
		_, err := s.conn.Exec(context.Background(), v.Query)
		log.Println(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s Secret) Save() {
	_, err := s.conn.Exec(context.Background(), s.config.InsertSecret, s.Key, s.Value)
	if err != nil {
		log.Fatal(err)
	}
}

func (s Secret) Remove() {
	_, err := s.conn.Exec(context.Background(), s.config.DeleteSecret, s.ID)
	if err != nil {
		log.Fatal(err)
	}
}

func (s Secret) List() error {
	rows, _ := s.conn.Query(context.Background(), s.config.SelectSecrets)
	for rows.Next() {
		err := rows.Scan(&s.ID, &s.Key, &s.Value)
		if err != nil {
			return err
		}
		fmt.Printf("ID: %v\n", s.ID)
		fmt.Printf("Key: %s\n", s.Key)
		fmt.Printf("Value: %s\n", s.Value)
		fmt.Printf("---\n")
	}
	return rows.Err()
}

func (s Secret) Update() error {
	_, err := s.conn.Exec(context.Background(), s.config.DeleteSecret, s.Key, s.Value, s.ID)
	return err
}

func PgTest() {
	connString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	connection, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("pgx.Connect", err)
	}
	rawData := FileAsString("test.yml")
	encStr, err := EncTest(rawData)
	if err != nil {
		log.Fatal(err)
	}
	config := ReadConfig()
	Secret := Secret{1, "postgres_user", encStr, config, connection}
	Secret.Print()
	Secret.Migrate()
	Secret.Save()
	Secret.List()
	//Secret.Remove()
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
