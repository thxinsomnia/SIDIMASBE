package middlewares

import (
	"net/http"

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
        tokenString, err := c.Cookie("token")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"Message": "Tidak Terverifikasi!! Harap Login Terlebih dahulu!!"})
            c.Abort()
            return
        }

        claims := &config.JWTClaims{}

        //parsing token
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return config.JWT_KEY, nil
        })

        if err != nil {
            v, _ := err.(*jwt.ValidationError)
            switch v.Errors {
            case jwt.ValidationErrorSignatureInvalid:
                c.JSON(http.StatusUnauthorized, gin.H{"Message": "Tidak Terverifikasi!! Harap Login Terlebih dahulu!!"})
                c.Abort()
                return

            case jwt.ValidationErrorExpired:
                c.JSON(http.StatusUnauthorized, gin.H{"Message": "Silahkan Login Ulang Sesi Sudah Kadaluarsa!!"})
                c.Abort()
                return

            default:
                c.JSON(http.StatusUnauthorized, gin.H{"Message": "Tidak Terverifikasi!! Harap Login Terlebih dahulu!!"})
                c.Abort()
                return
            }
        }

        if !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"Message": "Tidak Terverifikasi!! Harap Login Terlebih dahulu!!"})
            c.Abort()
            return
        }

        c.Next()
    }
}
