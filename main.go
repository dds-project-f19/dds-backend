package main

import (
	"dds-backend/routes"
)

func main() {
	router, conf, db, err := routes.MakeServer()
	defer db.Close()

	if err != nil {
		panic(err)
	}
	if err = router.Run(conf.Address); err != nil {
		panic(err)
	}
}
