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
	Scope    string
	Key      string
	Value    string
	Filepath string
	Conf     Config
	Conn     *pgx.Conn
}

// print the current secret helper function
func (s Secret) Print() {
	fmt.Println(s)
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

// make database migrations
func (s Secret) Migrate() {
	for _, v := range s.Conf.Queries {
		_, err := s.Conn.Exec(context.Background(), v.Query)
		//log.Println(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// save the current secret by value
func (s Secret) Save() {
	// log.Println(s.Conf.InsertSecret, s.Key, s.Value)
	// get encrypted data
	encValue, err := EncTest(s.Value)
	if err != nil {
		log.Fatal(err)
	}
	_, err = s.Conn.Exec(context.Background(), s.Conf.InsertSecret, s.Key, encValue)
	if err != nil {
		log.Fatal("Save ", err)
	}
}

// save the current secret by filename
func (s Secret) SaveFile() {
	rawData := FileAsString(s.Filepath)
	encValue, err := EncTest(rawData)
	if err != nil {
		log.Fatal(err)
	}
	_, err = s.Conn.Exec(context.Background(), s.Conf.InsertSecret, s.Key, encValue)
	if err != nil {
		log.Fatal("Save ", err)
	}
}

// select secret by key
func (s Secret) Get() error {
	rows, _ := s.Conn.Query(context.Background(), s.Conf.SelectSecret, s.Key)
	for rows.Next() {
		err := rows.Scan(&s.Key, &s.Value)
		if err != nil {
			return err
		}
		decStr, err := DecTest(s.Value)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Key: %s\n", s.Key)
		fmt.Printf("Value: %s\n", decStr)
		fmt.Printf("---\n")
	}
	return rows.Err()
}

// select secret by key
func (s Secret) Select() error {
	rows, _ := s.Conn.Query(context.Background(), s.Conf.SelectSecrets)
	for rows.Next() {
		err := rows.Scan(&s.Key, &s.Value)
		if err != nil {
			return err
		}
		fmt.Printf("Key: %s\n", s.Key)
		fmt.Printf("Value: %s...\n", s.Value[0:15])
		fmt.Printf("---\n")
	}
	return rows.Err()
}

// put secret value by key
func (s Secret) Update() error {
	_, err := s.Conn.Exec(context.Background(), s.Conf.DeleteSecret, s.Key, s.Value)
	return err
}

// remove secret by key
func (s Secret) Remove() {
	_, err := s.Conn.Exec(context.Background(), s.Conf.DeleteSecret, s.Key)
	if err != nil {
		log.Fatal(err)
	}
}
