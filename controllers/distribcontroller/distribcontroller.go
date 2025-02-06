package distribcontroller

import (
	"net/http"
	"SIDIMASBE/models"
	"github.com/gin-gonic/gin"
	"fmt"
	"time"
)

func DistribusiMenu(c *gin.Context) {
    var distribusiInput struct {
        IDMenu       int64  `json:"id_menu"`
        JumlahKirim  int64  `json:"jumlah_kirim"`
        AlamatTujuan string `json:"alamat_tujuan"`
    }

    if err := c.ShouldBindJSON(&distribusiInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
        return
    }

    tx := models.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"Message": "Terjadi kesalahan internal"})
        }
    }()

    // Cek apakah menu ada
    var menu models.Menu
    if err := tx.First(&menu, "id_menu = ?", distribusiInput.IDMenu).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"Message": "Menu tidak ditemukan"})
        return
    }

    // Cek apakah stok cukup
    if menu.Jumlah_porsi < distribusiInput.JumlahKirim {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{
            "Message": fmt.Sprintf("Stok menu tidak mencukupi. Tersedia %d, ingin mengirim %d.",
                menu.Jumlah_porsi, distribusiInput.JumlahKirim),
        })
        return
    }

    // Kurangi stok menu
    menu.Jumlah_porsi -= distribusiInput.JumlahKirim
    if err := tx.Save(&menu).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengupdate stok menu"})
        return
    }

    // Buat catatan distribusi
    distribusi := models.Distribusi{
        ID_menu:       menu.ID_menu,
        Nama_menu:     menu.Nama_menu,
        Jumlah_kirim:  distribusiInput.JumlahKirim,
        Alamat_tujuan: distribusiInput.AlamatTujuan,
        Status:        "Diterima",
        Tanggal_kirim: time.Now(),
    }

    if err := tx.Create(&distribusi).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mencatat distribusi"})
        return
    }

    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "Message": "Distribusi berhasil dicatat",
        "Data": distribusi,
        "SisaStokMenu": menu.Jumlah_porsi, // Menampilkan stok menu terbaru
    })
}

func AmbilDataDistribusi(c *gin.Context) {
    var distribusi []models.Distribusi

    if err := models.DB.Find(&distribusi).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengambil data distribusi"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "Message": "Data distribusi berhasil diambil",
        "Data": distribusi,
    })
}

