package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tochk/kirino/auth"
	"github.com/tochk/kirino/departments"
	"github.com/tochk/kirino/generator"
	"github.com/tochk/kirino/memorandums"
	"github.com/tochk/kirino/server"
	"github.com/tochk/kirino/users"
)

var (
	configFile  = flag.String("config", "conf.json", "Where to read the config from")
	servicePort = flag.Int("port", 4001, "Service port number")
)

func filesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	path := "." + r.URL.Path
	if f, err := os.Stat(path); err == nil && !f.IsDir() {
		http.ServeFile(w, r, path)
		return
	}
	http.NotFound(w, r)
}

func checkFolders() {
	if file, err := os.Open("userFiles"); err != nil {
		file.Close()
		if err = os.Mkdir("userFiles", 0644); err != nil {
			log.Fatal(err)
		}
		log.Println("Creating directory for documents")
	}
}

func main() {
	log.Println("Checking folder for documents")
	checkFolders()
	flag.Parse()
	if err := server.LoadConfig(*configFile); err != nil {
		log.Fatal(err)
	}
	log.Println("Config loaded from", *configFile)

	server.ConnectToDb()
	defer server.Core.Db.Close()
	log.Printf("Connected to database on %s", server.Config.DbHost)

	router := mux.NewRouter().StrictSlash(true)

	router.Methods("GET").PathPrefix("/static/").HandlerFunc(filesHandler)
	router.Methods("GET").PathPrefix("/userFiles/").HandlerFunc(filesHandler)

	router.HandleFunc("/", memorandums.FormsHandler).Methods("GET")
	router.HandleFunc("/{type}/", memorandums.FormsHandler).Methods("GET")
	router.HandleFunc("/generate/{type}/", generator.GenerateHandler).Methods("POST")
	router.HandleFunc("/generated/{type}/{token}/", generator.GeneratedHandler).Methods("GET")

	router.HandleFunc("/admin/{type}/", auth.Handler)

	router.HandleFunc("/wifi/memorandums/{action}/{num:[0-9]+}", memorandums.ListWifiHandler)
	router.HandleFunc("/wifi/users/{action}/{num:[0-9]+}", users.WifiUsersHandler)

	router.HandleFunc("/departments/{action}/{num:[0-9]+}", departments.Handler).Methods("GET")

	router.HandleFunc("/ethernet/memorandums/{action}/{num:[0-9]+}", memorandums.ListEthernetHandler)

	router.HandleFunc("/phone/memorandums/{action}/{num:[0-9]+}", memorandums.ListPhoneHandler)

	router.HandleFunc("/domain/memorandums/{action}/{num:[0-9]+}", memorandums.ListDomainHandler)

	router.HandleFunc("/mail/memorandums/{action}/{num:[0-9]+}", memorandums.ListMailHandler)

	port := strconv.Itoa(*servicePort)
	log.Println("Server started at port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}