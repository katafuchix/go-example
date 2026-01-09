// 全てを組み立てて実行する

package main

import (
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// --- 準備層 ---
	dsn := "root@tcp(127.0.0.1:3306)/stock?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 1. 具体的なリポジトリを作成
	repo := NewMySQLStockRepository(db)

	// --- Webサーバー層 ---
	// 2. ハンドラにリポジトリを「注入」してルーティング設定
	http.HandleFunc("/stocks", StockHandler(repo))

	log.Println("サーバー起動: http://localhost:8080/stocks")

	// 3. 起動！
	log.Fatal(http.ListenAndServe(":8080", nil))

	/*
		// DB接続
		dsn := "root@tcp(127.0.0.1:3306)/stock?parseTime=true"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}

		// 1. 具体的な「MySQL版の道具」を作る
		repo := NewMySQLStockRepository(db)

		// 2. 道具を使ってデータを保存（mainはSQLの中身を知らなくて良い！）
		//newStock := &Stock{Code: "A101", Name: "ボールペン", Price: 100}
		//repo.Save(newStock)

		// 3. 取得して表示
		stocks, _ := repo.FindAll()
		for _, s := range stocks {
			fmt.Printf("ID: %d, 商品名: %s\n", s.ID, s.Name)
		}
	*/
}
