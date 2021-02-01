package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"github.com/fabiofcferreira/novaengelapiscraper"
	"github.com/fabiofcferreira/novaengelapiscraper/http"
)

// Parse flags
var username string = ""
var password string = ""

var action int = 0
var outputFile string = ""

var loggedIn bool = false
var lastLoginTime time.Time

var auth *novaengelapiscraper.LoginAuthorization

var cmd int = -1

func printAllCommands() {
	fmt.Printf("1. Download products JSON.\n")
	fmt.Printf("2. Download all products' JSON and images.\n")
	fmt.Printf("3. Convert products JSON (fetches new copy) to Shopify-ready CSV.\n")
	fmt.Printf("4. Convert existing JSON to Shopify-ready CSV.\n")
	fmt.Printf("0. Quit.\n")
}

func login() error {
	// Login
	color.Cyan("Logging in...")
	authorization, err := http.Login(username, password)
	if err != nil {
		color.HiRed("Couldn't login.")
		loggedIn = false

		return err
	}

	loggedIn = true
	auth = authorization
	lastLoginTime = time.Now()
	return nil
}

func downloadProducts() error {
	color.Cyan("Fetching all products...")
	products, err := http.GetAllProducts(auth.Token)
	if err != nil {
		color.HiRed("Couldn't fetch all products.")
		return err
	}
	fmt.Printf("Fetched %s products.\n", color.HiGreenString("%d", len(*products)))

	// Save products in JSON file
	color.Cyan("Saving products to JSON file...")
	err = novaengelapiscraper.SaveProductsInJSONFile("products.json", products)
	if err != nil {
		color.HiRed("Couldn't save products in JSON file.")
		return err
	}

	return nil
}

func downloadProductsImages() error {
	color.Cyan("Fetching all products...")
	products, err := http.GetAllProducts(auth.Token)
	if err != nil {
		color.HiRed("Couldn't fetch all products.")
		return err
	}
	fmt.Printf("Fetched %s products.\n", color.HiGreenString("%d", len(*products)))

	// Fetch all images
	color.Cyan("Fetching images...")
	err = http.GetAllProductsImages(auth.Token, products)
	if err != nil {
		color.HiRed("Couldn't save products images.")
		return err
	}

	return nil
}

func generateShopifyCSV() error {
	color.Cyan("Fetching all products...")
	products, err := http.GetAllProducts(auth.Token)
	if err != nil {
		color.HiRed("Couldn't fetch all products.")
		return err
	}
	fmt.Printf("Fetched %s products.\n", color.HiGreenString("%d", len(*products)))

	os.Mkdir("downloaded", 0755)

	var productsCounted int = 0
	var productsNumber int = len(*products)
	var wg sync.WaitGroup

	partsSize := productsNumber / 16
	startIndex := 0
	endIndex := partsSize
	for i := 0; i < productsNumber; i += partsSize {
		wg.Add(1)
		go writeShopifyProductsCSV(&wg, (*products)[startIndex:endIndex], &productsCounted, &productsNumber)
		startIndex += partsSize
		endIndex += partsSize

		if endIndex >= productsNumber {
			endIndex = productsNumber - 1
		}
	}

	wg.Wait()

	// // Open volume pricing CSV file stream
	// file, err = os.Create("volume_pricing.csv")
	// if err != nil {
	// 	color.HiRed("Couldn't create volume_pricing.csv file.")
	// 	return err
	// }

	// // Create CSV writer
	// csvWriter = csv.NewWriter(file)
	// defer csvWriter.Flush()

	// // Add headers
	// headers = []string{"sku", "barcode", "wholesale_price_1", "from_quantity_1", "to_quantity_1", "wholesale_price_2", "from_quantity_2", "to_quantity_2", "wholesale_price_3", "from_quantity_3", "to_quantity_3", "wholesale_price_4", "from_quantity_4", "to_quantity_4", "wholesale_price_5", "from_quantity_5", "to_quantity_5", "increments"}
	// err = csvWriter.Write(headers)
	// if err != nil {
	// 	color.HiRed("Couldn't write volume pricing CSV file headers.")
	// }

	// productsCounted = 0
	// wg = sync.WaitGroup{}

	// partsSize = productsNumber / 16
	// startIndex = 0
	// endIndex = partsSize
	// for i := 0; i < productsNumber; i += partsSize {
	// 	wg.Add(1)
	// 	go writeShopifyVolumePricingCSV(&wg, csvWriter, (*products)[startIndex:endIndex], &productsCounted, &productsNumber)
	// 	startIndex += partsSize
	// 	endIndex += partsSize

	// 	if endIndex >= productsNumber {
	// 		endIndex = productsNumber - 1
	// 	}
	// }

	// wg.Wait()

	return nil
}

func writeShopifyProductsCSV(wg *sync.WaitGroup, products []novaengelapiscraper.Product, productsCounter *int, productsNumber *int) (string, error) {
	unixTimeStr := strconv.FormatInt(time.Now().Unix(), 10)

	// Generate unique hash
	var b strings.Builder
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	for i := 0; i < 8; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	// Create images folder
	productFilename := "product-" + b.String() + "-" + unixTimeStr

	// Open products CSV file stream
	file, err := os.Create(productFilename + ".csv")
	if err != nil {
		color.HiRed("Couldn't create " + productFilename + ".csv file.")
		return "", err
	}

	// Create CSV writer
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()

	// Add headers
	headers := []string{"Handle", "Title", "Body (HTML)", "Vendor", "Type", "Tags", "Published", "Option1 Name", "Option1 Value", "Option2 Name", "Option2 Value", "Option3 Name", "Option3 Value", "Variant SKU", "Variant Grams", "Variant Inventory Tracker", "Variant Inventory Qty", "Variant Inventory Policy", "Variant Fulfillment Service", "Variant Price", "Variant Compare At Price", "Variant Requires Shipping", "Variant Taxable", "Variant Barcode", "Image Src", "Image Position", "Image Alt Text", "Gift Card", "SEO Title", "SEO Description", "Google Shopping / Google Product Category", "Google Shopping / Gender", "Google Shopping / Age Group", "Google Shopping / MPN", "Google Shopping / AdWords Grouping", "Google Shopping / AdWords Labels", "Google Shopping / Condition", "Google Shopping / Custom Product", "Google Shopping / Custom Label 0", "Google Shopping / Custom Label 1", "Google Shopping / Custom Label 2", "Google Shopping / Custom Label 3", "Google Shopping / Custom Label 4", "Variant Image", "Variant Weight Unit", "Variant Tax Code", "Cost per item", "Status"}
	err = csvWriter.Write(headers)
	if err != nil {
		color.HiRed("Couldn't write products CSV file headers.")
	}

	for _, product := range products {
		// Login again if token expires
		if time.Now().After(lastLoginTime.Add(15 * time.Minute)) {
			login()
		}

		line := []string{}

		// Create handle
		handle := product.Description
		handle, _, _ = transform.String(transform.Chain(norm.NFD, transform.RemoveFunc(novaengelapiscraper.IsMn), norm.NFC), handle)
		handle = strings.ReplaceAll(handle, " ", "-")

		line = append(line, handle)

		// Add title
		line = append(line, product.Description)

		// Add Body
		line = append(line, "")

		// Add vendor
		line = append(line, product.BrandName)

		// Add type
		line = append(line, "")

		// Add categories (tags)
		categories := ""
		for i, category := range product.Families {
			if i > 1 && i <= len(product.Families)-1 {
				categories += ", "
			}

			categories += category
		}
		line = append(line, categories)

		// Add published
		line = append(line, "TRUE")

		// Add option 1
		line = append(line, "")
		line = append(line, "")

		// Add option 2
		line = append(line, "")
		line = append(line, "")

		// Add option 3
		line = append(line, "")
		line = append(line, "")

		// Add SKU, weight, inventory tracker, stock, policy and fulfillment mode
		line = append(line, strconv.Itoa(product.ID))
		line = append(line, fmt.Sprintf("%.3f", product.KGs))
		line = append(line, "shopify")
		line = append(line, strconv.Itoa(product.Stock))
		line = append(line, "deny")
		line = append(line, "manual")

		// Add price
		line = append(line, fmt.Sprintf("%.2f", (product.Price/0.9)*1.23))
		line = append(line, "")

		// Add shipping necessity
		line = append(line, "TRUE")

		// Add taxable
		line = append(line, "TRUE")

		// Add barcode (add first EAN code)
		line = append(line, product.EANs[0])

		// Add image source
		imageURL, err := http.GetProductImageURL(auth.Token, product.ID)
		if len(imageURL) == 0 {
			login()

			imageURL, err = http.GetProductImageURL(auth.Token, product.ID)
			if err != nil {
				fmt.Printf("Couldn't fetch product image URL after second try (ID: %d)\n", product.ID)
				continue
			}
		}

		if err != nil {
			fmt.Printf("Couldn't fetch product image URL (ID: %d)\n", product.ID)
		}

		line = append(line, imageURL)
		line = append(line, "1")
		line = append(line, "")

		// Add gift card
		line = append(line, "FALSE")

		// Add SEO title and description
		line = append(line, product.Description)
		line = append(line, "")

		// Add Google Shopping categories, gender, age group, MPN
		line = append(line, categories)
		line = append(line, product.Gender)
		line = append(line, "")
		line = append(line, "")
		line = append(line, "")
		line = append(line, "")
		line = append(line, "new")
		line = append(line, "FALSE")
		line = append(line, "")
		line = append(line, "")
		line = append(line, "")
		line = append(line, "")
		line = append(line, "")
		line = append(line, "")
		line = append(line, "kg")
		line = append(line, "")
		line = append(line, fmt.Sprintf("%.2f", product.Price))
		line = append(line, "active")

		// Count product
		(*productsCounter)++

		// Log
		// fmt.Printf("[%.2f%%] Fetched image %d of %d products.\n", (float32(*productsCounter)/float32(*productsNumber))*100, *productsCounter, *productsNumber)

		err = csvWriter.Write(line)
		if err != nil {
			color.HiRed("Couldn't write to CSV file: %v", err)
			break
		}

		csvWriter.Flush()
	}

	wg.Done()

	// Close file streams
	csvWriter.Flush()
	err = file.Close()
	if err != nil {
		color.Red("Couldn't close file stream")
		os.Exit(1)
	}

	// fmt.Println(err)

	color.Green("Finished writing " + productFilename + ".csv")

	// oldLocation := "./" + productFilename + ".csv"
	// newLocation := "./downloaded/" + productFilename + ".csv"
	// err = os.Rename(oldLocation, newLocation)
	// if err != nil {
	// 	log.Fatal(err)
	// 	os.Exit(1)
	// }

	return productFilename, nil
}

func writeShopifyVolumePricingCSV(wg *sync.WaitGroup, csvWriter *csv.Writer, products []novaengelapiscraper.Product, productsCounter *int, productsNumber *int) error {
	for _, product := range products {
		line := []string{}

		// Add SKU and barcode
		line = append(line, strconv.Itoa(product.ID))
		line = append(line, product.EANs[0])

		// Add price of 1 unit
		line = append(line, fmt.Sprintf("%.2f", product.Price1))
		line = append(line, "1")
		line = append(line, "1")

		// Add price of 3 units
		line = append(line, fmt.Sprintf("%.2f", product.Price3))
		line = append(line, "3")
		line = append(line, "3")

		// Add price of 12 units
		line = append(line, fmt.Sprintf("%.2f", product.Price12))
		line = append(line, "12")
		line = append(line, "12")

		// Add price of 48 units
		line = append(line, fmt.Sprintf("%.2f", product.Price48))
		line = append(line, "48")
		line = append(line, "48")

		// Add price of 120 unit
		line = append(line, fmt.Sprintf("%.2f", product.Price120))
		line = append(line, "120")
		line = append(line, "120")

		line = append(line, "0")

		err := csvWriter.Write(line)
		if err != nil {
			fmt.Println(line)
			color.HiRed("Couldn't write to CSV file: %v", err)
			break
		}
	}

	wg.Done()

	return nil
}

func main() {
	runtime.GOMAXPROCS(6)

	flag.StringVar(&username, "username", "", "Username used to login to Nova Engel's API")
	flag.StringVar(&password, "password", "", "Password used to login to Nova Engel's API")

	flag.IntVar(&action, "action", 0, "Available actions:\n1. Download products\n2. Download products and its images\n3. Generate Shopify CSV\n")
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

	login()

	if action == 1 {
		downloadProducts()
	} else if action == 2 {
		downloadProductsImages()
	} else if action == 3 {
		generateShopifyCSV()
	} else {
		color.HiYellow("Couldn't read command.")
		os.Exit(1)
	}
}
