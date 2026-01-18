package model

import "time"

type User struct {
	// 複数のタグを半角スペースで区切って並べる
	ID        uint      `json:"id" param:"id" gorm:"primaryKey;column:id"`
	Name      string    `json:"name"           gorm:"column:name"`
	CreatedAt time.Time `json:"created_at"    gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at"    gorm:"column:updated_at"`
}

/*
type User struct {
	//                                  [DB用]               [API用]         [Echoパス用]   [HTMLフォーム用]
	ID        uint      `gorm:"primaryKey"      json:"id"         param:"id"     form:"id"`
	Name      string    `gorm:"column:nickname" json:"name"                      form:"name"`
	Email     string    `gorm:"uniqueIndex"     json:"email"                     form:"email"`
	CreatedAt time.Time `gorm:"column:created"  json:"created_at"`
}
*/
