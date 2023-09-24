package user

type User struct {
	ID       string `gorm:"primaryKey"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
