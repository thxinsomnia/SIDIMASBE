package models

type User struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Username string `gorm:"type varchar(255)" json:"username"`
	Email    string `gorm:"type varchar(255)" json:"email"`
	Password string `gorm:"type varchar(255)" json:"password"`
	Nama_Supplier string `gorm:"type varchar(255)" json:"nama_supplier"`
	Penyuplai string `gorm:"type varchar(255)" json:"supplier"`
	Nama_Perusahaan string `gorm:"type varchar(255)" json:"nama_perusahaan"`
} 