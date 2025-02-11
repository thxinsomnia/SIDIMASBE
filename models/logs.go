package models

import "time"

type Log struct {
    ID_log          int64     `gorm:"primaryKey" json:"id_log"`
    ID_bahan        int64     `json:"id_bahan"`
    Tanggal         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"tanggal"`
    JumlahDigunakan int64     `json:"jumlah_digunakan"`
    SisaBahan        int64     `json:"sisa_bahan"`
}