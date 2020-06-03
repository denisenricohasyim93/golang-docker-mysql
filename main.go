package main

import (
    "database/sql"
    "log"
    "net/http"
    "encoding/json"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

type Employee struct {
    Id    int `form:"id" json:"id"`
    Name  string `form:"Name" json:"Name"`
    City string `form:"City" json:"City"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Employee
}

type ResponseShow struct {
        Status  int    `json:"status"`
        Message string `json:"message"`
        Data    Employee
}


func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    db, err := sql.Open(dbDriver, "root:supersecret@tcp(172.17.0.2:3306)/goblog")
    if err != nil {
        panic(err.Error())
    }
    return db
}

func Index(w http.ResponseWriter, r *http.Request) {
    var response Response
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM employee ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
    }
    emp := Employee{}
    res := []Employee{}
    for selDB.Next() {
        var id int
        var name, city string
        err = selDB.Scan(&id, &name, &city)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Name = name
        emp.City = city
        res = append(res, emp)
    }

    response.Status = 1
    response.Message = "Success"
    response.Data = res
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
    var response ResponseShow
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := Employee{}
    for selDB.Next() {
        var id int
        var name, city string
        err = selDB.Scan(&id, &name, &city)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Name = name
        emp.City = city
    }
    // tmpl.ExecuteTemplate(w, "Show", emp)

    response.Status = 1
    response.Message = "Success"
    response.Data = emp
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
    var u Employee

    if r.Body == nil {
        http.Error(w, "Please send a request body", 400)
        return
    }

    err := json.NewDecoder(r.Body).Decode(&u)

    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }
    fmt.Println(u.Name)
    fmt.Println(u.City)

    db := dbConn()
    if r.Method == "POST" {
        name := u.Name
        city := u.City
        insForm, err := db.Prepare("INSERT INTO employee(name, city) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name, city)
        log.Println("INSERT: Name: " + name + " | City: " + city)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main() {
    log.Println("Server started on: http://localhost:8080")
    http.HandleFunc("/", Index)
    http.HandleFunc("/show", Show)
    http.HandleFunc("/insert", Insert)
    http.ListenAndServe(":8080", nil)
}
