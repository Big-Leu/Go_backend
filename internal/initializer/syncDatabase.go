package initializer

import (
	"kubequntumblock/pkg/models"
     "fmt"
)
func SyncDatabase(){
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.EndPoint{})
	fmt.Println("Database synchronized successfully!")
}