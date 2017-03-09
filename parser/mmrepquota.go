package parser

import (
	"github.com/elastic/beats/libbeat/common"
)

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

// GetQuotaEvent turns the quota information into a MapStr
func GetQuotaEvent(quota *QuotaInfo) common.MapStr {
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

// ParseMmRepQuota converts the lines into the desired information
func ParseMmRepQuota(output string) ([]QuotaInfo, error) {

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
