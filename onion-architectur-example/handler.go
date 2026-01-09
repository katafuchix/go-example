package main

import (
	"fmt"
	"net/http"
)

// StockHandler はリポジトリを使ってHTMLを返す関数
func StockHandler(repo StockRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. リポジトリを使ってDBから全件取得
		stocks, err := repo.FindAll()
		if err != nil {
			http.Error(w, "データ取得失敗", http.StatusInternalServerError)
			return
		}

		// 2. ブラウザに表示（簡易的な表示）
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<h1>在庫一覧</h1><ul>")
		for _, s := range stocks {
			fmt.Fprintf(w, "<li>ID: %d | 商品名: %s </li>", s.ID, s.Name)
		}
		fmt.Fprint(w, "</ul>")
	}
}
