package authcontroller

import (

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"SIDIMASBE/config"

	"SIDIMASBE/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login godoc
//	@Summary		Login a user
//	@Description	login a user by taking a JSON input
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User			true	"User to login"
//	@Success		200		{object}	map[string]string	"Successfully logged in"
//	@Failure		400		{object}	map[string]string	"Bad Request"
//	@Router			/login [post]
func Login(c *gin.Context) {
    var userInput models.User
    if err := c.ShouldBindJSON(&userInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
        return
    }

    var user models.User
    if err := models.DB.Where("username = ?", userInput.Username).First(&user).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.JSON(http.StatusUnauthorized, gin.H{"Message": "Username atau Password Tidak Sesuai"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
            return
        }
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"Message": "Username atau Password Tidak Sesuai"})
        return
    }

    expTime := time.Now().Add(time.Minute * 1)
    claims := &config.JWTClaims{
        Username: user.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    "go-jwt-mux",
            ExpiresAt: jwt.NewNumericDate(expTime),
        },
    }

    tokenDeklarasi := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token, err := tokenDeklarasi.SignedString(config.JWT_KEY)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
        return
    }

    c.SetCookie("token", token, 3600, "/", "", false, true)
 c.JSON(http.StatusOK, gin.H{"Message": "Login Berhasil!", "Token": token})
}



// Register godoc
//	@Summary		Register a new user
//	@Description	register a new user by taking a JSON input
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User			true	"User to register"
//	@Success		200		{object}	map[string]string	"Successfully registered"
//	@Failure		400		{object}	map[string]string	"Bad Request"
//	@Router			/register [post]
func Register(c *gin.Context) {
    var userInput models.User
    if err := c.ShouldBindJSON(&userInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
        return
    }

    hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
    userInput.Password = string(hashPassword)

    if err := models.DB.Create(&userInput).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Message": "Pendaftaran Berhasil"})
}


// Logout godoc
//	@Summary		Logout user
//	@Description	clear JWT token from the cookie
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Successfully logged out"
//	@Security		Bearer
//	@Router			/logout [get]
func Logout(c *gin.Context) {
    c.SetCookie("token", "", -1, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{"Message": "Logout Berhasil!"})
}
