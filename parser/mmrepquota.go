package parser

import (
	"github.com/elastic/beats/v7/libbeat/common"
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

// ToMapStr turns the quota information into a common.MapStr
func (q *QuotaInfo) ToMapStr() common.MapStr {
	return common.MapStr{
		"filesystem":    q.filesystem,
		"fileset":       q.fileset,
		"kind":          q.kind,
		"entity":        q.entity,
		"block_usage":   q.blockUsage,
		"block_soft":    q.blockSoft,
		"block_hard":    q.blockHard,
		"block_doubt":   q.blockDoubt,
		"block_expired": q.blockGrace,
		"files_usage":   q.filesUsage,
		"files_soft":    q.filesSoft,
		"files_hard":    q.filesHard,
		"files_doubt":   q.filesDoubt,
		"files_expired": q.filesGrace,
	}
}

// UpdateDevice does not do anything, since we already have that information
func (q *QuotaInfo) UpdateDevice(device string) {}

func parseMmRepQuotaCallback(fields []string, fieldMap map[string]int) ParseResult {
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
	return &qi
}

// ParseMmRepQuota converts the lines into the desired information
func ParseMmRepQuota(output string) ([](QuotaInfo), error) {

	var prefixFieldlocation = 0
	var identifierFieldLocation = 1
	var headerFieldLocation = 2

	qs, _ := parseGpfsYOutput(prefixFieldlocation, identifierFieldLocation, headerFieldLocation, "mmrepquota", output, parseMmRepQuotaCallback)

	var quotaInfos = make([](QuotaInfo), 0, len(qs))
	for _, q := range qs {
		quotaInfos = append(quotaInfos, *(q.(*QuotaInfo)))
	}

	return quotaInfos, nil
}
