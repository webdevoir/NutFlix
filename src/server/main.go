package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"./controllers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "clone_development"
)

var db *sql.DB

func init() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})

	userController := controllers.UserController{}
	movieController := controllers.MovieController{}

	r.HandleFunc("/api/users/{id}", userController.Show(db)).Methods("GET")
	r.HandleFunc("/api/movies", movieController.Index(db)).Methods("GET")
	r.HandleFunc("/api/movies/{id}", movieController.Show(db)).Methods("GET")
	r.HandleFunc("/signup", userController.Signup(db)).Methods("POST")
	r.HandleFunc("/login", userController.Login(db)).Methods("POST")
	r.HandleFunc("/logout", userController.Logout())
	port := ":5000"

	fmt.Println("Listening on port " + port)
	http.ListenAndServe(port, handlers.CORS(headers, methods, origins)(r))
}
