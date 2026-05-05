# Metrics

Here is a summary of the metrics produced by `unaspoller`. They fall into 5 categories:
* Device metrics (cpu, temperature, memory, network traffic in/out)
* Summary metrics (counts of disks, drives, pools, etc)
* Per Pool metrics
* Per Disk metrics
* Per Drive metrics

A few definitions may help:
* Disk: Inside a UNAS you have a collection of physical disks (HDDs/SSDs)
* Raid Groups: One or more physical disks belong to a Raid Group
* Storage Pool: A storage pool is made up of one or more Raid Groups (but usually just one)
* Drive: Think of `Drive` as in `Network Drive` or `Network Share`. It has an optional quota. 

## Metrics produced by unaspoller

Device metrics:
```
# HELP unas_up Whether the UNAS is up or not TODO - fix me
# TYPE unas_up gauge
unas_up{firmware_version="5.0.17",model="UNASPRO",name="unas01",version="4.1.16"} 1
# HELP unas_uptime Uptime in seconds of the device
# TYPE unas_uptime gauge
unas_uptime 920112.163195206
# HELP unas_cpu_load CPU load of UNAS device
# TYPE unas_cpu_load gauge
unas_cpu_load 0.055
# HELP unas_cpu_temperature Temperature of CPU in deg C
# TYPE unas_cpu_temperature gauge
unas_cpu_temperature 63
# HELP unas_memory_avail Amount of memory available on UNAS device in KB
# TYPE unas_memory_avail gauge
unas_memory_avail 5.319808e+06
# HELP unas_memory_free Amount of memory free on UNAS device in KB
# TYPE unas_memory_free gauge
unas_memory_free 4.8592e+06
# HELP unas_memory_total Amount of total memory on UNAS device in KB
# TYPE unas_memory_total gauge
unas_memory_total 8.272256e+06
# HELP unas_network_receive_kbps Network traffic received in KBps
# TYPE unas_network_receive_kbps gauge
unas_network_receive_kbps 1.6277525717426689
# HELP unas_network_transmit_kbps Network traffic transmitted in KBps
# TYPE unas_network_transmit_kbps gauge
unas_network_transmit_kbps 34.81847702651475
```

Summary metrics:
```
# HELP unas_nos_disks Number of physical disks
# TYPE unas_nos_disks gauge
unas_nos_disks 7
# HELP unas_nos_disks_by_state Number of physical disks in each state
# TYPE unas_nos_disks_by_state gauge
unas_nos_disks_by_state{state="optimal"} 7
# HELP unas_nos_drives The number of drives by type
# TYPE unas_nos_drives gauge
unas_nos_drives{type="shared"} 4
# HELP unas_nos_pools Number of storage pools
# TYPE unas_nos_pools gauge
unas_nos_pools 2
# HELP unas_nos_pools_by_status Number of storage pools in each status
# TYPE unas_nos_pools_by_status gauge
unas_nos_pools_by_status{status="fullyOperational"} 2
```

Pool specific metrics (example of two different pools here):
```
# HELP unas_storage_pool_capacity Capacity of the storage pool in bytes
# TYPE unas_storage_pool_capacity gauge
unas_storage_pool_capacity{pool_id="aa460908-1e83-4acb-ab65-436913517d61",pool_number="1"} 2.974424236032e+12
unas_storage_pool_capacity{pool_id="d6212f95-4ea7-4874-89f6-5a1ce1e292fc",pool_number="2"} 3.99205466112e+12
# HELP unas_storage_pool_usage Usage of the storage pool in bytes
# TYPE unas_storage_pool_usage gauge
unas_storage_pool_usage{pool_id="aa460908-1e83-4acb-ab65-436913517d61",pool_number="1"} 2.41736548352e+11
unas_storage_pool_usage{pool_id="d6212f95-4ea7-4874-89f6-5a1ce1e292fc",pool_number="2"} 3.388354789376e+12
```

Per disk metrics (example of only two disks here):
```
# HELP unas_disk_bad_sector_count The count of bad sectors reported by the disk
# TYPE unas_disk_bad_sector_count gauge
unas_disk_bad_sector_count{serial="N8GUTKVY",slotId="2"} 0
unas_disk_bad_sector_count{serial="Z9CBFX9K",slotId="1"} 0
# HELP unas_disk_health_score The health score from 0=bad to 5=good
# TYPE unas_disk_health_score gauge
unas_disk_health_score{serial="N8GUTKVY",slotId="2"} 5
unas_disk_health_score{serial="Z9CBFX9K",slotId="1"} 5
# HELP unas_disk_power_on_hours The number of hours the disk has been powered on
# TYPE unas_disk_power_on_hours gauge
unas_disk_power_on_hours{serial="N8GUTKVY",slotId="2"} 42687
unas_disk_power_on_hours{serial="Z9CBFX9K",slotId="1"} 260
# HELP unas_disk_read_error_count The number of read errors
# TYPE unas_disk_read_error_count gauge
unas_disk_read_error_count{serial="N8GUTKVY",slotId="2"} 0
unas_disk_read_error_count{serial="Z9CBFX9K",slotId="1"} 0
# HELP unas_disk_read_error_rate The read error rate of the disk in whoknowswhat
# TYPE unas_disk_read_error_rate gauge
unas_disk_read_error_rate{serial="N8GUTKVY",slotId="2"} 0
unas_disk_read_error_rate{serial="Z9CBFX9K",slotId="1"} 4.24093e+06
# HELP unas_disk_read_kbps The read rate of the disk in kbps
# TYPE unas_disk_read_kbps gauge
unas_disk_read_kbps{serial="N8GUTKVY",slotId="2"} 0
unas_disk_read_kbps{serial="Z9CBFX9K",slotId="1"} 0
# HELP unas_disk_rpm The RPM of the disk
# TYPE unas_disk_rpm gauge
unas_disk_rpm{serial="N8GUTKVY",slotId="2"} 7200
unas_disk_rpm{serial="Z9CBFX9K",slotId="1"} 5900
# HELP unas_disk_size The size of the disk in bytes
# TYPE unas_disk_size gauge
unas_disk_size{serial="N8GUTKVY",slotId="2"} 4.000787030016e+12
unas_disk_size{serial="Z9CBFX9K",slotId="1"} 1.000204886016e+12
# HELP unas_disk_temperature The temperature of the disk in deg C
# TYPE unas_disk_temperature gauge
unas_disk_temperature{serial="N8GUTKVY",slotId="2"} 48
unas_disk_temperature{serial="Z9CBFX9K",slotId="1"} 40
# HELP unas_disk_uncorrectable_sector_count The count of uncorrectable sector errors reported by the disk
# TYPE unas_disk_uncorrectable_sector_count gauge
unas_disk_uncorrectable_sector_count{serial="N8GUTKVY",slotId="2"} 0
unas_disk_uncorrectable_sector_count{serial="Z9CBFX9K",slotId="1"} 0
# HELP unas_disk_write_kbps The write rate of the disk in kbps
# TYPE unas_disk_write_kbps gauge
unas_disk_write_kbps{serial="N8GUTKVY",slotId="2"} 0
unas_disk_write_kbps{serial="Z9CBFX9K",slotId="1"} 0
```

Per drive metrics:
```
# HELP unas_drive_quota The size of the drive quota in bytes
# TYPE unas_drive_quota gauge
unas_drive_quota{name="misc",poolId="aa460908-1e83-4acb-ab65-436913517d61"} -1
unas_drive_quota{name="vm_backups",poolId="aa460908-1e83-4acb-ab65-436913517d61"} -1
# HELP unas_drive_size The size of the drive in bytes
# TYPE unas_drive_size gauge
unas_drive_size{name="misc",poolId="aa460908-1e83-4acb-ab65-436913517d61"} 2.30161252352e+11
unas_drive_size{name="vm_backups",poolId="aa460908-1e83-4acb-ab65-436913517d61"} 9.411362816e+09
```

## Drive API details

Most of the Drive API details can be inferred from the [driveapi.go](driveapi.go) and [drivetypes.go](drivetypes.go) files.

The structs in [drivetypes.go](drivetypes.go) match the JSON structures to allow simple parsing. However we all know that APIs change, JSON can be a bit looser (in terms of conforming to standards) than expected, and things naturally drift over time.

## Validation

As the Drive API is undocumented I'm trying to be as strict as possible in verifying that I'm pulling the correct information from it.A

TODO - validation strategy, possible future changes, etc
