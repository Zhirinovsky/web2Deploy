package bin

import (
	"encoding/csv"
	"github.com/dnlo/struct2csv"
	"github.com/prongbang/excelx"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

type ProductExport struct {
	ID         int     `header:"ID" no:"1"`
	Name       string  `header:"Name" no:"2"`
	Price      float64 `header:"Price" no:"3"`
	Amount     int     `header:"Amount" no:"4"`
	Discount   int     `header:"Discount" no:"5"`
	ImageLink  string  `header:"ImageLink" no:"6"`
	CategoryID string  `header:"CategoryID" no:"7"`
	IsExist    bool    `header:"IsExist" no:"8"`
}

type CharacteristicExport struct {
	ID       int    `header:"ID" no:"1"`
	Name     string `header:"Name" no:"2"`
	Type     string `header:"Type" no:"3"`
	Relation int    `header:"Relation" no:"4"`
	IsExist  bool   `header:"IsExist" no:"5"`
}

type UserExport struct {
	ID         int    `header:"ID" no:"1"`
	Email      string `header:"Email" no:"2"`
	Password   string `header:"Password" no:"3"`
	Phone      string `header:"Phone" no:"4"`
	LastName   string `header:"LastName" no:"5"`
	Name       string `header:"Name" no:"6"`
	MiddleName string `header:"MiddleName" no:"7"`
	Gender     string `header:"Gender" no:"8"`
	RoleID     int    `header:"RoleID" no:"9"`
	IsExist    bool   `header:"IsExist" no:"10"`
}

type OrderExport struct {
	ID       int    `header:"ID" no:"1"`
	Date     string `header:"Date" no:"2"`
	Address  string `header:"Address" no:"3"`
	StatusID int    `header:"StatusID" no:"4"`
	UserID   int    `header:"UserID" no:"5"`
}

type PositionExport struct {
	ID            int     `header:"ID" no:"1"`
	CheckoutPrice float64 `header:"CheckoutPrice" no:"2"`
	Amount        int     `header:"Amount" no:"3"`
	OrderID       int     `header:"OrderID" no:"4"`
	ProductID     int     `header:"ProductID" no:"5"`
}

type StatusExport struct {
	ID      int    `header:"ID" no:"1"`
	Status  string `header:"Status" no:"2"`
	IsExist bool   `header:"IsExist" no:"3"`
}

type CategoryExport struct {
	ID       int    `header:"ID" no:"1"`
	Name     string `header:"Name" no:"2"`
	Relation int    `header:"Relation" no:"3"`
	IsExist  bool   `header:"IsExist" no:"4"`
}

type RoleExport struct {
	ID      int    `header:"ID" no:"1"`
	Name    string `header:"Name" no:"2"`
	IsExist bool   `header:"IsExist" no:"3"`
}

func Export(w http.ResponseWriter, r *http.Request) {
	SaveLog(log.Fields{
		"group":    "export/import",
		"subgroup": "export",
		"table":    r.FormValue("Table"),
		"datatype": r.FormValue("Type"),
	}, log.TraceLevel, "Export attempt...")
	var items []any
	var data [][]string
	switch r.FormValue("Table") {
	case "products":
		var products []ProductExport
		Request("/"+r.FormValue("Table"), "GET", ServerToken, nil, &products)
		for _, item := range products {
			items = append(items, item)
		}
	case "orders":
		var orders []OrderExport
		Request("/"+r.FormValue("Table"), "GET", ServerToken, nil, &orders)
		for _, item := range orders {
			items = append(items, item)
		}
	case "positions":
		var positions []PositionExport
		Request("/"+r.FormValue("Table"), "GET", ServerToken, nil, &positions)
		for _, item := range positions {
			items = append(items, item)
		}
	case "statuses":
		var statuses []StatusExport
		Request("/"+r.FormValue("Table"), "GET", ServerToken, nil, &statuses)
		for _, item := range statuses {
			items = append(items, item)
		}
	case "characteristics":
		var characteristics []CharacteristicExport
		Request("/"+r.FormValue("Table"), "GET", ServerToken, nil, &characteristics)
		for _, item := range characteristics {
			items = append(items, item)
		}
	case "categories":
		var categories []CategoryExport
		Request("/"+r.FormValue("Table"), "GET", ServerToken, nil, &categories)
		for _, item := range categories {
			items = append(items, item)
		}
	case "users":
		var users []UserExport
		Request("/"+r.FormValue("Table"), "GET", ServerToken, nil, &users)
		for _, item := range users {
			items = append(items, item)
		}
	case "roles":
		var roles []RoleExport
		Request("/"+r.FormValue("Table"), "GET", ServerToken, nil, &roles)
		for _, item := range roles {
			items = append(items, item)
		}
	}
	switch r.FormValue("Type") {
	case "csv":
		enc := struct2csv.New()
		row, err := enc.GetColNames(items[0])
		GlobalCheck(err)
		data = append(data, row)
		for _, item := range items {
			row, err = enc.GetRow(item)
			GlobalCheck(err)
			data = append(data, row)
		}
		file, err := os.Create(r.FormValue("Table") + ".csv")
		defer file.Close()
		GlobalCheck(err)
		writer := csv.NewWriter(file)
		defer writer.Flush()
		err = writer.WriteAll(data)
		GlobalCheck(err)
		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(r.FormValue("Table")+".csv"))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, r.FormValue("Table")+".csv")
		err = file.Close()
		GlobalCheck(err)
		err = os.Remove(r.FormValue("Table") + ".csv")
		GlobalCheck(err)
	case "xlsx":
		file, err := excelx.Convert(items)
		GlobalCheck(err)
		err = excelx.ResponseWriter(file, w, r.FormValue("Table")+".xlsx")
		GlobalCheck(err)
	}
	if GlobalError == "" {
		SaveLog(log.Fields{
			"group":    "export/import",
			"subgroup": "export",
			"table":    r.FormValue("Table"),
			"datatype": r.FormValue("Type"),
		}, log.InfoLevel, "Successful export")
	} else {
		SaveLog(log.Fields{
			"group":    "export/import",
			"subgroup": "export",
			"table":    r.FormValue("Table"),
			"datatype": r.FormValue("Type"),
			"error":    GlobalError,
		}, log.ErrorLevel, "Export error")
	}
}
