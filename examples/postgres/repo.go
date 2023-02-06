package postgres

import "database/sql"

type User struct {
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

func CreateUser(db *sql.DB, user User) error {
	_, err := db.Exec("insert into users (name, email) values ($1, $2)", user.Name, user.Email)
	return err
}
