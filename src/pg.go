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
func (s SecretCtx) Save(recordType string) {
	// get encrypted data
	encValue := EncryptString(s.CliSecret.Value)
	var encNotes string
	log.Println(s.CliSecret.Notes)
	if s.CliSecret.Notes != "" {
		encNotes = EncryptString(s.CliSecret.Notes)
	}
	s.CliSecret.Revision = 1
	_, err := s.Pool.Exec(
		ctx,
		s.Conf.InsertSecret,
		s.CliSecret.Key,
		s.CliSecret.Revision,
		encValue,
		s.CliSecret.Username,
		s.CliSecret.Uri,
		encNotes,
		recordType,
	)
	log.Println(s.CliSecret.Key)
	checkErr(err)
}

// select secret by key
func (s SecretCtx) Select() {
	var secrets []*Secret
	err := pgxscan.Select(ctx, s.Pool, &secrets, s.Conf.SelectSecrets)
	for _, v := range secrets {
		fmt.Printf("- [ %s ]\n", v.Key)
	}
	fmt.Printf("- Total items: [ %v ]\n", len(secrets))
	checkErr(err)
}

// select secret by key
func (s SecretCtx) Get() Secret {
	var dbSecret Secret
	err := pgxscan.Get(ctx, s.Pool, &dbSecret, s.Conf.SelectSecret, s.CliSecret.Key)
	checkErr(err)
	dbSecret.Value = DecryptString(dbSecret.Value)
	if dbSecret.Notes != "" {
		dbSecret.Notes = DecryptString(dbSecret.Notes)
	}
	fmt.Printf("- [ %s, %s, %s, %s, %s ]\n", dbSecret.Key, dbSecret.Value, dbSecret.Username, dbSecret.Uri, dbSecret.Notes)
	return dbSecret
}

// put secret value by key
func (s SecretCtx) Update() {
	dbSecret := s.Get()
	//fmt.Print(dbSecret)
	log.Printf("Change secret revision from [%v] to [%v]\n", dbSecret.Revision, dbSecret.Revision+1)
	switch {
	case s.CliSecret.Key != dbSecret.Key:
		log.Printf("Changed key from [%s] to [%s]\n", dbSecret.Key, s.CliSecret.Key)
		dbSecret.Key = s.CliSecret.Key

	case s.CliSecret.Value != "" && s.CliSecret.Value != dbSecret.Value:
		log.Printf("Changed value from [%s] to [%s]\n", dbSecret.Value, s.CliSecret.Value)
		dbSecret.Value = s.CliSecret.Value

	case s.CliSecret.Username != "" && s.CliSecret.Username != dbSecret.Username:
		log.Printf("Changed username from [%s] to [%s]\n", dbSecret.Username, s.CliSecret.Username)
		dbSecret.Username = s.CliSecret.Username

	case s.CliSecret.Uri != "" && s.CliSecret.Uri != dbSecret.Uri:
		log.Printf("Changed uri from [%s] to [%s]\n", dbSecret.Uri, s.CliSecret.Uri)
		dbSecret.Uri = s.CliSecret.Uri

	case s.CliSecret.Notes != "" && s.CliSecret.Notes != dbSecret.Notes:
		log.Printf("Changed notes from [%s] to [%s]\n", dbSecret.Notes, s.CliSecret.Notes)
		dbSecret.Notes = s.CliSecret.Notes

	case s.CliSecret.IsDeleted != dbSecret.IsDeleted:
		log.Printf("Changed is_deleted from [%t] to [%t]\n", dbSecret.IsDeleted, s.CliSecret.IsDeleted)
		dbSecret.IsDeleted = s.CliSecret.IsDeleted

	default:
		log.Println("No changes")
	}
	s.CliSecret = dbSecret
	// get encrypted data
	encValue := EncryptString(s.CliSecret.Value)
	var encNotes string
	log.Println(s.CliSecret.Notes)
	if s.CliSecret.Notes != "" {
		encNotes = EncryptString(s.CliSecret.Notes)
	}
	s.CliSecret.Revision = 1
	_, err := s.Pool.Exec(
		ctx,
		s.Conf.UpdateSecret,
		s.CliSecret.Key,
		s.CliSecret.Revision,
		encValue,
		s.CliSecret.Username,
		s.CliSecret.Uri,
		encNotes,
		s.CliSecret.IsDeleted,
	)
	log.Println(s.CliSecret.Key)
	checkErr(err)
	//s.WriteRevision()
}

func (s SecretCtx) WriteRevision() {

}

// drop secrets table
func (s SecretCtx) Drop() {
	_, err := s.Pool.Exec(ctx, s.Conf.DropSecrets)
	checkErr(err)
}

func (s SecretCtx) ImportBitwarden() {
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
