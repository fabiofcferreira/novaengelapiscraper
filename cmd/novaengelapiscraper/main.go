package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"

	"github.com/fabiofcferreira/novaengelapiscraper"
	"github.com/fabiofcferreira/novaengelapiscraper/http"
)

func main() {
	runtime.GOMAXPROCS(6)

	// Parse flags
	username := ""
	password := ""

	flag.StringVar(&username, "username", "", "Username used to login to Nova Engel's API.")
	flag.StringVar(&password, "password", "", "Password used to login to Nova Engel's API.")
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

	// Login
	color.Cyan("Logging in...")
	authorization, err := http.Login(username, password)
	if err != nil {
		color.HiRed("Couldn't login.")
		os.Exit(1)
	}

	token := authorization.Token
	fmt.Printf("Token: ")
	color.HiGreen(token)

	// Use token to get all products
	color.Cyan("Fetching all products...")
	products, err := http.GetAllProducts(token)
	if err != nil {
		color.HiRed("Couldn't fetch all products.")
		os.Exit(1)
	}
	fmt.Printf("Fetched %s products.\n", color.HiGreenString("%d", len(*products)))

	// Save products in JSON file
	color.Cyan("Saving products to JSON file...")
	err = novaengelapiscraper.SaveProductsInJSONFile("products.json", products)
	if err != nil {
		color.HiRed("Couldn't save products in JSON file.")
		os.Exit(1)
	}

	// Fetch all images
	color.Cyan("Fetching images...")
	err = http.GetAllProductsImages(token, products)
	if err != nil {
		color.HiRed("Couldn't save products images.")
		os.Exit(1)
	}

	os.Exit(0)
}
