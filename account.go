package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"git.hex.uz/bobur1/go-django-auth"
)

func (s *Srv) account(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		userID, err := takeIdFromToken(r)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Login failed!", http.StatusUnauthorized)
			return
		}

		var ua Account
		err = s.db.Get(&ua,
			`SELECT email 
						from users 
						where id = $1;`, userID)

		if err != nil {
			if err == sql.ErrNoRows {

			} else {
				log.Println(err)
				ResponseMessage(w, "Some problems during finding row", http.StatusBadRequest)
				return
			}
		}
		articles := []Articles{}
		err = s.db.Select(&articles,
			`SELECT * 
						from articles 
						where user_id = $1;`, userID)

		if err != nil {
			log.Println(err)
			ResponseMessage(w, "User not founded!", http.StatusBadRequest)
			return
		}

		fullInformation := &Account{
			Email:ua.Email,
			Articles:articles,
		}

		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Request failed!", http.StatusBadRequest)
			return
		}

		WriteJSON(w, fullInformation)

	case "PUT":
		userID, err := takeIdFromToken(r)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Login failed!", http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Login failed!", http.StatusBadRequest)
			return
		}

		var userData Users
		json.Unmarshal(body, &userData)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Not right format provided!", http.StatusBadRequest)
			return
		}
		var find Users
		err = s.db.Get(&find,
			`SELECT * 
						from users 
						where id = $1 ;`, userID)

		if err != nil {
			log.Println(err)
			ResponseMessage(w, "User is not exist!", http.StatusUnauthorized)
			return
		}

		var u Users
		err = s.db.Get(&u,
			`SELECT * 
						from users 
						where email = $1 ;`, userData.Email)

		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Some problems with server!", http.StatusBadRequest)
			return

		}
			if u.Id != userID {
				ResponseMessage(w, "Email address already registered! Please, choose another one!", http.StatusBadRequest)
				return
			}
		if userData.Password != "" {
			saltNum := dj_auth.PBKDFf2GenSalt()
			pass := dj_auth.PBKDF2PasswordString(userData.Password, "pbkdf2_sha256", 100000, saltNum)

			_, err = s.db.Exec(`UPDATE users
			  SET email = $1, password= $2
			  WHERE id = $3;`, userData.Email, pass, userID)

			if err != nil {
				log.Println(err)
				ResponseMessage(w, "Some problems while updating! 1 ", http.StatusBadGateway)
				return
			}
		}else {
			if u.Id != userID {
				_, err = s.db.Exec(`UPDATE users
			  SET email = $1,
			  WHERE id = $2;`, userData.Email, userID)

				if err != nil {
					log.Println(err)
					ResponseMessage(w, "Some problems while updating! 2 ", http.StatusBadGateway)
					return
				}
			}
		}
		var updatedUser UsersResponse
		err = s.db.Get(&updatedUser,
			`SELECT email 
						from users 
						where id = $1 ;`, userID)

		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Some problems with server!", http.StatusBadGateway)
			return
		}

		WriteJSON(w,updatedUser)
	default:
		ResponseMessage(w,"Not allowed method!", http.StatusConflict)
	}
}