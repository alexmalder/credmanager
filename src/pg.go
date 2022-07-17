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
	Key       string
	Revision  int
	Value     string
	Username  string
	Uri       string
	Notes     string
	Type      string
	IsDeleted bool
}

type SecretCtx struct {
	cliSecret Secret
	dbSecret  Secret
	Scope     string
	Filepath  string
	Conf      Config
	Pool      *pgxpool.Pool
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
	// get encrypted data
	encValue := EncryptString(s.cliSecret.Value)
	encNotes := EncryptString(s.cliSecret.Notes)
	s.cliSecret.Revision = 1
	_, err := s.Pool.Exec(
		ctx,
		s.Conf.InsertSecret,
		s.cliSecret.Key,
		s.cliSecret.Revision,
		encValue,
		s.cliSecret.Username,
		s.cliSecret.Uri,
		encNotes,
		TypeValue,
	)
	log.Println(s.cliSecret.Key)
	checkErr(err)
}

// save the current secret by filename
func (s SecretCtx) SaveFile() {
	encValue := EncryptString(FileAsString(s.Filepath))
	encNotes := EncryptString(s.cliSecret.Notes)
	s.cliSecret.Revision = 1
	_, err := s.Pool.Exec(
		ctx,
		s.Conf.InsertSecret,
		s.cliSecret.Key,
		s.cliSecret.Revision,
		encValue,
		s.cliSecret.Username,
		s.cliSecret.Uri,
		encNotes,
		TypeFile,
	)
	checkErr(err)
}

// select secret by key
func (s SecretCtx) Select() error {
	var secrets []*Secret
	err := pgxscan.Select(ctx, s.Pool, &secrets, s.Conf.SelectSecrets)
	for _, v := range secrets {
		fmt.Printf("- [ %s ]\n", v.Key)
	}
	fmt.Printf("- Total items: [ %v ]\n", len(secrets))
	checkErr(err)
	return err
}

// select secret by key
func (s SecretCtx) Get() Secret {
	err := pgxscan.Get(ctx, s.Pool, &s.dbSecret, s.Conf.SelectSecret, s.cliSecret.Key)
	checkErr(err)
	decValue := DecryptString(s.dbSecret.Value)
	decNotes := DecryptString(s.dbSecret.Notes)
	fmt.Printf("- [ %s, %s, %s, %s, %s ]\n", s.dbSecret.Key, decValue, s.dbSecret.Username, s.dbSecret.Uri, decNotes)
	return s.dbSecret
}

// put secret value by key
func (s SecretCtx) UpdateValue() error {
	dbSecret := s.Get()
	//fmt.Print(dbSecret)
	decValue := DecryptString(dbSecret.Value)
	decNotes := DecryptString(dbSecret.Notes)
	log.Printf("Change secret revision from [%v] to [%v]\n", dbSecret.Revision, dbSecret.Revision+1)
	switch {
	case s.cliSecret.Key != dbSecret.Key:
		log.Printf("Change secret key from [%s] to [%s]\n", dbSecret.Key, s.cliSecret.Key)
	case s.cliSecret.Value != "" && s.cliSecret.Value != decValue:
		log.Printf("Change secret value from [%s] to [%s]\n", decValue, s.cliSecret.Value)
	case s.cliSecret.Username != "" && s.cliSecret.Username != dbSecret.Username:
		log.Printf("Change secret username from [%s] to [%s]\n", dbSecret.Username, s.cliSecret.Username)
	case s.cliSecret.Uri != "" && s.cliSecret.Uri != dbSecret.Uri:
		log.Printf("Change secret uri from [%s] to [%s]\n", dbSecret.Uri, s.cliSecret.Uri)
	case s.cliSecret.Notes != "" && s.cliSecret.Notes != decNotes:
		log.Printf("Change secret notes from [%s] to [%s]\n", decNotes, s.cliSecret.Notes)
	/* is not implemented
	case s.cliSecret.Type == dbSecret.Type:
		log.Printf("Change secret Type from [%s] to [%s]\n", dbSecret.Type, s.cliSecret.Type)
	*/
	case s.cliSecret.IsDeleted != dbSecret.IsDeleted:
		log.Printf("Change secret is_deleted from [%t] to [%t]\n", dbSecret.IsDeleted, s.cliSecret.IsDeleted)
	default:
		log.Println("No changes")
	}
	return nil
}

// put secret file by key
func (s SecretCtx) UpdateFile() {
	encValue := EncryptString(FileAsString(s.Filepath))
	encNotes := EncryptString(s.cliSecret.Notes)
	_, err := s.Pool.Exec(
		ctx,
		s.Conf.DeleteSecret,
		s.cliSecret.Key,
		encValue,
		s.cliSecret.Username,
		s.cliSecret.Uri,
		encNotes,
	)
	checkErr(err)
}

// drop secrets table
func (s SecretCtx) Drop() {
	_, err := s.Pool.Exec(ctx, s.Conf.DropSecrets)
	checkErr(err)
}

func (s SecretCtx) ImportBitwarden() {
	payload := ReadJson(s.Filepath)
	for _, v := range payload {
		s.cliSecret.Key = v.Name
		s.cliSecret.Value = v.Login.Password
		s.cliSecret.Username = v.Login.Username
		s.cliSecret.Notes = EncryptString(v.Notes)
		//fmt.Printf("%s %s %s ", v.Name, v.Login.Username, v.Login.Password)
		if len(v.Login.Uris) == 1 {
			for _, u := range v.Login.Uris {
				//log.Printf("%s", u.Uri)
				s.cliSecret.Uri = u.Uri
			}
		}
		s.SaveValue()
	}
}
