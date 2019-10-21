# SimpleBlog


[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://git.hex.uz/bobu1/simpleBlog)

![enter image description here](https://atinysliceofkate.files.wordpress.com/2013/01/expired_domain_names.png)

**New ReadMe will be provided soon...**

SimpleBlog - простой бэкенд блога написанный на Golang.

  Имеет следущий функционал
  - [Registration](#registration)
  - [Login](#login)
  - [Account](#account) (show/edit)
  - [Article](#article) (show/read/write/edit/vote/delete)
  - [Comment](#comment) (write)

## [Registration]
Header
```
Content-Type	application/json
```
Body
```
{
	"email": "new@mail.com",
	"password" :"password"
}
```

![alt text](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/registration.PNG)  
## [Login]
Header
```
Content-Type	application/json
```
Body
```
{
	"email": "new@mail.com",
	"password" :"password"
}
```
Copy given token and past it in the header and from now our header will be:
Header
![Image](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/login.PNG) 

```
Content-Type	application/json
Authorization	Berear eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```
> Note: from now I will use this header.  
The exception will be spelled out

## [Account]
### Show
It is enough to send our header

![alt text](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/account.PNG)  
### Edit
Body
```
{
	"email": "new123@mail.com",
	"password" :"password222"
}
```
![alt text](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/accountEdit.PNG)  
## Article
It is not mandatory to add Auth token in the header for *show* and *read* article.
### Show
*Have limitation 10 articles*

![alt text](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/articles.PNG)  
### Read
Body (id of the article)
```
{
	"Id":"2"
}
```
Avg -- avarage rate (vote)

![alt text](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/articleRead.PNG)  
### Write
*We need to use different header.*
Header
```
Content-Type	multipart/form-data
```
And body part should be the same as shown below.

![alt text](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/articleWrite.PNG)  
### Edit
*We need to use different header.*
Header
```
Content-Type	multipart/form-data
```
And body part should be the same as shown below.
*deleteImg --- images (names) wich should be replaced(deleted)*

![alt text](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/articleEdit.PNG)  
### Vote
```
{
	"article_id" : 47,
	"mark": 5
}
```
mark -- rate; 
*deleteImg --- images (names) wich should be replaced(deleted)*

![alt text](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/articleVote.PNG)  
### Delete
```
{
	"Id" : "47"
}
```


![alt text](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/articleDelete.PNG)  
## [Comment]
### Write
```
{
	"article_id": 47,
	"text" : "Good article!"
}
```
![alt text](https://git.hex.uz/bobur1/simpleBlog/raw/branch/master/screens/commentWrite.PNG) 

### Todos

 - Write MORE Tests
 - Optimize code
 - More functionality

