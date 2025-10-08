package metrics_service

import (
	"context"
	"errors"
	"fmt"
	"goodhumored/wmi-metrics-client/internal"
	"goodhumored/wmi-metrics-client/internal/config"
	"goodhumored/wmi-metrics-client/internal/metrics"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type WMIClient interface {
	GetMetrics() (metrics.Metrics, error)
	GetSystemInfo() (metrics.SystemInfo, error)
}

type WSClient interface {
	Send(data any) error
}

type MetricsService struct {
	WmiClient WMIClient
	WsClient  WSClient
	cfg       config.Config
}

func New(WmiClient WMIClient, WsClient WSClient, cfg config.Config) *MetricsService {
	return &MetricsService{WmiClient, WsClient, cfg}
}

func (s *MetricsService) StartSendMetricsLoop(ctx context.Context) {
	err := s.SendHandshakeMessage()
	if err != nil {
		fmt.Println("Failure", err)
		return
	}
	metricsChan := make(chan metrics.Metrics, 10)

	// sending metrics
	go func() {
		for metrics := range metricsChan {
			err := s.WsClient.Send(metrics)
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure, websocket.CloseGoingAway) || errors.Is(err, websocket.ErrCloseSent) {
				fmt.Println("Server closed connection")
				return
			}
			if err != nil {
				fmt.Printf("Failed sending: %+V\n", err)
			}
		}
	}()

	// getting metrics
	for {
		select {
		case <-ctx.Done():
			close(metricsChan)
			return
		default:
			metrics, err := s.WmiClient.GetMetrics()
			if err != nil {
				fmt.Println("Error fetching metrics: ", err)
				return
			}
			fmt.Printf("Got metrics: %v\n", metrics)
			metricsChan <- metrics
			time.Sleep(time.Duration(s.cfg.MetricsReadPeriod) * time.Millisecond)
		}
	}
}

func (s *MetricsService) SendHandshakeMessage() error {
	systemInfo, err := s.WmiClient.GetSystemInfo()
	if err != nil {
		return fmt.Errorf("failed getting system info: %w", err)
	}

	hostname, _ := os.Hostname()
	handshake := metrics.ClientHandshake{ID: fmt.Sprintf("%s-%s", hostname, internal.GetMacAddress()), SystemInfo: systemInfo}

	err = s.WsClient.Send(handshake)
	if err != nil {
		return fmt.Errorf("failed sending handshake message: %w", err)
	}
	return nil
}
