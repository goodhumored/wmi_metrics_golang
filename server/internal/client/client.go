package client

type ConnectionStatus int

const (
	Connected    ConnectionStatus = iota // 0
	Disconnected                         // 1
)

type HealthStatus int

const (
	Healthy   HealthStatus = iota // 0
	Unhealthy                     // 1
	Uncertain                     // 2
)

// OSInfo holds information about the client's operating system
type OSInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Arch    string `json:"arch"`
}

// NetworkInfo holds information about the client's network
type NetworkInfo struct {
	IPAddress  string `json:"ip_address"`
	MACAddress string `json:"mac_address"`
	Hostname   string `json:"hostname"`
}

// SystemInfo holds detailed information about the client's system
type SystemInfo struct {
	OS      OSInfo      `json:"os"`
	Network NetworkInfo `json:"network"`
	CPU     string      `json:"cpu"`
	Memory  string      `json:"memory"`
	Disk    string      `json:"disk"`
}

// Client represents a known client with its status and health information
type Client struct {
	ID     string           `json:"id"`
	Status ConnectionStatus `json:"status"`
	Health HealthStatus     `json:"health"`
	Info   SystemInfo       `json:"info"`
}

type ClientHandshake struct {
	ID string `json:"id"`
	SystemInfo
}

type ClientMetrics struct {
	ID string `json:"id"`
	SystemInfo
}

func New(ID string, Status ConnectionStatus, Health HealthStatus, Info SystemInfo) *Client {
	return &Client{ID, Status, Health, Info}
}
