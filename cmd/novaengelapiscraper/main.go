package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/fabiofcferreira/novaengelapiscraper"
	"github.com/fabiofcferreira/novaengelapiscraper/http"
	"github.com/fatih/color"
)

// Parse flags
var username string = ""
var password string = ""

var action int = 0

func main() {
	runtime.GOMAXPROCS(6)

	flag.IntVar(&action, "action", 0, "Available actions:\n1. Generate products JSON\n2. Download products and its images\n3. Generate Shopify CSV\n")

	flag.StringVar(&password, "password", "", "Password used to login to Nova Engel's API")
	flag.StringVar(&username, "username", "", "Username used to login to Nova Engel's API")
	// flag.StringVar(&outputFile, "outputFile", "", "Output filename only")

	flag.Parse()

	// Validate username and password
	if len(username) == 0 {
		color.HiRed("Invalid username.")
		os.Exit(1)
	}

	if len(password) == 0 {
		color.HiRed("Invalid password.")
		os.Exit(1)
	}

	if action == 0 {
		color.HiRed("Couldn't read command")
		fmt.Println("Commands list:")
		flag.PrintDefaults()
	}

	auth := &novaengelapiscraper.LoginAuthorization{
		User:     username,
		Password: password,
	}

	switch action {
	case 1:
		http.GenerateProductsJSON(auth)
		break
	case 2:
		http.GenerateShopifyCSV(auth)
		break
	}
}
