// domain.go: データの形（構造体）を定義する
package main

// Domain: 在庫データそのもの
type Stock struct {
	ID   uint   `gorm:"primaryKey"` // 主キーとして認識させる
	Code string `gorm:"unique"`     // 重複を許さない設定
	Name string
	/*gorm.Model
	Code  string
	Name  string
	Price int*/
}
