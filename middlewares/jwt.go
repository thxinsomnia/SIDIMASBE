package middlewares

import (
	"net/http"
    "fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"SIDIMASBE/config"

)

// JWTVerif godoc
//	@Summary		Verify JWT token
//	@Description	verify JWT token from the cookie. Client should send "Cookie" header with the format "token=<JWT token>".
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			Cookie	header		string				true	"JWT token"
//	@Success		200		{object}	map[string]string	"Successfully verified"
//	@Failure		401		{object}	map[string]string	"Unauthorized"
//	@Router			/api [get]
func JWTVerif() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 🔥 Ambil token dari header Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"Message": "Token tidak ditemukan di Header!"})
            c.Abort()
            return
        }

        // 🔥 Pastikan format "Bearer <token>"
        var tokenString string
        if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
            tokenString = authHeader[7:]
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"Message": "Format Token Salah!"})
            c.Abort()
            return
        }

        claims := &config.JWTClaims{}

        // 🔥 Parsing token
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return config.JWT_KEY, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"Message": "Token tidak valid atau sudah expired"})
            c.Abort()
            return
        }

        // ✅ Debugging: Token berhasil diverifikasi
        fmt.Println("✅ Token Valid! Username:", claims.Username)

        c.Set("username", claims.Username) // Simpan username untuk request berikutnya
        c.Next()
    }
}
