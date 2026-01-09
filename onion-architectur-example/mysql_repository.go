// 実際にMySQLを叩く処理を書く

package main

import "gorm.io/gorm"

type mysqlStockRepository struct {
	db *gorm.DB
}

// 設計図（インターフェース）に基づいて実装
func (r *mysqlStockRepository) Save(stock *Stock) error {
	return r.db.Create(stock).Error
}

func (r *mysqlStockRepository) FindAll() ([]Stock, error) {
	var stocks []Stock
	err := r.db.Find(&stocks).Error
	return stocks, err
}

// 構造体を新しく作るための関数
func NewMySQLStockRepository(db *gorm.DB) StockRepository {
	return &mysqlStockRepository{db: db}
}
