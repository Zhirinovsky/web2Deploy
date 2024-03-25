package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"web2/bin"
	"web2/bin/pages"
)

func main() {
	//Настройка логов
	bin.ConfigLogs()
	bin.SaveLog(log.Fields{"group": "server"}, log.TraceLevel, "Server startup...")
	r := mux.NewRouter()
	//Страницы
	r.HandleFunc("/", pages.MainPage)
	r.HandleFunc("/login", pages.Login)
	r.HandleFunc("/registration", pages.Registration)
	r.HandleFunc("/admin", pages.Admin)
	r.HandleFunc("/profile", pages.Profile)
	//Методы сессии
	r.HandleFunc("/logout", bin.Logout)
	//Методы главной страницы
	r.HandleFunc("/search", pages.SearchProduct)
	r.HandleFunc("/category", pages.SetCategory)
	//Методы профиля
	r.HandleFunc("/personalisation", pages.ChangePersonalData)
	r.HandleFunc("/change-password", pages.ChangePassword)
	//Cud действия
	r.HandleFunc("/products", bin.CudProduct)
	r.HandleFunc("/characteristics", bin.CudCharacteristic)
	r.HandleFunc("/categories", bin.CudCategory)
	r.HandleFunc("/sets", bin.CudSet)
	r.HandleFunc("/images", bin.CudImages)
	r.HandleFunc("/statuses", bin.CudStatus)
	r.HandleFunc("/roles", bin.CudRole)
	r.HandleFunc("/users", bin.CudUser)
	r.HandleFunc("/orders", bin.CudOrder)
	r.HandleFunc("/positions", bin.CudPosition)
	r.HandleFunc("/cards", bin.CudCard)
	//Импорт и экспорт
	r.HandleFunc("/export", bin.Export)
	r.HandleFunc("/import", bin.Import)
	//Загрузка bootstrap, модальных окон, изображений и иконок
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images", http.FileServer(http.Dir("images"))))
	//Подключение api
	bin.ConnectAPI()
	//bin.ConnectDB()
	// Метрики
	bin.CreateMetrics()
	r.Handle("/metrics", promhttp.HandlerFor(bin.Reg, promhttp.HandlerOpts{Registry: bin.Reg}))
	// Установка CSRF токена
	csrf.Secure(false)
	CSRF := csrf.Protect([]byte("my-current-secret-key-5467"), csrf.SameSite(csrf.SameSiteStrictMode))
	// Https
	//err := httpscerts.Check("cert.pem", "key.pem")
	//if err != nil {
	//	err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:8081")
	//	if err != nil {
	//		log.Fatal("Ошибка: Не можем сгенерировать https сертификат.")
	//	}
	//}
	//err = http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", CSRF(r))
	bin.SaveLog(log.Fields{"group": "server"}, log.InfoLevel, "The server is running successfully")
	err := http.ListenAndServe(":8080", CSRF(r))
	if err != nil {
		bin.SaveLog(log.Fields{
			"group": "server",
			"error": err.Error(),
		}, log.ErrorLevel, "Error at server startup")
	}
}
