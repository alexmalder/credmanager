package src

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

// secret in database data structure
type Secret struct {
	Key      string
	Value    string
	Username string
	Uri      string
	Notes    string
}

type SecretCtx struct {
	Secret
	Scope    string
	Filepath string
	Conf     Config
	Pool     *pgxpool.Pool
}

var ctx = context.Background()

// print the current secret helper function
func (s Secret) Print() {
	fmt.Println(s)
}

// get connection string for postgresql database
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
func (s SecretCtx) Migrate() {
	for _, v := range s.Conf.Queries {
		_, err := s.Pool.Exec(context.Background(), v.Query)
		//log.Println(v)
		checkErr(err)
	}
}

// save the current secret by value
func (s SecretCtx) SaveValue() {
	fmt.Println(s.Conf.InsertSecret, s.Key, s.Value, s.Username, s.Uri, s.Notes)
	// get encrypted data
	encValue := EncryptString(s.Value)
	_, err := s.Pool.Exec(ctx, s.Conf.InsertSecret, s.Key, encValue, s.Username, s.Uri, s.Notes, TypeValue)
	log.Println(s.Key)
	checkErr(err)
}

// save the current secret by filename
func (s SecretCtx) SaveFile() {
	rawData := FileAsString(s.Filepath)
	encValue := EncryptString(rawData)
	_, err := s.Pool.Exec(ctx, s.Conf.InsertSecret, s.Key, encValue, s.Username, s.Uri, s.Notes, TypeFile)
	checkErr(err)
}

// select secret by key
func (s SecretCtx) Get() error {
	//var secret Secret
	err := pgxscan.Get(ctx, s.Pool, &s, s.Conf.SelectSecret, s.Key)
	checkErr(err)
	decStr := DecryptString(s.Value)
	fmt.Printf("Key: %s\n", s.Key)
	fmt.Printf("Value: %s\n", decStr)
	fmt.Printf("---\n")
	return err
}

// select secret by key
func (s SecretCtx) Select() error {
	var secrets []*Secret
	err := pgxscan.Select(ctx, s.Pool, &secrets, s.Conf.SelectSecrets)
	for _, v := range secrets {
		log.Println(v.Key)
	}
	checkErr(err)
	return err
}

// put secret value by key
func (s SecretCtx) UpdateValue() error {
	encValue := EncryptString(s.Value)
	encNotes := EncryptString(s.Notes)
	_, err := s.Pool.Exec(ctx, s.Conf.DeleteSecret, s.Key, encValue, s.Username, s.Uri, encNotes)
	return err
}

// put secret file by key
func (s SecretCtx) UpdateFile() error {
	encValue := EncryptString(FileAsString(s.Filepath))
	encNotes := EncryptString(s.Notes)
	_, err := s.Pool.Exec(ctx, s.Conf.DeleteSecret, s.Key, encValue, s.Username, s.Uri, encNotes)
	return err
}

// remove secret by key
func (s SecretCtx) Remove() {
	_, err := s.Pool.Exec(ctx, s.Conf.DeleteSecret, s.Key)
	checkErr(err)
}

// drop secrets table
func (s SecretCtx) Drop() {
	_, err := s.Pool.Exec(ctx, s.Conf.DropSecrets)
	checkErr(err)
}

func (s SecretCtx) ImportBitwarden() {
	payload := ReadJson()
	for _, v := range payload {
		s.Key = v.Name
		s.Value = v.Login.Password
		s.Username = v.Login.Username
		s.Notes = EncryptString(v.Notes)
		//fmt.Printf("%s %s %s ", v.Name, v.Login.Username, v.Login.Password)
		if len(v.Login.Uris) == 1 {
			for _, u := range v.Login.Uris {
				//log.Printf("%s", u.Uri)
				s.Uri = u.Uri
			}
		}
		s.SaveValue()
	}
}
