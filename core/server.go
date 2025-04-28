package main

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"gopkg.in/yaml.v3"
)

type IServer interface {
	Serve()
	NewServer(endpoints map[string]func(w http.ResponseWriter, r *http.Request)) *Server
	GetInitialized() bool
	GetEndpoints() map[string]func(w http.ResponseWriter, r *http.Request)
	GetPrefix() string
	SetInitialized(value bool)
	SetEndpoints(value map[string]func(w http.ResponseWriter, r *http.Request))
	SetPrefix(value string)
}

type Server struct {
	Port int `yaml:"port"`
	Prefix string `yaml:"prefix"`
	Endpoints map[string]func(w http.ResponseWriter, r *http.Request)
	Initialized bool
}

var _ IServer = (*Server)(nil)

func (server *Server) NewServer(endpoints map[string]func(w http.ResponseWriter, r *http.Request)) *Server {

	// -- Read YAML config
	file, err := os.ReadFile("./config/server.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(file, server); err != nil {
		log.Fatal(err)
	}

	// -- Validate YAML config values
	if server.Port < 0 {
		log.Fatal("\n[Obelisk] Port should not be less than 0.")
	}
	
	match, err := regexp.MatchString(`^$|^[a-z]+$`, server.Prefix)
	if err != nil {
		log.Fatal("\n[Obelisk] Prefix does not follow the pattern:\n		'<foo>'")
	}
	if !match {
		server.SetPrefix(server.Prefix)
		log.Fatal("\n[Obelisk] Prefix does not follow the pattern:\n		'<foo>'")
	}

	server.SetEndpoints(endpoints)
	server.SetInitialized(true)

	return server
}

func (server *Server) Serve() {
	if !server.Initialized {
		log.Fatal("\n[Obelisk] Server could not be initialized. Please take a look at server.yaml file and try again.")
	}

	pattern := `^$|^[a-z]+(/[a-z]+)*$`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}

	for endpoint, endpointFn := range server.Endpoints {
		if !regex.MatchString(endpoint) {
			log.Fatal("\n[Obelisk] Endpoint does not follow the pattern:\n		'<foo>/<bar>/<foo>'")
		}
		http.HandleFunc("/" + server.Prefix + "/" + endpoint, endpointFn)
	}
	log.Printf(`
   ____  __         ___      __  
  / __ \/ /_  ___  / (_)____/ /__
 / / / / __ \/ _ \/ / / ___/ //_/
/ /_/ / /_/ /  __/ / (__  ) ,<   
\____/_.___/\___/_/_/____/_/|_|  
                                 `)
	log.Printf("API Server started! Listening on :%d/%s", server.Port, server.Prefix)
	
	log.Fatal(http.ListenAndServe(":" + strconv.FormatInt(int64(server.Port), 10), nil))
}

func (server *Server) GetInitialized() bool {
	return server.Initialized
}

func (server *Server) GetEndpoints() map[string]func(w http.ResponseWriter, r *http.Request) {
	return server.Endpoints
}

func (server *Server) GetPrefix() string {
	return server.Prefix
}

func (server *Server) SetInitialized(value bool) {
	server.Initialized = value
}

func (server *Server) SetEndpoints(value map[string]func(w http.ResponseWriter, r *http.Request)) {
	server.Endpoints = value
}

func (server *Server) SetPrefix(value string) {
	server.Prefix = value
}