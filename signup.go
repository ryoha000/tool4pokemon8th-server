package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

type LoginRequestBody struct {
	Name string `json:"username" form:"username"`
	Pass string `json:"password" form:"password"`
}

type User struct {
	ID   int    `json:"id"  db:"id"`
	Name string `json:"name"  db:"name"`
	Pass string `json:"-"  db:"pass"`
}

func signup(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		w.Write([]byte("only POST"))
		fmt.Fprint(w, "This method allow only POST")
		return nil
	}
	_db, err := sqlx.Connect("mysql", "root:ryoha0216@/mydb")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("db is bad"))
		fmt.Fprint(w, "Cannot connect database")
		return err
	}
	db = _db

	if err = db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("db is bad"))
		fmt.Fprint(w, "Database tyotto umaku ittenai")
		return err
	}

	req := LoginRequestBody{}

	body, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Println(e.Error())
		return e
	}
	e = json.Unmarshal(body, &req)
	if e != nil {
		fmt.Println(e.Error())
		return e
	}
	// もう少し真面目にバリデーションするべき
	if req.Pass == "" || req.Name == "" {
		// エラーは真面目に返すべき
		fmt.Fprint(w, "form is empty")
		return nil
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Pass), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}

	// ユーザーの存在チェック
	var count int

	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE name=?", req.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("db is bad"))
		fmt.Fprint(w, "Database tyotto umaku ittenai")
		return err
	}

	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("already used this name"))
		fmt.Fprint(w, "already used this name")
		return nil
	}

	_, err = db.Exec("INSERT INTO users (name, pass) VALUES (?, ?)", req.Name, hashedPass)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("db is bad"))
		fmt.Fprint(w, "Database tyotto umaku ittenai")
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("your account is created!"))
	fmt.Fprint(w, "your account is created!")
	return nil
}
