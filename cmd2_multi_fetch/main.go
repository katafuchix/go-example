package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Stock struct {
	ID   int
	Code string
	Name string
}

func main() {
	fmt.Println("Hello, World!")

	//db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/stock?parseTime=true&loc=Asia%2FTokyo")

	// 設定を構造体で作る
	c := mysql.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "stock",
		ParseTime:            true,
		Loc:                  time.Local, // システムのタイムゾーンを使用
		AllowNativePasswords: true,
	}

	// FormatDSN() で正しい文字列に変換してくれる
	db, err := sql.Open("mysql", c.FormatDSN())

	panic_err(err)
	defer db.Close()

	fmt.Println("接続成功！")

	// 1. Query を使う（複数を返すクエリ）
	rows, err := db.Query("SELECT id, code, name FROM stocks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // rowsも忘れずに閉じる

	// 2. ループで1行ずつ取り出す
	for rows.Next() {
		var s Stock
		// 1行分のデータを変数に流し込む
		err := rows.Scan(&s.ID, &s.Code, &s.Name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, コード: %s, 銘柄名: %s\n", s.ID, s.Code, s.Name)
	}

	// 1. 格納用のスライス（可変長配列）を準備
	var stocks []Stock

	// 1. Query を使う（複数を返すクエリ）
	rows, err = db.Query("SELECT id, code, name FROM stocks")

	for rows.Next() {
		// 毎回クリーンな変数を定義（スコープ優先！）
		var s Stock

		err := rows.Scan(&s.ID, &s.Code, &s.Name)
		if err != nil {
			log.Fatal(err)
		}
		// 2. スライスに追加（append）
		stocks = append(stocks, s)
	}

	// ループの外で、まとめて処理ができる
	fmt.Printf("全部で %d 件のデータを取得しました\n", len(stocks))

	// 3. ループ中にエラーが起きていなかったか最後にチェック
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func panic_err(err error) {
	if err != nil {
		panic(err.Error())
	}
}
