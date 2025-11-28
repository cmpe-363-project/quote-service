package repository

import "gorm.io/gorm"

type Quote struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	Message  string `gorm:"type:text;not null"`
	AuthorID int    `gorm:"not null"`
}

func (Quote) TableName() string {
	return "quotes"
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Quote{},
	)
}
