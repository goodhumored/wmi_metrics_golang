package main

import (
	"encoding/json"
	"fmt"
	"goodhumored/wmi-metrics-server/internal/client"
	"goodhumored/wmi-metrics-server/internal/client/clients_repository"
	"goodhumored/wmi-metrics-server/internal/client/clients_service"
	"goodhumored/wmi-metrics-server/internal/client/metrics"
	"goodhumored/wmi-metrics-server/internal/client/metrics/metrics_repository"
	"goodhumored/wmi-metrics-server/internal/controllers/clients_controller"
	"goodhumored/wmi-metrics-server/internal/ws_server"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	// Инициализация репозиториев и сервисов
	metricsRepo := metrics_repository.NewMemoryMetricsRepository()
	clientRepo := clients_repository.New()
	clientService := clients_service.New(clientRepo, metricsRepo)

	// Инициализация контроллера
	clientsCtrl := clients_controller.New(clientService, clientRepo)

	// Настройка HTTP роутера
	router := mux.NewRouter()
	router.HandleFunc("/api/clients", clientsCtrl.GetAllClients).Methods("GET")
	router.HandleFunc("/api/clients/{id}", clientsCtrl.GetClient).Methods("GET")
	router.HandleFunc("/api/clients/{id}/metrics/latest", clientsCtrl.GetLatestMetrics).Methods("GET")
	router.HandleFunc("/api/clients/{id}/metrics", clientsCtrl.GetMetricsForPeriod).Methods("GET")

	// Настройка WebSocket сервера
	server := ws_server.New("/ws", router)
	server.AddHandler(func(conn *websocket.Conn, w http.ResponseWriter, r *http.Request) {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		var clientHandshake client.ClientHandshake
		if err := json.Unmarshal(msg, &clientHandshake); err != nil {
			fmt.Println("Unmarshal error:", err)
			return
		}
		cl := clientService.HandleFirstClientMessage(clientHandshake)
		defer clientService.HandleClientDisconnects(cl)

		fmt.Println("Client connected:", cl.ID)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}
			var metricsMessage metrics.Metrics
			if err := json.Unmarshal(msg, &metricsMessage); err != nil {
				fmt.Println("Unmarshal error:", err)
				continue
			}
			if err := clientService.HandleClientMetricsMessage(cl, metricsMessage); err != nil {
				fmt.Println("Handle metrics error:", err)
			}
		}
	})

	// Запуск сервера
	fmt.Println("Server starting on :8080")
	server.Start()
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Server error:", err)
	}
}
