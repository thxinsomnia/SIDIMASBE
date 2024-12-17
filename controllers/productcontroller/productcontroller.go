package productcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	
)

// Index godoc
//	@Summary		List products
//	@Security		Bearer
//	@Description	get list of products
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	object{ID=int, nama_barang=string, harga=int}	"Successfully retrieved list of products"
//	@Router			/api/products [get]
func Index(c *gin.Context) {
    data := []map[string]interface{}{
        {
            "id":          1,
            "nama barang": "Kaos Yoasobi",
            "harga":       215000,
        },

        {
            "id":          2,
            "nama barang": "Kaos Kenshi Yonezu",
            "harga":       278000,
        },

        {
            "id":          3,
            "nama barang": "Kaos One Ok Rock",
            "harga":       251000,
        },
    }

    c.JSON(http.StatusOK, data)
}