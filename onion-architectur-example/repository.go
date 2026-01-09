// repository.go: DBを操作する「道具」の設計図（インターフェース）を作る

package main

// Repository: 「在庫を保存する」「在庫を取得する」という機能の定義
// 具体的にMySQLを使うかどうかは、ここでは気にしません。
type StockRepository interface {
	Save(stock *Stock) error
	FindAll() ([]Stock, error)
}
