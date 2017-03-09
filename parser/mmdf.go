package parser

import "github.com/elastic/beats/libbeat/common"

// DfInfo contains the relevant information obtained in a single run of mmdf
type DfInfo struct {
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
	return nil, nil
}
