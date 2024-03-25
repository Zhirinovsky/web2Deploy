package bin

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"strconv"
	"web2/bin/models"
)

var Reg *prometheus.Registry

var (
	CurrentSessions prometheus.Gauge
	UsersAmount     prometheus.Gauge
	Products        *prometheus.CounterVec
	TotalFailures   *prometheus.CounterVec
	Failures        *prometheus.CounterVec
)

func CreateMetrics() {
	Reg = prometheus.NewRegistry()
	CurrentSessions = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "current_sessions_amount",
		Help: "Current amount of sessions.",
	})
	UsersAmount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "accounts_amount",
		Help: "Amount of existing accounts.",
	})
	Failures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors",
			Help: "Number of handled errors.",
		},
		[]string{"error"},
	)
	TotalFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fatal_errors",
			Help: "Number of fatal errors.",
		},
		[]string{"fatal_error"},
	)
	Products = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "products",
			Help: "Info about products.",
		},
		[]string{"name", "price", "amount", "category"},
	)
	Reg.MustRegister(CurrentSessions, UsersAmount, Failures, TotalFailures, Products)
	var products []models.Product
	Request("/products", "GET", ServerToken, nil, &products)
	for _, product := range products {
		Products.With(prometheus.Labels{"name": product.Name, "price": strconv.FormatFloat(product.Price, 'g', -1, 64), "amount": strconv.Itoa(product.Amount), "category": product.Category.Name}).Add(product.Price * float64(product.Amount))
	}
	SaveLog(log.Fields{
		"group": "server",
	}, log.InfoLevel, "Metrics loaded")
	UpdateMetrics()
}

func UpdateMetrics() {
	CurrentSessions.Set(float64(Client.DBSize(context.Background()).Val() - 1))
	var users []models.User
	Request("/users", "GET", ServerToken, nil, &users)
	UsersAmount.Set(float64(len(users)))
}
