package main

import (
	"github.com/andreluzz/poc/modules/core/routes"
	module "github.com/andreluzz/poc/modules/shared"
)

func main() {
	module.ListenAndServe(":3010", "../../cert.pem", "../../key.pem", routes.Routes())
}
