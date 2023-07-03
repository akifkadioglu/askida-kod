package database

import (
	"github.com/akifkadioglu/askida-kod/ent"
)

var client *ent.Client
var err error

type PostgreSQL struct{}
type MySQL struct{}

func Connection() {
	var database MySQL
	database.main()
}
func DBManager() *ent.Client {
	return client
}
