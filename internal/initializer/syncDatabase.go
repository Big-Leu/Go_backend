package initializer

import (
	"kubequntumblock/pkg/models"
     "fmt"
)
func SyncDatabase(){
	DB.AutoMigrate(&models.User{})
	fmt.Println("Database synchronized successfully!")
}