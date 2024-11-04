package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"kubequntumblock/internal/initializer"
	"kubequntumblock/models"
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)


func Token(c *gin.Context , user uint,goth_cooke string)(*gin.Context){

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"sub" : user,
		"exp" : time.Now().Add(time.Hour *24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECERET")))

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error" : err,
		})
		return c
	}
	c.SetSameSite(http.SameSiteDefaultMode) 

	c.SetCookie("_gothic_session", goth_cooke, 3600*24, "/", "localhost", false, true)
	c.SetCookie("Authorization", tokenString, 3600*24, "/", "localhost", false, true)
	return c
}

func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes)[:length], nil
}


func CreateUser(c *gin.Context, email string) (uint, error) {

    var existingUser models.User
    if err := initializer.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
        return uint(existingUser.ID), nil
    }

    password, err := generateRandomPassword(12)
    if err != nil {
        return 0, fmt.Errorf("failed to generate random password: %w", err)
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return 0, fmt.Errorf("failed to hash password: %w", err)
    }

    user := models.User{Email: email, Password: string(hash)}
    result := initializer.DB.Create(&user)
    if result.Error != nil {
        return 0, fmt.Errorf("failed to create user: %w", result.Error)
    }

    return uint(user.ID), nil
}
