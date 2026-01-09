package main

type Stock struct {
	Code string         `gorm:"column:code"`           //
	ID   int            `gorm:"column:id;primary_key"` //
	Memo sql.NullString `gorm:"column:memo"`           //
	Name string         `gorm:"column:name"`           //
}

// TableName sets the insert table name for this struct type
func (s *Stock) TableName() string {
	return "stocks"
}

