package http

import (
	"time"

	"github.com/fabiofcferreira/novaengelapiscraper"
	"github.com/fatih/color"
)

// GenerateProductsJSON generates JSON containing all products
func GenerateProductsJSON(auth *novaengelapiscraper.LoginAuthorization) error {
	if time.Now().After(auth.LastLoggedIn.Add(5 * time.Minute)) {
		newAuth, err := Login(auth.User, auth.Password)
		if err != nil {
			color.HiRed("Couldn't login.")
			return err
		}

		auth = newAuth
	}

	// Get all products
	products, err := GetAllProducts(auth)
	if err != nil {
		color.HiRed("Couldn't fetch products list.")
		return err
	}

	err = novaengelapiscraper.SaveProductsJSONFile("products", products)
	if err != nil {
		color.HiRed("Couldn't save products JSON file.")
		return err
	}

	return nil
}
