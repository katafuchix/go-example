package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// HTMLファイルを解析
		tmpl := template.Must(template.ParseFiles("index.html"))

		// HTMLに渡すデータ
		data := map[string]string{"Status": "稼働中"}

		// データを入れて表示
		tmpl.Execute(w, data)
	})

	http.ListenAndServe(":8080", nil)
}

/*package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 1. ルーティングの設定（どのURLにアクセスしたら、どの関数を動かすか）
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello! これはGoで作ったWebサーバーです。")
	})

	fmt.Println("サーバーを起動しました: http://localhost:8080")

	// 2. サーバーの起動（ポート番号 8080 で待ち受け）
	http.ListenAndServe(":8080", nil)
}
*/
