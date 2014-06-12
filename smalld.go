package main

import "os"
import "log"
import "fmt"
import ( 
	_ "github.com/lib/pq"
	"database/sql"
)
import "net/http"
import "net/url"
import "encoding/json"

var db *sql.DB //to share with our handlers

func SafeValues(v *url.Values) bool {
	log.Printf("safe %+v", v)
	//lat, lon, acc are floats, label is string safe for db insert
	return true //for now
}

func buildQuery(v *url.Values) (string) {
	point := fmt.Sprintf("POINT(%s %s)", v.Get("lon"), v.Get("lat"))
	query := fmt.Sprintf("select name from adminareas where st_contains(adminareas.geom, st_geomfromtext('%s', 4326));", point)
	return query
}

func LocationHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("handling url", req.URL)
	if req.Method == "GET" {
		if req.URL.RawQuery != "" {
		values, err := url.ParseQuery(req.URL.RawQuery)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(values)
			if SafeValues(&values) {
				q := buildQuery(&values)
				log.Println("query:", q)
				rows, err := db.Query(q)
				if err != nil {
					log.Print("db error",err)
				}
				var l []string
				for rows.Next() {
					var name string
					rows.Scan(&name)
					l = append(l, name)
				}

				m := make(map[string][]string)
				m["names"] = l
				j, err := json.Marshal(m)
				if err != nil {
					log.Fatal(err)
				}
				w.Write(j)
				return
			}
		} else {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}
		
	}
}

func EventsHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("handling /events request")
	if req.Method != "GET" {
		//die horribly with an error
	}
}

func main() {
	log.Println("smalld starting")
	db_connection := os.Getenv("SMALLD_DB_CONNECTION")
	url_base := os.Getenv("SMALLD_URL_BASE")
	options := os.Getenv("SMALLD_OPTIONS") //override command line flags 
	log.Println("SMALLD_DB_CONNECTION:", db_connection)
	log.Println("SMALLD_URL_BASE:", url_base)
	log.Println("SMALLD_OPTIONS:", options)
	var err error
	db, err = sql.Open("postgres", db_connection)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} 
	log.Println("connected to database")
	http.HandleFunc("/location", LocationHandler)
	log.Println("registered LocationHandler")
	http.ListenAndServe(":8000", nil)
}
