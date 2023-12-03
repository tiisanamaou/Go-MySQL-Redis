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

	w.Header().Set("JWT", JwtString)

	//w.WriteHeader(http.StatusNoContent) //204
	w.WriteHeader(http.StatusOK) //200
	// Responseにデータを格納
	buf, _ := json.Marshal(loginRequest)
	_, _ = w.Write(buf)

	// RedisにUserIDをKEYにしてJWTをValueとして登録する
	rdb := RedisConnect()
	RedisDataSet(rdb, loginRequest.UserID, JwtString)

	log.Println("サインイン完了")
}

// サインイン後のリクエスト処理
/*
JWTがあっている場合のみリクエストを受け付ける
*/
func SigninGET(w http.ResponseWriter, req *http.Request) {
	// Methodの確認
	if req.Method != "GET" {
		//w.WriteHeader(http.StatusBadRequest) //400
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}

	jwtHeader := req.Header.Get("JWT")
	// JWTの検証
	userid, err := JwtVerify(jwtHeader)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}

	db := MysqlConnect()
	// MySQLに接続できなかったなら503を返す
	if db == nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	defer db.Close()

	rdb := RedisConnect()
	//RedisDataSet(rdb, userid, jwtHeader)
	jwtRedis := RedisDataGet(rdb, userid)
	if jwtRedis == "" {
		fmt.Println("有効期限切れ:TTL")
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	if jwtRedis != jwtHeader {
		fmt.Println("Redisエラー")
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	// 問題なければRedis上のTTLを更新する
	RedisDataSet(rdb, userid, jwtHeader)

	// クエリパラメータの値を取得
	req.FormValue("search")
	QueryPara := req.FormValue("search")
	// 特定のユーザーのToDoを取得
	SearchUserToDo(db, userid, QueryPara)

	// ステータスコードの設定
	w.WriteHeader(http.StatusOK) //200
	// Responseにデータを格納
	buf, _ := json.Marshal(ResToDo)
	_, _ = w.Write(buf)

	fmt.Println("----------")
}

// サインイン後のリクエスト処理
/*
JWTがあっている場合のみリクエストを受け付ける
*/
func SigninPOST(w http.ResponseWriter, req *http.Request) {
	// Methodの確認
	if req.Method != "POST" {
		//w.WriteHeader(http.StatusBadRequest) //400
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}

	jwtHeader := req.Header.Get("JWT")
	// JWTの検証
	userid, err := JwtVerify(jwtHeader)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}

	db := MysqlConnect()
	// MySQLに接続できなかったなら503を返す
	if db == nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	defer db.Close()

	rdb := RedisConnect()
	//RedisDataSet(rdb, userid, jwtHeader)
	jwtRedis := RedisDataGet(rdb, userid)
	if jwtRedis == "" {
		fmt.Println("有効期限切れ:TTL")
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	if jwtRedis != jwtHeader {
		fmt.Println("Redisエラー")
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	// 問題なければRedis上のTTLを更新する
	RedisDataSet(rdb, userid, jwtHeader)

	POSTMethod(w, req)

	fmt.Println("----------")
}

// サインイン後のリクエスト処理
/*
JWTがあっている場合のみリクエストを受け付ける
*/
func SigninPUT(w http.ResponseWriter, req *http.Request) {
	// Methodの確認
	if req.Method != "PUT" {
		//w.WriteHeader(http.StatusBadRequest) //400
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}

	jwtHeader := req.Header.Get("JWT")
	// JWTの検証
	userid, err := JwtVerify(jwtHeader)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}

	db := MysqlConnect()
	// MySQLに接続できなかったなら503を返す
	if db == nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	defer db.Close()

	rdb := RedisConnect()
	//RedisDataSet(rdb, userid, jwtHeader)
	jwtRedis := RedisDataGet(rdb, userid)
	if jwtRedis == "" {
		fmt.Println("有効期限切れ:TTL")
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	if jwtRedis != jwtHeader {
		fmt.Println("Redisエラー")
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	// 問題なければRedis上のTTLを更新する
	RedisDataSet(rdb, userid, jwtHeader)

	PUTMethod(w, req)

	fmt.Println("----------")
}

// サインイン後のリクエスト処理
/*
JWTがあっている場合のみリクエストを受け付ける
*/
func SigninDELETE(w http.ResponseWriter, req *http.Request) {
	// Methodの確認
	if req.Method != "DELETE" {
		//w.WriteHeader(http.StatusBadRequest) //400
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}

	jwtHeader := req.Header.Get("JWT")
	// JWTの検証
	userid, err := JwtVerify(jwtHeader)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}

	db := MysqlConnect()
	// MySQLに接続できなかったなら503を返す
	if db == nil {
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	defer db.Close()

	rdb := RedisConnect()
	//RedisDataSet(rdb, userid, jwtHeader)
	jwtRedis := RedisDataGet(rdb, userid)
	if jwtRedis == "" {
		fmt.Println("有効期限切れ:TTL")
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	if jwtRedis != jwtHeader {
		fmt.Println("Redisエラー")
		w.WriteHeader(http.StatusServiceUnavailable) //503
		return
	}
	// 問題なければRedis上のTTLを更新する
	RedisDataSet(rdb, userid, jwtHeader)

	DELETEMethod(w, req)

	fmt.Println("----------")
}
