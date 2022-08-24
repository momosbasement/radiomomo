package db

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	model "github.com/momosbasement/radiomomo/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Connection struct {
	db *gorm.DB
}

func (c *Connection) Init(dsn string) {
	c.db, _ = gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	c.db.AutoMigrate(&model.Track{})

}

func (conn *Connection) GetConnection() *gorm.DB {
	return conn.db
}
