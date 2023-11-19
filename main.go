package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Film struct {
	Title    string
	Director string
}

func main() {
	fmt.Println("Go app...")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&Checkout{})
	db.AutoMigrate(&Order{})
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Cart{})

	service := UberService{db: db}
	handler := HttpHandler{service: &service}

	// handler function #1 - returns the index.html template, with film data
	// h1 := func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl := template.Must(template.ParseFiles("index.html"))
	// 	films := map[string][]Film{
	// 		"Films": {
	// 			{Title: "The Godfather", Director: "Francis Ford Coppola"},
	// 			{Title: "Blade Runner", Director: "Ridley Scott"},
	// 			{Title: "The Thing", Director: "John Carpenter"},
	// 		},
	// 	}
	// 	tmpl.Execute(w, films)
	// }

	// // handler function #2 - returns the template block with the newly added film, as an HTMX response
	// h2 := func(w http.ResponseWriter, r *http.Request) {
	// 	time.Sleep(1 * time.Second)
	// 	title := r.PostFormValue("title")
	// 	director := r.PostFormValue("director")
	// 	// htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
	// 	// tmpl, _ := template.New("t").Parse(htmlStr)
	// 	tmpl := template.Must(template.ParseFiles("index.html"))
	// 	tmpl.ExecuteTemplate(w, "film-list-element", Film{Title: title, Director: director})
	// }

	h3 := func(w http.ResponseWriter, r *http.Request) {
		// tmpl := template.Must(template.ParseFiles("index.html"))
		// films := map[string][]Film{
		// 	"Films": {
		// 		{Title: "The Godfather", Director: "Francis Ford Coppola"},
		// 		{Title: "Blade Runner", Director: "Ridley Scott"},
		// 		{Title: "The Thing", Director: "John Carpenter"},
		// 	},
		// }

		// tmpl.Execute(w, films)
		tmpl := template.Must(template.ParseFiles("products.html", "base.html"))
		films := map[string][]Film{
			"Films": {
				{Title: "The Godfather", Director: "Francis Ford Coppola"},
				{Title: "Blade Runner", Director: "Ridley Scott"},
				{Title: "The Thing", Director: "John Carpenter"},
			},
		}

		//fmt.Print(tmpl.Tree)
		tmpl.ExecuteTemplate(w, "base", films)

		// mpl.ExecuteTemplate(w, "film-list-element", map[string][]Film{})
	}

	h4 := func(w http.ResponseWriter, r *http.Request) {
		// time.Sleep(1 * time.Second)
		// title := r.PostFormValue("title")
		// director := r.PostFormValue("director")
		// htmlStr := fmt.Sprintf("<li class='list-group-item bg-primary text-white'>%s - %s</li>", title, director)
		// tmpl, _ := template.New("t").Parse(htmlStr)
		// tmpl := template.Must(template.ParseFiles("index.html"))
		// tmpl.ExecuteTemplate(w, "film-list-element", Film{Title: title, Director: director})
		fmt.Print("sdfasfafd")
		w.Write([]byte("<li class=\"list-group-item bg-primary text-white\">{{ .Title }} - {{ .Director }}</li>"))
	}

	// define handlers
	// http.HandleFunc("/", h1)
	// http.HandleFunc("/add-film/", h2)
	http.HandleFunc("/queijo", h3)
	http.HandleFunc("/queijo/add-film/ ", h4)
	http.HandleFunc("/accounts", handler.accountsHandler)

	fmt.Print("heeeey")

	log.Fatal(http.ListenAndServe(":8000", nil))

}

type Product struct {
	gorm.Model
	price int
	name  string
	slug  string
	sku   string
}

type Cart struct {
	gorm.Model
	products []Product
}

func (c Cart) AddProduct(p Product) {
	c.products = append(c.products, p)
}

func (c Cart) IsClosable() bool {
	return len(c.products) > 0
}

type Checkout struct {
	gorm.Model
	cart          Cart
	paymentMethod int
	address       string
	date          string
	status        int
}

func (c Checkout) IsClosable() bool {
	return c.cart.IsClosable() && c.address != "" && c.date != ""
}

func (c Checkout) MakeOrder() Order {
	c.status = CheckoutClosed
	return Order{checkout: c}
}

const (
	CreditCardPaymentMethod = iota
	DebitCardPaymentMethod
)

const (
	CheckoutOpen = iota
	CheckoutClosed
)

type Order struct {
	gorm.Model
	checkout Checkout
}

type Account struct {
	gorm.Model
	openCart Cart
	orders   []Order
	username string
	email    string
}

func makeAccount(username string, email string) Account {
	return Account{openCart: Cart{}, orders: make([]Order, 0), username: username, email: email}
}

type HttpHandler struct {
	service *UberService
}

func (h HttpHandler) accountsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		accounts := h.service.listAccounts()
		b, _ := json.Marshal(accounts)

		w.Write(b)
	}
}

func (h HttpHandler) cartHandler(w http.ResponseWriter, r *http.Request) {

}

func (h HttpHandler) cartCheckoutHandler(w http.ResponseWriter, r *http.Request) {

}

func (h HttpHandler) checkoutAddAdressHandler(w http.ResponseWriter, r *http.Request) {

}

func (h HttpHandler) checkoutAddPaymentMethodHandler(w http.ResponseWriter, r *http.Request) {

}

func (h HttpHandler) checkoutAddDateHandler(w http.ResponseWriter, r *http.Request) {

}

func (h HttpHandler) checkoutFinalizeHandler(w http.ResponseWriter, r *http.Request) {

}

func (h HttpHandler) checkoutHandler(w http.ResponseWriter, r *http.Request) {

}

func (h HttpHandler) productsHandler(w http.ResponseWriter, r *http.Request) {

}

type UberService struct {
	db *gorm.DB
}

func (self UberService) listAccounts() []Account {
	var acconts []Account
	self.db.Find(&acconts)
	return acconts
}

func (self UberService) cart(id int) Cart {
	var cart Cart
	self.db.First(&cart, id)
	return cart
}

func (self UberService) cartCheckout(id int) Checkout {
	var cart Cart
	self.db.First(&cart, id)

	checkout := Checkout{cart: cart}
	self.db.Create(&checkout)

	return checkout
}

func (self UberService) checkoutAddAdress(id int, adress string) {
	var checkout Checkout
	self.db.First(&checkout, id)

	checkout.address = adress

	self.db.Model(&checkout).Update("adress", adress)
}

func (self UberService) checkoutAddPaymentMethod() {

}

func (self UberService) checkoutAddDate() {

}

func (self UberService) checkoutFinalize() {

}

func (self UberService) checkout() {

}

func (self UberService) products() {

}

//clientes tem um carrinho vazio
//clientes adicionam produtos no carrinho
//clientes finalizam as compras e vão pro pagamento
//cliente escolhe a forma de pagamento, endereço e data da entrega
//cliente finaliza checkout e ordem é criada

//Get accounts/id

//Post accounts/id/carts/id/add-product
//{ product }

//Get accounts/id/carts/id

//Post accounts/id/carts/id/checkout

//Patch accouts/id/checkouts/id/addAdress
//Patch accouts/id/checkouts/id/addPaymentMethod
//Patch accouts/id/checkouts/id/addDate

//Post accouts/id/checkouts/id/finalize
