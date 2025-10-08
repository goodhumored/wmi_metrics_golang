package metrics

type CPUMetrics struct {
	LoadPercentage uint32 `json:"load"`
}

type RAMMetrics struct {
	FreePhysicalMemory     uint64 `json:"free_ram"`
	TotalVisibleMemorySize uint64 `json:"total_ram"`
}

type DiskMetrics struct {
	DeviceID  string `json:"device_id"`
	FreeSpace uint64 `json:"free"`
	Size      uint64 `json:"total"`
}

type Metrics struct {
	Disks     []DiskMetrics `json:"disks"`
	Ram       []RAMMetrics  `json:"os"`
	CPU       []CPUMetrics  `json:"cpu"`
	Timestamp int64         `json:"timestamp"`
}
