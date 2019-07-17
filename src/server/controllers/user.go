package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "time"

	"../models"
	userrepository "../repository"
	"../utils"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var users []models.User

// UserController struct
type UserController struct{}

// Claims struct
// type Claims struct {
// 	Username string `json:"username"`
// 	jwt.StandardClaims
// }

var jwtKey = []byte("my_secret_key")

// Login function to log the user in
func (c UserController) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := utils.Store.Get(r, "session")
		if err != nil {
			log.Println("error identifying session")
			return
		}

		var user models.User
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		password := user.Password

		userRepo := userrepository.UserRepository{}
		user, erro := userRepo.Login(db, user)

		log.Println(err)

		if erro != nil {
			if erro == sql.ErrNoRows {
				error.Message = "The user does not exist"
				utils.RespondWithError(w, http.StatusBadRequest, error)
				return
			} else {
				log.Fatal(err)
			}
		}

		hashedPassword := user.Password

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

		if err != nil {
			error.Message = "Invalid Password"
			utils.RespondWithError(w, http.StatusUnauthorized, error)
			return
		}

		session.Values["loggedin"] = "true"
		session.Values["email"] = user.Email
		session.Save(r, w)
		log.Print("User ", user.Email, " is authenticated")

		// expirationTime := time.Now().Add(120 * time.Minute)

		// claims := &Claims{
		// 	Username: user.Username,
		// 	StandardClaims: jwt.StandardClaims{
		// 		ExpiresAt: expirationTime.Unix(),
		// 	},
		// }

		// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// tokenString, err := token.SignedString(jwtKey)

		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }

		// http.SetCookie(w, &http.Cookie{
		// 	Name:    "token",
		// 	Value:   tokenString,
		// 	Path: "/",
		// 	Expires: expirationTime,
		// })

	}
}

// Signup for new users
func (c UserController) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)
		fmt.Printf("%+v\n", user)

		if user.Email == "" {
			error.Message = "Email is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing."
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)

		userRepo := userrepository.UserRepository{}
		user = userRepo.Signup(db, user)

		if err != nil {
			error.Message = "Server error."
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		user.Password = ""

		w.Header().Set("Content-Type", "application/json")
		utils.ResponseJSON(w, user)

		json.NewEncoder(w).Encode(user)
	}
}

// Show to display profile
func (c UserController) Show(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		params := mux.Vars(r)

		if r.Method != "GET" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		id := params["id"]
		fmt.Printf(id)

		if id == "" {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		row := db.QueryRow("SELECT * FROM users WHERE id=$1", id)

		newUser := user
		err := row.Scan(&newUser.ID, &newUser.Email, &newUser.Username, &newUser.Password)
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		} else if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
}

// Welcome will be the route for a signed in user
// func (c UserController) Welcome(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		c, err := r.Cookie("token")
// 		if err != nil {
// 			if err == http.ErrNoCookie {
// 				w.WriteHeader(http.StatusUnauthorized)
// 				return
// 			}
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		tknStr := c.Value

// 		claims := &Claims{}

// 		// Parse the JWT string and store the result in `claims`.
// 		// Note that we are passing the key in this method as well. This method will return an error
// 		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
// 		// or if the signature does not match
// 		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})
// 		if err != nil {
// 			if err == jwt.ErrSignatureInvalid {
// 				w.WriteHeader(http.StatusUnauthorized)
// 				return
// 			}
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		if !tkn.Valid {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}

// 		w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
// 	}
// }

//Logout logs the user out duh
func (c UserController) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := utils.Store.Get(r, "session")
		if err == nil { //If there is no error, then remove session
			if session.Values["loggedin"] != "false" {
				session.Values["loggedin"] = "false"
				session.Save(r, w)
			}
		}
		log.Print("User had been logged out!")
	}
}