package pages

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
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

func SetCategory(w http.ResponseWriter, r *http.Request) {
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
