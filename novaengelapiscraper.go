package novaengelapiscraper

import (
	"time"
	"unicode"
)

// LoginAuthorization is used for auth tokens
type LoginAuthorization struct {
	User         string
	Password     string
	Token        string
	LastLoggedIn time.Time
}

// Product is used to represent products
type Product struct {
	ID          int
	EANs        []string
	Description string
	Price       float32
	PVR         float32
	Stock       int
	BrandID     string
	BrandName   string
	Gender      string
	Families    []string
	KGs         float32

	Width  int
	Height int
	Depth  int
	VAT    float32

	Date string

	Content     string
	ProductLine string

	PriceQuantity []ProductPriceQuantity `json:",omitempty"`

	Price1   float32
	Price3   float32
	Price12  float32
	Price48  float32
	Price120 float32

	ItemID string
}

// ProductPriceQuantity is used to express price for each quantity
type ProductPriceQuantity struct {
	Quantity int
	Price    float32
}

// IsMn ...
func IsMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}
