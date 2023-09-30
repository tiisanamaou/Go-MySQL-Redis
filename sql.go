package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	//_ "github.com/go-sql-driver/mysql"
)

// MySQLと接続する関数
func MysqlConnect() *sql.DB {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	config := mysql.Config{
		DBName:               "go_db",
		User:                 "GoUser",
		Passwd:               "gogopass",
		Addr:                 "127.0.0.1:3306",
		Collation:            "utf8mb4_unicode_ci",
		ParseTime:            true,
		AllowNativePasswords: true,
		Loc:                  jst,
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		fmt.Println(err.Error())
	}
	//
	fmt.Println("----------")
	err = db.Ping()
	if err != nil {
		fmt.Println("データベース接続失敗")
		return nil
	} else {
		fmt.Println("データベース接続成功")
	}
	return db
}

// INSERT
/*
新規レコード作成
*/
func InsertTodo(db *sql.DB, id int, uuid string, title string) (int64, error) {
	//var progress int = 5
	res, err := db.Exec(
		"INSERT INTO "+TableName+" (ID, Title) VALUES (?, ?)",
		id,
		//uuid,
		title,
		//progress,
	)
	if err != nil {
		fmt.Println("INSERT USER ERROR db.Exe err")
		fmt.Println(err)
		return 0, err
	}
	//
	insetID, err := res.LastInsertId()
	if err != nil {
		fmt.Println("INSERT USER res.LastInsertId error")
		fmt.Println(err)
		return 0, err
	}
	fmt.Println("レコードを追加しました")
	return insetID, nil
}

// UPDATE
/*
指定されたUUIDのレコードのIDとTitleを更新する
*/
func UpdateToDo(db *sql.DB, id int, title string, progress int) error {
	if _, err := db.Exec(
		//"UPDATE todo123 SET id = ?, title = ? WHERE uuid = ?",
		"UPDATE "+TableName+" SET id = ?, title = ?, progress = ? WHERE uuid = ?",
		id,
		title,
		progress,
	); err != nil {
		fmt.Println(err)
		fmt.Println("UPDATE エラー")
		return err
	}
	return nil
}

// DELETE
/*
指定されたUUIDのレコードを削除する
*/
func DeleteTodo(db *sql.DB, id int) error {
	if _, err := db.Exec(
		//"DELETE FROM todo123 WHERE uuid = ?",
		"DELETE FROM "+TableName+" WHERE uuid = ?",
		id,
	); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// 指定したタイトルのレコード数を数える関数
/*
引数：検索する文字列
戻り値：レコード数を返す
*/
func RecordTitleCount(db *sql.DB, title string) (int, error) {
	var recordCount int
	err := db.QueryRow(
		//"SELECT COUNT(*) FROM todoTable WHERE title LIKE CONCAT('%', ?, '%')",
		"SELECT COUNT(*) FROM "+TableName+" WHERE title LIKE CONCAT('%', ?, '%')",
		title,
	).
		Scan(
			&recordCount,
		)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("sql.ErrNoRows:指定したタイトルのレコードがありません")
		return recordCount, err
	}
	if err != nil {
		fmt.Println(err)
		fmt.Println("db.QueryRow:指定したタイトルのレコードはありません")
		return recordCount, err
	}
	//
	fmt.Print("レコード数：")
	fmt.Println(recordCount)
	if recordCount == 0 {
		fmt.Println("指定されたタイトルのレコードはありません")
		return recordCount, nil
	}
	return recordCount, nil
}

// 指定したUUIDのレコードが何件あるか調べる
/*
引数：UUID
戻り値：UUIDのレコード数を返す、更新するUUIDのレコードがあるか確認用
*/
func UuidCount(db *sql.DB, id int) (int, error) {
	var recordCount int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM "+TableName+" WHERE uuid = ?",
		id,
	).
		Scan(
			&recordCount,
		)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("レコードがありません")
		return 0, nil
	}
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	fmt.Print("レコード数：")
	fmt.Println(recordCount)
	if recordCount == 0 {
		fmt.Println("指定されたUUIDのレコードはありません")
		return recordCount, nil
	}
	return recordCount, nil
}

// 指定したタイトルを含むレコードの内容を全件取得する
/*
タイトルを指定した場合、指定したタイトルを含むレコードを全件取得する
タイトルを指定しない場合、レコードを全件取得する
ワイルドカードで検索する関数
SELECT * FROM xxx WHERE xxx LIKE CONCAT('%', ?, '%')
*/
func SearchToDo(db *sql.DB, searchWord string) {
	rows, err := db.Query(
		"SELECT * FROM "+TableName+" WHERE title LIKE CONCAT('%', ?, '%')",
		searchWord,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	var todolist []ToDo
	//
	for rows.Next() {
		todo := &ToDo{}
		if err := rows.Scan(
			&todo.ID,
			//&todo.Uuid,
			&todo.Title,
			&todo.Status, //Add
			&todo.CreateAt,
			&todo.UpdateAt,
		); err != nil {
			fmt.Println(err)
			fmt.Println("Scanエラー")
			return
		}
		//
		todolist = append(todolist, ToDo{
			ID: todo.ID,
			//Uuid:     todo.Uuid,
			Title:    todo.Title,
			Status:   todo.Status, // Add
			CreateAt: todo.CreateAt,
			UpdateAt: todo.UpdateAt,
		})
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("rows.Nextエラー")
		fmt.Println(err)
	}
	//
	ResToDo.ToDos = todolist
}

// ユーザーテーブルにユーザーデータがあるか検索する
func SearchUserPassword(db *sql.DB, userid string) (string, error) {
	rows, err := db.Query(
		"SELECT * FROM "+UserTableName+" WHERE userid = ?",
		userid,
	)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//
	loginUser := &LoginRequest{}
	//
	for rows.Next() {
		if err := rows.Scan(
			&loginUser.UserID,
			&loginUser.Password,
		); err != nil {
			fmt.Println(err)
			fmt.Println("Scanエラー")
			return "", err
		}
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("rows.Nextエラー")
		fmt.Println(err)
	}
	//fmt.Println(loginUser.UserID)
	//fmt.Println(loginUser.Password)
	return loginUser.Password, nil
}

// ユーザーテーブルにユーザーデータを登録する
func UserSignUp(db *sql.DB, userid string, password string) error {
	res, err := db.Exec(
		"INSERT INTO "+UserTableName+" (userid, signin_password) VALUES (?, ?)",
		userid,
		password,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//
	fmt.Println(res)
	return nil
}
