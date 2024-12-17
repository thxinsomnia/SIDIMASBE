package models

type product struct {
	ID    int64   `gorm:"primaryKey" json:"id"`
	Nama  string  `gorm:"type varchar(255)" json:"nama"`
	Stok  int32   `gorm:"type int(5)" json:"stok"`
	Harga float32 `gorm:"type decimal(14,2)" json:"harga"`
}
