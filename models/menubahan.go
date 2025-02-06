package models

type MenuBahan struct {
    ID_menu         int64 `gorm:"primaryKey" json:"id_menu"`
    ID_bahan        int64 `gorm:"primaryKey" json:"id_bahan"`
    Kebutuhan int64 `gorm:"type:int" json:"kebutuhan"`

    // Relasi ke tabel bahan
    Materials Material `gorm:"foreignKey:ID_bahan;references:id_bahan" json:"bahan"`
}