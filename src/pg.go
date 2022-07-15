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
	Conf Config
	Conn   *pgx.Conn
}

// print the current secret helper function
func (s Secret) Print() {
	fmt.Println(s)
}

// make database migrations
func (s Secret) Migrate() {
	for _, v := range s.Conf.Queries {
		_, err := s.Conn.Exec(context.Background(), v.Query)
		log.Println(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// save the current secret
func (s Secret) Save() {
    log.Println(s.Conf.InsertSecret, s.Key, s.Value)
	_, err := s.Conn.Exec(context.Background(), s.Conf.InsertSecret, s.Key, s.Value)
	if err != nil {
		log.Fatal("Save ", err)
	}
}

// select secret list by key
func (s Secret) List() error {
	rows, _ := s.Conn.Query(context.Background(), s.Conf.SelectSecret, s.Key)
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
	_, err := s.Conn.Exec(context.Background(), s.Conf.DeleteSecret, s.Key)
	if err != nil {
		log.Fatal(err)
	}
}

// put secret value by key
func (s Secret) Update() error {
	_, err := s.Conn.Exec(context.Background(), s.Conf.DeleteSecret, s.Key, s.Value)
	return err
}

func ConnectionString() string {
    return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
}

// test postgres functions
func PgTest() {
	connection, err := pgx.Connect(context.Background(), ConnectionString())
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
