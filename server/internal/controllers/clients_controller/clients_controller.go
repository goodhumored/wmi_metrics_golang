package clients_controller

import (
	"encoding/json"
	"goodhumored/wmi-metrics-server/internal/client"
	"goodhumored/wmi-metrics-server/internal/client/metrics"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ClientsService interface {
	GetClientLatestMetrics(cl *client.Client) (metrics.Metrics, bool)
	GetClientMetricsForPeriod(cl *client.Client, from, to int64) ([]metrics.Metrics, error)
}

type ClientsRepository interface {
	GetClient(id string) (*client.Client, bool)
	GetAllClients() []*client.Client
}

type ClientsController struct {
	service    ClientsService
	repository ClientsRepository
}

func New(service ClientsService, repository ClientsRepository) *ClientsController {
	return &ClientsController{
		service:    service,
		repository: repository,
	}
}

// GetAllClients возвращает список всех клиентов
func (c *ClientsController) GetAllClients(w http.ResponseWriter, r *http.Request) {
	clients := c.repository.GetAllClients()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

// GetClient возвращает информацию о конкретном клиенте
func (c *ClientsController) GetClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["id"]

	client, exists := c.repository.GetClient(clientID)
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

// GetLatestMetrics возвращает последние метрики клиента
func (c *ClientsController) GetLatestMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["id"]

	client, exists := c.repository.GetClient(clientID)
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	metrics, found := c.service.GetClientLatestMetrics(client)
	if !found {
		http.Error(w, "No metrics found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetMetricsForPeriod возвращает метрики клиента за период
func (c *ClientsController) GetMetricsForPeriod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["id"]

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		http.Error(w, "Parameters 'from' and 'to' are required", http.StatusBadRequest)
		return
	}

	from, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'from' parameter", http.StatusBadRequest)
		return
	}

	to, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'to' parameter", http.StatusBadRequest)
		return
	}

	client, exists := c.repository.GetClient(clientID)
	if !exists {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	metricsData, err := c.service.GetClientMetricsForPeriod(client, from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metricsData)
}
