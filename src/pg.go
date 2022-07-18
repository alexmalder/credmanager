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
	CliSecret Secret
	DbSecret  Secret
	Scope     string
	Filepath  string
	Conf      Config
	Pool      *pgxpool.Pool
}

var ctx = context.Background()

// print the current secret helper function
func (s *Secret) Print() {
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
func (s *SecretCtx) Migrate() {
	for _, v := range s.Conf.Queries {
		_, err := s.Pool.Exec(context.Background(), v.Query)
		//log.Println(v)
		checkErr(err)
	}
}

// save the current secret by value
func (s *SecretCtx) Save(recordType string) {
	encValue := EncryptString(s.CliSecret.Value)
	encNotes := EncryptString(s.CliSecret.Notes)
	s.CliSecret.Type = recordType
	_, err := s.Pool.Exec(
		ctx,
		s.Conf.InsertSecret,
		s.CliSecret.Key,
		1,
		encValue,
		s.CliSecret.Username,
		s.CliSecret.Uri,
		encNotes,
		s.CliSecret.Type,
	)
	//log.Println(s.CliSecret.Key)
	checkErr(err)
	s.WriteRevision()
}

// select secret by key
func (s *SecretCtx) Select() {
	var secrets []*Secret
	err := pgxscan.Select(ctx, s.Pool, &secrets, s.Conf.SelectSecrets)
	for _, v := range secrets {
		fmt.Printf("- [ %s ]\n", v.Key)
	}
	fmt.Printf("- Total items: [ %v ]\n", len(secrets))
	checkErr(err)
}

// select secret by key
func (s *SecretCtx) Get() {
	err := pgxscan.Get(ctx, s.Pool, &s.DbSecret, s.Conf.SelectSecret, s.CliSecret.Key)
	checkErr(err)
	s.DbSecret.Value = DecryptString(s.DbSecret.Value)
	s.DbSecret.Notes = DecryptString(s.DbSecret.Notes)
	fmt.Printf(
		"- [ %s, %s, %s, %s, %s, %v ]\n",
		s.DbSecret.Key,
		s.DbSecret.Value,
		s.DbSecret.Username,
		s.DbSecret.Uri,
		s.DbSecret.Notes,
		s.DbSecret.Revision,
	)
}

// put secret value by key
func (s *SecretCtx) Update() {
	s.Get()
	log.Printf("Change secret revision from [%v] to [%v]\n", s.DbSecret.Revision, s.DbSecret.Revision+1)
	log.Println("Notes ", s.CliSecret.Notes, s.DbSecret.Notes)
	if s.CliSecret.Key != s.DbSecret.Key {
		log.Printf("Changed key from [%s] to [%s]\n", s.DbSecret.Key, s.CliSecret.Key)
		s.DbSecret.Key = s.CliSecret.Key
	}

	if s.CliSecret.Value != "" && s.CliSecret.Value != s.DbSecret.Value {
		log.Printf("Changed value from [%s] to [%s]\n", s.DbSecret.Value, s.CliSecret.Value)
		s.DbSecret.Value = s.CliSecret.Value
	}

	if s.CliSecret.Username != "" && s.CliSecret.Username != s.DbSecret.Username {
		log.Printf("Changed username from [%s] to [%s]\n", s.DbSecret.Username, s.CliSecret.Username)
		s.DbSecret.Username = s.CliSecret.Username
	}

	if s.CliSecret.Uri != "" && s.CliSecret.Uri != s.DbSecret.Uri {
		log.Printf("Changed uri from [%s] to [%s]\n", s.DbSecret.Uri, s.CliSecret.Uri)
		s.DbSecret.Uri = s.CliSecret.Uri
	}

	if s.CliSecret.Notes != "" && s.CliSecret.Notes != s.DbSecret.Notes {
		log.Printf("Changed notes from [%s] to [%s]\n", s.DbSecret.Notes, s.CliSecret.Notes)
		s.DbSecret.Notes = s.CliSecret.Notes
	}

	if s.CliSecret.IsDeleted != s.DbSecret.IsDeleted {
		log.Printf("Changed is_deleted from [%t] to [%t]\n", s.DbSecret.IsDeleted, s.CliSecret.IsDeleted)
		s.DbSecret.IsDeleted = s.CliSecret.IsDeleted
	}

	// get encrypted data
	encValue := EncryptString(s.DbSecret.Value)
	encNotes := EncryptString(s.DbSecret.Notes)
	s.DbSecret.Revision += 1
	_, err := s.Pool.Exec(
		ctx,
		s.Conf.UpdateSecret,
		s.DbSecret.Key,
		s.DbSecret.Revision,
		encValue,
		s.DbSecret.Username,
		s.DbSecret.Uri,
		encNotes,
		s.DbSecret.IsDeleted,
	)
	checkErr(err)
	s.WriteRevision()
}

func (s *SecretCtx) WriteRevision() {
	log.Println("s.DbSecret ", s.DbSecret)
	// get encrypted data
	encValue := EncryptString(s.DbSecret.Value)
	encNotes := EncryptString(s.DbSecret.Notes)
	_, err := s.Pool.Exec(
		ctx,
		s.Conf.InsertRevision,
		s.DbSecret.Key,
		s.DbSecret.Revision,
		encValue,
		s.DbSecret.Username,
		s.DbSecret.Uri,
		encNotes,
		s.DbSecret.Type,
		s.DbSecret.IsDeleted,
	)
	checkErr(err)
}

// drop secrets table
func (s *SecretCtx) Drop() {
	for _, v := range s.Conf.Drops {
		_, err := s.Pool.Exec(ctx, v.Query)
		checkErr(err)
	}
}

func (s *SecretCtx) ImportBitwarden() {
	payload := ReadJson(s.Filepath)
	for _, v := range payload {
		s.CliSecret.Key = v.Name
		s.CliSecret.Value = v.Login.Password
		s.CliSecret.Username = v.Login.Username
		s.CliSecret.Notes = v.Notes
		//fmt.Printf("%s %s %s ", v.Name, v.Login.Username, v.Login.Password)
		if len(v.Login.Uris) == 1 {
			for _, u := range v.Login.Uris {
				//log.Printf("%s", u.Uri)
				s.CliSecret.Uri = u.Uri
			}
		}
		s.Save(TypeValue)
	}
}
