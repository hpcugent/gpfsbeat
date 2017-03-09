package parser

// ParseMmLsFs returns the different devices in the GPFS cluster
func ParseMmLsFs(output string) ([]string, error) {
	var prefixFieldlocation = 0
	var identifierFieldLocation = 1
	var headerFieldLocation = 2

	ds, _ := parseGpfsYOutput(prefixFieldlocation, identifierFieldLocation, headerFieldLocation, "mmlsfs", output, parseMmLsFsCallback)

	// avoid duplicates
	var dsm map[string]bool
	var devices = make([]string, 0, len(ds))
	for _, device := range ds {
		d := device.(string) // explicit type conversion
		_, ok := dsm[d]
		if !ok {
			dsm[d] = true
			devices = append(devices, d)
		}
	}

	return devices, nil
}

// parseMmLsFsCallback returns the device name found in the fields
func parseMmLsFsCallback(fields []string, fieldMap map[string]int) interface{} {
	return fields[fieldMap["deviceName"]]
}
