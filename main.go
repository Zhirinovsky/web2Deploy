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
	r.HandleFunc("/", pages.MainPage)                 //Основная страница
	r.HandleFunc("/login", pages.Login)               //Страница авторизации
	r.HandleFunc("/registration", pages.Registration) //Страница регистрация
	r.HandleFunc("/admin", pages.Admin)               //Панель управления
	r.HandleFunc("/profile", pages.Profile)           //Профиль
	r.HandleFunc("/cart", pages.Cart)                 //Корзина
	//Методы сессии
	r.HandleFunc("/logout", bin.Logout) //Выход из учётной записи
	//Методы главной страницы
	r.HandleFunc("/main/search", pages.SearchProduct) //Метод поиска продукта на основной странице
	r.HandleFunc("/main/category", pages.Filter)      //Метод фильтрации продуктов по категориям
	r.HandleFunc("/main/add-cart", pages.AddCart)     //Метод добавления товара в корзину
	//Методы профиля
	r.HandleFunc("/profile/personalisation", pages.ChangePersonalData) //Метод изменения персональных данных
	r.HandleFunc("/profile/change-password", pages.ChangePassword)     //Метод смены пароля
	r.HandleFunc("/profile/card-create", pages.CreateCard)             //Метод создания скидочной карты
	r.HandleFunc("/profile/card-delete", pages.DeleteCard)             //Метод удаления скидочной карты
	//Методы корзины
	r.HandleFunc("/cart/remove-cart", pages.RemoveCart)     //Метод удаления товара из корзины
	r.HandleFunc("/cart/change-cart", pages.ChangeCart)     //Метод изменения количества товара в корзине
	r.HandleFunc("/cart/change-active", pages.ChangeActive) //Метод изменения статуса товара в заказе
	r.HandleFunc("/cart/make-order", pages.MakeOrder)       //Метод оформления заказа
	//Cud действия
	r.HandleFunc("/products", bin.CudProduct)               //Методы CRUD действий над продуктами
	r.HandleFunc("/characteristics", bin.CudCharacteristic) //Методы CRUD действий над характеристиками
	r.HandleFunc("/categories", bin.CudCategory)            //Методы CRUD действий над категориями
	r.HandleFunc("/sets", bin.CudSet)                       //Методы CRUD действий над наборами
	r.HandleFunc("/images", bin.CudImages)                  //Методы CRUD действий над изображениями
	r.HandleFunc("/statuses", bin.CudStatus)                //Методы CRUD действий над статусами
	r.HandleFunc("/roles", bin.CudRole)                     //Методы CRUD действий над ролями
	r.HandleFunc("/users", bin.CudUser)                     //Методы CRUD действий над пользователями
	r.HandleFunc("/orders", bin.CudOrder)                   //Методы CRUD действий над заказами
	r.HandleFunc("/positions", bin.CudPosition)             //Методы CRUD действий над позициямми заказов
	r.HandleFunc("/cards", bin.CudCard)                     //Методы CRUD действий над скидочными картами
	//Импорт и экспорт
	r.HandleFunc("/export", bin.Export) //Метод фильтрации продуктов по категориям
	r.HandleFunc("/import", bin.Import) //Метод фильтрации продуктов по категориям
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
