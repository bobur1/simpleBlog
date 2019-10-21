package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func cors(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Headers",
			"accept, accept-encoding, authorization, content-type, dnt, origin, user-agent, x-csrftoken, x-requested-with, Access-Control-Allow-Origin")
		w.Header().Set("Access-Control-Allow-Methods",
			"DELETE, GET, OPTIONS, PATCH, POST, PUT")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
		}
		h.ServeHTTP(w, r)
	})
}

type JWTData struct {
	// Standard claims are the standard jwt claims from the IETF standard
	// https://tools.ietf.org/html/rfc7519
	jwt.StandardClaims
	CustomClaims int `json:"user_id,omitempty"`
}

func (s *Srv) remove(fileName string) (error) {

	_, err := s.db.Exec(`Delete from uploads where path = $1`, fileName)
	if err != nil {
		panic(err)
	}
	os.Remove(fileName)

	return nil
}



//func WriteJSON(w http.ResponseWriter, msg interface{}) {
//	enc := json.NewEncoder(w)
//	enc.SetIndent("", "  ")
//	err := enc.Encode(msg)
//	if err != nil {
//		log.Println(err)
//	}
//}

func takeIdFromToken(r *http.Request) (int, error) {
	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken, " ")

	if len(authArr) != 2 {
		return 0, errors.New("Auth failed!")
	}

	jwtToken := authArr[1]

	claims, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		//if jwt.SigningMethodHS256 != token.Method {
		//	return nil, errors.New("Invalid signing algorithm")
		//}
		return []byte(Cfg.SRV.Secret), nil
	})

	if err != nil {
		log.Println(err)
		return 0, errors.New("Auth failed!")
	}

	data := claims.Claims.(*JWTData)

	userID := data.CustomClaims

	return userID, nil
}

func convertToArray(img string) ([]string) {
	img = strings.Replace(img, "[", "", -1)
	img = strings.Replace(img, "]", "", -1)
	img = strings.Replace(img, ",", " ", -1)
	img = strings.Replace(img, "'", "", -1)
	img = strings.Replace(img, "\"", "", -1)
	return strings.Fields(img)

}

func (s *Srv) getArticle(w http.ResponseWriter, articleId string) (SingleArticle, error) {

	var a SingleArticle
	log.Println("Listening id--->: ", articleId)
	err := s.db.Get(&a,
		`SELECT articles.*, avg(vote_article.mark) as vote
						from articles left outer join vote_article on vote_article.article_id = articles.id
						where articles.id=$1  GROUP BY articles.id`, articleId)

	if err != nil {
		log.Println(err)
		return a, err
	}
	up := []Uploads{}
	err = s.db.Select(&up,
		`SELECT path 
						from uploads 
						where article_id = $1;`, articleId)

	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return a, err
	}
	a.Images = up

	c := []Comment{}
	err = s.db.Select(&c,
		`SELECT comments.id, comments.article_id, comments.parent_id, comments.context, users.email as user_email 
						from comments join users on comments.user_id = users.id
						where article_id = $1;`, articleId)

	if err != nil && err != sql.ErrNoRows {
		log.Println(err)

		return a, err
	}
	a.Comments = c

	return a, nil
}


func (s *Srv) uploadHandler(w http.ResponseWriter, r *http.Request, articleId int) {

	m := r.MultipartForm

	//get the *fileheaders
	files := m.File["myfiles"]


	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Some problems while uploading files", http.StatusBadRequest)
			return
		}
		//create destination file making sure the path is writeable.
		var extension = filepath.Ext(files[i].Filename)
		randBytes := make([]byte, 8)
		rand.Read(randBytes)
		NewFileName := hex.EncodeToString(randBytes) + extension

		//if no DIRECTORY ---> make DIRECTORY
		if _, err := os.Stat(Cfg.SRV.Root); os.IsNotExist(err) {
			err = os.MkdirAll(Cfg.SRV.Root, 0755)
			if err != nil {
				log.Println(err)
			}
		}
		dst, err := os.Create(Cfg.SRV.Root + NewFileName)
		defer dst.Close()
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Some problems while creating destination folder", http.StatusBadRequest)
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			log.Println(err)
			ResponseMessage(w, "Some problems while uploading files in destination folder", http.StatusBadRequest)
			return
		}
		_, err = s.db.Exec(`INSERT INTO uploads (article_id,path)
								VALUES ($1, $2)`, articleId, Cfg.SRV.Root + NewFileName)
		if err != nil {
			log.Println(err)
			ResponseMessage(w, "Opps! Internal error", http.StatusBadRequest)
			return
		}
	}

}

func ResponseMessage(w http.ResponseWriter, message string, headerCode int) {
	var Resp ResponseMsg
	w.WriteHeader(headerCode)
	Resp.Message = message
	WriteJSON(w, Resp)
}
func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
func writeJSON(rw http.ResponseWriter, js []byte, code int) {

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	rw.Write(js)
}
func WriteJSON(rw http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(rw, js, http.StatusOK)
}
func WriteIndentJSON(rw http.ResponseWriter, data interface{}) {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(rw, js, http.StatusOK)
}
