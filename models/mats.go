package models

import "time"

type Material struct {
	ID_bahan      int64  `gorm:"primaryKey" json:"id_bahan"`
	Nama_bahan    string `gorm:"type varchar(255)" json:"nama_bahan"`
	Jumlah        int64 `gorm:"type varchar(255)" json:"jumlah"`
	Asal_bahan    string `gorm:"type varchar(255)" json:"asal_bahan"`
	Kategori      string `gorm:"type varchar(255)" json:"kategori"`
	Status        string `gorm:"type varchar(255)" json:"status"`
	Tanggal       time.Time `gorm:"type timestamp" json:"tanggal"`

	//input foreign key untuk di tabel
	ID_supplier int64 `gorm:"type:int;not null" json:"id_supplier"`
	Supplier Supplier `gorm:"foreignKey:ID_supplier;references:id_supplier" json:"supplier"`
}