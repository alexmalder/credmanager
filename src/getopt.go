package src

import (
	"fmt"
	"log"
	"os"

	"github.com/kesselborn/go-getopt"
)

const (
	ScopeCreate     string = "create"
	ScopeCreateFile        = "create-file"
	ScopeGet               = "get"
	ScopePut               = "put"
	ScopeDelete            = "delete"
)

// get command line arguments
func Getopts() Secret {
	sco := getopt.SubCommandOptions{
		getopt.Options{
			"global description",
			getopt.Definitions{
				{"config|c", "key of a new secret", getopt.Optional | getopt.ExampleIsDefault, "config.yml"},
			},
		},
		getopt.SubCommands{
			ScopeCreate: {
				"create key-value pair as string",
				getopt.Definitions{
					{"key|k", "key of a new secret", getopt.Required, ""},
					{"value|v", "value of a new secret", getopt.Required, ""},
				},
			},
			ScopeCreateFile: {
				"create key-value pair as file",
				getopt.Definitions{
					{"key|k", "key of a new secret", getopt.Required, ""},
					{"file|f", "file of a new secret", getopt.Required, ""},
				},
			},
			ScopeGet: {
				"get secret by key",
				getopt.Definitions{
					{"key|k", "key of a new secret", getopt.Required, ""},
				},
			},
			ScopePut: {
				"put secret by key",
				getopt.Definitions{
					{"key|k", "key of a new secret", getopt.Required, ""},
					{"value|v", "value of a new secret", getopt.Required, ""},
				},
			},
			ScopeDelete: {
				"delete secret by key",
				getopt.Definitions{
					{"key|k", "key of a new secret", getopt.Required, ""},
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
	var request Secret
	request.Scope = scope
	for k, v := range options {
		//log.Println(k, v.String)
		switch {
		case k == "key":
			request.Key = v.String
		case k == "value":
			request.Value = v.String
		case k == "file":
			request.File = v.String
		}
	}
	return request
}
