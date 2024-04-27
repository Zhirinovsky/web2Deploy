package pages

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
	"web2/bin"
	"web2/bin/models"
	"web2/bin/structures"
)

type CartPageObject struct {
	BaseObject      structures.BaseObject
	Carts           []models.Cart
	Products        []models.Product
	Characteristics []models.Characteristic
}

func Cart(w http.ResponseWriter, r *http.Request) {
	bin.RefreshAPI()
	pageObject := getCartData(bin.GetCurrentUser(w, r))
	pageObject.BaseObject.CurrentUser = bin.GetCurrentUser(w, r)
	pageObject.BaseObject.CsrfField = csrf.TemplateField(r)
	if pageObject.BaseObject.ErrorStr == "" {
		pageObject.BaseObject.ErrorStr = bin.GetError(&bin.GlobalError)
	}
	pageObject.BaseObject.MessageStr = bin.GlobalMessage
	bin.GlobalMessage = ""
	bin.Refresh(w, r)
	t, err := template.New("cart.html").Funcs(sprig.FuncMap()).ParseFiles("views/cart.html")
	template.Must(t.ParseGlob("views/templates/*.html"))
	template.Must(t.ParseGlob("views/templates/modals/*.html"))
	template.Must(t.ParseGlob("views/templates/tables/*.html"))
	bin.NullErrorCheck(&pageObject.BaseObject.ErrorStr, err)
	err = t.Execute(w, pageObject)
	bin.CheckErr(err)
}

func getCartData(user models.User) CartPageObject {
	var pageObject CartPageObject
	var carts []models.Cart
	response := bin.Request("/products", "GET", bin.ServerToken, nil, &pageObject.Products)
	bin.ResponseCheck(response, "/products", "GET")
	response = bin.Request("/characteristics", "GET", bin.ServerToken, nil, &pageObject.Characteristics)
	bin.ResponseCheck(response, "/characteristics", "GET")
	response = bin.Request("/carts", "GET", bin.ServerToken, nil, &carts)
	bin.ResponseCheck(response, "/carts", "GET")
	for _, cart := range carts {
		if cart.UserID == user.ID {
			pageObject.Carts = append(pageObject.Carts, cart)
		}
	}
	sort.Slice(pageObject.Carts, func(i, j int) bool {
		return pageObject.Carts[i].ID < pageObject.Carts[j].ID
	})
	return pageObject
}

func RemoveCart(w http.ResponseWriter, r *http.Request) {
	var result map[string]string
	response := bin.Request("/carts/"+r.FormValue("Id"), "DELETE", bin.GetCurrentToken(w, r), nil, &result)
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		bin.GlobalError = result["message"]
	}
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func ChangeCart(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	amount, _ := strconv.Atoi(r.FormValue("Amount"))
	active, _ := strconv.ParseBool(r.FormValue("Active"))
	var product models.Product
	response := bin.Request("/products/"+r.FormValue("ProductId"), "GET", bin.GetCurrentToken(w, r), nil, &product)
	bin.ResponseCheck(response, "/products/"+r.FormValue("ProductId"), "GET")
	cart := models.Cart{
		ID:        id,
		Amount:    amount,
		Active:    active,
		ProductID: product.ID,
		UserID:    bin.GetCurrentUser(w, r).ID,
	}
	var result map[string]string
	if r.FormValue("Operation") == "Add" {
		if product.Amount > int(amount) {
			cart.Amount = cart.Amount + 1
			response = bin.Request("/carts/"+r.FormValue("Id"), "PUT", bin.GetCurrentToken(w, r), cart, &result)
			if !(response.StatusCode == 200 || response.StatusCode == 201) {
				bin.GlobalError = result["message"]
			}
		} else {
			bin.GlobalMessage = "Нельзя повысить кол-во товара в позиции, иначе его количество превысит его количество на складе"
		}
	} else {
		if amount <= 1 {
			response = bin.Request("/carts/"+r.FormValue("Id"), "DELETE", bin.GetCurrentToken(w, r), nil, &result)
			if !(response.StatusCode == 200 || response.StatusCode == 201) {
				bin.GlobalError = result["message"]
			}
		} else {
			cart.Amount = cart.Amount - 1
			response = bin.Request("/carts/"+r.FormValue("Id"), "PUT", bin.GetCurrentToken(w, r), cart, &result)
			if !(response.StatusCode == 200 || response.StatusCode == 201) {
				bin.GlobalError = result["message"]
			}
		}
	}
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func ChangeActive(w http.ResponseWriter, r *http.Request) {
	var cart models.Cart
	response := bin.Request("/carts/"+r.FormValue("Id"), "GET", bin.GetCurrentToken(w, r), nil, &cart)
	bin.ResponseCheck(response, "/carts/"+r.FormValue("Id"), "GET")
	cart.Active = !cart.Active
	var result map[string]string
	response = bin.Request("/carts/"+r.FormValue("Id"), "PUT", bin.GetCurrentToken(w, r), cart, &result)
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		bin.GlobalError = result["message"]
	}
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func MakeOrder(w http.ResponseWriter, r *http.Request) {
	user := bin.GetCurrentUser(w, r)
	var result map[string]any
	order := models.Order{
		Date:     time.Now().Format(time.DateTime),
		Address:  r.FormValue("Address"),
		StatusID: 1,
		UserID:   user.ID,
	}
	requestOrder := log.Fields{
		"group":     "ordering",
		"email":     user.Email,
		"date":      order.Date,
		"address":   order.Address,
		"status_id": order.StatusID,
		"user_id":   order.UserID,
	}
	bin.SaveLog(requestOrder, log.TraceLevel, "Sending order creation request...")
	response := bin.Request("/orders", "POST", bin.GetCurrentToken(w, r), order, &result)
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		bin.GlobalError = result["message"].(string)
		requestOrder["error"] = result["message"]
		bin.SaveLog(requestOrder, log.ErrorLevel, "Failed order creation")
	} else {
		bin.SaveLog(requestOrder, log.InfoLevel, "Successful order creation")
		var carts []models.Cart
		response = bin.Request("/carts", "GET", bin.GetCurrentToken(w, r), nil, &carts)
		bin.ResponseCheck(response, "/carts", "GET")
		var content []models.Cart
		for _, cart := range carts {
			if cart.UserID == user.ID && cart.Active {
				content = append(content, cart)
			}
		}
		id := int(result["ID"].(float64))
		for _, cart := range content {
			position := models.Position{
				CheckoutPrice: cart.Product.Price,
				Amount:        cart.Amount,
				OrderID:       id,
				ProductID:     cart.ProductID,
				Order:         models.Order{},
				Product:       models.Product{},
			}
			requestPosition := log.Fields{
				"group":          "ordering",
				"email":          user.Email,
				"checkout_price": position.CheckoutPrice,
				"amount":         position.Amount,
				"order_id":       position.OrderID,
				"product_id":     position.ProductID,
			}
			response = bin.Request("/positions", "POST", bin.GetCurrentToken(w, r), position, &result)
			if !(response.StatusCode == 200 || response.StatusCode == 201) {
				bin.GlobalError = result["message"].(string)
				requestPosition["error"] = result["message"]
				bin.SaveLog(requestPosition, log.ErrorLevel, "Failed order position creation")
			} else {
				bin.SaveLog(requestPosition, log.InfoLevel, "Successful order position creation")
				requestCart := log.Fields{
					"group":      "ordering",
					"email":      user.Email,
					"id":         cart.ID,
					"active":     cart.Active,
					"amount":     cart.Amount,
					"user_id":    cart.UserID,
					"product_id": cart.ProductID,
				}
				cartId := strconv.Itoa(cart.ID)
				response = bin.Request("/carts/"+cartId, "DELETE", bin.GetCurrentToken(w, r), nil, &result)
				if !(response.StatusCode == 200 || response.StatusCode == 201) {
					bin.GlobalError = result["message"].(string)
					requestPosition["error"] = result["message"]
					bin.SaveLog(requestCart, log.ErrorLevel, "Failed cart position deletion")
				} else {
					bin.SaveLog(requestCart, log.InfoLevel, "Successful cart position deletion")
				}
			}
		}
		bin.GlobalMessage = "Заказ успешно создан"
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
