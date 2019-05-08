package main

import (
	"github.com/andreluzz/poc/modules/financials/routes"
	module "github.com/andreluzz/poc/modules/shared"
)

func main() {
	module.ListenAndServe(":3040", routes.Routes())
}
