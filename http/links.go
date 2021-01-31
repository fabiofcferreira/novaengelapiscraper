package http

// Links is the list of all Nova Engel API endpoints
var Links = map[string]string{
	"login":        "https://b2b.novaengel.com/api/login",
	"products":     "https://b2b.novaengel.com/api/products/availables/$TOKEN$/pt",
	"productImage": "https://b2b.novaengel.com/api/products/image/$TOKEN$/$ID$",
}
