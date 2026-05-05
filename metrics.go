package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	// /drive/api/v2/systems/device-info
	upGauge  *prometheus.GaugeVec
	memFree  prometheus.Gauge
	memTotal prometheus.Gauge
	memAvail prometheus.Gauge
	cpuLoad  prometheus.Gauge
	cpuTemp  prometheus.Gauge
	// /drive/api/v2/storage
	nosPools         prometheus.Gauge
	nosPoolsByStatus *prometheus.GaugeVec
	nosDisks         prometheus.Gauge
	nosDisksByState  *prometheus.GaugeVec
	// - per Pool info
	poolCapacity *prometheus.GaugeVec
	poolUsage    *prometheus.GaugeVec
	// - per Disk info
	diskRPM                      *prometheus.GaugeVec
	diskSize                     *prometheus.GaugeVec
	diskTemp                     *prometheus.GaugeVec
	diskPowerOnHours             *prometheus.GaugeVec // This is a counter as it can reset if the disk is removed/replaced
	diskBadSectorCount           *prometheus.GaugeVec // This is a counter as it can reset if the disk is removed/replaced
	diskUncorrectableSectorCount *prometheus.GaugeVec // This is a counter as it can reset if the disk is removed/replaced
	diskReadErrorRate            *prometheus.GaugeVec
	diskSmartReadErrorCount      *prometheus.GaugeVec // This is a counter as it can reset if the disk is removed/replaced
	diskReadKBPS                 *prometheus.GaugeVec
	diskWriteKBPS                *prometheus.GaugeVec
	diskHealthScore              *prometheus.GaugeVec
	// /proxy/users/drive/api/v2/drives
	nosDrives *prometheus.GaugeVec
	// per Drive info
	driveQuota *prometheus.GaugeVec
	driveUsage *prometheus.GaugeVec
}

func (u *UNAS) newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		// Device metrics
		upGauge: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "up",
				Help:      "Whether the UNAS is up or not TODO - fix me",
			},
			[]string{"name", "model", "version", "firmware_version"},
		),
		memFree: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "memory_free",
				Help:      "Amount of memory free on UNAS device in KB",
			}),
		memTotal: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "memory_total",
				Help:      "Amount of total memory on UNAS device in KB",
			}),
		memAvail: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "memory_avail",
				Help:      "Amount of memory available on UNAS device in KB",
			}),
		cpuLoad: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "cpu_load",
				Help:      "CPU load of UNAS device",
			}),
		cpuTemp: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "cpu_temperature",
				Help:      "Temperature of CPU in deg C",
			}),
		// Device summary metrics
		nosPools: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "nos_pools",
				Help:      "Number of storage pools",
			}),
		nosPoolsByStatus: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "nos_pools_by_status",                    // Don't like that Pools use "status" and Disks use "state"
				Help:      "Number of storage pools in each status", // But that's the API not a decision by me
			},
			[]string{"status"},
		),
		nosDisks: promauto.With(reg).NewGauge(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "nos_disks",
				Help:      "Number of physical disks",
			}),
		nosDisksByState: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "nos_disks_by_state",
				Help:      "Number of physical disks in each state",
			},
			[]string{"state"},
		),
		// Per pool metrics
		poolCapacity: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "storage_pool_capacity",
				Help:      "Capacity of the storage pool in bytes",
			},
			[]string{"pool_number", "pool_id"},
		),
		poolUsage: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "storage_pool_usage",
				Help:      "Usage of the storage pool in bytes",
			},
			[]string{"pool_number", "pool_id"},
		),
		// Per disk metrics
		diskRPM: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_rpm",
				Help:      "The RPM of the disk",
			},
			[]string{"slotId", "serial"},
		),
		diskSize: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_size",
				Help:      "The size of the disk in bytes",
			},
			[]string{"slotId", "serial"},
		),
		diskTemp: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_temperature",
				Help:      "The temperature of the disk in deg C",
			},
			[]string{"slotId", "serial"},
		),
		// We use a GaugeVec for these counters as we cannot be sure that the counters will not reset
		// especially if a disk is replaced (and Serials aren't unique)
		// I suppose Serial + Model is unlikely not to be unique but you never know
		// Also we can't set a counter to an arbitrary value with promauto
		diskPowerOnHours: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_power_on_hours",
				Help:      "The number of hours the disk has been powered on",
			},
			[]string{"slotId", "serial"},
		),
		diskBadSectorCount: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_bad_sector_count",
				Help:      "The count of bad sectors reported by the disk",
			},
			[]string{"slotId", "serial"},
		),
		diskUncorrectableSectorCount: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_uncorrectable_sector_count",
				Help:      "The count of uncorrectable sector errors reported by the disk",
			},
			[]string{"slotId", "serial"},
		),
		diskReadErrorRate: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_read_error_rate",
				Help:      "The read error rate of the disk in whoknowswhat",
			},
			[]string{"slotId", "serial"},
		),
		diskSmartReadErrorCount: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_read_error_count",
				Help:      "The number of read errors",
			},
			[]string{"slotId", "serial"},
		),
		diskReadKBPS: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_read_kbps",
				Help:      "The read rate of the disk in kbps",
			},
			[]string{"slotId", "serial"},
		),
		diskWriteKBPS: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_write_kbps",
				Help:      "The write rate of the disk in kbps",
			},
			[]string{"slotId", "serial"},
		),
		diskHealthScore: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "disk_health_score",
				Help:      "The health score from 0=bad to 5=good",
			},
			[]string{"slotId", "serial"},
		),
		// Drives metrics
		nosDrives: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "nos_drives",
				Help:      "The number of drives by type",
			},
			[]string{"type"},
		),
		// Per drive metrics
		driveUsage: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "drive_size",
				Help:      "The size of the drive in bytes",
			},
			[]string{"name", "poolId"},
		),
		driveQuota: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: u.c.flagMetricPrefix,
				Name:      "drive_quota",
				Help:      "The size of the drive quota in bytes",
			},
			[]string{"name", "poolId"},
		),
	}
	return m
}
