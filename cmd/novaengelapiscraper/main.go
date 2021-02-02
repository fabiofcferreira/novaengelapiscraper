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

var hostSecure bool = false
var hostAddress string = ""

// AssetsHost is the host's address
var AssetsHost string = ""

func main() {
	runtime.GOMAXPROCS(6)

	flag.IntVar(&action, "action", 0, "Available actions:\n1. Generate products JSON\n2. Generate Shopify CSV\n3. Download all products images\n4. Generate Shopify CSV and add images download link\n")

	flag.StringVar(&password, "password", "", "Password used to login to Nova Engel's API")
	flag.StringVar(&username, "username", "", "Username used to login to Nova Engel's API")

	flag.BoolVar(&hostSecure, "hostSecure", false, "Host is secure")
	flag.StringVar(&hostAddress, "hostAddress", "0.0.0.0", "Host address")
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
	case 3:
		http.GetAllProductsImages(auth)
		break
	case 4:
		schema := ""
		if hostSecure {
			schema = "https"
		} else {
			schema = "http"
		}

		http.GenerageShopifyCSVWithImages(auth, schema, hostAddress)
		break
	}
}
