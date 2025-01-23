package suppliercontroller

import (
	"fmt"
	"strconv"

	"net/http"
	"github.com/gin-gonic/gin"
	"SIDIMASBE/models"
)


//ambil semua supplier
func GetSupplier(c *gin.Context) {
	var supplier []models.Supplier

	// Fetch all supplier from the database
	if err := models.DB.Find(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Products": supplier})
}

// GetProductByID retrieves a supplier by its ID from the database
func GetSupplierByID(c *gin.Context) {
	supplierID := c.Param("id") // Assuming the supplier ID is passed as a URL parameter

	var supplier models.Supplier

	// Find the supplier by ID
	if err := models.DB.First(&supplier, supplierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Supplier Tidak Ditemukan!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Supplier": supplier})
}

// func Addsupplier(c *gin.Context) {
//     var userInput models.Supplier
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
func UpdateSupplier(c *gin.Context) {
	var userInput models.Supplier
	supplierID := c.Param("id") // Assuming the supplier ID is passed as a URL parameter

	// Bind the JSON input to the userInput struct
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	// Find the supplier by ID
	var supplier models.Supplier
	if err := models.DB.First(&supplier, supplierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Supplier Tidak Ditemukan"})
		return
	}

	// Update the supplier fields
	supplier.Nama_supplier = userInput.Nama_supplier
	supplier.Alamat = userInput.Alamat
	supplier.Kontak = userInput.Kontak
	supplier.Sertifikasi = userInput.Sertifikasi
	supplier.Verifikasi = userInput.Verifikasi
	// Add other fields as necessary

	// Save the updated supplier to the database
	if err := models.DB.Save(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Update Data Supplier Berhasil"})
}

// DeleteProduct deletes a supplier from the database
func DeleteSupplier(c *gin.Context) {
	supplierID := c.Param("id") // Assuming the supplier ID is passed as a URL parameter

	// Find the supplier by ID
	var supplier models.Supplier
	if err := models.DB.First(&supplier, supplierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Supplier Tidak Ditemukan!"})
		return
	}

	// Delete the supplier from the database
	if err := models.DB.Delete(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Hapus Data Supplier Berhasil"})
}

func GenerateUserID() (string, error) {
	var lastUserID string
	err := models.DB.Model(&models.Supplier{}).Order("id_supplier desc").Limit(1).Pluck("id_supplier", &lastUserID).Error
	if err != nil {
		return "", err // Jika tidak ada supplier sebelumnya, ID pertama bisa dimulai dari angka 177013001
	}

	// Mengambil angka urut dari ID terakhir
	if len(lastUserID) < 9 {
		return "177013001", nil // Jika ID pertama kali, mulai dari 177013001
	}

	// Ambil bagian angka urut dari ID terakhir (3 digit terakhir)
	lastIDNumeric, err := strconv.Atoi(lastUserID[6:])
	if err != nil {
		return "", err
	}

	// Tambahkan 1 ke ID terakhir
	newIDNumeric := lastIDNumeric + 1
	newID := fmt.Sprintf("177013%03d", newIDNumeric) // Format 3 digit angka

	return newID, nil
	
}

func Addsupplier(c *gin.Context) {
	var userInput models.Supplier

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
	c.JSON(http.StatusOK, gin.H{"Message": "Tambah Data Supplier Berhasil"})
}
