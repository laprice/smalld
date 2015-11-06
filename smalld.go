package main

import "os"
import "log"
import "fmt"
import "strconv"
import (
	"database/sql"
	_ "github.com/lib/pq"
)
import "net/http"
import "net/url"
import "encoding/json"

var db *sql.DB //to share with our handlers

func safeValues(v *url.Values) bool {
	log.Printf("safe %+v", v)
	return true //for now
}

func makePoint(v *url.Values) string {
	point := fmt.Sprintf("POINT(%s %s)", v.Get("lon"), v.Get("lat"))
	return point
}

func recordlocations(v *url.Values) {
	p := makePoint(v)
	txn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	label := fmt.Sprintf("%s", v.Get("label"))
	acc, err := strconv.ParseFloat(v.Get("acc"), 64)
	if err != nil {
		log.Fatal(err)
	}
	_, err = txn.Exec("insert into locations ( label, acc, geom ) values ( $1, $2, ST_PointFromText( $3, 4326) )", label, acc, p)
	if err != nil {
		log.Fatal(err)
	}
	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// LocationHandler is the main entry point for smalld 
// it receives the get request parses the location data from it
// and logs the values to the location table.
func LocationHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("handling url", req.URL)
	if req.Method == "GET" {
		if req.URL.RawQuery != "" {
			values, err := url.ParseQuery(req.URL.RawQuery)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(values)
			if safeValues(&values) {
				p := makePoint(&values)
				go recordlocations(&values)
				log.Println("point:", p)
				q := "select name from adminareas where st_contains(adminareas.geom, st_geomfromtext( $1 , 4326))"
				rows, err := db.Query(q, p)
				if err != nil {
					log.Print("db error", err)
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
				h := w.Header()
				h.Add("Access-Control-Allow-Origin", "*")
				w.Write(j)
				return
			}
		} else {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

	}
}

func main() {
	log.Println("smalld starting")
	dbConnection := os.Getenv("SMALLD_DB_CONNECTION")
	urlBase := os.Getenv("SMALLD_URL_BASE")
	listenAddress := os.Getenv("SMALLD_LISTEN_ADDRESS")
	options := os.Getenv("SMALLD_OPTIONS") //override command line flags
	log.Println("SMALLD_DB_CONNECTION:", dbConnection)
	log.Println("SMALLD_URL_BASE:", urlBase)
	log.Println("SMALLD_LISTEN_ADDRESS", listenAddress)
	log.Println("SMALLD_OPTIONS:", options)
	var err error
	db, err = sql.Open("postgres", dbConnection)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to database")
	http.HandleFunc("/location", LocationHandler)
	log.Println("registered LocationHandler")
	http.ListenAndServe(listenAddress, nil)
}
