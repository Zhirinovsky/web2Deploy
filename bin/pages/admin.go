package pages

import (
	"bufio"
	"encoding/json"
	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"web2/bin"
	"web2/bin/models"
	"web2/bin/structures"
)

type ChartItem struct {
	Y     int
	X     int
	Label string
}

type Log struct {
	Level string
	Group string
	Time  string
	Msg   string
	Data  string
}

type AdminPageObject struct {
	BaseObject      structures.BaseObject
	Products        []models.Product
	Categories      []models.Category
	Characteristics []models.Characteristic
	Roles           []models.Role
	Users           []models.User
	Statuses        []models.Status
	Orders          []models.Order
	Positions       []models.Position
	Images          []models.Image
	PieData         []ChartItem
	GraphData       []ChartItem
	ColumnData      []ChartItem
	Logs            []Log
}

func Admin(w http.ResponseWriter, r *http.Request) {
	bin.RefreshAPI()
	pageObject := getAdminData()
	if pageObject.BaseObject.ErrorStr == "" {
		pageObject.BaseObject.ErrorStr = bin.GetError(&bin.GlobalError)
	}
	pageObject.BaseObject.CurrentUser = bin.GetCurrentUser(w, r)
	pageObject.BaseObject.CsrfField = csrf.TemplateField(r)
	t, err := template.New("admin.html").Funcs(sprig.FuncMap()).ParseFiles("views/admin.html")
	bin.CheckErr(err)
	template.Must(t.ParseGlob("views/templates/*.html"))
	template.Must(t.ParseGlob("views/templates/tables/*.html"))
	template.Must(t.ParseGlob("views/templates/modals/*.html"))
	bin.Refresh(w, r)
	err = t.Execute(w, pageObject)
	bin.CheckErr(err)
}

func getAdminData() AdminPageObject {
	var pageObject AdminPageObject
	response := bin.Request("/products", "GET", bin.ServerToken, nil, &pageObject.Products)
	bin.ResponseCheck(response, "/products", "GET")
	response = bin.Request("/characteristics", "GET", bin.ServerToken, nil, &pageObject.Characteristics)
	bin.ResponseCheck(response, "/characteristics", "GET")
	response = bin.Request("/categories", "GET", bin.ServerToken, nil, &pageObject.Categories)
	bin.ResponseCheck(response, "/categories", "GET")
	response = bin.Request("/roles", "GET", bin.ServerToken, nil, &pageObject.Roles)
	bin.ResponseCheck(response, "/roles", "GET")
	response = bin.Request("/users", "GET", bin.ServerToken, nil, &pageObject.Users)
	bin.ResponseCheck(response, "/users", "GET")
	response = bin.Request("/statuses", "GET", bin.ServerToken, nil, &pageObject.Statuses)
	bin.ResponseCheck(response, "/statuses", "GET")
	response = bin.Request("/orders", "GET", bin.ServerToken, nil, &pageObject.Orders)
	bin.ResponseCheck(response, "/orders", "GET")
	response = bin.Request("/positions", "GET", bin.ServerToken, nil, &pageObject.Positions)
	bin.ResponseCheck(response, "/positions", "GET")
	dir, err := os.Open("./images")
	bin.NullErrorCheck(&pageObject.BaseObject.ErrorStr, err)
	files, err := dir.Readdir(-1)
	bin.NullErrorCheck(&pageObject.BaseObject.ErrorStr, err)
	for _, file := range files {
		image := models.Image{
			Name: file.Name(),
			Path: "/images/" + file.Name(),
			Type: "file",
		}
		pageObject.Images = append(pageObject.Images, image)
	}
	getPieData(&pageObject)
	getGraphData(&pageObject)
	getColumnData(&pageObject)
	getLogs(&pageObject)
	return pageObject
}

func getPieData(pageObject *AdminPageObject) {
	var buffer []ChartItem
	for _, category := range pageObject.Categories {
		if category.Relation != 0 {
			item := ChartItem{
				Y:     0,
				Label: category.Name,
			}
			buffer = append(buffer, item)
		}
	}
	for _, order := range pageObject.Orders {
		for _, position := range order.Positions {
			for i := range buffer {
				if buffer[i].Label == position.Product.Category.Name {
					buffer[i].Y += position.Amount * int(position.CheckoutPrice)
				}
			}
		}
	}
	for i := range buffer {
		if buffer[i].Y != 0 {
			pageObject.PieData = append(pageObject.PieData, buffer[i])
		}
	}
}

func getGraphData(pageObject *AdminPageObject) {
	for _, category := range pageObject.Categories {
		if category.Relation != 0 {
			for i := 0; i < 12; i++ {
				item := ChartItem{
					Y:     0,
					X:     i,
					Label: category.Name,
				}
				pageObject.GraphData = append(pageObject.GraphData, item)
			}
		}
	}
	for _, order := range pageObject.Orders {
		for _, position := range order.Positions {
			for i := range pageObject.GraphData {
				currentDate, _ := strconv.Atoi(order.Date[5:7])
				if pageObject.GraphData[i].Label == position.Product.Category.Name && currentDate == pageObject.GraphData[i].X {
					pageObject.GraphData[i].Y += position.Amount * int(position.CheckoutPrice)
				}
			}
		}
	}
}

func getColumnData(pageObject *AdminPageObject) {
	for _, product := range pageObject.Products {
		item := ChartItem{
			Y:     product.Amount,
			Label: product.Name,
		}
		pageObject.ColumnData = append(pageObject.ColumnData, item)
	}
}

func getLogs(pageObject *AdminPageObject) {
	file, err := os.Open("data/logs.log")
	bin.NullErrorCheck(&pageObject.BaseObject.ErrorStr, err)
	defer file.Close()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := []byte(sc.Text())
		var Map map[string]any
		err = json.Unmarshal(line, &Map)
		bin.NullErrorCheck(&pageObject.BaseObject.ErrorStr, err)
		var log Log
		for k, v := range Map {
			switch k {
			case "group":
				log.Group = v.(string)
			case "level":
				log.Level = v.(string)
			case "msg":
				log.Msg = v.(string)
			case "time":
				log.Time = v.(string)
			default:
				switch reflect.TypeOf(v).Kind() {
				case reflect.Int:
					v = strconv.Itoa(v.(int))
				case reflect.Float64:
					v = strconv.FormatFloat(v.(float64), 'f', -1, 64)
				}
				if log.Data != "" {
					log.Data += ", " + k + "=" + v.(string)
				} else {
					log.Data = k + "=" + v.(string)
				}
			}
		}
		pageObject.Logs = append(pageObject.Logs, log)
	}
}
