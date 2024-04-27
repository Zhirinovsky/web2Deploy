package pages

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/csrf"
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"web2/bin"
	"web2/bin/models"
	"web2/bin/structures"
)

var SearchText string
var FilterCategories []models.Category
var FilterCharacteristicSets []models.Characteristic
var FilterCharacteristicNumbers []FilterCharacteristicInt

type MainPageObject struct {
	BaseObject      structures.BaseObject
	Products        []models.Product
	Categories      []models.Category
	Characteristics []models.Characteristic
}

type FilterCharacteristicInt struct {
	Characteristic models.Characteristic
	ValueMax       float64
	ValueMin       float64
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
			if len(FilterCategories) == 0 {
				if len(FilterCharacteristicSets) == 0 {
					if len(FilterCharacteristicNumbers) == 0 {
						Products = append(Products, product)
					} else {
						for _, characteristic := range FilterCharacteristicNumbers {
							for _, set := range product.Sets {
								if characteristic.Characteristic.Name == set.Characteristic.Name && characteristic.ValueMin <= set.Value && characteristic.ValueMax >= set.Value {
									Products = append(Products, product)
									goto exit
								}
							}
						}
					}
				} else {
					for _, characteristic := range FilterCharacteristicSets {
						for _, set := range product.Sets {
							if characteristic.Name == set.Characteristic.Name {
								if len(FilterCharacteristicNumbers) == 0 {
									Products = append(Products, product)
									goto exit
								} else {
									for _, characteristic := range FilterCharacteristicNumbers {
										for _, set := range product.Sets {
											if characteristic.Characteristic.Name == set.Characteristic.Name && characteristic.ValueMin <= set.Value && characteristic.ValueMax >= set.Value {
												Products = append(Products, product)
												goto exit
											}
										}
									}
								}
							}
						}
					}
				}
			} else {
				for _, category := range FilterCategories {
					if category.Name == product.Category.Name {
						if len(FilterCharacteristicSets) == 0 {
							if len(FilterCharacteristicNumbers) == 0 {
								Products = append(Products, product)
								goto exit
							} else {
								for _, characteristic := range FilterCharacteristicNumbers {
									for _, set := range product.Sets {
										if characteristic.Characteristic.Name == set.Characteristic.Name && characteristic.ValueMin <= set.Value && characteristic.ValueMax >= set.Value {
											Products = append(Products, product)
											goto exit
										}
									}
								}
							}
						} else {
							for _, characteristic := range FilterCharacteristicSets {
								for _, set := range product.Sets {
									if characteristic.Name == set.Characteristic.Name {
										if len(FilterCharacteristicNumbers) == 0 {
											Products = append(Products, product)
											goto exit
										} else {
											for _, characteristic := range FilterCharacteristicNumbers {
												for _, set := range product.Sets {
													if characteristic.Characteristic.Name == set.Characteristic.Name && characteristic.ValueMin <= set.Value && characteristic.ValueMax >= set.Value {
														Products = append(Products, product)
														goto exit
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	exit:
	}
	SearchText = ""
	FilterCategories = nil
	FilterCharacteristicSets = nil
	FilterCharacteristicNumbers = nil
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
			FilterCategories = append(FilterCategories, category)
		}
	}
	var characteristics []models.Characteristic
	response = bin.Request("/characteristics", "GET", bin.ServerToken, nil, &characteristics)
	bin.ResponseCheck(response, "/characteristics", "GET")
	for _, characteristic := range characteristics {
		if r.FormValue(characteristic.Name) == "1" {
			FilterCharacteristicSets = append(FilterCharacteristicSets, characteristic)
		}
		valueMin, errMin := strconv.ParseFloat(r.FormValue(characteristic.Name+"|Min"), 64)
		valueMax, errMax := strconv.ParseFloat(r.FormValue(characteristic.Name+"|Max"), 64)
		if errMin == nil || errMax == nil {
			if errMin != nil {
				valueMin = -math.MaxFloat64
			}
			if errMax != nil {
				valueMax = math.MaxFloat64
			}
			characteristicNumber := FilterCharacteristicInt{
				Characteristic: characteristic,
				ValueMax:       valueMax,
				ValueMin:       valueMin,
			}
			FilterCharacteristicNumbers = append(FilterCharacteristicNumbers, characteristicNumber)
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
