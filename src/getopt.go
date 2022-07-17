package src

import (
	"fmt"
	"log"
	"os"

	"github.com/kesselborn/go-getopt"
)

const (
	ScopeCreateValue     string = "create-value"
	ScopeCreateFile             = "create-file"
	ScopeSelect                 = "select"
	ScopeGet                    = "get"
	ScopePutValue               = "put-value"
	ScopePutFile                = "put-file"
	ScopeDelete                 = "delete"
	ScopeDrop                   = "drop"
	ScopeImportBitwarden        = "import-bitwarden"
	TypeFile                    = "file"
	TypeValue                   = "value"
)

// get command line arguments
func Getopts() SecretCtx {
	sco := getopt.SubCommandOptions{
		getopt.Options{
			"global description",
			getopt.Definitions{
				{"config|c", "key of a new secret", getopt.Optional | getopt.ExampleIsDefault, "config.yml"},
			},
		},
		getopt.SubCommands{
			ScopeImportBitwarden: {
				"import bitwarden json file",
				getopt.Definitions{
					{"file|f", "file of a new secret", getopt.Required, ""},
				},
			},
			ScopeCreateValue: {
				"create key-value pair as string",
				getopt.Definitions{
					{"key|k", "key in secret", getopt.Required, ""},
					{"value|v", "value in secret", getopt.Required, ""},
					{"username", "username in secret", getopt.Optional | getopt.ExampleIsDefault, ""},
					{"uri", "uri in secret", getopt.Optional | getopt.ExampleIsDefault, ""},
					{"notes|n", "notes of the new secret", getopt.Optional | getopt.ExampleIsDefault, ""},
				},
			},
			ScopeCreateFile: {
				"create key-value pair as file",
				getopt.Definitions{
					{"key|k", "key of a new secret", getopt.Required, ""},
					{"file|f", "file of a new secret", getopt.Required, ""},
					{"notes|n", "notes of the new secret", getopt.Optional | getopt.ExampleIsDefault, ""},
				},
			},
			ScopeSelect: {
				"select secrets",
				getopt.Definitions{
					{"organization|o", "organization id", getopt.Optional | getopt.ExampleIsDefault, ""},
				},
			},
			ScopeGet: {
				"get secret by key",
				getopt.Definitions{
					{"key|k", "key of a new secret", getopt.Required, ""},
				},
			},
			ScopePutValue: {
				"put secret by key",
				getopt.Definitions{
					{"key|k", "key in secret", getopt.Required, ""},
					{"value|v", "value in secret", getopt.Required, ""},
					{"username", "username in secret", getopt.Optional | getopt.ExampleIsDefault, ""},
					{"uri", "uri in secret", getopt.Optional | getopt.ExampleIsDefault, ""},
					{"notes|n", "notes of the new secret", getopt.Optional | getopt.ExampleIsDefault, ""},
				},
			},
			ScopePutFile: {
				"put secret by key",
				getopt.Definitions{
					{"key|k", "key in secret", getopt.Required, ""},
					{"value|v", "value in secret", getopt.Required, ""},
					{"notes|n", "notes of the new secret", getopt.Optional | getopt.ExampleIsDefault, ""},
				},
			},
			ScopeDelete: {
				"delete secret by key",
				getopt.Definitions{
					{"key|k", "key of a new secret", getopt.Required, ""},
				},
			},
			ScopeDrop: {
				"drop secrets table",
				getopt.Definitions{
					{"organization|o", "organization id", getopt.Optional | getopt.ExampleIsDefault, ""},
				},
			},
		},
	}

	scope, options, _, _, e := sco.ParseCommandLine()

	help, wantsHelp := options["help"]

	if e != nil || wantsHelp {
		exit_code := 0
		switch {
		case wantsHelp && help.String == "usage":
			fmt.Print(sco.Usage())
		case wantsHelp && help.String == "help":
			fmt.Print(sco.Help())
		default:
			fmt.Println("**** Error: ", e.Error(), "\n", sco.Help())
			exit_code = e.ErrorCode
		}
		os.Exit(exit_code)
	}
	if scope == "*" {
		log.Fatal(sco.Help())
	}

	//fmt.Printf("scope: %s\n", scope)
	//fmt.Printf("options: %#v\n", options)
	var request SecretCtx
	request.Scope = scope
	log.Println(scope)
	for k, v := range options {
		//log.Println(k, v.String)
		switch {
		case k == "key":
			request.Key = v.String
		case k == "value":
			request.Value = v.String
		case k == "username":
			request.Username = v.String
		case k == "uri":
			request.Uri = v.String
		case k == "notes":
			request.Notes = v.String
		case k == "file":
			request.Filepath = v.String
		}
	}
	return request
}
