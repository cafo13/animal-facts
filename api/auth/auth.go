package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/cafo13/animal-facts/api/database"
	jwt "github.com/dgrijalva/jwt-go"
)

type AuthHandler interface {
	VerifyLogin(username string, password string) (string, error)
	GenerateToken(userId uint) (string, error)
	GetUserById(uid uint) (*database.User, error)
}

type AuthDataHandler struct {
	Handler database.DatabaseHandler
}

func NewAuthHandler(databaseHandler database.DatabaseHandler) AuthHandler {
	return AuthDataHandler{Handler: databaseHandler}
}

func (adh AuthDataHandler) VerifyLogin(username string, password string) (string, error) {
	user := &database.User{}
	user.Database = *adh.Handler.GetDatabase()
	user.Username = username
	user.Password = password
	err := user.CheckLogin()
	if err != nil {
		return "", err
	}

	token, err := adh.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (adh AuthDataHandler) GenerateToken(userId uint) (string, error) {
	tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func (adh AuthDataHandler) GetUserById(uid uint) (*database.User, error) {
	var user *database.User
	user.Database = *adh.Handler.GetDatabase()
	user.ID = uid
	err := user.GetById()
	if err != nil {
		return nil, err
	}

	return user, nil
}
