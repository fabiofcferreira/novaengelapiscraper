package novaengelapiscraper

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
)

// SaveProductsJSONFile writes a JSON file to the destination given.
// The list of all records are list in JSON file
func SaveProductsJSONFile(filename string, list *[]Product) error {
	file, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		color.HiRed("Couldn't marshal products' JSON.")
		return err
	}

	if err = ioutil.WriteFile(filename+".json", file, 0644); err != nil {
		color.HiRed("Couldn't write JSON to destination file %s.", filename)
		return err
	}

	return nil
}

// SaveCSVFile writes a CSV file to the destination given.
// Headers and lines are used as arguments and are written all at once
func SaveCSVFile(filename string, headers []string, lines [][]string) error {
	file, err := os.Create(filename + ".csv")
	if err != nil {
		color.HiRed("Couldn't open CSV file %s.\n", filename)
		return err
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.Flush()

	// Write headers
	err = csvWriter.Write(headers)
	if err != nil {
		color.HiRed("Couldn't write headers into CSV file %s.\n", filename)
		return err
	}

	// Write all lines at once
	err = csvWriter.WriteAll(lines)
	if err != nil {
		color.HiRed("Couldn't write all lines into CSV file %s.\n", filename)
		return err
	}

	return nil
}
