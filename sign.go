package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// 新規アカウント登録
func SignUp(w http.ResponseWriter, req *http.Request) {
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
	fmt.Println("Sign Up")
	//Request Body 取得
	contentLength := req.ContentLength
	contentBody := make([]byte, contentLength)
	req.Body.Read(contentBody)
	//
	var loginRequest LoginRequest
	err := json.Unmarshal(contentBody, &loginRequest)
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
	//
	err = UserSignUp(db, loginRequest.UserID, loginRequest.Password)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound) //404
		return
	}
	//
	//w.WriteHeader(http.StatusNoContent) //204
	w.WriteHeader(http.StatusOK) //200
	// Responseにデータを格納
	buf, _ := json.Marshal(loginRequest)
	_, _ = w.Write(buf)
	log.Println("アカウント作成完了")
}

// サインイン処理
func SignIn(w http.ResponseWriter, req *http.Request) {
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
	//Request Body 取得
	contentLength := req.ContentLength
	contentBody := make([]byte, contentLength)
	req.Body.Read(contentBody)
	//
	fmt.Println("Sign In")
	//
	var loginRequest LoginRequest
	err := json.Unmarshal(contentBody, &loginRequest)
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
	//
	userPassword, err := SearchUserPassword(db, loginRequest.UserID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusNotFound) //404
		return
	}
	//fmt.Println(userPassword)
	if loginRequest.Password != userPassword {
		log.Println("パスワード不一致")
		w.WriteHeader(http.StatusNotFound) //404
		return
	}
	log.Println("パスワード一致")
	JwtGenerater(loginRequest.UserID)
	JwtVerify(JwtString)
	//
	//w.Header().Set("JWT", JwtString)
	//
	//w.WriteHeader(http.StatusNoContent) //204
	w.WriteHeader(http.StatusOK) //200
	// Responseにデータを格納
	buf, _ := json.Marshal(loginRequest)
	_, _ = w.Write(buf)
	log.Println("サインイン完了")
}

// サインイン後のリクエスト処理
