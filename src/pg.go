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
	// base fields
	Key      string
	Value    string
	Username string
	Uri      string
	Notes    string
	// helper fields
	Scope    string
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
		checkErr(err)
	}
}

// save the current secret by value
func (s Secret) SaveValue() {
	fmt.Println(s.Conf.InsertSecret, s.Key, s.Value, s.Username, s.Uri, s.Notes)
	// get encrypted data
	encValue := EncryptString(s.Value)
	_, err := s.Conn.Exec(
		context.Background(),
		s.Conf.InsertSecret,
		s.Key,
		encValue,
		s.Username,
		s.Uri,
		s.Notes,
        "value",
	)
    log.Println(s.Key)
	checkErr(err)
}

// save the current secret by filename
func (s Secret) SaveFile() {
	rawData := FileAsString(s.Filepath)
	encValue := EncryptString(rawData)
	_, err := s.Conn.Exec(
        context.Background(), 
        s.Conf.InsertSecret, 
		s.Key,
		encValue,
		s.Username,
		s.Uri,
		s.Notes,
        "file",
    )
	checkErr(err)
}

// select secret by key
func (s Secret) Get() error {
	rows, _ := s.Conn.Query(context.Background(), s.Conf.SelectSecret, s.Key)
	for rows.Next() {
		err := rows.Scan(&s.Key, &s.Value, &s.Username, &s.Uri, &s.Notes)
		if err != nil {
			return err
		}
		decStr := DecryptString(s.Value)
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
		err := rows.Scan(&s.Key, &s.Value, &s.Username, &s.Uri, &s.Notes)
		checkErr(err)
		fmt.Printf("Key: %s\n", s.Key)
		fmt.Printf("Value: %s...\n", s.Value[0:15])
		fmt.Printf("---\n")
	}
	return rows.Err()
}

// put secret value by key
func (s Secret) UpdateValue() error {
	encValue := EncryptString(s.Value)
	_, err := s.Conn.Exec(
        context.Background(), 
        s.Conf.DeleteSecret, 
		s.Key,
		encValue,
		s.Username,
		s.Uri,
		s.Notes,
    )
	return err
}

// put secret file by key
func (s Secret) UpdateFile() error {
	rawData := FileAsString(s.Filepath)
	encValue := EncryptString(rawData)
	_, err := s.Conn.Exec(
        context.Background(), 
        s.Conf.DeleteSecret, 
		s.Key,
		encValue,
		s.Username,
		s.Uri,
		s.Notes,
    )
	return err
}

// remove secret by key
func (s Secret) Remove() {
	_, err := s.Conn.Exec(context.Background(), s.Conf.DeleteSecret, s.Key)
	checkErr(err)
}

// drop secrets table
func (s Secret) Drop() {
	_, err := s.Conn.Exec(context.Background(), s.Conf.DropSecrets)
	checkErr(err)
}
