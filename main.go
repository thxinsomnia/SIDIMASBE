package handler// ganti ke main kalau mau di run local

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"SIDIMASBE/controllers/distribcontroller"
	"SIDIMASBE/controllers/menucontroller"
	"SIDIMASBE/controllers/authcontroller"
	"SIDIMASBE/controllers/matscontroller"
	"SIDIMASBE/controllers/suppliercontroller"
	_ "SIDIMASBE/docs"
	"SIDIMASBE/middlewares"
	"SIDIMASBE/models"
)

func init() {
	// Connect to Database
	models.ConnectDatabase()
}

// Handler for deployment - Menerima request dan menangani routing dengan CORS
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set up Gin router di dalam handler
	router := gin.Default()

	router.Use(middlewares.CORSMiddleware())

	// Register routes setelah router diinisialisasi
	router.POST("/login", authcontroller.Login)
	router.POST("/register", authcontroller.Register)
	router.GET("/logout", authcontroller.Logout)

	api := router.Group("/api")
	api.Use(middlewares.JWTVerif())
	api.GET("/supl", suppliercontroller.GetSupplier)
	api.GET("/supl/:id", suppliercontroller.GetSupplierByID)
	api.POST("/asupl", suppliercontroller.Addsupplier)
	api.PUT("/esupl/:id", suppliercontroller.UpdateSupplier)
	api.DELETE("/dsupl/:id", suppliercontroller.DeleteSupplier)
	api.POST("/addbahan", matscontroller.Addbahan)
	api.POST("/addstok", matscontroller.TambahStokBahan)
	api.GET("/bahan", matscontroller.GetAllBahan)
	api.GET("/bahan/:id", matscontroller.GetBahanByID)
	api.PUT("/ebahan/:id", matscontroller.EditBahan)
	api.DELETE("/dbahan/:id", matscontroller.HapusBahan)
	api.POST("/addmenu", menucontroller.BuatMenu)
	api.GET("/menu", menucontroller.AmbilDataMenu)
	api.GET("/menu/:id", menucontroller.AmbilDataMenuID)
	api.PUT("/emenu/:id", menucontroller.EditMenu)
	api.DELETE("/dmenu/:id", menucontroller.HapusMenu)
	api.POST("/bmenu", menucontroller.BuatPorsiMenu)
	api.POST("/adddistrib", distribcontroller.DistribusiMenu)
	api.GET("/distrib", distribcontroller.AmbilDataDistribusi)

	// Serve the Swagger UI at /swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Handle HTTP request
	router.ServeHTTP(w, r)
}

func main() {
	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Server is running on port %s\n", port)

	// Gunakan handler untuk menangani request HTTP
	http.HandleFunc("/", Handler) // Memetakan path "/" ke Handler

	// Mulai server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
