package bootstrap

import (
	"log"

	"github.com/qxsugar/bill/api/config"
	"github.com/qxsugar/bill/api/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	cfg := config.GetDBConfig()
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)

	if err := db.AutoMigrate(
		&model.User{},
		&model.Room{},
		&model.RoomMember{},
		&model.Transaction{},
		&model.RoomLog{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	DB = db
	log.Println("database connected successfully")
}
