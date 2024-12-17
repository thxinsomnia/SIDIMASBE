package productcontroller

import (

	"net/http"
	"github.com/gin-gonic/gin"
	"SIDIMASBE/models"
)


// GetProducts retrieves all products from the database
func GetProduct(c *gin.Context) {
	var products []models.Product

	// Fetch all products from the database
	if err := models.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Products": products})
}

// GetProductByID retrieves a product by its ID from the database
func GetProductByID(c *gin.Context) {
	productID := c.Param("id") // Assuming the product ID is passed as a URL parameter

	var product models.Product

	// Find the product by ID
	if err := models.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Product": product})
}

func Addproduct(c *gin.Context) {
    var userInput models.Product
    if err := c.ShouldBindJSON(&userInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
        return
    }

    if err := models.DB.Create(&userInput).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Message": "Tambah Produk Berhasil"})
}

// UpdateProduct updates an existing product in the database
func UpdateProduct(c *gin.Context) {
	var userInput models.Product
	productID := c.Param("id") // Assuming the product ID is passed as a URL parameter

	// Bind the JSON input to the userInput struct
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	// Find the product by ID
	var product models.Product
	if err := models.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Product not found"})
		return
	}

	// Update the product fields
	product.Nama = userInput.Nama
	product.Stok = userInput.Stok
	product.Harga = userInput.Harga
	// Add other fields as necessary

	// Save the updated product to the database
	if err := models.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Update Produk Berhasil"})
}

// DeleteProduct deletes a product from the database
func DeleteProduct(c *gin.Context) {
	productID := c.Param("id") // Assuming the product ID is passed as a URL parameter

	// Find the product by ID
	var product models.Product
	if err := models.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Product not found"})
		return
	}

	// Delete the product from the database
	if err := models.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Delete Produk Berhasil"})
}