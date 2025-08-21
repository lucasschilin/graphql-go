package users

import (
	"database/sql"
	"log"

	database "github.com/lucasschilin/hackernews-graphql-go/internal/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

func (u *User) Create() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Users (Username, Password) VALUES (?,?)")
	print(stmt)
	if err != nil {
		log.Fatal()
	}

	hashedPassword, _ := HashPassword(u.Password)
	res, err := stmt.Exec(u.Username, hashedPassword)
	if err != nil {
		log.Fatal()
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return id
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetByLinkID(linkID string) *User {
	stmt, err := database.Db.Prepare("SELECT U.ID, U.Username FROM Users U JOIN Links L ON L.UserID = U.ID WHERE L.ID = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := stmt.QueryRow(linkID)

	var user User
	err = row.Scan(&user.ID, &user.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return nil
	}

	return &user
}

func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("select ID from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}

func (u *User) Authenticate() bool {
	stmt, err := database.Db.Prepare("SELECT Password FROM Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}

	row := stmt.QueryRow(u.Username)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
		return false
	}

	return CheckPasswordHash(u.Password, hashedPassword)
}

func GetAll() []User {
	stmt, err := database.Db.Prepare("SELECT ID, Username FROM Users")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal()
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User

		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return users

}
