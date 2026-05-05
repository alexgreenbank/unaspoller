package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// OK, each API endpoint has a function to unmarshal it and to extract and populate the metrics
// Optionally we also have functions to validate:
// a) the unmarshalled object via validateStrict()
// b) a generic JSON object if the unmarshal operation fails via validateLoose()
//
// I'd also expect at some point that the API changes will mean that we'll need to parse the
// metric data out the generic JSON object rather than the stricter types. But not yet.
type UNASDriveAPIDef struct {
	url            string
	unmarshal      func(*UNAS, []byte) (error, any)
	metrics        func(*UNAS, any) error
	validateStrict func(*UNAS, any) error
	validateLoose  func(*UNAS, any) error // No need (YET) to implement this - but we will see...
}

func (u *UNAS) registerAPIDef(url string, apidef *UNASDriveAPIDef) error {
	if _, exists := u.apidefs[url]; exists {
		return fmt.Errorf("attempting to re-register APIDef with url=[%s]", url)
	}
	u.c.log.Debugf("registering apidef [%s]\n", url)
	u.apidefs[url] = apidef
	return nil
}

func (u *UNAS) registerAPIDefs() error {
	var apidefs = []UNASDriveAPIDef{
		{"/proxy/drive/api/v2/storage", (*UNAS).driveAPIV2StorageUnmarshal, (*UNAS).driveAPIV2StorageMetrics, (*UNAS).driveAPIV2StorageValidateStrict, nil},
		{"/proxy/drive/api/v2/systems/device-info", (*UNAS).driveAPIV2SystemsDeviceInfoUnmarshal, (*UNAS).driveAPIV2SystemsDeviceInfoMetrics, (*UNAS).driveAPIV2SystemsDeviceInfoValidateStrict, nil},
		{"/proxy/users/drive/api/v2/drives", (*UNAS).driveAPIV2DrivesUnmarshal, (*UNAS).driveAPIV2DrivesMetrics, (*UNAS).driveAPIV2DrivesValidateStrict, nil},
		{"/proxy/drive/api/v2/systems/network-io", (*UNAS).driveAPIV2SystemsNetworkIOUnmarshal, (*UNAS).driveAPIV2SystemsNetworkIOMetrics, (*UNAS).driveAPIV2SystemsNetworkIOValidateStrict, nil},
	}
	for _, a := range apidefs {
		err := u.registerAPIDef(a.url, &a)
		if err != nil {
			return err
		}
	}
	return nil
}

// Test to see if a string value is in a list of expected values
func (u *UNAS) expectString(ok *bool, val string, expected []string, where string) {
	for _, v := range expected {
		if val == v {
			return
		}
	}
	u.c.log.Errorf("validation error: unseen value of [%s] for [%s]", val, where)
	*ok = false
}

// Test to see if an int value is in a list of expected values
func (u *UNAS) expectInt(ok *bool, val int, expected []int, where string) {
	for _, v := range expected {
		if val == v {
			return
		}
	}
	u.c.log.Errorf("validation error: unseen value of [%d] for [%s]", val, where)
	*ok = false
}

// Test to see if an int64 value is in a list of expected values
func (u *UNAS) expectInt64(ok *bool, val int64, expected []int64, where string) {
	for _, v := range expected {
		if val == v {
			return
		}
	}
	u.c.log.Errorf("validation error: unseen value of [%d] for [%s]", val, where)
	*ok = false
}

// Test to see if an int value is within a range
func (u *UNAS) expectIntRange(ok *bool, val int, rangeMin, rangeMax int, where string) {
	if val >= rangeMin && val <= rangeMax {
		return
	}
	u.c.log.Errorf("validation error: value [%d] outside range [%d,%d] for [%s]", val, rangeMin, rangeMax, where)
	*ok = false
}

// Test to see if an int64 value is within a range
func (u *UNAS) expectInt64Range(ok *bool, val int64, rangeMin, rangeMax int64, where string) {
	if val >= rangeMin && val <= rangeMax {
		return
	}
	u.c.log.Errorf("validation error: value [%d] outside range [%d,%d] for [%s]", val, rangeMin, rangeMax, where)
	*ok = false
}

// Test to see if a float64 value is within a range
// TODO - probably needs an epsilon tolerance test
func (u *UNAS) expectFloat64Range(ok *bool, val float64, rangeMin, rangeMax float64, where string) {
	if val >= rangeMin && val <= rangeMax {
		return
	}
	u.c.log.Errorf("validation error: value [%f] outside range [%f,%f] for [%s]", val, rangeMin, rangeMax, where)
	*ok = false
}

///////////////////////////////////////////////////////////////////////////////

func (u *UNAS) doDriveAPIDef(apiPath string) error {
	if _, exists := u.apidefs[apiPath]; !exists {
		return fmt.Errorf("cannot find APIDef with url=[%s]", apiPath)
	}
	apidef := u.apidefs[apiPath]

	// Perform the request to get the body data
	body, err := u.doGetRequest(apiPath)
	if err != nil {
		return fmt.Errorf("api %s failed: %s", apiPath, err)
	}

	// Log the raw JSON if that option is set
	if u.c.optLogRawJSON {
		u.c.log.Debugf("RAWJSON:%s:[%s]", apiPath, string(body))
	}

	if apidef.unmarshal == nil {
		u.c.log.Infof("no unmarshal function for %s - nothing we can do", apiPath)
		return nil
	}

	// Attempt to unmarshal the object
	u.c.log.Debugf("unmarshalling %s", apiPath)
	err, obj := apidef.unmarshal(u, body)
	if err != nil {
		u.c.log.Errorf("error unmarshalling %s: %s", apiPath, err)
		// log the json
		u.c.log.Debug("failed to unmarshal data for url=[%s] body=[%s]", apiPath, body)
	} else {
		u.c.log.Debugf("unmarshalled %s ok", apiPath)
	}
	// Try to validate it
	if obj != nil {
		// We managed to unmarshal the JSON so validate it strictly if we have a func
		if apidef.validateStrict != nil {
			if err = apidef.validateStrict(u, obj); err != nil {
				u.c.log.Errorf("error during strict validation %s: %s", apiPath, err)
			}
		}
	} else {
		// We didn't manage to unmarshal the JSON so validate it as best we can if we have a loose func
		if apidef.validateLoose != nil {
			var anyObj any
			err := json.Unmarshal(body, &anyObj)
			if err != nil {
				u.c.log.Errorf("error unmarshalling %s to generic object: %s", apiPath, err)
			} else {
				if err = apidef.validateLoose(u, anyObj); err != nil {
					u.c.log.Errorf("error during loose validation %s: %s", apiPath, err)
				}
			}
		}
	}
	// Extract metrics if it exists
	if obj != nil {
		if apidef.metrics != nil {
			if err := apidef.metrics(u, obj); err != nil {
				u.c.log.Errorf("error extracting metrics for %s: %s", apiPath, err)
			}
		} else {
			u.c.log.Debugf("metrics extraction %s ok", apiPath)
		}
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// API Specific parts below here

// //////////////////////////////////////////////////////////////////////////////
// /proxy/drive/api/v2/systems/network-io
func (u *UNAS) driveAPIV2SystemsNetworkIOUnmarshal(body []byte) (error, any) {
	var foo DriveApiV2SystemsNetworkIO
	err := json.Unmarshal(body, &foo)
	if err != nil {
		u.c.log.Errorf("failed to Unmarshal DriveApiV2SystemsNetworkIO: %w", err)
		return err, nil
	}
	return nil, foo
}

func (u *UNAS) driveAPIV2SystemsNetworkIOMetrics(obj any) error {
	var foo DriveApiV2SystemsNetworkIO

	foo = obj.(DriveApiV2SystemsNetworkIO)

	u.m.netReceiveKBPS.Set(foo.ReceiveKBPS)
	u.m.netTransmitKBPS.Set(foo.TransmitKBPS)

	return nil
}

func (u *UNAS) driveAPIV2SystemsNetworkIOValidateStrict(obj any) error {
	// var foo DriveApiV2SystemsNetworkIO
	ok := true

	// foo = obj.(DriveApiV2SystemsNetworkIO)

	if !ok {
		return fmt.Errorf("errors during strict validation of /proxy/users/drive/api/v2/drives")
	}
	return nil
}

// //////////////////////////////////////////////////////////////////////////////
// /proxy/users/drive/api/v2/drives
func (u *UNAS) driveAPIV2DrivesUnmarshal(body []byte) (error, any) {
	var foo DriveApiV2Drives
	err := json.Unmarshal(body, &foo)
	if err != nil {
		u.c.log.Errorf("failed to Unmarshal DriveApiV2Drives: %w", err)
		return err, nil
	}
	return nil, foo
}

func (u *UNAS) driveAPIV2DrivesMetrics(obj any) error {
	var foo DriveApiV2Drives

	foo = obj.(DriveApiV2Drives)

	driveTypeCount := make(map[string]int)
	for _, ddata := range foo.Drives {
		u.m.driveUsage.WithLabelValues(ddata.Name, ddata.StoragePoolId).Set(float64(ddata.Usage))
		u.m.driveQuota.WithLabelValues(ddata.Name, ddata.StoragePoolId).Set(float64(ddata.Quota))
		driveTypeCount[ddata.Type]++
	}

	for k, v := range driveTypeCount {
		u.m.nosDrives.WithLabelValues(k).Set(float64(v))
	}

	return nil
}

func (u *UNAS) driveAPIV2DrivesValidateStrict(obj any) error {
	var foo DriveApiV2Drives
	ok := true

	foo = obj.(DriveApiV2Drives)

	for _, ddata := range foo.Drives {
		u.expectString(&ok, ddata.Type, []string{"shared"}, "DriveApiV2DrivesDrive.Type")
		u.expectString(&ok, ddata.Status, []string{"active"}, "DriveApiV2DrivesDrive.Status")
		u.expectString(&ok, ddata.DataSync, []string{""}, "DriveApiV2DrivesDrive.DataSync")
		u.expectString(&ok, ddata.RecordSize, []string{""}, "DriveApiV2DrivesDrive.RecordSize")
		u.expectString(&ok, ddata.CompressionLevel, []string{""}, "DriveApiV2DrivesDrive.CompressionLevel")
		u.expectString(&ok, ddata.Deduplication, []string{""}, "DriveApiV2DrivesDrive.Deduplication")
		u.expectString(&ok, ddata.Role, []string{"admin"}, "DriveApiV2DrivesDrive.Role")
		u.expectString(&ok, ddata.Protections.EncryptionStatus, []string{"unencrypted"}, "DriveApiV2DrivesDriveProtecitons.EncryptionStatus")
	}

	if !ok {
		return fmt.Errorf("errors during strict validation of /proxy/users/drive/api/v2/drives")
	}
	return nil
}

// //////////////////////////////////////////////////////////////////////////////
// /proxy/drive/api/v2/systems/device-info
func (u *UNAS) driveAPIV2SystemsDeviceInfoUnmarshal(body []byte) (error, any) {
	var foo DriveApiV2SystemsDeviceInfo
	err := json.Unmarshal(body, &foo)
	if err != nil {
		u.c.log.Errorf("failed to Unmarshal DriveApiV2SystemsDeviceInfo: %w", err)
		return err, nil
	}
	return nil, foo
}

func (u *UNAS) driveAPIV2SystemsDeviceInfoMetrics(obj any) error {
	var foo DriveApiV2SystemsDeviceInfo

	foo = obj.(DriveApiV2SystemsDeviceInfo)

	u.m.upGauge.WithLabelValues(foo.Name, foo.Model, foo.Version, foo.FirmwareVersion).Set(1.0)

	// Calculate uptime from foo.StartupTime (e.g. "2026-04-24T22:20:27Z")
	startupTime, err := time.Parse("2006-01-02T15:04:05Z", foo.StartupTime)
	if err == nil {
		uptime := time.Since(startupTime).Seconds()
		u.m.uptime.Set(uptime)
	} else {
		u.c.log.Debugf("failed to parse DriveApiV2SystemsDeviceInfo.StartupTime value=[%s]: %w", foo.StartupTime, err)
	}

	// TODO - foo.Status? ("STATE_RUNNING")

	u.m.memFree.Set(float64(foo.Memory.Free))
	u.m.memTotal.Set(float64(foo.Memory.Total))
	u.m.memAvail.Set(float64(foo.Memory.Available))
	u.m.cpuLoad.Set(foo.CPU.CurrentLoad)
	u.m.cpuTemp.Set(float64(foo.CPU.Temperature))

	return nil
}

func (u *UNAS) driveAPIV2SystemsDeviceInfoValidateStrict(obj any) error {
	var foo DriveApiV2SystemsDeviceInfo
	ok := true

	foo = obj.(DriveApiV2SystemsDeviceInfo)

	for _, ni := range foo.NetworkInterfaces {
		u.expectString(&ok, ni.Interface, []string{"ethernet", "sfp+"}, "DriveApiV2SystemsDeviceInfo.Status")
		u.expectString(&ok, ni.MaxSpeed, []string{"GbE", "10 GbE"}, "DriveApiV2SystemsDeviceInfo.MaxSpeed")
		u.expectString(&ok, ni.LinkSpeed, []string{"GbE", ""}, "DriveApiV2SystemsDeviceInfo.LinkSpeed")
	}

	u.expectString(&ok, foo.Model, []string{"UNASPRO"}, "DriveApiV2SystemsDeviceInfo.Model")
	u.expectString(&ok, foo.Model, []string{"UNASPRO"}, "DriveApiV2SystemsDeviceInfo.Model")
	u.expectIntRange(&ok, foo.Memory.Total, 8000000, 8500000, "DriveApiV2SystemsDeviceInfo.Memory.Total")
	u.expectFloat64Range(&ok, foo.CPU.CurrentLoad, 0.0, 2.0, "DriveApiV2SystemsDeviceInfo.CPU.CurrentLoad")
	u.expectIntRange(&ok, foo.CPU.Temperature, 30, 80, "DriveApiV2SystemsDeviceInfo.CPU.Temperature")

	if !ok {
		return fmt.Errorf("errors during strict validation of /proxy/drive/api/v2/systems/device-info")
	}
	return nil
}

// //////////////////////////////////////////////////////////////////////////////
// /proxy/drive/api/v2/storage
func (u *UNAS) driveAPIV2StorageUnmarshal(body []byte) (error, any) {
	var foo DriveApiV2Storage
	err := json.Unmarshal(body, &foo)
	if err != nil {
		u.c.log.Errorf("failed to Unmarshal DriveApiV2Storage: %w", err)
		return err, nil
	}
	return nil, foo
}

func (u *UNAS) driveAPIV2StorageMetrics(obj any) error {
	var foo DriveApiV2Storage

	foo = obj.(DriveApiV2Storage)
	// Summary
	u.m.nosPools.Set(float64(len(foo.Pools)))
	u.m.nosDisks.Set(float64(len(foo.Disks)))
	// pools data
	poolStatusCount := make(map[string]int)
	for _, pdata := range foo.Pools {
		pNumber := fmt.Sprintf("%d", pdata.Number)
		u.m.poolCapacity.WithLabelValues(pNumber, pdata.Id).Set(float64(pdata.Capacity))
		u.m.poolUsage.WithLabelValues(pNumber, pdata.Id).Set(float64(pdata.Usage))
		poolStatusCount[pdata.Status]++
	}
	// Add a metric for each unique pool Status
	for pscK, pscV := range poolStatusCount {
		u.m.nosPoolsByStatus.WithLabelValues(pscK).Set(float64(pscV))
	}

	diskStateCount := make(map[string]int)
	// disks data
	for _, ddata := range foo.Disks {
		u.m.diskRPM.WithLabelValues(ddata.SlotId, ddata.Serial).Set(float64(ddata.RPM))
		u.m.diskSize.WithLabelValues(ddata.SlotId, ddata.Serial).Set(float64(ddata.Size))
		u.m.diskTemp.WithLabelValues(ddata.SlotId, ddata.Serial).Set(float64(ddata.Temperature))
		u.m.diskPowerOnHours.WithLabelValues(ddata.SlotId, ddata.Serial).Set(float64(ddata.PowerOnHours))
		u.m.diskBadSectorCount.WithLabelValues(ddata.SlotId, ddata.Serial).Set(float64(ddata.BadSectorCount))
		u.m.diskUncorrectableSectorCount.WithLabelValues(ddata.SlotId, ddata.Serial).Set(float64(ddata.UncorrectableSectorCount))
		u.m.diskReadErrorRate.WithLabelValues(ddata.SlotId, ddata.Serial).Set(float64(ddata.ReadErrorRate))
		u.m.diskSmartReadErrorCount.WithLabelValues(ddata.SlotId, ddata.Serial).Set(float64(ddata.SmartReadErrorCount))
		u.m.diskReadKBPS.WithLabelValues(ddata.SlotId, ddata.Serial).Set(ddata.ReadKBPS)
		u.m.diskWriteKBPS.WithLabelValues(ddata.SlotId, ddata.Serial).Set(ddata.WriteKBPS)
		u.m.diskHealthScore.WithLabelValues(ddata.SlotId, ddata.Serial).Set(float64(ddata.HealthScore))
		diskStateCount[ddata.State]++
	}
	// Add a metric for each unique disk State
	for dscK, dscV := range diskStateCount {
		u.m.nosDisksByState.WithLabelValues(dscK).Set(float64(dscV))
	}

	return nil
}

func (u *UNAS) driveAPIV2StorageValidateStrict(obj any) error {
	var foo DriveApiV2Storage
	ok := true

	foo = obj.(DriveApiV2Storage)

	// Check pool data
	for _, pdata := range foo.Pools {
		u.expectString(&ok, pdata.PreferLevel, []string{"raid5"}, "DriveApiV2Storage.Pools.PreferLevel")
		u.expectString(&ok, pdata.Type, []string{"lvm"}, "DriveApiV2Storage.Pools.Type")
		u.expectString(&ok, pdata.Status, []string{"fullyOperational", "repairing", "noDataProtectionYet"}, "DriveApiV2Storage.Pools.Status")
		for _, rgdata := range pdata.RaidGroups {
			u.expectString(&ok, rgdata.RemnantReason, []string{""}, "DriveApiV2Storage.Pools.RaidGroups.RemnantReason")
			u.expectString(&ok, rgdata.CurrentLevel, []string{"raid5", "raid1"}, "DriveApiV2Storage.Pools.RaidGroups.CurrentLevel")
			u.expectString(&ok, rgdata.ConfigLevel, []string{"raid5"}, "DriveApiV2Storage.Pools.RaidGroups.ConfigLevel")
			u.expectIntRange(&ok, rgdata.CurrentProtection, 0, 1, "DriveApiV2Storage.Pools.RaidGroups.CurrentProtection")
			u.expectIntRange(&ok, rgdata.ExpectedProtection, 0, 1, "DriveApiV2Storage.Pools.RaidGroups.ExpectedProtection")
			u.expectIntRange(&ok, rgdata.Progress, 0, 100, "DriveApiV2Storage.Pools.RaidGroups.Progress")
			u.expectIntRange(&ok, rgdata.Estimate, 0, 30000, "DriveApiV2Storage.Pools.RaidGroups.Estimate")
		}
		u.expectString(&ok, pdata.InitializingStatus, []string{"successful"}, "DriveApiV2Storage.Pools.InitializingStatus")
	}

	// Check disk data
	for _, ddata := range foo.Disks {
		u.expectString(&ok, ddata.Type, []string{"HDD", ""}, "DriveApiV2Storage.Disks.Type")
		u.expectString(&ok, ddata.State, []string{"optimal", "empty", "repairing", "scanning"}, "DriveApiV2Storage.Disks.State")
		u.expectInt(&ok, ddata.RPM, []int{0, 5900, 7200}, "DriveApiV2Storage.Disks.RPM")
		u.expectString(&ok, ddata.Sata, []string{"SATA 3.1", ""}, "DriveApiV2Storage.Disks.Sata")
		u.expectString(&ok, ddata.Ata, []string{"ATA8-ACS", "ACS-3", "ACS-2,", ""}, "DriveApiV2Storage.Disks.Ata")
		u.expectString(&ok, ddata.NvmeVersion, []string{""}, "DriveApiV2Storage.Disks.NvmeVersion")
		u.expectString(&ok, ddata.SectorFormat, []string{"512E", ""}, "DriveApiV2Storage.Disks.SectorFormat")
		if ddata.Temperature != 0 {
			u.expectIntRange(&ok, ddata.Temperature, 20, 70, "DriveApiV2Storage.Disks.Temperature")
		}
		// not doing counters
		u.expectIntRange(&ok, len(ddata.RiskReasons), 0, 0, "DriveApiV2Storage.Disks.RiskReasons.Len")
		u.expectIntRange(&ok, len(ddata.IncompatibleReasons), 0, 0, "DriveApiV2Storage.Disks.IncompatibleReasons.Len")
		u.expectIntRange(&ok, ddata.HealthScore, 0, 5, "DriveApiV2Storage.Disks.HealthScore")
	}

	u.expectIntRange(&ok, len(foo.CacheSlots), 0, 0, "DriveApiV2Storage.CacheSlots.Len")
	if !ok {
		return fmt.Errorf("errors during strict validation of /proxy/drive/api/v2/storage")
	}
	return nil
}
