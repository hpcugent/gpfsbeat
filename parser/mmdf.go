package parser

import (
	"github.com/elastic/beats/v7/libbeat/common"
)

// MmDfNSDInfo represents the `nsd` output line information
type MmDfNSDInfo struct {
	device                  string
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

// ToMapStr turns the nsd information into a common.MapStr
func (m *MmDfNSDInfo) ToMapStr() common.MapStr {
	return common.MapStr{
		"device":                    m.device,
		"version":                   m.version,
		"nsd_name":                  m.nsdname,
		"storage_pool":              m.storagePool,
		"disk_size":                 m.diskSize,
		"failure_group":             m.failureGroup,
		"metadata":                  m.metadata,
		"data":                      m.data,
		"free_blocks":               m.freeBlocks,
		"free_blocks_percentage":    m.freeBlocksPercentage,
		"free_fragments":            m.freeFragments,
		"free_fragments_percentage": m.freeFragmentsPercentage,
		"info_type":                 "nsd",
	}
}

// UpdateDevice sets the device name
func (m *MmDfNSDInfo) UpdateDevice(device string) {
	m.device = device
}

// MmDfPoolTotalInfo represent the `poolTotal` output line information
type MmDfPoolTotalInfo struct {
	device                  string
	version                 int64
	poolName                string
	poolSize                int64
	freeBlocks              int64
	freeBlocksPercentage    int64
	freeFragments           int64
	freeFragmentsPercentage int64
	maxDiskSize             int64
}

// ToMapStr turns the pool total information into a common.MapStr
func (m *MmDfPoolTotalInfo) ToMapStr() common.MapStr {
	return common.MapStr{
		"device":                    m.device,
		"version":                   m.version,
		"pool_name":                 m.poolName,
		"pool_size":                 m.poolSize,
		"free_blocks":               m.freeBlocks,
		"free_blocks_percentage":    m.freeBlocksPercentage,
		"free_fragments":            m.freeFragments,
		"free_fragments_percentage": m.freeFragmentsPercentage,
		"max_disk_size":             m.maxDiskSize,
		"info_type":                 "pooltotal",
	}
}

// UpdateDevice sets the device name
func (m *MmDfPoolTotalInfo) UpdateDevice(device string) {
	m.device = device
}

// MmDfFsTotalInfo represents the `fstotal` output line information
type MmDfFsTotalInfo struct {
	device                  string
	version                 int64
	fsSize                  int64
	freeBlocks              int64
	freeBlocksPercentage    int64
	freeFragments           int64
	freeFragmentsPercentage int64
}

// ToMapStr turns the fs total information into a common.MapStr
func (m *MmDfFsTotalInfo) ToMapStr() common.MapStr {
	return common.MapStr{
		"device":                    m.device,
		"version":                   m.version,
		"fs_size":                   m.fsSize,
		"free_blocks":               m.freeBlocks,
		"free_blocks_percentage":    m.freeBlocksPercentage,
		"free_fragments":            m.freeFragments,
		"free_fragments_percentage": m.freeFragmentsPercentage,
		"info_type":                 "fstotal",
	}

}

// UpdateDevice sets the device name
func (m *MmDfFsTotalInfo) UpdateDevice(device string) {
	m.device = device
}

// MmDfInodeInfo represents the `inode` ouput line information
type MmDfInodeInfo struct {
	device          string
	version         int64
	usedInodes      int64
	freeInodes      int64
	allocatedInodes int64
	maxInodes       int64
}

// ToMapStr turns the inode information into a common.MapStr
func (m *MmDfInodeInfo) ToMapStr() common.MapStr {
	return common.MapStr{
		"device":           m.device,
		"version":          m.version,
		"used_inodes":      m.usedInodes,
		"free_inodes":      m.freeInodes,
		"allocated_inodes": m.allocatedInodes,
		"max_inodex":       m.maxInodes,
		"info_type":        "inodes",
	}
}

// UpdateDevice sets the device name
func (m *MmDfInodeInfo) UpdateDevice(device string) {
	m.device = device
}

func parseMmDfCallback(fields []string, fieldMap map[string]int) ParseResult {

	var identifierFieldLocation = 1

	switch fields[identifierFieldLocation] {
	case "nsd":
		return &MmDfNSDInfo{
			version:                 parseCertainInt(fields[fieldMap["version"]]),
			nsdname:                 fields[fieldMap["nsdName"]],
			storagePool:             fields[fieldMap["storagePool"]],
			diskSize:                parseCertainInt(fields[fieldMap["diskSize"]]),
			failureGroup:            parseCertainInt(fields[fieldMap["failureGroup"]]),
			metadata:                fields[fieldMap["metadata"]] == "Yes",
			data:                    fields[fieldMap["data"]] == "Yes",
			freeBlocks:              parseCertainInt(fields[fieldMap["freeBlocks"]]),
			freeBlocksPercentage:    parseCertainInt(fields[fieldMap["freeBlocksPct"]]),
			freeFragments:           parseCertainInt(fields[fieldMap["freeFragments"]]),
			freeFragmentsPercentage: parseCertainInt(fields[fieldMap["freeFragmentsPct"]]),
			diskAvailableForAlloc:   fields[fieldMap["diskAvailableForAlloc"]],
		}
	case "poolTotal":
		return &MmDfPoolTotalInfo{
			version:                 parseCertainInt(fields[fieldMap["version"]]),
			poolName:                fields[fieldMap["poolName"]],
			poolSize:                parseCertainInt(fields[fieldMap["poolSize"]]),
			freeBlocks:              parseCertainInt(fields[fieldMap["freeBlocks"]]),
			freeBlocksPercentage:    parseCertainInt(fields[fieldMap["freeBlocksPct"]]),
			freeFragments:           parseCertainInt(fields[fieldMap["freeFragments"]]),
			freeFragmentsPercentage: parseCertainInt(fields[fieldMap["freeFragmentsPct"]]),
			maxDiskSize:             parseCertainInt(fields[fieldMap["maxDiskSize"]]),
		}
	case "fsTotal":
		return &MmDfFsTotalInfo{
			version:                 parseCertainInt(fields[fieldMap["version"]]),
			fsSize:                  parseCertainInt(fields[fieldMap["fsSize"]]),
			freeBlocks:              parseCertainInt(fields[fieldMap["freeBlocks"]]),
			freeBlocksPercentage:    parseCertainInt(fields[fieldMap["freeBlocksPct"]]),
			freeFragments:           parseCertainInt(fields[fieldMap["freeFragments"]]),
			freeFragmentsPercentage: parseCertainInt(fields[fieldMap["freeFragmentsPct"]]),
		}
	case "inode":
		return &MmDfInodeInfo{
			version:         parseCertainInt(fields[fieldMap["version"]]),
			usedInodes:      parseCertainInt(fields[fieldMap["usedInodes"]]),
			freeInodes:      parseCertainInt(fields[fieldMap["freeInodes"]]),
			allocatedInodes: parseCertainInt(fields[fieldMap["allocatedInodes"]]),
			maxInodes:       parseCertainInt(fields[fieldMap["maxInodes"]]),
		}
	}
	return nil
}

// ParseMmDf converts the lines in the output string into the desired information
func ParseMmDf(device string, output string) ([]ParseResult, error) {

	var prefixFieldlocation = 0
	var identifierFieldLocation = 1
	var headerFieldLocation = 2

	mmdfs, _ := parseGpfsYOutput(prefixFieldlocation, identifierFieldLocation, headerFieldLocation, "mmdf", output, parseMmDfCallback)

	var dfs = make([]ParseResult, 0, len(mmdfs))
	for _, info := range mmdfs {
		if info == nil {
			continue // line could not be parsed
		}
		info.UpdateDevice(device)
		dfs = append(dfs, info)
	}

	return dfs, nil
}
