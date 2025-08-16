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
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title, Address) VALUES (?,?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(link.Title, link.Address)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	log.Printf("Row wiith ID = %d inserted", id)
	return id
}
