package src

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

// secret in database data structure
type Secret struct {
	Scope  string
	Key    string
	Value  string
	File   string
	config Config
	conn   *pgx.Conn
}

// print the current secret helper function
func (s Secret) Print() {
	fmt.Println(s)
}

// make database migrations
func (s Secret) Migrate() {
	for _, v := range s.config.Queries {
		_, err := s.conn.Exec(context.Background(), v.Query)
		log.Println(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// save the current secret
func (s Secret) Save() {
	_, err := s.conn.Exec(context.Background(), s.config.InsertSecret, s.Key, s.Value)
	if err != nil {
		log.Fatal(err)
	}
}

// select secret list by key
func (s Secret) List(key string) error {
	rows, _ := s.conn.Query(context.Background(), s.config.SelectSecret, key)
	for rows.Next() {
		err := rows.Scan(&s.Key, &s.Value)
		if err != nil {
			return err
		}
		fmt.Printf("Key: %s\n", s.Key)
		fmt.Printf("Value: %s\n", s.Value)
		fmt.Printf("---\n")
	}
	return rows.Err()
}

// remove secret by key
func (s Secret) Remove() {
	_, err := s.conn.Exec(context.Background(), s.config.DeleteSecret, s.Key)
	if err != nil {
		log.Fatal(err)
	}
}

// put secret value by key
func (s Secret) Update() error {
	_, err := s.conn.Exec(context.Background(), s.config.DeleteSecret, s.Key, s.Value)
	return err
}

// test postgres functions
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
	Secret := Secret{ScopeCreate, "postgres_user", encStr, "", config, connection}
	Secret.Print()
	Secret.Migrate()
	Secret.Save()
	Secret.List()
	//Secret.Remove()
}
