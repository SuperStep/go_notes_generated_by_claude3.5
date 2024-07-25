package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Name    string
	Content string
	UserID  uint
}
