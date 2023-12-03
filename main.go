package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ToDo struct {
	ID       int       `json:"id"`
	UserID   string    `json:"userid"`
	Title    string    `json:"title"`
	Status   int       `json:"status"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}

type LoginRequest struct {
	UserID   string `json:"UserID"`
	Password string `json:"Password"`
}

type UserRequest struct {
	ID       int       `json:"id"`
	Password string    `json:"Password"`
	Title    string    `json:"title"`
	Status   int       `json:"status"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
}

type Response struct {
	ToDos []ToDo `json:"todos"`
}

var ResToDo Response

// 接続するテーブル名を記載、SQL文で使用する
const TableName string = "todoTable"
const UserTableName string = "userTable"

func main() {
	fmt.Println("Go MySQL")

	http.HandleFunc("/get", GETMethod)
	http.HandleFunc("/post", POSTMethod)
	http.HandleFunc("/put", PUTMethod)
	http.HandleFunc("/delete", DELETEMethod)
	http.HandleFunc("/signup", SignUp)
	http.HandleFunc("/signin", SignIn)
	http.HandleFunc("/signin-get", SigninGET)
	http.HandleFunc("/signin-post", SigninPOST)
	http.HandleFunc("/signin-put", SigninPUT)
	http.HandleFunc("/signin-delete", SigninDELETE)
	http.ListenAndServe(":8080", nil)
}

// HTTP GET Method, READ
func GETMethod(w http.ResponseWriter, req *http.Request) {
	// Methodの確認
	if req.Method != "GET" {
		//w.WriteHeader(http.StatusBadRequest) //400
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	// レスポンスヘッダーの設定
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// CROSエラーが出ないようにする設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//
	fmt.Println("GET Method Search")
	// クエリパラメータを抽出
	req.FormValue("search")
	fmt.Println(req.FormValue("search"))
	QueryPara := req.FormValue("search")
	// MySQL 接続
	db := MysqlConnect()
	// MySQLに接続できなかったなら503を返す
	if db == nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	defer db.Close()
	// クエリパラメータで検索
	recordCount, err := RecordTitleCount(db, QueryPara)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	// 指定したタイトルのレコードがなかったら404を返す
	if recordCount == 0 {
		w.WriteHeader(http.StatusNotFound) //404
		fmt.Println("404エラー")
		return
	}
	// 指定したタイトルのレコードを取得する
	SearchToDo(db, QueryPara)
	// ステータスコードの設定
	w.WriteHeader(http.StatusOK) //200
	// Responseにデータを格納
	buf, _ := json.Marshal(ResToDo)
	_, _ = w.Write(buf)
	log.Println("取得完了")
}

// HTTP POST Method, CREATE
func POSTMethod(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	// レスポンスヘッダーの設定
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// CROSエラーが出ないようにする設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//
	fmt.Println("POST Method")
	//
	// POST Request body data
	contentLength := req.ContentLength
	contentBody := make([]byte, contentLength)
	req.Body.Read(contentBody)
	//
	var todo ToDo
	err := json.Unmarshal(contentBody, &todo)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	//
	// MySQLに接続できなかったなら503を返す
	db := MysqlConnect()
	if db == nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	defer db.Close()
	// uuid
	//uu := UuidGenerate()
	_, err = InsertTodo(db, todo.UserID, todo.Title, 50)
	//
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		fmt.Println("id が重複しています")
		return
	}
	// ステータスコード
	w.WriteHeader(http.StatusNoContent) //204
	log.Println("作成完了")
}

// HTTP PUT Method, UPDATE
func PUTMethod(w http.ResponseWriter, req *http.Request) {
	if req.Method != "PUT" {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	// レスポンスヘッダーの設定
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// CROSエラーが出ないようにする設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//
	fmt.Println("PUT Method")
	//Request Body 取得
	contentLength := req.ContentLength
	contentBody := make([]byte, contentLength)
	req.Body.Read(contentBody)
	//
	var todo ToDo
	err := json.Unmarshal(contentBody, &todo)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound) //404
		return
	}
	// MySQLに接続できなかったなら503を返す
	db := MysqlConnect()
	if db == nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	defer db.Close()

	// 指定したUUIDのレコードがあるか確認
	// なかったら404を返す
	//recordCount, _ := UuidCount(db, todo.ID)
	//if recordCount == 0 {
	//	w.WriteHeader(http.StatusNotFound) //404
	//	return
	//}

	// SQLにデータをUPDATEする
	// エラーが発生したら503を返す
	err = UpdateToDo(db, todo.ID, todo.Title, todo.Status)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	w.WriteHeader(http.StatusNoContent) //204
	log.Println("更新完了")
}

// HTTP DELETE Method, DELETE
func DELETEMethod(w http.ResponseWriter, req *http.Request) {
	if req.Method != "DELETE" {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}
	// レスポンスヘッダーの設定
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// CROSエラーが出ないようにする設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//
	fmt.Println("DELETE Method")
	//Request Body 取得
	contentLength := req.ContentLength
	contentBody := make([]byte, contentLength)
	req.Body.Read(contentBody)
	//
	var todo ToDo
	err := json.Unmarshal(contentBody, &todo)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound) //404
		return
	}
	//
	// MySQLに接続できなかったなら503を返す
	db := MysqlConnect()
	if db == nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	defer db.Close()

	// 指定したUUIDのレコードがあるか確認
	//recordCount, _ := UuidCount(db, todo.ID)
	//if recordCount == 0 {
	//	w.WriteHeader(http.StatusNotFound) //404
	//	return
	//}

	// 指定UUIDのレコードを削除
	DeleteTodo(db, todo.ID)

	// ステータスコード
	w.WriteHeader(http.StatusNoContent) //204
	//fmt.Println("UUID:", todo.Uuid)
	log.Println("削除完了")
}

// UUID を作成する
func UuidGenerate() string {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	uu := u.String()
	fmt.Println(uu)
	return uu
}
