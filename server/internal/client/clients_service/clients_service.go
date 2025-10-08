package clients_service

import (
	"fmt"
	"goodhumored/wmi-metrics-server/internal/client"
	"goodhumored/wmi-metrics-server/internal/client/metrics"
)

type MetricsRepository interface {
	StoreMetrics(clientID string, metrics metrics.Metrics) error
	GetLatestMetrics(clientID string) (metrics.Metrics, bool)
	GetMetricsForPeriod(clientID string, from, to int64) ([]metrics.Metrics, error)
}

type ClientsService struct {
	repo        ClientsRepository
	metricsRepo MetricsRepository
}

func New(repo ClientsRepository, metricsRepo MetricsRepository) *ClientsService {
	return &ClientsService{repo: repo, metricsRepo: metricsRepo}
}

func (s *ClientsService) HandleFirstClientMessage(handshake client.ClientHandshake) *client.Client {
	cl, exists := s.repo.GetClient(handshake.ID)
	if !exists {
		cl = client.New(handshake.ID, client.Connected, client.Uncertain, handshake.SystemInfo)
		cl.Health = client.Uncertain
		s.repo.AddClient(cl)
	} else {
		cl.Status = client.Connected
		cl.Health = client.Uncertain
		cl.Info = handshake.SystemInfo
		s.repo.UpdateClient(cl)
	}
	return cl
}

func (s *ClientsService) HandleClientMetricsMessage(cl *client.Client, metricsMessage metrics.Metrics) error {
	fmt.Println("Received metrics from client", cl.ID)
	cl.Health = client.Healthy
	if err := s.repo.UpdateClient(cl); err != nil {
		return err
	}
	return s.metricsRepo.StoreMetrics(cl.ID, metricsMessage)
}

func (s *ClientsService) GetClientLatestMetrics(cl *client.Client) (metrics.Metrics, bool) {
	return s.metricsRepo.GetLatestMetrics(cl.ID)
}

func (s *ClientsService) GetClientMetricsForPeriod(cl *client.Client, from, to int64) ([]metrics.Metrics, error) {
	return s.metricsRepo.GetMetricsForPeriod(cl.ID, from, to)
}

func (s *ClientsService) HandleClientDisconnects(cl *client.Client) error {
	cl.Status = client.Disconnected
	cl.Health = client.Uncertain
	return s.repo.UpdateClient(cl)
}
