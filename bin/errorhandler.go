package bin

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		TotalFailures.With(prometheus.Labels{"fatal_error": err.Error()}).Inc()
		SaveLog(log.Fields{
			"group": "errors",
			"error": err,
		}, log.PanicLevel, err.Error())
	}
}

func NullErrorCheck(errorStr *string, error error) {
	if error != nil {
		SaveLog(log.Fields{
			"group": "errors",
			"error": error,
		}, log.ErrorLevel, error.Error())
		Failures.With(prometheus.Labels{"error": error.Error()}).Inc()
		if *errorStr == "" {
			*errorStr = error.Error()
		}
	}
}

func GetError(errorStr *string) string {
	buffer := *errorStr
	*errorStr = ""
	return buffer
}

func GlobalCheck(error error) {
	if error != nil {
		SaveLog(log.Fields{
			"group": "errors",
			"error": error,
		}, log.ErrorLevel, error.Error())
		if GlobalError == "" {
			GlobalError = error.Error()
		}
	}
}

func ResponseCheck(response *http.Response, url string, requestType string) {
	if response.StatusCode != 200 {
		if GlobalError == "" {
			GlobalError = "Ошибка при получении данных таблицы из Api по адресу: " + url
		}
		SaveLog(log.Fields{
			"group":        "errors",
			"error":        GlobalError,
			"url":          url,
			"type-request": requestType,
		}, log.ErrorLevel, "Error while retrieving table data")
	}
}
