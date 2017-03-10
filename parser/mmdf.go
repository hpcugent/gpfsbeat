package parser

import "github.com/elastic/beats/libbeat/common"

// DfInfo contains the relevant information obtained in a single run of mmdf
type DfInfo interface {
}

// MmDfNSDInfo represents the `nsd` output line information
type MmDfNSDInfo struct {
	version                 int64
	nsdname                 string
	storagePool             string
	diskSize                int64
	failureGroup            int64
	metadata                bool
	data                    bool
	freeBlocks              int64
	freeBlocksPercentage    int64
	freeFragments           int64
	freeFragmentsPercentage int64
	diskAvailableForAlloc   string // no idea what this should be
}

// MmDfPoolTotalInfo represent the `poolTotal` output line information
type MmDfPoolTotalInfo struct {
	version                 int64
	poolName                string
	poolSize                int64
	freeBlocks              int64
	freeBlocksPercentage    int64
	freeFragments           int64
	freeFragmentsPercentage int64
	maxDiskSize             int64
}

// MmDfFsTotalInfo represents the `fstotal` output line information
type MmDfFsTotalInfo struct {
	version                 int64
	fsSize                  int64
	freeBlocks              int64
	freeBlocksPercentage    int64
	freeFragments           int64
	freeFragmentsPercentage int64
}

// MmDfInodeInfo represents the `inode` ouput line information
type MmDfInodeInfo struct {
	version         int64
	usedInodes      int64
	freeInodes      int64
	allocatedInodes int64
	maxInodes       int64
}

// GetDfEvent turns the mmdf information into a MapStr
func GetDfEvent(dfinfo *DfInfo) common.MapStr {
	return common.MapStr{}

}

func parseMmDfInfoCallback(fields []string, fieldMap map[string]int) interface{} {

	return nil
}

// ParseMmDf converts the lines in the output string into the desired information
func ParseMmDf(output string) (DfInfo, error) {

	var prefixFieldlocation = 0
	var identifierFieldLocation = 1
	var headerFieldLocation = 2

	mmdfs, _ := parseGpfsYOutput(prefixFieldlocation, identifierFieldLocation, headerFieldLocation, "mmrepquota", output, parseMmRepQuotaCallback)

	var quotaInfos = make([]QuotaInfo, 0, len(qs))
	for _, q := range qs {
		quotaInfos = append(quotaInfos, q.(QuotaInfo))
	}

	return quotaInfos, nil	return nil, nil
}
