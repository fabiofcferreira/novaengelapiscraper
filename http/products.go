package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fabiofcferreira/novaengelapiscraper"
	"github.com/fatih/color"
)

var productsNumber int = 0
var imagesFetchedCounter int = 0

var noImageProductUIDs = []int{}

// GetAllProducts fetches all products
func GetAllProducts(token string) (*[]novaengelapiscraper.Product, error) {
	products := &[]novaengelapiscraper.Product{}
	url := strings.ReplaceAll(Links["products"], "$TOKEN$", token)

	resp, err := http.Get(url)
	if err != nil {
		color.HiRed("Couldn't perform login request.")
		return nil, err
	}
	defer resp.Body.Close()

	// Parse request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.HiRed("Couldn't parse response JSON.")
		return nil, err
	}

	// Unmarshal JSON
	err = json.Unmarshal(body, products)
	if err != nil {
		color.HiRed("Couldn't parse response JSON.")
		return nil, err
	}

	for index, product := range *products {
		product.Price1 = product.PriceQuantity[0].Price
		product.Price3 = product.PriceQuantity[1].Price
		product.Price12 = product.PriceQuantity[2].Price
		product.Price48 = product.PriceQuantity[3].Price
		product.Price120 = product.PriceQuantity[4].Price

		product.PriceQuantity = []novaengelapiscraper.ProductPriceQuantity{}

		(*products)[index] = product
	}

	return products, err
}

// GetProductImageURL fetches the product temporary link
func GetProductImageURL(token string, id int) (string, error) {
	var url string = ""
	url = strings.ReplaceAll(Links["productImage"], "$TOKEN$", token)
	url = strings.ReplaceAll(url, "$ID$", strconv.Itoa(id))

	resp, err := http.Get(url)
	if err != nil {
		color.HiRed("Couldn't get product image URL.")
		return "", err
	}
	defer resp.Body.Close()

	// Parse request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.HiRed("Couldn't parse response JSON.")
		return "", err
	}

	// Fetch image data
	imageURL := strings.ReplaceAll(string(body), `"`, "")
	if len(imageURL) == 0 {
		return "", err
	}

	return imageURL, nil
}

// GetAllProductsImages fetches all product images and saves them locally
func GetAllProductsImages(token string, products *[]novaengelapiscraper.Product) error {
	productsNumber = len(*products)

	start := time.Now().Unix()

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
	imagesFolderName := "images-" + b.String()
	os.Mkdir(imagesFolderName, 0755)

	fmt.Printf("Hopefully there are %d images.\n", productsNumber)

	// Start fetching all images
	var wg sync.WaitGroup

	partsSize := productsNumber / 16
	startIndex := 0
	endIndex := partsSize

	for i := 0; i < productsNumber; i += partsSize {
		wg.Add(1)
		go GetProductsImagesAsync(&wg, token, imagesFolderName, (*products)[startIndex:endIndex])
		startIndex += partsSize
		endIndex += partsSize

		if endIndex >= productsNumber {
			endIndex = productsNumber - 1
		}
	}

	wg.Wait()

	end := time.Now().Unix()
	duration := end - start

	fmt.Printf("%d seconds have passed. %s images were fetched of %s products.", uint(duration), color.HiGreenString("%d", imagesFetchedCounter), color.HiGreenString("%d", productsNumber))

	file, err := json.MarshalIndent(noImageProductUIDs, "", "  ")
	if err != nil {
		color.HiRed("Couldn't save list of products with no image.")
	} else {
		if err = ioutil.WriteFile("productsWithoutImage.json", file, 0644); err != nil {
			color.HiRed("Couldn't write JSON to destination file.")
		}
	}

	return nil
}

// GetProductsImagesAsync fetches products images and saves them inside a specified folder
func GetProductsImagesAsync(wg *sync.WaitGroup, token string, folderName string, products []novaengelapiscraper.Product) {
	defer wg.Done()

	url := ""
	for _, product := range products {
		url = strings.ReplaceAll(Links["productImage"], "$TOKEN$", token)
		url = strings.ReplaceAll(url, "$ID$", strconv.Itoa(product.ID))

		resp, err := http.Get(url)
		if err != nil {
			color.HiRed("Couldn't get product image URL.")
			continue
		}
		defer resp.Body.Close()

		// Parse request
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			noImageProductUIDs = append(noImageProductUIDs, product.ID)
			color.HiRed("Couldn't parse response JSON.")
			continue
		}

		// Fetch image data
		imageURL := strings.ReplaceAll(string(body), `"`, "")
		if len(imageURL) == 0 {
			noImageProductUIDs = append(noImageProductUIDs, product.ID)
			continue
		}

		resp, err = http.Get(imageURL)
		if err != nil {
			noImageProductUIDs = append(noImageProductUIDs, product.ID)
			color.HiRed("Couldn't fetch product image data.")
			continue
		}
		defer resp.Body.Close()

		// Open a file write stream
		for _, eanCode := range product.EANs {
			file, err := os.Create(path.Join(folderName, eanCode+".jpg"))
			if err != nil {
				noImageProductUIDs = append(noImageProductUIDs, product.ID)
				color.HiRed("Couldn't open file write stream.")
				continue
			}
			defer file.Close()

			// Write to file
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				noImageProductUIDs = append(noImageProductUIDs, product.ID)
				color.HiRed("Couldn't save to image file.")
				continue
			}
		}

		imagesFetchedCounter++
		fmt.Printf("[%.2f%%] Fetched image %d of %d products.\n", (float32(imagesFetchedCounter)/float32(productsNumber))*100, imagesFetchedCounter, productsNumber)
	}
}
