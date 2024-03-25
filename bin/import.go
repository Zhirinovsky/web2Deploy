package bin

import (
	"encoding/csv"
	"github.com/prongbang/excelx"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

func Import(w http.ResponseWriter, r *http.Request) {
	SaveLog(log.Fields{
		"group":    "export/import",
		"subgroup": "import",
		"table":    r.FormValue("Table"),
	}, log.TraceLevel, "Import attempt...")
	var MaxUploadSize int64 = 1024 * 1024 * 15
	err := r.ParseMultipartForm(32 << 20)
	GlobalCheck(err)
	if err == nil {
		files := r.MultipartForm.File["File"]
		for _, fileHeader := range files {
			if fileHeader.Size > MaxUploadSize {
				GlobalError = "Загружаемый файл слишком большой. Пожалуйста, выбирите файл размером меньше 15MB"
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			file, err := fileHeader.Open()
			GlobalCheck(err)
			if err != nil {
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			defer file.Close()
			buff := make([]byte, 512)
			_, err = file.Read(buff)
			GlobalCheck(err)
			if err != nil {
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			filetype := http.DetectContentType(buff)
			if filetype != "application/octet-stream" && filetype != "application/zip" {
				GlobalError = "Формат импортируемого файла должен быть csv или xlsx."
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			_, err = file.Seek(0, io.SeekStart)
			GlobalCheck(err)
			if err != nil {
				http.Redirect(w, r, "/admin", http.StatusSeeOther)
				return
			}
			switch filetype {
			case "application/octet-stream":
				reader := csv.NewReader(file)
				info, err := reader.ReadAll()
				if err != nil {
					GlobalError = "Не удаётся считать содержимое файла"
					http.Redirect(w, r, "/admin", http.StatusSeeOther)
					return
				}
				for i := 1; i < len(info); i++ {
					switch r.FormValue("Table") {
					case "products":
						id, _ := strconv.Atoi(info[i][0])
						price, _ := strconv.ParseFloat(info[i][2], 64)
						amount, _ := strconv.Atoi(info[i][3])
						discount, _ := strconv.Atoi(info[i][4])
						exist, _ := strconv.ParseBool(info[i][7])
						product := ProductExport{
							ID:         id,
							Name:       info[i][1],
							Price:      price,
							Amount:     amount,
							Discount:   discount,
							ImageLink:  info[i][5],
							CategoryID: info[i][6],
							IsExist:    exist,
						}
						Request("/products", "POST", GetCurrentToken(w, r), product, nil)
					case "orders":
						id, _ := strconv.Atoi(info[i][0])
						statusid, _ := strconv.Atoi(info[i][3])
						userid, _ := strconv.Atoi(info[i][4])
						order := OrderExport{
							ID:       id,
							Date:     info[i][1],
							Address:  info[i][2],
							StatusID: statusid,
							UserID:   userid,
						}
						Request("/orders", "POST", GetCurrentToken(w, r), order, nil)
					case "positions":
						id, _ := strconv.Atoi(info[i][0])
						checkoutprice, _ := strconv.ParseFloat(info[i][1], 64)
						amount, _ := strconv.Atoi(info[i][2])
						orderid, _ := strconv.Atoi(info[i][3])
						productid, _ := strconv.Atoi(info[i][4])
						position := PositionExport{
							ID:            id,
							CheckoutPrice: checkoutprice,
							Amount:        amount,
							OrderID:       orderid,
							ProductID:     productid,
						}
						Request("/positions", "POST", GetCurrentToken(w, r), position, nil)
					case "statuses":
						id, _ := strconv.Atoi(info[i][0])
						exist, _ := strconv.ParseBool(info[i][2])
						status := StatusExport{
							ID:      id,
							Status:  info[i][1],
							IsExist: exist,
						}
						Request("/statuses", "POST", GetCurrentToken(w, r), status, nil)
					case "characteristics":
						id, _ := strconv.Atoi(info[i][0])
						relation, _ := strconv.Atoi(info[i][3])
						exist, _ := strconv.ParseBool(info[i][4])
						characteristic := CharacteristicExport{
							ID:       id,
							Name:     info[i][1],
							Type:     info[i][2],
							Relation: relation,
							IsExist:  exist,
						}
						Request("/characteristics", "POST", GetCurrentToken(w, r), characteristic, nil)
					case "categories":
						id, _ := strconv.Atoi(info[i][0])
						relation, _ := strconv.Atoi(info[i][2])
						exist, _ := strconv.ParseBool(info[i][3])
						category := CategoryExport{
							ID:       id,
							Name:     info[i][1],
							Relation: relation,
							IsExist:  exist,
						}
						Request("/categories", "POST", GetCurrentToken(w, r), category, nil)
					case "users":
						id, _ := strconv.Atoi(info[i][0])
						roleid, _ := strconv.Atoi(info[i][8])
						exist, _ := strconv.ParseBool(info[i][9])
						user := UserExport{
							ID:         id,
							Email:      info[i][1],
							Password:   info[i][2],
							Phone:      info[i][3],
							LastName:   info[i][4],
							Name:       info[i][5],
							MiddleName: info[i][6],
							Gender:     info[i][7],
							RoleID:     roleid,
							IsExist:    exist,
						}
						Request("/users", "POST", GetCurrentToken(w, r), user, nil)
					case "roles":
						id, _ := strconv.Atoi(info[i][0])
						exist, _ := strconv.ParseBool(info[i][2])
						role := RoleExport{
							ID:      id,
							Name:    info[i][1],
							IsExist: exist,
						}
						Request("/roles", "POST", GetCurrentToken(w, r), role, nil)
					}
				}
			case "application/zip":
				switch r.FormValue("Table") {
				case "products":
					items, err := excelx.Parse[ProductExport](file)
					GlobalCheck(err)
					for _, item := range items {
						Request("/products", "POST", GetCurrentToken(w, r), item, nil)
					}
				case "orders":
					items, err := excelx.Parse[OrderExport](file)
					GlobalCheck(err)
					for _, item := range items {
						Request("/orders", "POST", GetCurrentToken(w, r), item, nil)
					}
				case "positions":
					items, err := excelx.Parse[PositionExport](file)
					GlobalCheck(err)
					for _, item := range items {
						Request("/positions", "POST", GetCurrentToken(w, r), item, nil)
					}
				case "statuses":
					items, err := excelx.Parse[StatusExport](file)
					GlobalCheck(err)
					for _, item := range items {
						Request("/statuses", "POST", GetCurrentToken(w, r), item, nil)
					}
				case "characteristics":
					items, err := excelx.Parse[CharacteristicExport](file)
					GlobalCheck(err)
					for _, item := range items {
						Request("/characteristics", "POST", GetCurrentToken(w, r), item, nil)
					}
				case "categories":
					items, err := excelx.Parse[CategoryExport](file)
					GlobalCheck(err)
					for _, item := range items {
						Request("/categories", "POST", GetCurrentToken(w, r), item, nil)
					}
				case "users":
					items, err := excelx.Parse[UserExport](file)
					GlobalCheck(err)
					for _, item := range items {
						Request("/users", "POST", GetCurrentToken(w, r), item, nil)
					}
				case "roles":
					items, err := excelx.Parse[RoleExport](file)
					GlobalCheck(err)
					for _, item := range items {
						Request("/roles", "POST", GetCurrentToken(w, r), item, nil)
					}
				}
			}
		}
	}
	if GlobalError == "" {
		SaveLog(log.Fields{
			"group":    "export/import",
			"subgroup": "import",
			"table":    r.FormValue("Table"),
		}, log.InfoLevel, "Successful import")
	} else {
		SaveLog(log.Fields{
			"group":    "export/import",
			"subgroup": "import",
			"table":    r.FormValue("Table"),
			"error":    GlobalError,
		}, log.ErrorLevel, "Import error")
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
