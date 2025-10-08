package wmi_client

import (
	"fmt"
	"github.com/yusufpapurcu/wmi"
	"goodhumored/wmi-metrics-client/internal/metrics"
)

type WMIClient struct {
}

func New() *WMIClient {
	return &WMIClient{}
}

func (c WMIClient) GetMetrics() (metrics.Metrics, error) {
	metrics := metrics.Metrics{}
	disks, err := c.GetDiskStatus()
	if err != nil {
		return metrics, err
	}
	os, err := c.GetRAMStatus()
	if err != nil {
		return metrics, err
	}
	cpu, err := c.GetCPUStatus()
	if err != nil {
		return metrics, err
	}
	metrics.Disks = disks
	metrics.OS = os
	metrics.Proc = cpu
	return metrics, nil
}

func (c WMIClient) GetDiskStatus() ([]metrics.Disk, error) {
	var disks []metrics.Disk
	err := wmi.Query("SELECT DeviceID, FreeSpace, Size FROM Win32_LogicalDisk WHERE DriveType=3", &disks)
	if err != nil {
		return []metrics.Disk{}, fmt.Errorf("Disk query failed: %w", err)
	}
	if len(disks) == 0 {
		return []metrics.Disk{}, fmt.Errorf("no disks found")
	}
	return disks, nil
}

func (c WMIClient) GetRAMStatus() ([]metrics.RAM, error) {
	var os []metrics.RAM
	err := wmi.Query("SELECT FreePhysicalMemory, TotalVisibleMemorySize FROM Win32_OperatingSystem", &os)
	if err != nil {
		return []metrics.RAM{}, fmt.Errorf("RAM query failed: %w", err)
	}
	if len(os) == 0 {
		return []metrics.RAM{}, fmt.Errorf("no OS metrics found")
	}
	return os, nil
}

func (c WMIClient) GetCPUStatus() ([]metrics.CPU, error) {
	var cpu []metrics.CPU
	err := wmi.Query("SELECT LoadPercentage FROM Win32_Processor WHERE DeviceID='CPU0'", &cpu)
	if err != nil {
		return []metrics.CPU{}, fmt.Errorf("CPU query failed: %w", err)
	}
	if len(cpu) == 0 {
		return []metrics.CPU{}, fmt.Errorf("no CPU metrics found")
	}
	return cpu, nil
}

// GetSystemInfo retrieves system information using WMI
func (c *WMIClient) GetSystemInfo() (metrics.SystemInfo, error) {
	var info metrics.SystemInfo
	var osInfo []Win32_OperatingSystem
	var netInfo []Win32_NetworkAdapterConfiguration
	var cpuInfo []Win32_Processor
	var diskInfo []Win32_LogicalDisk

	// Query OS info
	err := wmi.Query("SELECT Caption, Version, OSArchitecture, TotalVisibleMemorySize FROM Win32_OperatingSystem", &osInfo)
	if err != nil {
		return info, fmt.Errorf("failed to query OS info: %v", err)
	}
	if len(osInfo) > 0 {
		info.OS.Name = osInfo[0].Caption
		info.OS.Version = osInfo[0].Version
		info.OS.Arch = osInfo[0].OSArchitecture
		memoryMB := osInfo[0].TotalVisibleMemorySize / 1024 // Convert KB to GB
		info.Memory = fmt.Sprintf("%.2f GB", float64(memoryMB)/1024)
	}

	// Query CPU info
	err = wmi.Query("SELECT Name FROM Win32_Processor", &cpuInfo)
	if err != nil {
		return info, fmt.Errorf("failed to query CPU info: %v", err)
	}
	if len(cpuInfo) > 0 {
		info.CPU = cpuInfo[0].Name
	}

	// Query Disk info (Total size of C: drive in GB)
	err = wmi.Query("SELECT Size FROM Win32_LogicalDisk WHERE DeviceID='C:'", &diskInfo)
	if err != nil {
		return info, fmt.Errorf("failed to query disk info: %v", err)
	}
	if len(diskInfo) > 0 {
		diskGB := diskInfo[0].Size / (1024 * 1024 * 1024) // Convert bytes to GB
		info.Disk = fmt.Sprintf("%.2f GB", float64(diskGB))
	}

	// Query Network info
	err = wmi.Query("SELECT IPAddress, MACAddress, DNSHostName, IPEnabled FROM Win32_NetworkAdapterConfiguration WHERE IPEnabled=TRUE", &netInfo)
	if err != nil {
		return info, fmt.Errorf("failed to query network info: %v", err)
	}
	if len(netInfo) > 0 {
		if len(netInfo[0].IPAddress) > 0 {
			info.Network.IPAddress = netInfo[0].IPAddress[0] // Take first IP
		}
		info.Network.MACAddress = netInfo[0].MACAddress
		info.Network.Hostname = netInfo[0].DNSHostName
	}

	return info, nil
}

// Win32_OperatingSystem represents the WMI class for OS information
type Win32_OperatingSystem struct {
	Caption                string
	Version                string
	OSArchitecture         string
	TotalVisibleMemorySize uint64
}

// Win32_NetworkAdapterConfiguration represents the WMI class for network information
type Win32_NetworkAdapterConfiguration struct {
	IPAddress   []string
	MACAddress  string
	DNSHostName string
	IPEnabled   bool
}

// Win32_Processor represents the WMI class for CPU information
type Win32_Processor struct {
	Name string
}

// Win32_LogicalDisk represents the WMI class for disk information
type Win32_LogicalDisk struct {
	Size     uint64
	DeviceID string
}
