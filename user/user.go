package user

import (
	"net/mail"
	"time"
)

type User struct {
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Birthdate time.Time `json:"birthdate"`
	Password  string    `json:"password"`
}

type UserInterface interface {
	IsValid() bool
	GetAge() int
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// func function() bool {
// 	return false
// }

func (u *User) IsValid() bool {
	// if !function() {
	// 	return false
	// }

	if u.Firstname == "" ||
		u.Lastname == "" ||
		u.Email == "" ||
		!isValidEmail(u.Email) ||
		!isValidDate(u.Birthdate.Format("2006-01-02")) ||
		u.GetAge() < 13 {
		return false
	}

	return true
}

func (u *User) GetAge() int {
	return time.Now().Year() - u.Birthdate.Year()
}
