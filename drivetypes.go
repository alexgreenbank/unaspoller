package main

// /proxy/drive/api/v2/systems/disk-stats
// UI has different durations it can query
// 1M: ?start=1775452607&end=1778066207&interval=21600	period=30d+6h	interval = 6h
// 1W: ?start=1777457812&end=1778066212&interval=3600	period=7d+1h	interval = 1h
// 1D: ?start=1777978915&end=1778066215&interval=900	period=24h+15m	interval = 15m
// 1H: ?start=1778061721&end=1778066221&interval=900	period=1h+15m	interval = 15m
type DriveApiV2SystemsDiskStatsSeriesDisk struct {
	SlotId         string    `json:"slotId"`
	Type           string    `json:"type"`
	Serial         string    `json:"serial"`
	Temperatures   []float64 `json:"temperatures"`
	ReadKBPS       []float64 `json:"readKBPS"`
	WriteKBPS      []float64 `json:"writeKBPS"`
	BadSectorCount []int     `json:"badSectorCount"`
}

type DriveApiV2SystemsDiskStatsSeriesCacheSlot struct {
	// TODO - NO IDEA - NEED DATA
}

type DriveApiV2SystemsDiskStatsSeries struct {
	Disks      []DriveApiV2SystemsDiskStatsSeriesDisk      `json:"disks"`
	CacheSlots []DriveApiV2SystemsDiskStatsSeriesCacheSlot `json:"cacheSlots"`
}

type DriveApiV2SystemsDiskStatsWindow struct {
	Samples  int `json:"samples"`
	Start    int `json:"start"`
	End      int `json:"end"`
	Interval int `json:"interval"`
}

type DriveApiV2SystemsDiskStats struct {
	Series DriveApiV2SystemsDiskStatsSeries `json:"series"`
	Window DriveApiV2SystemsDiskStatsWindow `json:"window"`
}

// /proxy/drive/api/v2/systems/network-io
type DriveApiV2SystemsNetworkIO struct {
	ReceiveKBPS  float64 `json:"receiveKBPS"`
	TransmitKBPS float64 `json:"transmitKBPS"`
	Timestamp    string  `json:"timestamp"`
}

// /proxy/users/drive/api/v2/drives
type DriveApiV2DrivesDriveProtection struct {
	EncryptionStatus    string `json:"encryptionStatus"`
	RemoteBackupEnabled bool   `json:"remoteBackupEnabled"`
	SnapshotEnabled     bool   `json:"snapshotEnabled"`
}

type DriveApiV2DrivesDrive struct {
	Id                       string                          `json:"id"`
	Type                     string                          `json:"type"`
	Name                     string                          `json:"name"`
	Status                   string                          `json:"status"`
	StoragePoolId            string                          `json:"storagePoolId"`
	DataSync                 string                          `json:"dataSync"`
	RecordSize               string                          `json:"recordSize"`
	CompressionLevel         string                          `json:"compressionLevel"`
	Deduplication            string                          `json:"deduplication"`
	DeduplicationEverEnabled bool                            `json:"deduplicationEverEnabled"`
	Quota                    int64                           `json:"quota"`
	Usage                    int64                           `json:"usage"`
	Role                     string                          `json:"role"`
	Protections              DriveApiV2DrivesDriveProtection `json:"protections"`
	MemberCount              int                             `json:"memberCount"`
}

type DriveApiV2Drives struct {
	Drives []DriveApiV2DrivesDrive `json:"drives"`
}

type DriveApiV2NetworkInterface struct {
	Interface     string `json:"interface"`
	InterfaceName string `json:"interfaceName"`
	Connected     bool   `json:"connected"`
	MaxSpeed      string `json:"maxSpeed"`
	LinkSpeed     string `json:"linkSpeed"`
	Address       string `json:"address,omitempty"`
	MAC           string `json:"mac,omitempty"`
}

type DriveApiV2CPU struct {
	CurrentLoad float64 `json:"currentLoad"`
	Temperature int     `json:"temperature"` // float64? - only ever seen ints
}

type DriveApiV2Memory struct {
	Free      int `json:"free"`
	Total     int `json:"total"`
	Available int `json:"available"`
}

type DriveApiV2SystemsDeviceInfo struct {
	NetworkInterfaces []DriveApiV2NetworkInterface `json:""`
	Usbs              *string                      `json:"usbs"`
	Version           string                       `json:"version"`
	Name              string                       `json:"name"`
	Model             string                       `json:"model"`
	StartupTime       string                       `json:"startupTime"`
	Memory            DriveApiV2Memory             `json:"memory"`
	CPU               DriveApiV2CPU                `json:"cpu"`
	FirmwareVersion   string                       `json:"firmwareVersion"`
	Status            string                       `json:"status"`
	SfpAggregation    bool                         `json:"sfpAggregation"`
}

type DriveApiV2IncompatibleReason struct {
	// ??? No examples to work from
}

type DriveApiV2RiskReason struct {
	// ??? No examples to work from
}

type DriveApiV2RaidGroup struct {
	Number              int    `json:"number"`
	Id                  string `json:"id"`
	RemnantReason       string `json:"remnantReason"`
	IsSSDCache          bool   `json:"isSSDCache"`
	CurrentLevel        string `json:"currentLevel"`
	ConfigLevel         string `json:"configLevel"`
	CurrentProtection   int    `json:"currentProtection"`
	ExpectedProtection  int    `json:"expectedProtection"`
	RecommendedDiskSize int64  `json:"recommendedDiskSize"`
	Progress            int    `json:"progress"`
	Estimate            int    `json:"estimate"`
}

type DriveApiV2Pool struct {
	Number             int                   `json:"number"`
	Id                 string                `json:"id"`
	PreferLevel        string                `json:"preferLevel"`
	Type               string                `json:"type"`
	Status             string                `json:"status"`
	Capacity           int64                 `json:"capacity"`
	Usage              int64                 `json:"usage"`
	activeRaidGroupId  string                `json:"activeRaidGroupId"`
	RaidGroups         []DriveApiV2RaidGroup `json:"raidGroups"`
	InitializingStatus string                `json:"initializingStatus"`
}

type DriveApiV2Disk struct {
	SlotId                   string                         `json:"slotId"`
	Location                 string                         `json:"location"`
	PoolId                   string                         `json:"poolId"`
	RaidGroupId              string                         `json:"raidGroupId"`
	MetadataGroupId          string                         `json:"metadataGroupId"`
	IsGlobalHotSpare         bool                           `json:"isGlobalHotSpare"`
	IsLocalHotSpare          bool                           `json:"isLocalHotSpare"`
	Type                     string                         `json:"type"`
	State                    string                         `json:"state"`
	RPM                      int                            `json:"rpm"`
	Model                    string                         `json:"model"`
	Size                     int64                          `json:"size"`
	Sata                     string                         `json:"sata"`
	Ata                      string                         `json:"ata"`
	NvmeVersion              string                         `json:"nvmeVersion"`
	Firmware                 string                         `json:"firmware"`
	SectorFormat             string                         `json:"sectorFormat"`
	Serial                   string                         `json:"serial"`
	Temperature              int                            `json:"temperature"` // float64? - only ever seen ints
	PowerOnHours             int                            `json:"powerOnHours"`
	BadSectorCount           int                            `json:"badSectorCount"`           // int64?
	UncorrectableSectorCount int                            `json:"uncorrectableSectorCount"` // int64?
	ReadErrorRate            int                            `json:"readErrorRate"`
	SmartReadErrorCount      int                            `json:"smartReadErrorCount"` // int64?
	RiskReasons              []DriveApiV2RiskReason         `json:"riskReasons"`
	IncompatibleReasons      []DriveApiV2IncompatibleReason `json:"incompatibleReasons"`
	ReadKBPS                 float64                        `json:"readKBPS"`
	WriteKBPS                float64                        `json:"writeKBPS"`
	SmartTestSupported       bool                           `json:"smartTestSupported"`
	HealthScore              int                            `json:"healthScore"`
}

type DriveApiV2CacheSlot struct {
}

type DriveApiV2Storage struct {
	Pools      []DriveApiV2Pool      `json:"pools"`
	Disks      []DriveApiV2Disk      `json:"disks"`
	CacheSlots []DriveApiV2CacheSlot `json:"cacheSlots"`
	Expansions *string               `json:"expansions"`
}
