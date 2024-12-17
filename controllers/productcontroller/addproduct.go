package productcontroller

import (

	"net/http"
	"github.com/gin-gonic/gin"
	"SIDIMASBE/models"
)

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