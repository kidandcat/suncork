package main

import (
	"fmt"
	"os"
)

var SK_STRIPE_KEY string

func main() {
	if os.Getenv("ENV") != "prod" {
		os.Setenv("ENV", "dev")
	}

	fmt.Println("Initializing in", os.Getenv("ENV"), "mode")
	server := loadLetsencript()

	initSessionDB()
	initDB()
	initProductTables()
	initOrderTables()
	initPaymentTables()
	loadConfiguration()
	loadTemplates()

	setupRoutes(server)

	if os.Getenv("ENV") == "prod" {
		SK_STRIPE_KEY = "sk_test_wz7prNIhQcE8myjeqnEn6oT5"
		fmt.Println(":: - - - - Server Ready - - - - ::")
		server.ListenAndServeTLS("", "")
	} else {
		SK_STRIPE_KEY = "sk_test_wz7prNIhQcE8myjeqnEn6oT5"
		fmt.Println(":: - - - - Server Ready - - - - ::")
		server.ListenAndServe()
	}
}
