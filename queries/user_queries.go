package queries

import (
	"omokogo/models"
	"omokogo/utils"
)

func CheckUserName(name string) bool {
	user := models.User{}
	err := db.Get(
		&user, "SELECT id, name, email FROM user WHERE name = $1;",
		name,
	)
	return err == nil
}

func CheckUserEmail(email string) bool {
	user := models.User{}
	err := db.Get(
		&user, "SELECT id, name, email FROM user WHERE email = $1;",
		email,
	)
	return err == nil
}

func CreateUser(req models.RegisterReq) (models.User, error) {
	id := utils.CreateId()
	hashedPassword := utils.Hash(req.Password)
	_, err := db.Exec(
		"INSERT INTO user (id, name, password, email) VALUES($1, $2, $3, $4);",
		id, req.Name, hashedPassword, req.Email,
	)
	if err != nil {
		return models.User{}, err
	}
	user := models.User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	return user, err
}

func LoginUser(req models.LoginReq) (models.User, error) {
	hashedPassword := utils.Hash(req.Password)
	user := models.User{}
	err := db.Get(
		&user, "SELECT id, name, email FROM user WHERE name = $1 AND password = $2;",
		req.Name, hashedPassword,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, err
}
