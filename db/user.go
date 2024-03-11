package db

import "time"

type User struct {
	Id        int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserType  int    `gorm:"type:tinyint" json:"user_type,omitempty"`
	Name      string `gorm:"size:100" json:"name,omitempty"`
	Address   string `gorm:"type:varchar(250)" json:"address"`
	Email     string `gorm:"type:varchar(100)" json:"email,omitempty"`
	Phone     string `gorm:"type:varchar(20); unique" json:"phone,omitempty"`
	Password  string `gorm:"type:varchar(200)" json:"password,omitempty"`
	IsDeleted bool   `gorm:"type:bool" json:"payment,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
