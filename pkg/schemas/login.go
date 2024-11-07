package schemas

type SignupSchemas struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type EndPoint struct {
	ID   uint   `gorm:"primaryKey"`
	EndpointType string `json:"endpointtype"`
	EndpointName  string `json:"endpointname"`
}