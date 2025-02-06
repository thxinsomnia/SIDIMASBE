package models

type Menu struct {
	ID_menu      int64  `gorm:"primaryKey" json:"id_menu"`
	Nama_menu    string `gorm:"type varchar(255)" json:"nama_menu"`
	Deskripsi    string `gorm:"type varchar(255)" json:"deskripsi"`
	Jumlah_porsi int64  `gorm:"type varchar(255)" json:"jumlah_porsi"`

	Bahan []MenuBahan `gorm:"foreignKey:ID_menu" json:"bahan"`
}

type MenuBahanItem struct {
    ID_bahan  int64 `json:"id_bahan"`
    Kebutuhan int64 `json:"kebutuhan"`
}

