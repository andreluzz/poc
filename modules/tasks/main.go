package main

import (
	module "github.com/andreluzz/poc/modules/shared"
	"github.com/andreluzz/poc/modules/tasks/routes"
)

func main() {
	module.ListenAndServe(":3020", "../../cert.pem", "../../key.pem", routes.Routes())
}
