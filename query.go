package main

import "github.com/jmoiron/sqlx"

var db *sqlx.DB

func InitDB() error {
	dbc, err := sqlx.Connect("sqlite3", Env.DB_URL)
	if err != nil {
		return err
	}
	db = dbc
	return nil
}

func CheckUserName(name string) bool {
	user := User{}
	err := db.Get(
		&user, "SELECT id, name, email FROM user WHERE name = $1;",
		name,
	)
	return err == nil
}

func CheckUserEmail(email string) bool {
	user := User{}
	err := db.Get(
		&user, "SELECT id, name, email FROM user WHERE email = $1;",
		email,
	)
	return err == nil
}

func CreateUser(req RegisterReq) (User, error) {
	id := CreateId()
	hashedPassword := Hash(req.Password)
	_, err := db.Exec(
		"INSERT INTO user (id, name, password, email) VALUES($1, $2, $3, $4);",
		id, req.Name, hashedPassword, req.Email,
	)
	if err != nil {
		return User{}, err
	}
	user := User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	}
	return user, err
}

func LoginUser(req LoginReq) (User, error) {
	hashedPassword := Hash(req.Password)
	user := User{}
	err := db.Get(
		&user, "SELECT id, name, email FROM user WHERE name = $1 AND password = $2;",
		req.Name, hashedPassword,
	)
	if err != nil {
		return User{}, err
	}
	return user, err
}
