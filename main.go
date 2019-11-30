package main

import (
	"dds-backend/routes"
	"dds-backend/services"
)

func main() {
	router, conf, db, err := routes.MakeServer()
	defer db.Close()

	if err != nil {
		panic(err)
	}

	go services.LaunchBot() // TODO enable

	if err = router.Run(conf.Address); err != nil {
		panic(err)
	}
}
