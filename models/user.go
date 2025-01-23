package models

type User struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type varchar(255)" json:"username"`
	Email    string `gorm:"type varchar(255)" json:"email"`
	Password string `gorm:"type varchar(255)" json:"password"`
	Role string `gorm:"type varchar(255)" json:"role"`
} 