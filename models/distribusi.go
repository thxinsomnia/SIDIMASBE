package models

import "time"

type Distribusi struct {
    ID_distribusi  int64     `gorm:"primaryKey" json:"id_distrib"`
    ID_menu        int64     `json:"id_menu"`
    Nama_menu      string    `json:"nama_menu"`
    Jumlah_kirim   int64     `json:"jumlah_kirim"`
    Alamat_tujuan  string    `json:"alamat_tujuan"`
    Status         string    `json:"status"`
    Tanggal_kirim  time.Time `json:"tanggal_kirim"`
}