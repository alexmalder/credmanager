package src

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v4"
	"gopkg.in/yaml.v3"
)

type Secret struct {
	ID    int
	Key   string
	Value string
	conn  *pgx.Conn
}

type Config struct {
	Queries []struct {
		Query string
	}
}

func readConfig() Config {
	yfile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	data := Config{}
	err2 := yaml.Unmarshal(yfile, &data)
	if err2 != nil {
		log.Fatal(err2)
	}
	//for _, v := range data.Queries {fmt.Printf("%s\n", v.Query)}
	return data
}

func (s Secret) Print() {
	fmt.Println(s)
}

func (s Secret) Migrate(config Config) {
	for _, v := range config.Queries {
		_, err := s.conn.Exec(context.Background(), v.Query)
		log.Println(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s Secret) Save() {
	_, err := s.conn.Exec(context.Background(), "insert into secrets (key, value) values($1, $2)", s.Key, s.Value)
	if err != nil {
		log.Fatal(err)
	}
}

func (s Secret) Remove() {
	_, err := s.conn.Exec(context.Background(), "delete from secrets where id=$1", s.ID)
	if err != nil {
		log.Fatal(err)
	}
}

func ConnString() string {
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

func PgTest() {
	conn, err := pgx.Connect(context.Background(), ConnString())
	if err != nil {
		log.Fatal("pgx.Connect", err)
	}
	rawData := FileAsString("test.yml")
	encStr, err := EncTest(rawData)
	if err != nil {
		log.Fatal(err)
	}
	Secret := Secret{1, "postgres_user", encStr, conn}
	Secret.Print()
	config := readConfig()
	Secret.Migrate(config)
	Secret.Save()
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
