package types

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswordLength  = 8
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" valid:"email"`
	Password  string `json:"password"`
}

func IsValidPassword(encpw, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw))
	return err == nil
}

func (params CreateUserParams) Validate() map[string]string {
	errors := make(map[string]string)
	if len(params.FirstName) < minFirstNameLength {
		errors["firstName"] = fmt.Sprintf("firstName must be at least %d characters long", minFirstNameLength)
	}
	if len(params.LastName) < minLastNameLength {
		errors["lastName"] = fmt.Sprintf("lastName must be at least %d characters long", minLastNameLength)
	}
	if len(params.Password) < minPasswordLength {
		errors["password"] = fmt.Sprintf("password must be at least %d characters long", minPasswordLength)
	}
	if err := params.ValidateEmail(); err != nil {
		errors["email"] = err.Error()
	}
	return errors
}

func (params CreateUserParams) ValidateEmail() error {
	if !govalidator.IsEmail(params.Email) {
		return fmt.Errorf("email is invalid")
	}
	return nil
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (params UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(params.FirstName) > 0 {
		m["firstName"] = params.FirstName
	}
	if len(params.LastName) > 0 {
		m["lastName"] = params.LastName
	}
	return m
}

type User struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

func NewUserFromParams(params *CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
