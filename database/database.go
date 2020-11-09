package database
import (
		"github.com/jinzhu/gorm"
	Gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gogam/database/model"
)

func InitDB() *Gorm.DB {
	
	//db, err := gorm.Open("sqlite3", ":memory:")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Character{})
	db.AutoMigrate(&model.Ability{})
	return db
}