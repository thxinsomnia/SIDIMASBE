package middlewares

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "time"
)

func CORSMiddleware() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "https://sidimas.vercel.app/"}, // Bisa disesuaikan dengan domain yang diizinkan
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length","Set-Cookie"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
