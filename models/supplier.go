package models

type Supplier struct {
	ID_supplier    int64   `gorm:"primaryKey" json:"id_supplier"`
	Nama_supplier  string  `gorm:"type varchar(255)" json:"nama_supplier"`
	Alamat  string  `gorm:"type varchar(255)" json:"alamat"`
	Kontak string `gorm:"type varchar(255)" json:"kontak"`
	Sertifikasi string `gorm:"type varchar(255)" json:"sertifikasi"`
	Verifikasi string `gorm:"type varchar(255)" json:"verifikasi"`

	Materials []Material `gorm:"foreignKey:ID_supplier;references:ID_supplier" json:"bahans"`
	
}
