package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (s *Srv) writeComment(w http.ResponseWriter, r *http.Request) {
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

		var userData IncomningComment
		json.Unmarshal(body, &userData)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Not correct type of data!", http.StatusUnprocessableEntity)
			return
		}

		var a SingleArticle

		err = s.db.Get(&a,
			`SELECT id 
						from articles 
						where id=$1`, articleId)
		//if nothing ---> error

		if err != nil {
			log.Println(err)
			ResponseMessage(w, "No such article anymore!", http.StatusUnauthorized)
			return
		}
		artId, err := strconv.Atoi(articleId)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "No such article anymore!", http.StatusUnprocessableEntity)
			return
		}

		_, err = s.db.Exec(`INSERT INTO comments(article_id,user_id, parent_id, context)
								VALUES ($1, $2, $3, $4)`, artId, userID, userData.ParentId, userData.Context)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Opps! Some problem with server", http.StatusBadGateway)
			return
		}

		a, err = s.getArticle(w, articleId)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Problem with finding article", http.StatusBadRequest)
			return
		}
		WriteJSON(w, a)
	default:
		http.Error(w, "Not allowed request type!", http.StatusMisdirectedRequest)
	}
}
