package pages

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"web2/bin"
	"web2/bin/models"
	"web2/bin/structures"
)

var SearchText string
var SetCategories []models.Category

type MainPageObject struct {
	BaseObject      structures.BaseObject
	Products        []models.Product
	Categories      []models.Category
	Characteristics []models.Characteristic
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		bin.RefreshAPI()
		pageObject := getData()
		pageObject.BaseObject.CurrentUser = bin.GetCurrentUser(w, r)
		pageObject.BaseObject.CsrfField = csrf.TemplateField(r)
		if pageObject.BaseObject.ErrorStr == "" {
			pageObject.BaseObject.ErrorStr = bin.GetError(&bin.GlobalError)
		}
		pageObject.BaseObject.MessageStr = bin.GlobalMessage
		bin.GlobalMessage = ""
		bin.Refresh(w, r)
		t, err := template.New("main.html").Funcs(sprig.FuncMap()).ParseFiles("views/main.html")
		template.Must(t.ParseGlob("views/templates/*.html"))
		template.Must(t.ParseGlob("views/templates/modals/*.html"))
		bin.NullErrorCheck(&pageObject.BaseObject.ErrorStr, err)
		err = t.Execute(w, pageObject)
		bin.CheckErr(err)
	}
}

func getData() MainPageObject {
	var pageObject MainPageObject
	var unfilteredProducts []models.Product
	response := bin.Request("/products", "GET", bin.ServerToken, nil, &unfilteredProducts)
	bin.ResponseCheck(response, "/products", "GET")
	pageObject.Products = searchHandle(unfilteredProducts)
	response = bin.Request("/characteristics", "GET", bin.ServerToken, nil, &pageObject.Characteristics)
	bin.ResponseCheck(response, "/characteristics", "GET")
	response = bin.Request("/categories", "GET", bin.ServerToken, nil, &pageObject.Categories)
	bin.ResponseCheck(response, "/categories", "GET")
	return pageObject
}

func searchHandle(unfilteredProducts []models.Product) []models.Product {
	var Products []models.Product
	for _, product := range unfilteredProducts {
		if strings.Contains(product.Name, SearchText) {
			if len(SetCategories) == 0 {
				Products = append(Products, product)
			} else {
				for _, category := range SetCategories {
					if category.Name == product.Category.Name {
						Products = append(Products, product)
						break
					}
				}
			}
		}
	}
	SearchText = ""
	SetCategories = nil
	return Products
}

func SearchProduct(w http.ResponseWriter, r *http.Request) {
	SearchText = r.FormValue("SearchText")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Filter(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	response := bin.Request("/categories", "GET", bin.ServerToken, nil, &categories)
	bin.ResponseCheck(response, "/categories", "GET")
	for _, category := range categories {
		if r.FormValue(category.Name) == "1" {
			SetCategories = append(SetCategories, category)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AddCart(w http.ResponseWriter, r *http.Request) {
	idProduct, _ := strconv.ParseInt(r.FormValue("Id"), 0, 10)
	user := bin.GetCurrentUser(w, r)
	if user.ID != 0 {
		var carts []models.Cart
		response := bin.Request("/carts", "GET", bin.ServerToken, nil, &carts)
		bin.ResponseCheck(response, "/carts", "GET")
		check := true
		for _, cart := range carts {
			if cart.UserID == user.ID && cart.ProductID == int(idProduct) {
				check = false
			}
		}
		if check {
			var product models.Product
			response = bin.Request("/products/"+r.FormValue("Id"), "GET", bin.ServerToken, nil, &product)
			bin.ResponseCheck(response, "/products/"+r.FormValue("Id"), "GET")
			if product.Amount > 0 {
				cart := models.Cart{
					ID:        0,
					Amount:    1,
					Active:    true,
					ProductID: int(idProduct),
					UserID:    user.ID,
				}
				var result map[string]string
				response = bin.Request("/carts", "POST", bin.GetCurrentToken(w, r), cart, &result)
				if !(response.StatusCode == 200 || response.StatusCode == 201) {
					bin.GlobalError = result["message"]
				} else {
					bin.GlobalMessage = "Товар успешно добавлен в корзину"
				}
			}
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
