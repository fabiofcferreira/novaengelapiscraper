package novaengelapiparser

import (
	"encoding/json"
	"io/ioutil"

	"github.com/fatih/color"
)

// SaveProductsInJSONFile saves products list in JSON file
func SaveProductsInJSONFile(filePath string, products *[]Product) error {
	file, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		color.HiRed("Couldn't marshal products' JSON.")
		return err
	}

	if err = ioutil.WriteFile(filePath, file, 0644); err != nil {
		color.HiRed("Couldn't write JSON to destination file %s.", filePath)
	}

	return err
}
