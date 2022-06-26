package main

import (
	"backend/app/cmd"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

// Opts is the app start options
type Opts struct {
	AppServerCmd cmd.ServerCommand `command:"app-server"`
}

//func main() {
//	var opts Opts
//	p := flags.NewParser(&opts, flags.Default)
//	p.CommandHandler = func(command flags.Commander, args []string) error {
//		// commands implements CommonOptionsCommander to allow passing set of extra options defined for all commands
//		c := command.(cmd.CommonOptionsCommander)
//		err := c.Execute(args)
//		if err != nil {
//			fmt.Printf("[ERROR] failed with %+v", err)
//		}
//		return err
//	}
//
//	if _, err := p.Parse(); err != nil {
//		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
//			os.Exit(0)
//		} else {
//			os.Exit(1)
//		}
//	}
//}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000" // and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are // encountered during parsing the application will be terminated.
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", home)
	r.Get("/carousel", carousel)
	//r.Post("/create", createSnippet)

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, r)
	log.Fatal(err)
}
