package links

import (
	"log"

	database "github.com/lucasschilin/hackernews-graphql-go/internal/pkg/db/mysql"
	"github.com/lucasschilin/hackernews-graphql-go/internal/users"
)

type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

func (link Link) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title, Address, UserID) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(link.Title, link.Address, link.User.ID)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	log.Printf("Row with ID = %d inserted", id)
	return id
}

func GetAll() []Link {
	stmt, err := database.Db.Prepare("SELECT L.id, L.title, L.address, U.id, U.username FROM Links L JOIN Users U ON L.UserID = U.id")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal()
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		var user users.User
		err := rows.Scan(&link.ID, &link.Title, &link.Address, &user.ID, &user.Username)
		if err != nil {
			log.Fatal(err)
		}

		link.User = &user

		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return links
}
