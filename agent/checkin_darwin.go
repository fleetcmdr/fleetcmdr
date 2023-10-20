package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
	"time"

	"howett.net/plist"
)

func (d *agentDaemon) checkin() {

	var data checkinData

	data.ID = d.ID
	data.Version = d.version
	//data.Serial = d.getSystemData().SPHardwareDataType[0].SerialNumber

	b := &bytes.Buffer{}
	ge := gob.NewEncoder(b)
	err := ge.Encode(data)
	if checkError(err) {
		return
	}

	resp, err := d.hc.Post(fmt.Sprintf("%s/%s", d.cmdHost, checkinURL), "application/octet-stream", b)
	if checkError(err) {
		return
	}
	defer resp.Body.Close()

}

func (d systemData) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintln(&b, d)
	return b.Bytes(), nil
}

func (d *agentDaemon) getSystemData() AppleSystemProfilerOutput {
	return d.systemData.(AppleSystemProfilerOutput)
}

func (d *agentDaemon) sendSystemData() {
	var data systemData
	data.ID = d.ID
	//start := time.Now()
	log.Printf("Reading system data...")
	var err error
	d.systemData, err = readSystemData()
	if checkError(err) {
		return
	}
	data.Payload = d.systemData

	//log.Printf("Got system data (took %s): %+v", time.Since(start).String(), d.getSystemData())
	//log.Printf("Got serial: %s", d.getSystemData().SPHardwareDataType[0].SerialNumber)
	//log.Printf("System data: %+v", d.getSystemData())

	b := &bytes.Buffer{}
	gob.Register(data)
	ge := gob.NewEncoder(b)
	err = ge.Encode(data)
	if checkError(err) {
		return
	}

	resp, err := d.hc.Post(fmt.Sprintf("%s/%s", d.cmdHost, systemDataURL), "application/octet-stream", b)
	if checkError(err) {
		return
	}
	defer resp.Body.Close()
}

type SPApplications struct {
	Name         string    `plist:"_name"`
	ArchKind     string    `plist:"arch_kind"`
	LastModified time.Time `plist:"lastModified"`
	ObtainedFrom string    `plist:"obtained_from"`
	Path         string    `plist:"path"`
	SignedBy     []string  `plist:"signed_by"`
	Version      string    `plist:"version"`
}

type SPConfigurationProfile struct {
	Items []struct {
		Items []struct {
			Name                             string `plist:"_name"`
			SpconfigprofilePayloadData       string `plist:"spconfigprofile_payload_data"`
			SpconfigprofilePayloadIdentifier string `plist:"spconfigprofile_payload_identifier"`
			SpconfigprofilePayloadUUID       string `plist:"spconfigprofile_payload_uuid"`
			SpconfigprofilePayloadVersion    int    `plist:"spconfigprofile_payload_version"`
		} `plist:"_items"`
		Name                             string `plist:"_name"`
		SpconfigprofileDescription       string `plist:"spconfigprofile_description"`
		SpconfigprofileInstallDate       string `plist:"spconfigprofile_install_date"`
		SpconfigprofileOrganization      string `plist:"spconfigprofile_organization"`
		SpconfigprofileProfileIdentifier string `plist:"spconfigprofile_profile_identifier"`
		SpconfigprofileProfileUUID       string `plist:"spconfigprofile_profile_uuid"`
		SpconfigprofileRemovalDisallowed string `plist:"spconfigprofile_RemovalDisallowed"`
		SpconfigprofileVerificationState string `plist:"spconfigprofile_verification_state"`
		SpconfigprofileVersion           int    `plist:"spconfigprofile_version"`
	} `plist:"_items"`
	Name string `plist:"_name"`
}

type SPDisplays struct {
	Name                          string `plist:"_name"`
	SpdisplaysMtlgpufamilysupport string `plist:"spdisplays_mtlgpufamilysupport"`
	SpdisplaysNdrvs               []struct {
		Name                          string `plist:"_name"`
		SpdisplaysDisplayProductID    string `plist:"_spdisplays_display-product-id"`
		SpdisplaysDisplaySerialNumber string `plist:"_spdisplays_display-serial-number"`
		SpdisplaysDisplayVendorID     string `plist:"_spdisplays_display-vendor-id"`
		SpdisplaysDisplayWeek         string `plist:"_spdisplays_display-week"`
		SpdisplaysDisplayYear         string `plist:"_spdisplays_display-year"`
		SpdisplaysDisplayID           string `plist:"_spdisplays_displayID"`
		SpdisplaysPixels              string `plist:"_spdisplays_pixels"`
		SpdisplaysResolution          string `plist:"_spdisplays_resolution"`
		SpdisplaysAmbientBrightness   string `plist:"spdisplays_ambient_brightness,omitempty"`
		SpdisplaysMain                string `plist:"spdisplays_main,omitempty"`
		SpdisplaysMirror              string `plist:"spdisplays_mirror"`
		SpdisplaysOnline              string `plist:"spdisplays_online"`
		SpdisplaysPixelresolution     string `plist:"spdisplays_pixelresolution"`
		SpdisplaysResolution0         string `plist:"spdisplays_resolution,omitempty"`
		SpdisplaysRotation            string `plist:"spdisplays_rotation,omitempty"`
		SpdisplaysTelevision          string `plist:"spdisplays_television,omitempty"`
		SpdisplaysConnectionType      string `plist:"spdisplays_connection_type,omitempty"`
		SpdisplaysDisplayType         string `plist:"spdisplays_display_type,omitempty"`
	} `plist:"spdisplays_ndrvs"`
	SpdisplaysVendor string `plist:"spdisplays_vendor"`
	SppciBus         string `plist:"sppci_bus"`
	SppciCores       string `plist:"sppci_cores"`
	SppciDeviceType  string `plist:"sppci_device_type"`
	SppciModel       string `plist:"sppci_model"`
}

type SPDisabledSoftware struct {
	Name         string `plist:"_name"`
	DisabledDate string `plist:"disabledDate"`
	Reason       string `plist:"reason"`
	Version      string `plist:"version"`
}

type SPEthernet struct {
	Name                     string `plist:"_name"`
	SpethernetAvbSupport     string `plist:"spethernet_avb_support"`
	SpethernetBSDDeviceName  string `plist:"spethernet_BSD_Device_Name"`
	SpethernetBus            string `plist:"spethernet_bus"`
	SpethernetDriver         string `plist:"spethernet_driver"`
	SpethernetMacAddress     string `plist:"spethernet_mac_address"`
	SpethernetProductName    string `plist:"spethernet_product_name"`
	SpethernetProductID      string `plist:"spethernet_product-id"`
	SpethernetUsbDeviceSpeed string `plist:"spethernet_usb_device_speed"`
	SpethernetVendorName     string `plist:"spethernet_vendor_name"`
	SpethernetVendorID       string `plist:"spethernet_vendor-id"`
}

type SPFirewall struct {
	Name                     string `plist:"_name"`
	SpfirewallGlobalstate    string `plist:"spfirewall_globalstate"`
	SpfirewallLoggingenabled string `plist:"spfirewall_loggingenabled"`
	SpfirewallStealthenabled string `plist:"spfirewall_stealthenabled"`
}

type SPInstallHistory struct {
	Name           string    `plist:"_name"`
	InstallDate    time.Time `plist:"install_date"`
	InstallVersion string    `plist:"install_version"`
	PackageSource  string    `plist:"package_source"`
}

type SPMemory struct {
	DimmManufacturer string `plist:"dimm_manufacturer"`
	DimmType         string `plist:"dimm_type"`
	SPMemory         string `plist:"SPMemory"`
}

type SPNetwork struct {
	Name     string `plist:"_name"`
	Ethernet struct {
		MACAddress   string `plist:"MAC Address"`
		MediaOptions []any  `plist:"MediaOptions"`
		MediaSubType string `plist:"MediaSubType"`
	} `plist:"Ethernet"`
	Hardware  string `plist:"hardware"`
	Interface string `plist:"interface"`
	IPv4      struct {
		ConfigMethod string `plist:"ConfigMethod"`
	} `plist:"IPv4"`
	IPv6 struct {
		ConfigMethod string `plist:"ConfigMethod"`
	} `plist:"IPv6"`
	Proxies struct {
		ExceptionsList []string `plist:"ExceptionsList"`
		FTPPassive     string   `plist:"FTPPassive"`
	} `plist:"Proxies"`
	SpnetworkServiceOrder int    `plist:"spnetwork_service_order"`
	Type                  string `plist:"type"`
}

type SPNetworkVolume struct {
	Name                       string `plist:"_name"`
	SpnetworkvolumeAutomounted string `plist:"spnetworkvolume_automounted"`
	SpnetworkvolumeFsmtnonname string `plist:"spnetworkvolume_fsmtnonname"`
	SpnetworkvolumeFstypename  string `plist:"spnetworkvolume_fstypename"`
	SpnetworkvolumeMntfromname string `plist:"spnetworkvolume_mntfromname"`
}

type SPNVMe struct {
	Items []struct {
		Name              string `plist:"_name"`
		BsdName           string `plist:"bsd_name"`
		DetachableDrive   string `plist:"detachable_drive"`
		DeviceModel       string `plist:"device_model"`
		DeviceRevision    string `plist:"device_revision"`
		DeviceSerial      string `plist:"device_serial"`
		PartitionMapType  string `plist:"partition_map_type"`
		RemovableMedia    string `plist:"removable_media"`
		Size              string `plist:"size"`
		SizeInBytes       int64  `plist:"size_in_bytes"`
		SmartStatus       string `plist:"smart_status"`
		SpnvmeTrimSupport string `plist:"spnvme_trim_support"`
		Volumes           []struct {
			Name        string `plist:"_name"`
			BsdName     string `plist:"bsd_name"`
			Iocontent   string `plist:"iocontent"`
			Size        string `plist:"size"`
			SizeInBytes int    `plist:"size_in_bytes"`
		} `plist:"volumes"`
	} `plist:"_items"`
	Name string `plist:"_name"`
}

type SPPower struct {
	Name                     string `plist:"_name"`
	SppowerBatteryChargeInfo struct {
		SppowerBatteryAtWarnLevel   string `plist:"sppower_battery_at_warn_level"`
		SppowerBatteryFullyCharged  string `plist:"sppower_battery_fully_charged"`
		SppowerBatteryIsCharging    string `plist:"sppower_battery_is_charging"`
		SppowerBatteryStateOfCharge int    `plist:"sppower_battery_state_of_charge"`
	} `plist:"sppower_battery_charge_info,omitempty"`
	SppowerBatteryHealthInfo struct {
		SppowerBatteryCycleCount            int    `plist:"sppower_battery_cycle_count"`
		SppowerBatteryHealth                string `plist:"sppower_battery_health"`
		SppowerBatteryHealthMaximumCapacity string `plist:"sppower_battery_health_maximum_capacity"`
	} `plist:"sppower_battery_health_info,omitempty"`
	SppowerBatteryModelInfo struct {
		PackLotCode                    string `plist:"Pack Lot Code"`
		PCBLotCode                     string `plist:"PCB Lot Code"`
		SppowerBatteryCellRevision     string `plist:"sppower_battery_cell_revision"`
		SppowerBatteryDeviceName       string `plist:"sppower_battery_device_name"`
		SppowerBatteryFirmwareVersion  string `plist:"sppower_battery_firmware_version"`
		SppowerBatteryHardwareRevision string `plist:"sppower_battery_hardware_revision"`
		SppowerBatterySerialNumber     string `plist:"sppower_battery_serial_number"`
	} `plist:"sppower_battery_model_info,omitempty"`
	ACPower struct {
		CurrentPowerSource                     string `plist:"Current Power Source"`
		DiskSleepTimer                         int    `plist:"Disk Sleep Timer"`
		DisplaySleepTimer                      int    `plist:"Display Sleep Timer"`
		HibernateMode                          int    `plist:"Hibernate Mode"`
		HighPowerMode                          int    `plist:"HighPowerMode"`
		LowPowerMode                           int    `plist:"LowPowerMode"`
		PrioritizeNetworkReachabilityOverSleep int    `plist:"PrioritizeNetworkReachabilityOverSleep"`
		SleepOnPowerButton                     string `plist:"Sleep On Power Button"`
		SystemSleepTimer                       int    `plist:"System Sleep Timer"`
		WakeOnLAN                              string `plist:"Wake On LAN"`
	} `plist:"AC Power,omitempty"`
	BatteryPower struct {
		DiskSleepTimer                         int    `plist:"Disk Sleep Timer"`
		DisplaySleepTimer                      int    `plist:"Display Sleep Timer"`
		HibernateMode                          int    `plist:"Hibernate Mode"`
		HighPowerMode                          int    `plist:"HighPowerMode"`
		LowPowerMode                           int    `plist:"LowPowerMode"`
		PrioritizeNetworkReachabilityOverSleep int    `plist:"PrioritizeNetworkReachabilityOverSleep"`
		ReduceBrightness                       string `plist:"ReduceBrightness"`
		SleepOnPowerButton                     string `plist:"Sleep On Power Button"`
		SystemSleepTimer                       int    `plist:"System Sleep Timer"`
		WakeOnLAN                              string `plist:"Wake On LAN"`
	} `plist:"Battery Power,omitempty"`
	SppowerUpsInstalled             string `plist:"sppower_ups_installed,omitempty"`
	SppowerAcChargerFamily          string `plist:"sppower_ac_charger_family,omitempty"`
	SppowerAcChargerFirmwareVersion string `plist:"sppower_ac_charger_firmware_version,omitempty"`
	SppowerAcChargerHardwareVersion string `plist:"sppower_ac_charger_hardware_version,omitempty"`
	SppowerAcChargerID              string `plist:"sppower_ac_charger_ID,omitempty"`
	SppowerAcChargerManufacturer    string `plist:"sppower_ac_charger_manufacturer,omitempty"`
	SppowerAcChargerName            string `plist:"sppower_ac_charger_name,omitempty"`
	SppowerAcChargerSerialNumber    string `plist:"sppower_ac_charger_serial_number,omitempty"`
	SppowerAcChargerWatts           string `plist:"sppower_ac_charger_watts,omitempty"`
	SppowerBatteryChargerConnected  string `plist:"sppower_battery_charger_connected,omitempty"`
	SppowerBatteryIsCharging        string `plist:"sppower_battery_is_charging,omitempty"`
	Items                           []struct {
		Items []struct {
			AppPID      int       `plist:"appPID"`
			Eventtype   string    `plist:"eventtype"`
			Scheduledby string    `plist:"scheduledby"`
			Time        time.Time `plist:"time"`
			UserVisible bool      `plist:"UserVisible"`
		} `plist:"_items"`
		Name string `plist:"_name"`
	} `plist:"_items,omitempty"`
}

type SPPrefPane struct {
	Name                 string `plist:"_name"`
	SpprefpaneBundlePath string `plist:"spprefpane_bundlePath"`
	SpprefpaneIdentifier string `plist:"spprefpane_identifier"`
	SpprefpaneIsVisible  string `plist:"spprefpane_isVisible"`
	SpprefpaneKind       string `plist:"spprefpane_kind"`
	SpprefpaneSupport    string `plist:"spprefpane_support"`
	SpprefpaneVersion    string `plist:"spprefpane_version"`
}

type SPPrinters struct {
	Cupsversion string `plist:"cupsversion"`
	Status      string `plist:"status"`
}

type SPSecureElement struct {
	CtlFw              string `plist:"ctl_fw"`
	CtlHw              string `plist:"ctl_hw"`
	CtlInfo            string `plist:"ctl_info"`
	CtlMw              string `plist:"ctl_mw"`
	SeDevice           string `plist:"se_device"`
	SeFw               string `plist:"se_fw"`
	SeHw               string `plist:"se_hw"`
	SeID               string `plist:"se_id"`
	SeInRestrictedMode string `plist:"se_in_restricted_mode"`
	SeInfo             string `plist:"se_info"`
	SeOsVersion        string `plist:"se_os_version"`
	SePlt              string `plist:"se_plt"`
	SeProdSigned       string `plist:"se_prod_signed"`
}

type SPSoftware struct {
	Name            string `plist:"_name"`
	BootMode        string `plist:"boot_mode"`
	BootVolume      string `plist:"boot_volume"`
	KernelVersion   string `plist:"kernel_version"`
	LocalHostName   string `plist:"local_host_name"`
	OsVersion       string `plist:"os_version"`
	SecureVM        string `plist:"secure_vm"`
	SystemIntegrity string `plist:"system_integrity"`
	Uptime          string `plist:"uptime"`
	UserName        string `plist:"user_name"`
}

type SPSPI struct {
	Name          string `plist:"_name"`
	AProductID    string `plist:"a_product_id"`
	BVendorID     string `plist:"b_vendor_id"`
	CStfwVersion  string `plist:"c_stfw_version"`
	DSerialNum    string `plist:"d_serial_num"`
	FManufacturer string `plist:"f_manufacturer"`
	GLocationID   string `plist:"g_location_id"`
	HMtfwVersion  string `plist:"h_mtfw_version"`
}

type SPStorage struct {
	Name             string `plist:"_name"`
	BsdName          string `plist:"bsd_name"`
	FileSystem       string `plist:"file_system"`
	FreeSpaceInBytes int64  `plist:"free_space_in_bytes"`
	IgnoreOwnership  string `plist:"ignore_ownership"`
	MountPoint       string `plist:"mount_point"`
	PhysicalDrive    struct {
		DeviceName       string `plist:"device_name"`
		IsInternalDisk   string `plist:"is_internal_disk"`
		MediaName        string `plist:"media_name"`
		MediumType       string `plist:"medium_type"`
		PartitionMapType string `plist:"partition_map_type"`
		Protocol         string `plist:"protocol"`
		SmartStatus      string `plist:"smart_status"`
	} `plist:"physical_drive"`
	SizeInBytes int64  `plist:"size_in_bytes"`
	VolumeUUID  string `plist:"volume_uuid"`
	Writable    string `plist:"writable"`
}

type SPThunderbolt struct {
	Name           string `plist:"_name"`
	DeviceNameKey  string `plist:"device_name_key"`
	DomainUUIDKey  string `plist:"domain_uuid_key"`
	Receptacle1Tag struct {
		CurrentLinkWidthKey string `plist:"current_link_width_key"`
		CurrentSpeedKey     string `plist:"current_speed_key"`
		LinkStatusKey       string `plist:"link_status_key"`
		ReceptacleIDKey     string `plist:"receptacle_id_key"`
		ReceptacleStatusKey string `plist:"receptacle_status_key"`
	} `plist:"receptacle_1_tag"`
	RouteStringKey string `plist:"route_string_key"`
	SwitchUIDKey   string `plist:"switch_uid_key"`
	VendorNameKey  string `plist:"vendor_name_key"`
}

type SPUniversalAccess struct {
	Name         string `plist:"_name"`
	Contrast     string `plist:"contrast"`
	CursorMag    string `plist:"cursor_mag"`
	Display      string `plist:"display"`
	FlashScreen  string `plist:"flash_screen"`
	KeyboardZoom string `plist:"keyboardZoom"`
	MouseKeys    string `plist:"mouse_keys"`
	ScrollZoom   string `plist:"scrollZoom"`
	SlowKeys     string `plist:"slow_keys"`
	StickyKeys   string `plist:"sticky_keys"`
	Voiceover    string `plist:"voiceover"`
	ZoomMode     string `plist:"zoomMode"`
}

type SPUSB []struct {
	Items []struct {
		Name             string `plist:"_name"`
		BcdDevice        string `plist:"bcd_device"`
		BusPower         string `plist:"bus_power"`
		BusPowerUsed     string `plist:"bus_power_used"`
		DeviceSpeed      string `plist:"device_speed"`
		ExtraCurrentUsed string `plist:"extra_current_used"`
		LocationID       string `plist:"location_id"`
		Manufacturer     string `plist:"manufacturer"`
		ProductID        string `plist:"product_id"`
		SerialNum        string `plist:"serial_num"`
		VendorID         string `plist:"vendor_id"`
	} `plist:"_items"`
	Name           string `plist:"_name"`
	HostController string `plist:"host_controller"`
}

type AppleSystemProfilerOutput struct {
	SPApplicationsDataType         []SPApplications `xml:"SPApplicationsDataType"`
	SPConfigurationProfileDataType []SPConfigurationProfile
	SPDisabledSoftwareDataType     []SPDisabledSoftware
	SPDisplaysDataType             []SPDisplays
	SPEthernetDataType             []SPEthernet
	SPFirewallDataType             []SPFirewall
	SPHardwareDataType             []SPHardware `plist:"SPHardwareDataType"`
	SPInstallHistoryDataType       []SPInstallHistory
	SPMemoryDataType               []SPMemory
	SPNetworkDataType              []SPNetwork
	SPNetworkVolumeDataType        []SPNetworkVolume
	SPNVMeDataType                 []SPNVMe
	SPPowerDataType                []SPPower
	SPPrefPaneDataType             []SPPrefPane
	SPPrintersDataType             []SPPrinters
	SPSecureElementDataType        []SPSecureElement
	SPSoftwareDataType             []SPSoftware
	SPSPIDataType                  []SPSPI
	SPStorageDataType              []SPStorage
	SPThunderboltDataType          []SPThunderbolt
	SPUniversalAccessDataType      []SPUniversalAccess
	SPUSBDataType                  []SPUSB
}

type SPHardware struct {
	Items struct {
		Name                 string `plist:"_name"`
		ActivationLockStatus string `plist:"activation_lock_status"`
		BootRomVersion       string `plist:"boot_rom_version"`
		ChipType             string `plist:"chip_type"`
		MachineModel         string `plist:"machine_model"`
		MachineName          string `plist:"machine_name"`
		ModelNumber          string `plist:"model_number"`
		NumberProcessors     string `plist:"number_processors"`
		OsLoaderVersion      string `plist:"os_loader_version"`
		PhysicalMemory       string `plist:"physical_memory"`
		PlatformUUID         string `plist:"platform_UUID"`
		ProvisioningUDID     string `plist:"provisioning_UDID"`
		SerialNumber         string `plist:"serial_number"`
	} `plist:"_items"`
}

type AppleSystemProfilerOutput2 struct {
	//SPApplicationsDataType []SPApplications `plist:"SPApplicationsDataType"`

	SPHardwareDataType []SPHardware `plist:"SPHardwareDataType"`
}

type thing []struct {
	SPHardwareDataType SPHardware
}

//func (a AppleSystemProfilerOutput) UnmarshalJSON(b []byte) error {
//	f := make(map[string]any)
//	err := json.Unmarshal(b, &f)
//	if checkError(err) {
//		return err
//	}
//
//	a.SPHardwareDataType = (f["SPHardwareDataType"]).([]SPHardware)
//
//	return nil
//}

func readSystemData() (AppleSystemProfilerOutput, error) {

	desirous := []string{
		"SPHardware",
		"SPApplications",
		"SPConfigurationProfile",
		"SPDisabledSoftware",
		"SPDisplays",
		"SPEthernet",
		"SPFirewall",
		"SPHardware",
		"SPInstallHistory",
		"SPMemory",
		"SPNetwork",
		"SPNetworkVolume",
		"SPNVMe",
		"SPPower",
		"SPPrefPane",
		"SPPrinters",
		"SPSecureElement",
		"SPSoftware",
		"SPStorage",
		"SPThunderbolt",
		"SPUniversalAccess",
		"SPUSB",
	}

	xmlData, err := run(fmt.Sprintf("/usr/sbin/system_profiler -xml %s", strings.Join(desirous, " ")))
	if checkError(err) {
		return AppleSystemProfilerOutput{}, err
	}

	//log.Printf("XML: %s", xmlData)

	//jsonData, err := os.ReadFile("systemprofiler.json")
	//if checkError(err) {
	//	return AppleSystemProfilerOutput{}, err
	//}

	var aspo AppleSystemProfilerOutput

	//var aspo2 interface{}

	//var stuff thing

	_, err = plist.Unmarshal([]byte(xmlData), &aspo)
	if checkError(err) {
		return aspo, err
	}

	//genericMap := make(map[string]any)

	//err = aspo.UnmarshalJSON([]byte(jsonData))
	//if checkError(err) {
	//	return aspo, err
	//}

	//err = json.Unmarshal([]byte(jsonData), &aspo)
	//if checkError(err) {
	//	return AppleSystemProfilerOutput{}, err
	//}

	log.Printf("Stuff: %#v", aspo)

	//for _, v := range aspo.([]interface{}) {
	//	vv := v.(map[string]interface{})
	//	for k, vvv := range vv {
	//		log.Printf("K: %v, V: %v", k, vvv)
	//		plist.Unmarshal()
	//	}
	//	log.Printf("Data: %#v", vv["_items"].([]interface{}))
	//}

	//config := mapstructure.DecoderConfig{WeaklyTypedInput: true, Result: &aspo.SPHardwareDataType}

	//decoder, err := mapstructure.NewDecoder(&config)
	//if checkError(err) {
	//	return aspo, err
	//}
	//err = decoder.Decode(genericMap["SPHardwareDataType"])
	//err = mapstructure.Decode(genericMap, &aspo)
	//if checkError(err) {
	//	return aspo, err
	//}

	//log.Printf("Converted: %#v", aspo.SPHardwareDataType)

	return aspo, nil
}
