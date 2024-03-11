package db

import (
	"fmt"

	cfg "github.com/Safayet-Shawn/Parking-Lot/config"
	"github.com/golang-jwt/jwt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDB() {
	var err error
	arg := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True",
		cfg.Username(),
		cfg.Password(),
		cfg.Host(),
		cfg.DB(),
	)
	DB, err = gorm.Open("mysql", arg)
	if err != nil {
		panic(err)
	}
}
func GetDB() *gorm.DB {
	return DB
}

type Token struct {
	UserId   int
	UserType int
	Name     string
	Address  string
	jwt.StandardClaims
}
