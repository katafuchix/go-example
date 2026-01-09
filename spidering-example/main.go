package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {
	// 1. コレクター（クローラーの本体）を初期化
	c := colly.NewCollector(
		// 自分のドメイン以外に行かないように制限（スパイダリングの基本）
		colly.AllowedDomains("go.dev"),
	)

	// 2. HTMLの特定の要素（ここでは <a> タグ）を見つけた時の処理
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Printf("リンクを発見: %v\n", link)

		// 発見したリンクをさらに巡回する（スパイダリング）
		e.Request.Visit(link)
	})

	// 3. リクエスト送信前の処理
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("訪問中:", r.URL.String())
	})

	// 4. 最初にアクセスするURLを指定
	c.Visit("https://go.dev/")
}
