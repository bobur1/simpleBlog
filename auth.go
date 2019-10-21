package main

import (
	"database/sql"
	"encoding/json"
	"git.hex.uz/bobur1/go-django-auth"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (s *Srv) registration(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		ResponseMessage(w, "Auth failed!", http.StatusUnauthorized)


	}

	var userData Users
	json.Unmarshal(body, &userData)
	log.Println("email: ", userData.Email)
	log.Println("password: ", userData.Password)

	if userData.Email != "" && userData.Password != "" {

		var u Users
		err = s.db.Get(&u,
			`SELECT * 
						from users 
						where email = $1 ;`, userData.Email)

		if err == sql.ErrNoRows {
			saltNum := dj_auth.PBKDFf2GenSalt()

			pass:= dj_auth.PBKDF2PasswordString(userData.Password, "pbkdf2_sha256", 100000, saltNum)

			var id int
			err = s.db.QueryRow(`INSERT INTO users(email, password)
								VALUES ($1, $2) RETURNING users.id`, userData.Email, pass).Scan(&id)

			if err != nil {
				log.Println(err)
				ResponseMessage(w, "Auth failed!", http.StatusBadGateway)
			}
			var rU UsersResponse
			err = s.db.Get(&rU,
				`SELECT  email
						from users 
						where id = $1;`, id)
			//jsonData, err := json.Marshal(struct {
			//	Success string `json:"success"`
			//}{
			//	"Created new User",
			//})
			//
			//if err != nil {
			//	log.Println(err)
			//	http.Error(w, "Request failed!", http.StatusUnauthorized)
			//}
			if err != nil {
				log.Println(err)
				ResponseMessage(w, "Error while finding user!", http.StatusBadGateway)
				return

			}
			WriteJSON(w, rU)
		}else if err != nil {
			ResponseMessage(w, "Some problems with server, please, try again later", http.StatusUnauthorized)
		}else{
			ResponseMessage(w, "Auth failed! The email exists", http.StatusUnauthorized)
		}

	}
}

func (s *Srv) login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		ResponseMessage(w, "Login failed!", http.StatusUnauthorized)
	}

	var userData Users
	json.Unmarshal(body, &userData)

	if userData.Email != "" && userData.Password != "" {
		var u Users
		err = s.db.Get(&u,
			`SELECT * 
						from users 
						where email = $1;`, userData.Email)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Login failed! No such email", http.StatusUnauthorized)
			return
		}
		// Comparing the password with the hash

		if !dj_auth.PBKDF2CheckPassword(u.Password, userData.Password) {
			ResponseMessage(w, "Login failed! 1", http.StatusUnauthorized)
		} else {


			claims := JWTData{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour).Unix(),
				},

				CustomClaims: u.Id,
			}


			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(Cfg.SRV.Secret))
			if err != nil {
				log.Println(err)
				ResponseMessage(w, "Login failed! 2", http.StatusUnauthorized)
			}

			json, err := json.Marshal(struct {
				Token string `json:"access"`
			}{
				tokenString,
			})

			if err != nil {
				log.Println(err)
				ResponseMessage(w, "Login failed! 3", http.StatusUnauthorized)
				return
			}

			w.Write(json)
		}
	} else {
		ResponseMessage(w, "Login failed! 4", http.StatusUnauthorized)
		return
	}
}


