package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "gopkg.in/dancannon/gorethink.v2"
)

type user struct {
	Email    string `json:"email" gorethink:"email"`
	Password string `json:"password" gorethink:"password"`
}

var session *db.Session

func init() {
	var err error
	session, err = db.Connect(db.ConnectOpts{
		Addresses:     []string{"db1:28015", "db2:28015"},
		Database:      "bourbaki",
		DiscoverHosts: true,
	})
	// session, err = db.Connect(db.ConnectOpts{
	// 	Address:  "localhost:28015",
	// 	Database: "bourbaki",
	// 	MaxIdle:  10,
	// 	MaxOpen:  10,
	// })
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	http.HandleFunc("/", handler) // each request calls handler
	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
}

// handler echoes the Path component of the request URL r.
func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var u user
		err := decoder.Decode(&u)
		if err != nil {
			log.Fatalln(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
		} else {
			db.Table("user").Insert(u).Run(session)
			fmt.Fprintf(w, "Hello from api.")
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such method")
	}
}
