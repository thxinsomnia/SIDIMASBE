package productcontroller

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"SIDIMASBE/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProducts retrieves all materials from the database
func GetMaterial(c *gin.Context) {
	var materials []models.Material

	// Fetch all materials from the database
	if err := models.DB.Find(&materials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Material": materials})
}

// GetProductByID retrieves a supplier by its ID from the database
func GetMaterialByID(c *gin.Context) {
	supplierID := c.Param("id") // Assuming the supplier ID is passed as a URL parameter

	var supplier models.Material

	// Find the supplier by ID
	if err := models.DB.First(&supplier, supplierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Material Tidak Ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Material": supplier})
}

// func Addmaterial(c *gin.Context) {
//     var userInput models.Material
//     if err := c.ShouldBindJSON(&userInput); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
//         return
//     }

//     if err := models.DB.Create(&userInput).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"Message": "Tambah Data Suppplier Berhasil"})
// }

// UpdateProduct updates an existing supplier in the database
func UpdateMaterial(c *gin.Context) {
	var userInput models.Material
	materialID := c.Param("id") // Assuming the supplier ID is passed as a URL parameter

	// Bind the JSON input to the userInput struct
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	// Find the supplier by ID
	var material models.Material
	if err := models.DB.First(&material, materialID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Material Tidak Ditemukan"})
		return
	}

	// Update the material fields
	material.Nama_bahan = userInput.Nama_bahan
	material.Jumlah = userInput.Jumlah
	material.Asal_bahan = userInput.Asal_bahan
	material.Kategori = userInput.Kategori
	material.Status = userInput.Status
	material.Tanggal = userInput.Tanggal
	// Add other fields as necessary

	// Save the updated material to the database
	if err := models.DB.Save(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Update Data Material Berhasil"})
}

// DeleteProduct deletes a material from the database
func DeleteMaterial(c *gin.Context) {
	materialID := c.Param("id") // Assuming the supplier ID is passed as a URL parameter

	// Find the supplier by ID
	var material models.Material
	if err := models.DB.First(&material, materialID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Material Tidak Ditemukan!"})
		return
	}

	// Delete the supplier from the database
	if err := models.DB.Delete(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Hapus Data Material Berhasil"})
}

func GenerateUserID() (string, error) {
	var lastUserID string
	err := models.DB.Model(&models.Material{}).Order("id_bahan desc").Limit(1).Pluck("id_bahan", &lastUserID).Error
	if err != nil {
		return "", err // Jika tidak ada supplier sebelumnya, ID pertama bisa dimulai dari angka 177013001
	}

	// Mengambil angka urut dari ID terakhir
	if len(lastUserID) < 9 {
		return "290220001", nil // Jika ID pertama kali, mulai dari 177013001
	}

	// Ambil bagian angka urut dari ID terakhir (3 digit terakhir)
	lastIDNumeric, err := strconv.Atoi(lastUserID[6:])
	if err != nil {
		return "", err
	}

	// Tambahkan 1 ke ID terakhir
	newIDNumeric := lastIDNumeric + 1
	newID := fmt.Sprintf("290220%03d", newIDNumeric) // Format 3 digit angka

	return newID, nil
	
}

func Addmaterial(c *gin.Context) {
	var userInput models.Material

	// Generate User ID baru dengan format yang diinginkan
	userID, err := GenerateUserID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to generate user ID"})
		return
	}

	userID64, err := strconv.ParseInt(userID, 10, 64)  // Mengonversi string ke int64
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Failed to convert user ID to int64"})
		return
	}

	// Assign ID yang baru ke userInput
	
	userInput.ID_supplier = userID64

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	// Simpan data supplier ke database
	if err := models.DB.Create(&userInput).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	// Respons sukses
	c.JSON(http.StatusOK, gin.H{"Message": "Tambah Data Material Berhasil"})
}



func checkSupplierExists(db *gorm.DB, supplierID int64) (bool, error) {
	var supplier models.Material
	result := db.First(&supplier, supplierID)
	if result.Error != nil {
		// Material tidak ditemukan
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	// Material ditemukan
	return true, nil
}

func addBahan(c *gin.Context) {
	var material models.Material

	// Bind JSON ke struct Material
	if err := c.ShouldBindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	// Validasi bahwa ID_supplier ada di tabel Material
	isExists, err := checkSupplierExists(models.DB, material.ID_supplier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error checking supplier"})
		return
	}
	if !isExists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Material does not exist"})
		return
	}

	// Simpan data Material ke database
	if err := models.DB.Create(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material berhasil ditambahkan"})
}