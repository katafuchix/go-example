package main

import (
	"fmt"
	"log"
	"time"

	// 1. 本家ドライバーに「mysqlDriver」という別名をつける
	// パッケージエイリアス（別名）
	mysqlDriver "github.com/go-sql-driver/mysql"
	// GORM本体
	"gorm.io/gorm"
	// GORM用のMySQLドライバ（これ自体に github.com/go-sql-driver/mysql が含まれています）
	"gorm.io/driver/mysql"
)

// GORMで使う構造体（タグを使って設定を追加できます）
type Stock struct {
	ID    uint   `gorm:"primaryKey"` // 主キーとして認識させる
	Code  string `gorm:"unique"`     // 重複を許さない設定
	Name  string
	Price int
}

func main() {
	// 1. 前と同じように Config を設定
	cfg := mysqlDriver.Config{
		User:                 "root",
		Passwd:               "",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "stock",
		ParseTime:            true,
		Loc:                  time.Local,
		AllowNativePasswords: true,
	}

	// 接続（GORM専用の書き方）
	//dsn := "root@tcp(127.0.0.1:3306)/stock?parseTime=true"
	//db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{})

	if err != nil {
		log.Fatal("データベース接続失敗:", err)
	}

	// 1. データの追加 (INSERT)
	//newStock := Stock{Code: "A100", Name: "消しゴム"}
	//db.Create(&newStock)

	// 2. データの取得 (SELECT)
	var s Stock
	db.First(&s, 1) // IDが1のデータを取得してsに入れる（Scan不要！）

	fmt.Println(s.ID)
	fmt.Println(s.Name)
	fmt.Println(s.Code)
}
