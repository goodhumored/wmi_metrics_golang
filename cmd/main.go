package main

import (
	"context"
	"fmt"
	"goodhumored/wmi-metrics-client/internal/config"
	"goodhumored/wmi-metrics-client/internal/metrics_service"
	"goodhumored/wmi-metrics-client/internal/wmi_client"
	"goodhumored/wmi-metrics-client/internal/ws_client"
	"os"
	"sync"
)

func main() {
	config := config.GetConfig()
	wmiClient := wmi_client.New()
	wsClient := ws_client.New(config.ServerUrl)
	metricsService := metrics_service.New(wmiClient, wsClient, config)

	err := wsClient.Connect()
	if err != nil {
		fmt.Printf("Failed connecting: %v\n", err)
		panic(err)
	}
	defer wsClient.Close()

	var wg sync.WaitGroup

	sigChan := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigChan
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsService.StartSendMetricsLoop(ctx)
	}()
	wg.Wait()
	fmt.Println("Normal shutdown")
}
