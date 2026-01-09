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

	// 取得したいデータのID
	targetID := 1
	var s Stock

	// クエリの実行と変数へのマッピング
	// Scanに渡す順番は、SELECTで指定した順番と同じにする必要があります
	err = db.QueryRow("SELECT id, code, name FROM stocks WHERE id = ?", targetID).Scan(&s.ID, &s.Code, &s.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("対象のデータが見つかりませんでした")
		} else {
			log.Fatal(err)
		}
		return
	}

	fmt.Printf("取得結果: ID=%d, 名前=%s\n", s.ID, s.Name)
}

func panic_err(err error) {
	if err != nil {
		panic(err.Error())
	}
}
