package parser

import "github.com/elastic/beats/libbeat/common"

// MmLsFsInfo contains the relevant information from a single mmlsfs run
type MmLsFsInfo struct {
	deviceName string
}

// ToMapStr turns the filesystem info into a common.MapStr
func (m *MmLsFsInfo) ToMapStr() common.MapStr {
	return common.MapStr{
		"device_name": m.deviceName,
	}
}

// UpdateDevice does not do anything, since we already have that information
func (m *MmLsFsInfo) UpdateDevice(device string) {}

// ParseMmLsFs returns the different devices in the GPFS cluster
func ParseMmLsFs(output string) ([]string, error) {
	var prefixFieldlocation = 0
	var identifierFieldLocation = 1
	var headerFieldLocation = 2

	ds, _ := parseGpfsYOutput(prefixFieldlocation, identifierFieldLocation, headerFieldLocation, "mmlsfs", output, parseMmLsFsCallback)

	// avoid duplicates
	var dsm = make(map[string]bool)
	var devices = make([]string, 0, len(ds))
	for _, device := range ds {
		d := device.(*MmLsFsInfo).deviceName // explicit type conversion
		_, ok := dsm[d]
		if !ok {
			dsm[d] = true
			devices = append(devices, d)
		}
	}

	return devices, nil
}

// parseMmLsFsCallback returns the device name found in the fields
func parseMmLsFsCallback(fields []string, fieldMap map[string]int) ParseResult {
	return &MmLsFsInfo{deviceName: fields[fieldMap["deviceName"]]}
}
