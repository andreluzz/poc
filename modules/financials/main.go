package main

import (
	"github.com/poc/modules"
	"github.com/poc/modules/financials/routes"
)

func main() {
	modules.ListenAndServe(":3040", "../../cert.pem", "../../key.pem", routes.Routes())
}
