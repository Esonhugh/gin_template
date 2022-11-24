package database

import (
	"gorm.io/gorm"
)

var GlobalDatabase GlobalDB

type GlobalDB struct {
	MainDB *gorm.DB
}

func Init() error {
	return nil
}
