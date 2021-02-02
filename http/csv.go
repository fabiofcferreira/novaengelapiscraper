package http

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fabiofcferreira/novaengelapiscraper"
	"github.com/fatih/color"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// GenerateShopifyCSV generates a Shopify-ready CSV with all the products with
// a valid image.
// It includes all the basic details needed for a product to be added to
// a shopify store instance.
func GenerateShopifyCSV(auth *novaengelapiscraper.LoginAuthorization) error {
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

	// Create headers and lines list
	headers := []string{"Handle", "Title", "Body (HTML)", "Vendor", "Type", "Tags", "Published", "Option1 Name", "Option1 Value", "Option2 Name", "Option2 Value", "Option3 Name", "Option3 Value", "Variant SKU", "Variant Grams", "Variant Inventory Tracker", "Variant Inventory Qty", "Variant Inventory Policy", "Variant Fulfillment Service", "Variant Price", "Variant Compare At Price", "Variant Requires Shipping", "Variant Taxable", "Variant Barcode", "Image Src", "Image Position", "Image Alt Text", "Gift Card", "SEO Title", "SEO Description", "Google Shopping / Google Product Category", "Google Shopping / Gender", "Google Shopping / Age Group", "Google Shopping / MPN", "Google Shopping / AdWords Grouping", "Google Shopping / AdWords Labels", "Google Shopping / Condition", "Google Shopping / Custom Product", "Google Shopping / Custom Label 0", "Google Shopping / Custom Label 1", "Google Shopping / Custom Label 2", "Google Shopping / Custom Label 3", "Google Shopping / Custom Label 4", "Variant Image", "Variant Weight Unit", "Variant Tax Code", "Cost per item", "Status"}
	lines := [][]string{{}}

	for _, product := range *products {
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
		line = append(line, "")
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

		lines = append(lines, line)
	}

	err = novaengelapiscraper.SaveCSVFile("products_shopify", headers, lines)
	if err != nil {
		color.HiRed("Couldn't save Shopify CSV file.\n")
		return err
	}

	return nil
}

// GenerageShopifyCSVWithImages generates a Shopify-ready CSV
func GenerageShopifyCSVWithImages(auth *novaengelapiscraper.LoginAuthorization, schema string, assetHost string) error {
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

	// Get all products images
	imageFolderName, err := GetAllProductsImagesUsingList(auth, products)
	if err != nil {
		color.HiRed("Couldn't fetch products images.")
		return err
	}

	// Create headers and lines list
	headers := []string{"Handle", "Title", "Body (HTML)", "Vendor", "Type", "Tags", "Published", "Option1 Name", "Option1 Value", "Option2 Name", "Option2 Value", "Option3 Name", "Option3 Value", "Variant SKU", "Variant Grams", "Variant Inventory Tracker", "Variant Inventory Qty", "Variant Inventory Policy", "Variant Fulfillment Service", "Variant Price", "Variant Compare At Price", "Variant Requires Shipping", "Variant Taxable", "Variant Barcode", "Image Src", "Image Position", "Image Alt Text", "Gift Card", "SEO Title", "SEO Description", "Google Shopping / Google Product Category", "Google Shopping / Gender", "Google Shopping / Age Group", "Google Shopping / MPN", "Google Shopping / AdWords Grouping", "Google Shopping / AdWords Labels", "Google Shopping / Condition", "Google Shopping / Custom Product", "Google Shopping / Custom Label 0", "Google Shopping / Custom Label 1", "Google Shopping / Custom Label 2", "Google Shopping / Custom Label 3", "Google Shopping / Custom Label 4", "Variant Image", "Variant Weight Unit", "Variant Tax Code", "Cost per item", "Status"}
	lines := [][]string{{}}

	for _, product := range *products {
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
		line = append(line, schema+"://"+assetHost+"/"+imageFolderName+"/"+product.EANs[0]+".jpg")
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

		lines = append(lines, line)
	}

	err = novaengelapiscraper.SaveCSVFile("products_shopify", headers, lines)
	if err != nil {
		color.HiRed("Couldn't save Shopify CSV file.\n")
		return err
	}

	return nil
}
