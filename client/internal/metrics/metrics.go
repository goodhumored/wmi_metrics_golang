package metrics

type CPU struct {
	LoadPercentage uint32 `json:"load"`
}

type RAM struct {
	FreePhysicalMemory     uint64 `json:"free_ram"`
	TotalVisibleMemorySize uint64 `json:"total_ram"`
}

type Disk struct {
	DeviceID  string `json:"device_id"`
	FreeSpace uint64 `json:"free"`
	Size      uint64 `json:"total"`
}

type Metrics struct {
	Disks []Disk `json:"disks"`
	OS    []RAM  `json:"ram"`
	Proc  []CPU  `json:"cpu"`
}

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

type ClientHandshake struct {
	ID string `json:"id"`
	SystemInfo
}

type ClientMetrics struct {
	ID string `json:"id"`
	SystemInfo
}
