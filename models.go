package main

type Users struct {
	Id       int    `db:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}
type UsersResponse struct {
	Email    string `db:"email" json:"email"`
}
type Account struct {
	Email    string     `db:"email"`
	Articles []Articles   `json:"articles" db:"article"`
}
type Articles struct {
	Id int `db:"id"`
	//ToDo:: int
	UserId string `db:"user_id"`
	Title  string `db:"title"`
	Text   string `db:"context"`
}

type SingleArticle struct {
	Id       int             `db:"id" json:"id"`
	UserId   int             `db:"user_id" json:"user_id"`
	Title    string          `db:"title" json:"title"`
	Text     string          `db:"context" json:"text"`
	Avg      JsonNullFloat64 `db:"vote" json:"avg"`
	Vote     `db:"vote_article"`
	Images   []Uploads `json:"images"`
	Comments []Comment `json:"comments"`
}

type IncomningComment struct {
	ParentId  int    `json:"parent_id"`
	ArticleId int    `json:"article_id"`
	Context   string `json:"text"`
}

type Uploads struct {
	//not used ---> comment if u need
	//Id 				int 	`db:"id"`
	//ArticleId 		int 	`db:"article_id"`
	Img string `json:"img" db:"path"`
}
type UploadsDelete struct {
	Id        int    `db:"id"`
	ArticleId int    `db:"article_id"`
	Img       string `db:"path"`
}

type Comment struct {
	Id        int           `json:"id" db:"id"`
	UserEmail string        `json:"user_email" db:"user_email"`
	ArticleId int           `json:"article_id" db:"article_id"`
	ParentId  JsonNullInt `json:"parent_id" db:"parent_id"`
	Context   string        `json:"context" db:"context"`
}

type Vote struct {
	Id        int `json:"id" db:"id"`
	ArticleId int `json:"article_id" db:"article_id"`
	UserId    int `json:"user_id"  db:"user_id"`
	Mark      int `json:"mark"  db:"mark"`
}

type ResponseMsg struct {
	Message string `json:"message"`
}