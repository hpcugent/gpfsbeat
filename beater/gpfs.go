package beater

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

var debugf = logp.MakeDebug("gpfs")

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

// MmRepQuota is a wrapper around the mmrepquota command
func (bt *Gpfsbeat) MmRepQuota() ([]QuotaInfo, error) {

	var quotas []QuotaInfo
	for _, filesystem := range bt.config.Filesystem {

		logp.Info("Running mmrepquota for filesystem %s", filesystem)

		cmd := exec.Command(bt.config.MMRepQuotaCommand, "-Y", filesystem) // TODO: pass arguments
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			logp.Err("Command mmrepquota did not run correctly for filesystem %s! Aborting.", filesystem)
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

// ParseMmRepQuota converts the lines into the desired information
func parseMmRepQuota(output string) ([]QuotaInfo, error) {

	lines := strings.Split(output, "\n")
	fieldMap := parseMmRepQuotaHeader(lines[0])

	var qs = make([]QuotaInfo, 0, 100000)

	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		fields := strings.Split(line, ":")
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
		qs = append(qs, qi)
	}
	return qs, nil
}

// parseMmRepQuotaHeader builds a map of the field names and the corresponding index
func parseMmRepQuotaHeader(header string) map[string]int {

	var m = make(map[string]int)

	for i, s := range strings.Split(header, ":") {
		if s == "" {
			continue
		}
		logp.Info("Currently processing header entry %d: %s", i, s)
		m[s] = i
	}
	logp.Info("All headers processed")

	return m
}
