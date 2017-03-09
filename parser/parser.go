package parser

import (
	"strconv"
	"strings"

	"github.com/elastic/beats/libbeat/logp"
)

type parseCallBack func([]string, map[string]int) interface{}

// parseCertainInt parses a string into an integer value, and discards the error
func parseCertainInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logp.Err("Oops, could not parse an int from %s", s)
		panic(err)
	}
	return v
}

// parseMmRepQuotaHeader builds a map of the field names and the corresponding index
func parseGpfsHeaderFields(fields []string) (m map[string]int) {

	m = make(map[string]int)
	for i, s := range fields {
		if s == "" {
			continue
		}
		logp.Info("Currently processing header entry %d: %s", i, s)
		m[s] = i
	}
	logp.Info("All headers fields processed")

	return
}

// updateHeaderMap will update (or add) an entry in the headerMap corresponding to the identifier field
// If the fields do not indicate they are derived from a header line, nothing is done
// The function returns true if we updated the headerMap, false otherwise
func updateHeaderMap(headerFieldLocation int, identifierFieldLocation int, headerMap map[string](map[string]int), fields []string) bool {
	identifier := fields[identifierFieldLocation]
	if fields[headerFieldLocation] == "HEADER" {
		fieldMap := parseGpfsHeaderFields(fields)
		headerMap[identifier] = fieldMap
		return true
	}
	return false
}

// parseGpfsYOutput parses data produced by a GPFS command using the -Y flag
func parseGpfsYOutput(
	prefixFieldlocation int,
	identifierFieldLocation int,
	headerFieldLocation int,
	prefix string,
	output string,
	fn parseCallBack) ([]interface{}, error) {

	lines := strings.Split(output, "\n")
	var headerMap = make(map[string](map[string]int))

	result := make([]interface{}, 0, len(lines))

	for _, line := range lines {

		// ignore empty lines
		if line == "" {
			continue
		}

		// we need to get the fields since we need to know the identifier in order to look up the correct field names
		fields := strings.Split(line, ":")

		// ignore all weird lines :)
		if fields[prefixFieldlocation] != prefix {
			continue
		}

		// there may be multiple HEADER lines so we need to gather them here (which is ugly, granted)
		// we then also already have the line identifier, so no need to get it from the HEADER parsing
		if updateHeaderMap(headerFieldLocation, identifierFieldLocation, headerMap, fields) {
			continue // we updated the map, so we can skip the remainder for this header line
		}
		identifier := fields[identifierFieldLocation]
		fieldMap := headerMap[identifier]
		info := fn(fields, fieldMap)
		result = append(result, info)
	}

	return result, nil
}
