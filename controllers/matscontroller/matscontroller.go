package matscontroller

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"SIDIMASBE/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

//fungsi ini untuk mendapatkan semua data bahan
func GetAllBahan(c *gin.Context) {
    var materials []models.Material

    if err := models.DB.Preload("Supplier").Find(&materials).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Materials": materials})
}

//fungsi ini untuk mendapatkan data bahan hanya berdasarkan id
func GetBahanByID(c *gin.Context) {
    materialID := c.Param("id") // Ambil ID dari URL

    var material models.Material
    if err := models.DB.Preload("Supplier").First(&material, materialID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"Message": "Material tidak ditemukan!"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Material": material})
}

//fungsi ini untuk edit bahan
func EditBahan(c *gin.Context) {
    materialID := c.Param("id") // ID dari URL
    var userInput models.Material

    // Bind JSON ke struct
    if err := c.ShouldBindJSON(&userInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid request data: " + err.Error()})
        return
    }

    var material models.Material
    if err := models.DB.First(&material, materialID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"Message": "Material tidak ditemukan"})
        return
    }

    // Update data bahan
    material.Nama_bahan = userInput.Nama_bahan
    material.Jumlah = userInput.Jumlah
    material.Asal_bahan = userInput.Asal_bahan
    material.Kategori = userInput.Kategori
    material.Status = userInput.Status
    material.Tanggal = userInput.Tanggal
    material.ID_supplier = userInput.ID_supplier

    if err := models.DB.Save(&material).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Message": "Update Data Material Berhasil"})
}

//fungsi ini utnuk hapus bahan dari database
func HapusBahan(c *gin.Context) {
    materialID := c.Param("id") // Ambil ID dari URL

    var material models.Material
    if err := models.DB.First(&material, materialID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"Message": "Material tidak ditemukan!"})
        return
    }

    if err := models.DB.Delete(&material).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Message": "Hapus Data Material Berhasil"})
}

// Fungsi Ini Untuk Cek Apakah Supplier Ada PAda Data
func checkSupplierExists(db *gorm.DB, supplierID int64) (bool, error) {
    var supplier models.Supplier // Pastikan kita mengecek di tabel Supplier, bukan Material
    result := db.First(&supplier, "id_supplier = ?", supplierID)
    
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return false, nil // Supplier tidak ditemukan
        }
        return false, result.Error
    }
    return true, nil // Supplier ditemukan
}

//Fungsi Ini Untuk Menmabhak Data
func Addbahan(c *gin.Context) {
	var userInput models.Material

    if err := c.ShouldBindJSON(&userInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
        return
    }

	if userInput.ID_bahan == 0 {
        newID, err := GenerateMatsID()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to generate Material ID"})
            return
        }
        userInput.ID_bahan, _ = strconv.ParseInt(newID, 10, 64) // Konversi ke int64
    }

    // Cek apakah supplier dengan ID tersebut ada
    isExists, err := checkSupplierExists(models.DB, userInput.ID_supplier)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Error checking supplier"})
        return
    }
    if !isExists {
        c.JSON(http.StatusBadRequest, gin.H{"Message": "Supplier tidak ditemukan!"})
        return
    }

    // Simpan data material ke database
    if err := models.DB.Create(&userInput).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Message": "Tambah Data Material Berhasil"})
}

//fungsi ini untuk custom id
func GenerateMatsID() (string, error) {
    var lastMatsID string

    // Ambil ID terakhir dari database
    err := models.DB.Model(&models.Material{}).
        Select("id_bahan").
        Order("id_bahan DESC").
        Limit(1).
        Pluck("id_bahan", &lastMatsID).
        Error

    if err != nil || lastMatsID == "" {
        // Jika tidak ada data sebelumnya, mulai dari 3399001
        return "3399001", nil
    }

    // Pastikan ID terakhir memiliki panjang yang sesuai
    if len(lastMatsID) < 7 {
        return "3399001", nil
    }

    // Ambil bagian angka urut (3 digit terakhir)
    lastIDNumeric, err := strconv.Atoi(lastMatsID[4:]) // Ambil setelah "3399"
    if err != nil {
        return "", err
    }

    // Tambahkan 1 ke ID terakhir
    newIDNumeric := lastIDNumeric + 1
    newID := fmt.Sprintf("3399%03d", newIDNumeric) // Format: "3399XXX"

    return newID, nil
}

func TambahStokBahan(c *gin.Context) {
	var input struct {
		ID_bahan int64 `json:"id_bahan"`
		Jumlah   int64 `json:"jumlah"`
	}

	// Bind JSON input ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	// Cek apakah bahan dengan ID tersebut ada
	var material models.Material
	if err := models.DB.First(&material, input.ID_bahan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Material tidak ditemukan"})
		return
	}

	// Update stok bahan
	material.Jumlah += input.Jumlah

	if err := models.DB.Save(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menambahkan stok"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stok berhasil ditambahkan", "data": material})
}

