package novaengelapiscraper

// LoginCredentials is used to login in Nova Engel API
type LoginCredentials struct {
	User     string
	Password string
}

// LoginAuthorization is used for auth tokens
type LoginAuthorization struct {
	Token string
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

	PriceQuantity []ProductPriceQuantity

	ItemID string
}

// ProductPriceQuantity is used to express price for each quantity
type ProductPriceQuantity struct {
	Quantity int
	Price    float32
}
