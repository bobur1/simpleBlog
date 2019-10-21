package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (s *Srv) voteArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		userID, err := takeIdFromToken(r)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Login failed!", http.StatusUnauthorized)
			return
		}
		queryValues := r.URL.Query()
		articleId := queryValues.Get("id")
		log.Println("Article: " + articleId)
		if articleId == "" {
			ResponseMessage(w, "No such article!", http.StatusMisdirectedRequest)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Login failed!", http.StatusUnauthorized)
			return
		}

		var userData Vote
		json.Unmarshal(body, &userData)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Not correct type of data!", http.StatusUnprocessableEntity)
			return
		}

		// let's find if there any article in db
		var a SingleArticle

		err = s.db.Get(&a,
			`SELECT id 
						from articles 
						where id=$1`, articleId)
		//if nothing ---> error

		if err != nil {
			http.Error(w, "No such article anymore!", http.StatusUnauthorized)
		}

		artId, err := strconv.Atoi(articleId)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Not correct type of data!", http.StatusUnprocessableEntity)
			return
		}
		var v Vote

		err = s.db.Get(&v,
			`SELECT id 
						from vote_article 
						where user_id=$1 and article_id =$2`, userID, artId)
		if err == sql.ErrNoRows {
			_, err := s.db.Exec(`INSERT INTO vote_article (article_id,user_id, mark)
								VALUES ($1, $2, $3)`, artId, userID, userData.Mark)

			if err != nil {
				log.Println(err)
				ResponseMessage(w, "Opps! Some problem in server side", http.StatusBadGateway)
				return
			}

			a, err = s.getArticle(w, articleId)
			if err != nil {
				log.Println(err)
				ResponseMessage(w, "Problem with finding article", http.StatusBadRequest)
				return
			}
			WriteJSON(w, a)
		} else if err != nil {
			ResponseMessage(w, "Some problems with voting this article!", http.StatusBadRequest)
			return
		} else {
			ResponseMessage(w, "You have already voted", http.StatusBadRequest)
			return
		}
	default:
		ResponseMessage(w, "Not allowed request type!", http.StatusMisdirectedRequest)

	}
}
