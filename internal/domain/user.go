package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deleted   bool      `json:"-"`
}

func NewUser(email, firstName, lastName, password string) User {
	return User{
		ID:        uuid.NewString(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Deleted:   false,
	}
}

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}
