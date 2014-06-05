package main

import "os"
import "log"
import ( 
	_ "github.com/lib/pq"
	"database/sql"
)
import "net/http"

func LocationHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("handling url", req.URL)
} 

func main() {
	log.Println("smalld starting")
	db_connection := os.Getenv("SMALLD_DB_CONNECTION")
	url_base := os.Getenv("SMALLD_URL_BASE")
	options := os.Getenv("SMALLD_OPTIONS") //override command line flags 
	log.Println("SMALLD_DB_CONNECTION:", db_connection)
	log.Println("SMALLD_URL_BASE:", url_base)
	log.Println("SMALLD_OPTIONS:", options)
	db, err := sql.Open("postgres", db_connection)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} 
	log.Println("connected to database")
	http.HandleFunc("/location", LocationHandler)
	log.Println("registered LocationHandler")
	http.ListenAndServe(":8000", nil)
}
