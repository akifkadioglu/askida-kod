package main

import (
	"github.com/akifkadioglu/askida-kod/database"
	"github.com/akifkadioglu/askida-kod/env"
	"github.com/akifkadioglu/askida-kod/route"
	"github.com/akifkadioglu/askida-kod/utils"
)

func main() {
	env.InitEnv(env.LOCAL)
	database.Connection()
	utils.InitTokenAuth()
	route.InitServices()
}
