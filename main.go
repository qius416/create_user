package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	db "gopkg.in/dancannon/gorethink.v2"
)

type user struct {
	Email    string `json:"email" gorethink:"email"`
	Password string `json:"password" gorethink:"password"`
}

type auth struct {
	Token string `json:"token"`
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
	// 	Address:  "192.168.99.100:28015",
	// 	Database: "bourbaki",
	// 	MaxIdle:  10,
	// 	MaxOpen:  10,
	// })
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	http.HandleFunc("/signup", signupHandler) // each request calls handler
	http.HandleFunc("/login", loginHandler)   // each request calls handler
	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
	//log.Fatal(http.ListenAndServe("0.0.0.0:7000", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var u user
		err := decoder.Decode(&u)
		if err != nil {
			log.Fatalln(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, err.Error())
		} else {
			res, err := db.Table("user").Get(u.Email).Run(session)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, err.Error())
				return
			}
			defer res.Close()

			if res.IsNil() {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "User not found")
				return
			}

			var myuser map[string]interface{}
			err = res.One(&myuser)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, err.Error())
				return
			}
			token, err := makeToken()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, err.Error())
			} else {
				p := auth{Token: token}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(p)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such method")
	}
}

// handler echoes the Path component of the request URL r.
func signupHandler(w http.ResponseWriter, r *http.Request) {
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

func makeToken() (tokenString string, err error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["role"] = "user"
	token.Claims["exp"] = time.Now().Add(time.Second * 30).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err = token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return
}
