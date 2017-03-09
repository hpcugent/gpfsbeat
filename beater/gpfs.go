package beater

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

var debugf = logp.MakeDebug("gpfs")
var mmrepquotaTimeOut = 5 * 1000 * time.Millisecond
var mmlsfsTimeout = 1 * 1000 * time.Millisecond

// QuotaInfo contains the information of a single entry produced by mmrepquota
type QuotaInfo struct {
	filesystem string
	fileset    string
	kind       string
	entity     string
	blockUsage int64
	blockSoft  int64
	blockHard  int64
	blockDoubt int64
	blockGrace string
	filesUsage int64
	filesSoft  int64
	filesHard  int64
	filesDoubt int64
	filesGrace string
}

type parseCallBack func([]string, map[string]int) interface{}

// MmLsFs returns an array of the devices known to the GPFS cluster
func (bt *Gpfsbeat) MmLsFs() ([]string, error) {
	// get the filesystems from mmlsfs
	ctx, cancel := context.WithTimeout(context.Background(), mmlsfsTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bt.config.MMLsFsCommand, "all", "-Y")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		logp.Err("Command %s did not run correctly! Aborting! Error: %s", bt.config.MMLsFsCommand, err)
		panic(err)
	}

	devices, err := parseMmLsFs(out.String())
	if err != nil {
		var nope []string
		return nope, errors.New("mmlsfs info could not be parsed")
	}

	return devices, nil
}

// MmRepQuota is a wrapper around the mmrepquota command
func (bt *Gpfsbeat) MmRepQuota() ([]QuotaInfo, error) {

	var quotas []QuotaInfo

	for _, device := range bt.config.Devices {

		logp.Info("Running mmrepquota for device %s", device)

		ctx, cancel := context.WithTimeout(context.Background(), mmrepquotaTimeOut)
		defer cancel()

		cmd := exec.CommandContext(ctx, bt.config.MMRepQuotaCommand, "-Y", device)
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			logp.Err("Command mmrepquota did not run correctly for device %s! Aborting. Error: %s", device, err)
			var nope []QuotaInfo
			return nope, errors.New("mmrepquota failed")
		}

		var qs []QuotaInfo
		qs, err = parseMmRepQuota(out.String())
		if err != nil {
			var nope []QuotaInfo
			return nope, errors.New("mmrepquota info could not be parsed")
		}
		quotas = append(quotas, qs...)
	}
	return quotas, nil
}

// GetQuotaEvent turns the quota information into a MapStr
func (bt *Gpfsbeat) GetQuotaEvent(quota *QuotaInfo) common.MapStr {
	return common.MapStr{
		"filesystem":    quota.filesystem,
		"fileset":       quota.fileset,
		"kind":          quota.kind,
		"entity":        quota.entity,
		"block_usage":   quota.blockUsage,
		"block_soft":    quota.blockSoft,
		"block_hard":    quota.blockHard,
		"block_doubt":   quota.blockDoubt,
		"block_expired": quota.blockGrace,
		"files_usage":   quota.filesUsage,
		"files_soft":    quota.filesSoft,
		"files_hard":    quota.filesHard,
		"files_doubt":   quota.filesDoubt,
		"files_expired": quota.filesGrace,
	}
}

// parseCertainInt parses a string into an integer value, and discards the error
func parseCertainInt(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		logp.Err("Oops, could not parse an int from %s", s)
		panic(err)
	}
	return v
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

// parseMmLsFs returns the different devices in the GPFS cluster
func parseMmLsFs(output string) ([]string, error) {
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

// ParseMmRepQuota converts the lines into the desired information
func parseMmRepQuota(output string) ([]QuotaInfo, error) {

	var prefixFieldlocation = 0
	var identifierFieldLocation = 1
	var headerFieldLocation = 2

	qs, _ := parseGpfsYOutput(prefixFieldlocation, identifierFieldLocation, headerFieldLocation, "mmrepquota", output, parseMmRepQuotaCallback)

	var quotaInfos = make([]QuotaInfo, 0, len(qs))
	for _, q := range qs {
		quotaInfos = append(quotaInfos, q.(QuotaInfo))
	}

	return quotaInfos, nil
}

func parseMmRepQuotaCallback(fields []string, fieldMap map[string]int) interface{} {
	qi := QuotaInfo{
		filesystem: fields[fieldMap["filesystemName"]],
		fileset:    fields[fieldMap["filesetname"]],
		kind:       fields[fieldMap["quotaType"]],
		entity:     fields[fieldMap["name"]],
		blockUsage: parseCertainInt(fields[fieldMap["blockUsage"]]),
		blockSoft:  parseCertainInt(fields[fieldMap["blockQuota"]]),
		blockHard:  parseCertainInt(fields[fieldMap["blockLimit"]]),
		blockDoubt: parseCertainInt(fields[fieldMap["blockInDoubt"]]),
		blockGrace: fields[fieldMap["blockGrace"]],
		filesUsage: parseCertainInt(fields[fieldMap["filesUsage"]]),
		filesSoft:  parseCertainInt(fields[fieldMap["filesQuota"]]),
		filesHard:  parseCertainInt(fields[fieldMap["filesLimit"]]),
		filesDoubt: parseCertainInt(fields[fieldMap["filesInDoubt"]]),
		filesGrace: fields[fieldMap["filesGrace"]],
	}
	if qi.kind == "FILESET" {
		qi.fileset = qi.entity // filesets have no name, and we need to have a link between FILESET and USR quota
	}
	return qi
}

// parseMmLsFsCallback returns the device name found in the fields
func parseMmLsFsCallback(fields []string, fieldMap map[string]int) interface{} {
	return fields[fieldMap["deviceName"]]
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
