package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func (s *Srv) articleHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		queryValues := r.URL.Query()
		articleId := queryValues.Get("id")
		log.Println("Article: " + articleId)
		if articleId == "" {

			a := []Articles{}
			//ToDo: pagination will be more practicel
			err := s.db.Select(&a,
				`SELECT * 
						from articles 
						limit 10`)

			if err != nil {
				log.Println(err)
				ResponseMessage(w, "Opps! Some problems in server", http.StatusBadGateway)
				return
			}
			WriteJSON(w, a)
		} else {

			a,err := s.getArticle(w, articleId)
			if err != nil {
				log.Println(err)
				ResponseMessage(w, "Problem with finding article", http.StatusBadRequest)
				return
			}
			WriteJSON(w, a)
		}
	case "POST":
		err := r.ParseMultipartForm(100000)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Post request failed!", http.StatusMisdirectedRequest)
			return
		}

		userId, err := takeIdFromToken(r)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Login failed!", http.StatusUnauthorized)
			return
		}
		notInitialized := ""
		title := r.FormValue("title")
		if title == "" {
			notInitialized += "Title,"
		}
		text := r.FormValue("text")
		if text == "" {
			notInitialized += "Text,"
		}
		log.Println("Listening user id--->: ", userId)
		if notInitialized != ""{
			ResponseMessage(w, notInitialized+" not initialized", http.StatusUnprocessableEntity)
			return
		}
		var articleId int
		err = s.db.QueryRow(`INSERT INTO articles(user_id, title, context)
								VALUES ($1, $2, $3) RETURNING articles.id`, userId, title, text).Scan(&articleId)

		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Opps! Internal error", http.StatusBadGateway)
			return
		}

		art := strconv.Itoa(articleId)
		log.Println("art-->: ", art)
		s.uploadHandler(w, r, articleId)
		a,err := s.getArticle(w,art)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Problem with finding article", http.StatusBadRequest)
			return
		}
		WriteJSON(w, a)

	case "PUT":
		queryValues := r.URL.Query()
		articleId := queryValues.Get("id")
		log.Println("Article: " + articleId)
		if articleId == "" {
			ResponseMessage(w, "Put request failed! No article id!", http.StatusMisdirectedRequest)
			return
		}
		err := r.ParseMultipartForm(100000)

		if err != nil {
			ResponseMessage(w, "Put request failed!", http.StatusMisdirectedRequest)
			return
		}

		userId, err := takeIdFromToken(r)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Login failed!", http.StatusUnauthorized)
			return
		}
		log.Println("Listening user id--->: ", userId)
		notInitialized := ""
		title := r.FormValue("title")
		if title == "" {
			notInitialized += "Title,"
		}
		text := r.FormValue("text")
		if text == "" {
			notInitialized += "Text,"
		}
		log.Println("Listening user id--->: ", userId)
		if notInitialized != ""{
			ResponseMessage(w, notInitialized+" not initialized", http.StatusUnprocessableEntity)
			return
		}
		//ToDo::Make method which should delete all images
		//deleteImg := r.FormValue("deleteImg")
		//log.Println("Delete Img--->: ", deleteImg)
		//if deleteImg != "" {
		//	imgArray := convertToArray(deleteImg)
		//
		//	//var imgArray []int
		//	//err = json.Unmarshal([]byte(deleteImg), &imgArray)
		//	//if err != nil {
		//	//	log.Fatal(err)
		//	//}
		//	for i := 0; i < len(imgArray); i++ {
		//		log.Println("Deleting Img--->: ", imgArray[i])
		//		err = s.remove(imgArray[i])
		//		if err != nil {
		//			http.Error(w, "Some problems during removing image!", http.StatusConflict)
		//		}
		//	}
		//	log.Println("Delete Img--->: ", imgArray)
		//}
		log.Println("Title--->: ", title)
		log.Println("text--->: ", text)

		art, err := strconv.Atoi(articleId)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Article id must be integer", http.StatusUnprocessableEntity)
			return
		}

		log.Println("Article id--->: ", art)
		log.Println("User id--->: ", userId)
		_, err = s.db.Exec(`UPDATE articles
								SET title = $1, context= $2
								WHERE id = $3 and user_id = $4;`, title, text, art, userId)

		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Article is not exist", http.StatusBadRequest)
			return
		}
		s.uploadHandler(w, r, art)
		a,err := s.getArticle(w, articleId)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Problem with finding article", http.StatusBadRequest)
			return
		}
		WriteJSON(w, a)
	case "DELETE":

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
			ResponseMessage(w, "No article id!", http.StatusMisdirectedRequest)
			return
		}
		// let's find if there any article in db
		var a SingleArticle

		err = s.db.Get(&a,
			`SELECT * 
						from articles 
						where id=$1 and user_id = $2`, articleId, userID)
		log.Println("article id: ", a.Id)
		//if nothing ---> error
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "No access!", http.StatusUnauthorized)
			return
		}

		up := []UploadsDelete{}
		err = s.db.Select(&up,
			`SELECT  *
						from uploads 
						where article_id=$1 group by article_id, id`, a.Id)
		//if nothing ---> go next

		if err == nil {
			// remove all images first

			for _, im := range up {
				log.Println("im : ", im.Img)
				s.remove(im.Img)
				//err := s.remove(im.Img)

				if err != nil {
					log.Println(err)
					ResponseMessage(w, "Some problems during removing image!", http.StatusConflict)
					return
				}
			}
			_,err = s.db.Exec(`Delete from uploads where article_id = $1`, a.Id)
			if err != nil{
				log.Println(err)
				ResponseMessage(w, "Server side problem", http.StatusBadGateway)
				return
			}
		} else if err == sql.ErrNoRows {
			log.Println("No images")

		}else{
			log.Println(err)
			ResponseMessage(w, "Server side problem", http.StatusBadGateway)
			return
		}

		_, err = s.db.Exec(`Delete from articles where id = $1`, a.Id)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Server side problem", http.StatusBadGateway)
			return
		}

		ResponseMessage(w, "Deleted!", http.StatusOK)
		return
	default:
		ResponseMessage(w, "Not allowed request type!", http.StatusMisdirectedRequest)
	}

}

func (s *Srv) deleteImage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":

		userID, err := takeIdFromToken(r)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Login failed!", http.StatusUnauthorized)
			return
		}
		queryValues := r.URL.Query()
		path := queryValues.Get("id")
		log.Println("Uploads: " + path)
		if path == "" {
			ResponseMessage(w, "No such article!", http.StatusMisdirectedRequest)
			return
		}

		//body, err := ioutil.ReadAll(r.Body)
		//if err != nil {
		//	log.Println(err)
		//	ResponseMessage(w, "Login failed!", http.StatusUnauthorized)
		//	return
		//}
		//
		//var path string
		//json.Unmarshal(body, &path)
		//if err != nil {
		//	log.Println(err)
		//	ResponseMessage(w, "Not correct type of data!", http.StatusUnprocessableEntity)
		//	return
		//}

		var a Uploads
		path = strings.Replace(path, "\"", "", -1)
		err = s.db.Get(&a,
			`SELECT uploads.path 
						from articles 
						join uploads on articles.id=uploads.article_id
						where articles.user_id=$1 and uploads.path=$2`, userID,path)
		//if nothing ---> error

		if err != nil {
			log.Println(err)
			ResponseMessage(w, "No such image!"+path, http.StatusUnauthorized)
			return
		}

		err = s.remove(path)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Some problems during removing image!", http.StatusConflict)
			return
		}
		_,err = s.db.Exec(`Delete from uploads where path = $1`, path)
		if err != nil{
			log.Println(err)
			ResponseMessage(w, "Server side problem", http.StatusBadGateway)
			return
		}
		ResponseMessage(w, "Success! Deleted", http.StatusOK)

	default:
		ResponseMessage(w, "Not allowed request type!", http.StatusMisdirectedRequest)
	}
}

func (s *Srv) deleteAllImages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":

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
		artId, err := strconv.Atoi(articleId)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "No such article!", http.StatusUnprocessableEntity)
			return
		}

		up := []UploadsDelete{}
		err = s.db.Select(&up,
			`SELECT  uploads.path
						from uploads join articles on uploads.article_id = articles.id
						where articles.id=$1 and articles.user_id =$2`, artId, userID)
		//if nothing ---> error

		if err != nil {
			log.Println(err)
			ResponseMessage(w, "No such article!", http.StatusUnauthorized)
			return
		}
		for _, im := range up {
			log.Println("im : ", im.Img)
			s.remove(im.Img)
			//err := s.remove(im.Img)

			if err != nil {
				log.Println(err)
				ResponseMessage(w, "Some problems during removing image!", http.StatusConflict)
				return
			}
		}
		_,err = s.db.Exec(`Delete from uploads where article_id = $1`, artId)
		if err != nil{
			log.Println(err)
			ResponseMessage(w, "Server side problem", http.StatusBadGateway)
			return
		}

		ResponseMessage(w, "Success! Deleted all", http.StatusOK)


	default:
		ResponseMessage(w, "Not allowed request type!", http.StatusMisdirectedRequest)
	}
}

