package main

import (
	"github.com/andreluzz/poc/modules/resources/routes"
	module "github.com/andreluzz/poc/modules/shared"
)

func main() {
	module.ListenAndServe(":3030", routes.Routes())
}
