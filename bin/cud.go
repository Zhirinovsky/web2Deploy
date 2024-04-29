package bin

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"web2/bin/models"
)

func CudProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	price, _ := strconv.ParseFloat(r.FormValue("Price"), 64)
	amount, _ := strconv.Atoi(r.FormValue("Amount"))
	discount, _ := strconv.Atoi(r.FormValue("Discount"))
	product := models.Product{
		ID:         id,
		Name:       r.FormValue("Name"),
		Price:      price,
		Amount:     amount,
		Discount:   discount,
		ImageLink:  r.FormValue("Image"),
		CategoryID: r.FormValue("Category.Id"),
		IsExist:    true,
	}
	var result map[string]string
	var response *http.Response
	requestProduct := log.Fields{
		"group":       "crud",
		"method":      r.FormValue("Method"),
		"table":       "product",
		"id":          product.ID,
		"name":        product.Name,
		"price":       product.Price,
		"amount":      product.Amount,
		"discount":    product.Discount,
		"image_link":  product.ImageLink,
		"category_id": product.CategoryID,
	}
	SaveLog(requestProduct, log.TraceLevel, "Sending database change request...")
	switch r.FormValue("Method") {
	case "POST":
		response = Request("/products", "POST", GetCurrentToken(w, r), product, &result)
	case "PUT":
		response = Request("/products/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), product, &result)
	case "DELETE":
		if r.FormValue("Type") == "Physical" {
			response = Request("/products/"+r.FormValue("Id"), "DELETE", GetCurrentToken(w, r), nil, &result)
		} else {
			var item models.Product
			Request("/products/"+r.FormValue("Id"), "GET", GetCurrentToken(w, r), nil, &item)
			item.IsExist = false
			response = Request("/products/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), item, &result)
		}
	}
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		GlobalError = result["message"]
		requestProduct["error"] = result["message"]
		SaveLog(requestProduct, log.ErrorLevel, "Failed database change")
	} else {
		SaveLog(requestProduct, log.InfoLevel, "Successful database change")
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func CudSet(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	var set models.Set
	set.ID = id
	if r.FormValue("ValueSet"+r.FormValue("Characteristic.Id")) == "" {
		set.CharacteristicID, _ = strconv.Atoi(r.FormValue("Characteristic.Id"))
		set.Value, _ = strconv.ParseFloat(r.FormValue("ValueInt"), 64)
	} else {
		set.CharacteristicID, _ = strconv.Atoi(r.FormValue("ValueSet" + r.FormValue("Characteristic.Id")))
	}
	var result map[string]string
	var response *http.Response
	requestSet := log.Fields{
		"group":             "crud",
		"method":            r.FormValue("Method"),
		"table":             "set",
		"id":                set.ID,
		"characteristic_id": set.CharacteristicID,
		"product_id":        set.ProductID,
	}
	SaveLog(requestSet, log.TraceLevel, "Sending database change request...")
	check := true
	if r.FormValue("Method") == "POST" {
		var product models.Product
		response = Request("/products/"+r.FormValue("Product.Id"), "GET", GetCurrentToken(w, r), nil, &product)
		ResponseCheck(response, "/products", "GET")
		for i := range product.Sets {
			if product.Sets[i].CharacteristicID == set.CharacteristicID {
				check = false
			}
		}
		if !check {
			GlobalError = "Данная характеристика уже присвоена продукту"
			requestSet["error"] = "This characteristic has already been assigned to the product"
			SaveLog(requestSet, log.ErrorLevel, "Failed database change")
		}
	}
	if check {
		switch r.FormValue("Method") {
		case "POST":
			set.ProductID, _ = strconv.Atoi(r.FormValue("Product.Id"))
			response = Request("/sets", "POST", GetCurrentToken(w, r), set, &result)
		case "PUT":
			var setP models.Set
			Request("/sets/"+r.FormValue("Id"), "GET", GetCurrentToken(w, r), nil, &setP)
			set.ProductID = setP.ProductID
			response = Request("/sets/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), set, &result)
		case "DELETE":
			response = Request("/sets/"+r.FormValue("Id"), "DELETE", GetCurrentToken(w, r), nil, &result)
		}
		if !(response.StatusCode == 200 || response.StatusCode == 201) {
			GlobalError = result["message"]
			requestSet["error"] = result["message"]
			SaveLog(requestSet, log.ErrorLevel, "Failed database change")
		} else {
			SaveLog(requestSet, log.InfoLevel, "Successful database change")
		}
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func CudCharacteristic(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	relation, _ := strconv.Atoi(r.FormValue("Relation"))
	_type := r.FormValue("Type")
	if _type == "" {
		_type = "END"
	}
	characteristic := models.Characteristic{
		ID:       id,
		Name:     r.FormValue("Name"),
		Type:     _type,
		Relation: relation,
		IsExist:  true,
	}
	var result map[string]string
	var response *http.Response
	requestCharacteristic := log.Fields{
		"group":    "crud",
		"method":   r.FormValue("Method"),
		"table":    "characteristic",
		"id":       characteristic.ID,
		"name":     characteristic.Name,
		"type":     characteristic.Type,
		"relation": characteristic.Relation,
	}
	SaveLog(requestCharacteristic, log.TraceLevel, "Sending database change request...")
	switch r.FormValue("Method") {
	case "POST":
		response = Request("/characteristics", "POST", GetCurrentToken(w, r), characteristic, &result)
	case "PUT":
		response = Request("/characteristics/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), characteristic, &result)
	case "DELETE":
		if r.FormValue("Type") == "Physical" {
			response = Request("/characteristics/"+r.FormValue("Id"), "DELETE", GetCurrentToken(w, r), nil, &result)
		} else {
			var item models.Characteristic
			Request("/characteristics/"+r.FormValue("Id"), "GET", GetCurrentToken(w, r), nil, &item)
			item.IsExist = false
			response = Request("/characteristics/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), item, &result)
		}
	}
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		GlobalError = result["message"]
		requestCharacteristic["error"] = result["message"]
		SaveLog(requestCharacteristic, log.ErrorLevel, "Failed database change")
	} else {
		SaveLog(requestCharacteristic, log.InfoLevel, "Successful database change")
	}
	http.Redirect(w, r, "/admin#nav-characteristic", http.StatusSeeOther)
}

func CudCategory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	relation, _ := strconv.Atoi(r.FormValue("Relation"))
	category := models.Category{
		ID:       id,
		Name:     r.FormValue("Name"),
		Relation: relation,
		IsExist:  true,
	}
	var result map[string]string
	var response *http.Response
	requestCategory := log.Fields{
		"group":    "crud",
		"method":   r.FormValue("Method"),
		"table":    "category",
		"id":       category.ID,
		"name":     category.Name,
		"relation": category.Relation,
	}
	SaveLog(requestCategory, log.TraceLevel, "Sending database change request...")
	switch r.FormValue("Method") {
	case "POST":
		response = Request("/categories", "POST", GetCurrentToken(w, r), category, &result)
	case "PUT":
		response = Request("/categories/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), category, &result)
	case "DELETE":
		if r.FormValue("Type") == "Physical" {
			response = Request("/categories/"+r.FormValue("Id"), "DELETE", GetCurrentToken(w, r), nil, &result)
		} else {
			var item models.Category
			Request("/categories/"+r.FormValue("Id"), "GET", GetCurrentToken(w, r), nil, &item)
			item.IsExist = false
			response = Request("/categories/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), item, &result)
		}
	}
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		GlobalError = result["message"]
		requestCategory["error"] = result["message"]
		SaveLog(requestCategory, log.ErrorLevel, "Failed database change")
	} else {
		SaveLog(requestCategory, log.InfoLevel, "Successful database change")
	}
	http.Redirect(w, r, "/admin#nav-category", http.StatusSeeOther)
}

func CudImages(w http.ResponseWriter, r *http.Request) {
	requestImage := log.Fields{
		"group":  "crud",
		"method": r.FormValue("Method"),
		"table":  "images",
	}
	switch r.FormValue("Method") {
	case "POST":
		SaveLog(requestImage, log.TraceLevel, "Sending image to server...")
		var MaxUploadSize int64 = 1024 * 1024 * 15
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			GlobalError = err.Error()
			requestImage["error"] = err.Error()
			SaveLog(requestImage, log.ErrorLevel, "Failed image submission")
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		}
		files := r.MultipartForm.File["File"]
		for _, fileHeader := range files {
			if fileHeader.Size > MaxUploadSize {
				GlobalError = "Загружаемый файл слишком большой: " + strconv.Itoa(int(fileHeader.Size)) + ". Пожалуйста, выбирите файл размером меньше 15MB."
				requestImage["error"] = "The size of uploaded file is too big: " + strconv.Itoa(int(fileHeader.Size))
				SaveLog(requestImage, log.ErrorLevel, "Failed image submission")
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			file, err := fileHeader.Open()
			if err != nil {
				GlobalError = err.Error()
				requestImage["error"] = err.Error()
				SaveLog(requestImage, log.ErrorLevel, "Failed image submission")
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			defer file.Close()
			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				GlobalError = err.Error()
				requestImage["error"] = err.Error()
				SaveLog(requestImage, log.ErrorLevel, "Failed image submission")
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			filetype := http.DetectContentType(buff)
			if filetype != "image/jpeg" && filetype != "image/png" {
				GlobalError = "Формат загружаемого файла не разрешён. Пожалуйста, выбирите JPEG или PNG изображения"
				requestImage["error"] = "Unauthorised format of uploaded file."
				SaveLog(requestImage, log.ErrorLevel, "Failed image submission")
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				GlobalError = err.Error()
				requestImage["error"] = err.Error()
				SaveLog(requestImage, log.ErrorLevel, "Failed image submission")
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			var f *os.File
			f, err = os.Create(fmt.Sprintf("./images/%s", fileHeader.Filename))
			if err != nil {
				GlobalError = err.Error()
				requestImage["error"] = err.Error()
				SaveLog(requestImage, log.ErrorLevel, "Failed image submission")
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			defer f.Close()
			_, err = io.Copy(f, file)
			if err != nil {
				GlobalError = err.Error()
				requestImage["error"] = err.Error()
				SaveLog(requestImage, log.ErrorLevel, "Failed image submission")
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			SaveLog(requestImage, log.InfoLevel, "Successful image submission")
		}
	case "DELETE":
		SaveLog(requestImage, log.TraceLevel, "Removing image from server...")
		err := os.Remove("." + r.FormValue("Path"))
		if err != nil {
			GlobalError = err.Error()
			requestImage["error"] = err.Error()
			SaveLog(requestImage, log.ErrorLevel, "Failed image submission")
		}
		SaveLog(requestImage, log.InfoLevel, "Successful image submission")
	}
	http.Redirect(w, r, "/admin#nav-image", http.StatusSeeOther)
}

func CudStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	status := models.Status{
		ID:      id,
		Status:  r.FormValue("Status"),
		IsExist: true,
	}
	var result map[string]string
	var response *http.Response
	requestStatus := log.Fields{
		"group":  "crud",
		"method": r.FormValue("Method"),
		"table":  "status",
		"id":     status.ID,
		"status": status.Status,
	}
	SaveLog(requestStatus, log.TraceLevel, "Sending database change request...")
	switch r.FormValue("Method") {
	case "POST":
		response = Request("/statuses", "POST", GetCurrentToken(w, r), status, &result)
	case "PUT":
		response = Request("/statuses/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), status, &result)
	case "DELETE":
		if r.FormValue("Type") == "Physical" {
			response = Request("/statuses/"+r.FormValue("Id"), "DELETE", GetCurrentToken(w, r), nil, &result)
		} else {
			var item models.Status
			Request("/statuses/"+r.FormValue("Id"), "GET", GetCurrentToken(w, r), nil, &item)
			item.IsExist = false
			response = Request("/statuses/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), item, &result)
		}
	}
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		GlobalError = result["message"]
		requestStatus["error"] = result["message"]
		SaveLog(requestStatus, log.ErrorLevel, "Failed database change")
	} else {
		SaveLog(requestStatus, log.InfoLevel, "Successful database change")
	}
	http.Redirect(w, r, "/admin#nav-status", http.StatusSeeOther)
}

func CudRole(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	role := models.Role{
		ID:      id,
		Name:    r.FormValue("Name"),
		IsExist: true,
	}
	var result map[string]string
	var response *http.Response
	requestRole := log.Fields{
		"group":  "crud",
		"method": r.FormValue("Method"),
		"table":  "role",
		"id":     role.ID,
		"name":   role.Name,
	}
	SaveLog(requestRole, log.TraceLevel, "Sending database change request...")
	switch r.FormValue("Method") {
	case "POST":
		response = Request("/roles", "POST", GetCurrentToken(w, r), role, &result)
	case "PUT":
		response = Request("/roles/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), role, &result)
	case "DELETE":
		if r.FormValue("Type") == "Physical" {
			response = Request("/roles/"+r.FormValue("Id"), "DELETE", GetCurrentToken(w, r), nil, &result)
		} else {
			var item models.Role
			Request("/roles/"+r.FormValue("Id"), "GET", GetCurrentToken(w, r), nil, &item)
			item.IsExist = false
			response = Request("/roles/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), item, &result)
		}
	}
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		GlobalError = result["message"]
		requestRole["error"] = result["message"]
		SaveLog(requestRole, log.ErrorLevel, "Failed database change")
	} else {
		SaveLog(requestRole, log.InfoLevel, "Successful database change")
	}
	http.Redirect(w, r, "/admin#nav-role", http.StatusSeeOther)
}

func CudUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	role, _ := strconv.Atoi(r.FormValue("Role.Id"))
	var password string
	if r.FormValue("Password") == "" {
		var user models.User
		Request("/users/"+r.FormValue("Id"), "GET", GetCurrentToken(w, r), nil, &user)
		password = user.Password
	} else {
		password = r.FormValue("Password")
	}
	user := models.User{
		ID:         id,
		Email:      r.FormValue("Email"),
		Password:   password,
		Phone:      r.FormValue("Phone"),
		LastName:   r.FormValue("LastName"),
		Name:       r.FormValue("Name"),
		MiddleName: r.FormValue("MiddleName"),
		Gender:     r.FormValue("Gender"),
		RoleID:     role,
		IsExist:    true,
	}
	requestUser := log.Fields{
		"group":       "crud",
		"method":      r.FormValue("Method"),
		"table":       "user",
		"id":          user.ID,
		"email":       user.Email,
		"phone":       user.Phone,
		"last_name":   user.LastName,
		"name":        user.Name,
		"middle_name": user.MiddleName,
		"gender":      user.Gender,
		"role_id":     user.RoleID,
	}
	SaveLog(requestUser, log.TraceLevel, "Sending database change request...")
	var result map[string]string
	var response *http.Response
	switch r.FormValue("Method") {
	case "POST":
		response = Request("/users", "POST", GetCurrentToken(w, r), user, &result)
	case "PUT":
		response = Request("/users/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), user, &result)
	case "DELETE":
		if r.FormValue("Type") == "Physical" {
			response = Request("/users/"+r.FormValue("Id"), "DELETE", GetCurrentToken(w, r), nil, &result)
		} else {
			var item models.User
			Request("/users/"+r.FormValue("Id"), "GET", GetCurrentToken(w, r), nil, &item)
			item.IsExist = false
			response = Request("/users/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), item, &result)
		}
	}
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		GlobalError = result["message"]
		requestUser["error"] = result["message"]
		SaveLog(requestUser, log.ErrorLevel, "Failed database change")
	} else {
		SaveLog(requestUser, log.InfoLevel, "Successful database change")
	}
	http.Redirect(w, r, "/admin#nav-user", http.StatusSeeOther)
}

func CudCard(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	discount, _ := strconv.Atoi(r.FormValue("Discount"))
	card := models.Card{
		ID:       id,
		Date:     r.FormValue("Date"),
		Discount: discount,
	}
	var result map[string]string
	var response *http.Response
	requestCard := log.Fields{
		"group":    "crud",
		"method":   r.FormValue("Method"),
		"table":    "card",
		"id":       card.ID,
		"date":     card.Date,
		"discount": card.Discount,
	}
	SaveLog(requestCard, log.TraceLevel, "Sending database change request...")
	switch r.FormValue("Method") {
	case "POST":
		card.Date = time.Now().Format(time.DateTime)
		response = Request("/cards", "POST", GetCurrentToken(w, r), card, &result)
	case "PUT":
		response = Request("/cards/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), card, &result)
	case "DELETE":
		response = Request("/cards/"+r.FormValue("Id"), "DELETE", GetCurrentToken(w, r), nil, &result)
	}
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		GlobalError = result["message"]
		requestCard["error"] = result["message"]
		SaveLog(requestCard, log.ErrorLevel, "Failed database change")
	} else {
		SaveLog(requestCard, log.InfoLevel, "Successful database change")
	}
	http.Redirect(w, r, "/admin#nav-user", http.StatusSeeOther)
}

func CudOrder(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	status, _ := strconv.Atoi(r.FormValue("Status.Id"))
	user, _ := strconv.Atoi(r.FormValue("User.Id"))
	date := r.FormValue("Date")
	if date == "" {
		date = time.Now().Format(time.DateTime)
	}
	order := models.Order{
		ID:       id,
		Date:     date,
		Address:  r.FormValue("Address"),
		StatusID: status,
		UserID:   user,
	}
	var result map[string]string
	var response *http.Response
	requestOrder := log.Fields{
		"group":     "crud",
		"method":    r.FormValue("Method"),
		"table":     "order",
		"id":        order.ID,
		"date":      order.Date,
		"address":   order.Address,
		"status_id": order.StatusID,
		"user_id":   order.UserID,
	}
	SaveLog(requestOrder, log.TraceLevel, "Sending database change request...")
	switch r.FormValue("Method") {
	case "POST":
		response = Request("/orders", "POST", GetCurrentToken(w, r), order, &result)
	case "PUT":
		response = Request("/orders/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), order, &result)
	case "DELETE":
		response = Request("/orders/"+r.FormValue("Id"), "DELETE", GetCurrentToken(w, r), nil, &result)
	}
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		GlobalError = result["message"]
		requestOrder["error"] = result["message"]
		SaveLog(requestOrder, log.ErrorLevel, "Failed database change")
	} else {
		SaveLog(requestOrder, log.InfoLevel, "Successful database change")
	}
	http.Redirect(w, r, "/admin#nav-order", http.StatusSeeOther)
}

func CudPosition(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("Id"))
	price, _ := strconv.ParseFloat(r.FormValue("Price"), 64)
	amount, _ := strconv.Atoi(r.FormValue("Amount"))
	order, _ := strconv.Atoi(r.FormValue("Order.Id"))
	product, _ := strconv.Atoi(r.FormValue("Product.Id"))
	position := models.Position{
		ID:            id,
		CheckoutPrice: price,
		Amount:        amount,
		OrderID:       order,
		ProductID:     product,
	}
	var result map[string]string
	var response *http.Response
	requestPosition := log.Fields{
		"group":          "crud",
		"method":         r.FormValue("Method"),
		"table":          "position",
		"id":             position.ID,
		"checkout_price": position.CheckoutPrice,
		"amount":         position.Amount,
		"order_id":       position.OrderID,
		"product_id":     position.ProductID,
	}
	SaveLog(requestPosition, log.TraceLevel, "Sending database change request...")
	switch r.FormValue("Method") {
	case "POST":
		response = Request("/positions", "POST", GetCurrentToken(w, r), position, &result)
	case "PUT":
		response = Request("/positions/"+r.FormValue("Id"), "PUT", GetCurrentToken(w, r), position, &result)
	case "DELETE":
		response = Request("/positions/"+r.FormValue("Id"), "DELETE", GetCurrentToken(w, r), nil, &result)
	}
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		GlobalError = result["message"]
		requestPosition["error"] = result["message"]
		SaveLog(requestPosition, log.ErrorLevel, "Failed database change")
	} else {
		SaveLog(requestPosition, log.InfoLevel, "Successful database change")
	}
	http.Redirect(w, r, "/admin#nav-order", http.StatusSeeOther)
}
