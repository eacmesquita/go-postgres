package main

import (
  "database/sql"
  "fmt"
  "os"
  _ "github.com/lib/pq"
  "net/http"
  "log"
  "encoding/json"
  "github.com/gorilla/mux"
)

const apiVersion = "v1"

var (
  host     = os.Getenv("DATABASE_SERVICE_NAME")
  port     = 5432
  user     = os.Getenv("POSTGRESQL_USER")
  password = os.Getenv("POSTGRESQL_PASSWORD")
  dbname   = os.Getenv("POSTGRESQL_DATABASE")
)

type Article struct {
  Id int `json:"id"`
  Title string `json:"title"`
  Description string `json:"desc"`
  Content string `json:"content"`
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
  fmt.Println("Endpoint Hit: returnAllArticles")
  sqlStatement := "select * from articles;"
  db := dbConn()
  var articles []Article
  rows,err := db.Query(sqlStatement)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  for rows.Next(){
    var article Article
    err := rows.Scan(&article.Id, &article.Title, &article.Description, &article.Content)
    articles = append(articles,article)
    if err != nil {
      panic(err)
    }
  }
  encoder := json.NewEncoder(w)
  //encoder.SetIndent("","    ")
  encoder.Encode(articles)
}

func homePage(w http.ResponseWriter, r *http.Request){
  fmt.Fprintf(w, "Welcome to the HomePage! This is version %s", apiVersion)
  fmt.Println("Endpoint Hit: homePage")
}

func createArticle(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var article Article
  err := decoder.Decode(&article)
  if err != nil {
    panic(err)
  }
  log.Println("Received:", article)
  sqlStatement := fmt.Sprintf("insert into articles(title,description,content) values($1,$2,$3) returning id")
  db := dbConn()
  var id int
  err = db.QueryRow(sqlStatement,article.Title, article.Description, article.Content).Scan(&id)
  if err != nil {
    panic(err)
  }
  fmt.Println("Created article with id: ", id)
}

func deleteArticleById(w http.ResponseWriter, r *http.Request ){
  id := mux.Vars(r)["id"]
  sqlStatement := fmt.Sprintf("delete from articles where id = %s", id)
  db := dbConn()
  res, err := db.Exec(sqlStatement)
  if err != nil {
    panic(err)
  }
  _, err = res.RowsAffected()
  if err != nil {
   w.WriteHeader(http.StatusBadRequest)
   w.Write([]byte("There was a problem"))
  }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Deleted the article")) 
}

func updateArticleById(w http.ResponseWriter, r *http.Request){
  id := mux.Vars(r)["id"]
  decoder := json.NewDecoder(r.Body)
  var article Article
  err := decoder.Decode(&article)
  fmt.Print("Article received: ", article)
  if err != nil {
    panic(err)
  }
  sqlStatement := "update articles set title=$1,description=$2,content=$3 where id=$4"
  db := dbConn()
  _, err = db.Exec(sqlStatement, article.Title,article.Description,article.Content,id)
  if err != nil {
    panic(err)
  }
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Updated the article"))
}

func handleRequests() {
  router := mux.NewRouter()
  router.HandleFunc("/",homePage).Methods(http.MethodGet)
  router.HandleFunc("/articles/", createArticle).Methods(http.MethodPost)
  router.HandleFunc("/articles/", returnAllArticles).Methods(http.MethodGet)
  router.HandleFunc("/articles/{id:[0-9]+}", updateArticleById).Methods(http.MethodPut)
  router.HandleFunc("/articles/{id:[0-9]+}", deleteArticleById).Methods(http.MethodDelete)
  log.Fatal(http.ListenAndServe(":8080", router))
}

func dbConn() *sql.DB{
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
    fmt.Println("Db string is: ", psqlInfo)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  return db
}

func main() {
  handleRequests()
}
