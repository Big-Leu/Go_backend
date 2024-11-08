package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
	Email string `gorm:"unique"`
	Password string

}
type EndPoint struct {
    gorm.Model
	EndpointType string 
	EndpointName string `gorm:"unique"`
	EndpointRoute string `gorm:"unique"`
}