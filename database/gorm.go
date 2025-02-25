package database

import (
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetGormDB() (*gorm.DB, error) {
	hostname := os.Getenv("RDS_HOSTNAME")
	port := os.Getenv("RDS_PORT")
	username := os.Getenv("RDS_USERNAME")
	password := os.Getenv("RDS_PASSWORD")
	dbname := os.Getenv("RDS_DB_NAME")

	dsn := username + ":" + password + "@tcp(" + hostname + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=UTC"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
}
