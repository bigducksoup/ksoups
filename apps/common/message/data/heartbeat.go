package data

// PerformanceMetrics is a struct that contains the performance metrics of the system
type PerformanceMetrics struct {
	// CPU Metrics
	CPUUserUtilization   float64 `json:"cpu_user_utilization"`
	CPUSystemUtilization float64 `json:"cpu_system_utilization"`
	CPUIdleUtilization   float64 `json:"cpu_idle_utilization"`
	CPUWaitUtilization   float64 `json:"cpu_wait_utilization"`

	// Memory Metrics
	MemoryAvailable uint64  `json:"memory_available"`
	MemoryCached    uint64  `json:"memory_cached"`
	SwapUtilization float64 `json:"swap_utilization"`

	// Disk Metrics
	DiskReadIOPS         uint64  `json:"disk_read_iops"`
	DiskWriteIOPS        uint64  `json:"disk_write_iops"`
	DiskReadThroughput   uint64  `json:"disk_read_throughput"`  // Bytes/s
	DiskWriteThroughput  uint64  `json:"disk_write_throughput"` // Bytes/s
	DiskUtilization      float64 `json:"disk_utilization"`
	DiskSpaceUtilization float64 `json:"disk_space_utilization"`

	// Network Metrics
	NetworkReceiveBPS  uint64  `json:"network_receive_bps"`
	NetworkTransmitBPS uint64  `json:"network_transmit_bps"`
	NetworkErrorRate   float64 `json:"network_error_rate"`
	NetworkDropRate    float64 `json:"network_drop_rate"`

	// Process Metrics
	ProcessCount          uint64  `json:"process_count"`
	ZombieProcessCount    uint64  `json:"zombie_process_count"`
	ProcessCPUUtilization float64 `json:"process_cpu_utilization"`
	ProcessMemUtilization float64 `json:"process_mem_utilization"`

	// System Load
	SystemLoad1  float64 `json:"system_load_1"`
	SystemLoad5  float64 `json:"system_load_5"`
	SystemLoad15 float64 `json:"system_load_15"`

	// I/O Metrics
	IOWaitTime    uint64  `json:"io_wait_time"`
	IOUtilization float64 `json:"io_utilization"`

	// Filesystem Metrics
	FSInodesFree       uint64  `json:"fs_inodes_free"`
	FSSpaceUtilization float64 `json:"fs_space_utilization"`

	// System Parameters
	ContextSwitchRate uint64 `json:"context_switch_rate"`
	InterruptCount    uint64 `json:"interrupt_count"`

	// Errors and Exceptions
	KernelErrors    uint64 `json:"kernel_errors"`
	OOMErrors       uint64 `json:"oom_errors"`
	InterruptErrors uint64 `json:"interrupt_errors"`
}
