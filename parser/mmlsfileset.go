package parser

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
)

// MmLsFilesetInfo contains relevant information about GPFS filesets. For now, we ignore AFM info, as we do no use it.
type MmLsFilesetInfo struct {
	device            string
	version           int64
	filesystemName    string
	filesetName       string
	ID                int64
	rootInode         int64
	status            string
	path              string
	parentID          int64
	created           time.Time
	inodes            int64
	dataInKB          int64
	comment           string
	filesetMode       string
	inodeSpace        int64
	isInodeSpaceOwner bool
	maxInodes         int64
	allocInodes       int64
	inodeSpaceMask    int64
	snapID            int64
	permChangeFlag    string
	freeInodes        int64
}

// ToMapStr returns the fileset information in a common.MapStr
func (m *MmLsFilesetInfo) ToMapStr() common.MapStr {
	return common.MapStr{
		"version":             m.version,
		"filesystem_name":     m.filesystemName,
		"fileset_name":        m.filesetName,
		"ID":                  m.ID,
		"root_inode":          m.rootInode,
		"status":              m.status,
		"path":                m.path,
		"parent_ID":           m.parentID,
		"created":             m.created,
		"inodes":              m.inodes,
		"data_in_KB":          m.dataInKB,
		"comment":             m.comment,
		"fileset_mode":        m.filesetMode,
		"inode_space":         m.inodeSpace,
		"isInode_space_owner": m.isInodeSpaceOwner,
		"max_inodes":          m.maxInodes,
		"alloc_inodes":        m.allocInodes,
		"inode_space_mask":    m.inodeSpaceMask,
		"snap_ID":             m.snapID,
		"perm_change_flag":    m.permChangeFlag,
		"free_inodes":         m.freeInodes,
	}

}

// UpdateDevice sets the device name (this should be the same as filesystemName)
func (m *MmLsFilesetInfo) UpdateDevice(device string) {
	m.device = device
}

func parseMmLsFilesetCallback(fields []string, fieldMap map[string]int) ParseResult {

	creation_time, _ := time.Parse("Mon Jan 2 15%3A04%3A05 2006", fields[fieldMap["created"]])

	return &MmLsFilesetInfo{
		version:           parseCertainInt(fields[fieldMap["version"]]),
		filesystemName:    fields[fieldMap["filesystemName"]],
		filesetName:       fields[fieldMap["filesetName"]],
		ID:                parseCertainInt(fields[fieldMap["id"]]),
		rootInode:         parseCertainInt(fields[fieldMap["rootInode"]]),
		status:            fields[fieldMap["status"]],
		path:              fields[fieldMap["path"]],
		parentID:          parseCertainInt(fields[fieldMap["parentId"]]),
		created:           creation_time,
		inodes:            parseCertainInt(fields[fieldMap["inodes"]]),
		dataInKB:          parseCertainInt(fields[fieldMap["dataInKB"]]),
		comment:           fields[fieldMap["comment"]],
		filesetMode:       fields[fieldMap["filesetMode"]],
		inodeSpace:        parseCertainInt(fields[fieldMap["inodesSpace"]]),
		isInodeSpaceOwner: fields[fieldMap["isInodeSpaceOwner"]] == "1",
		maxInodes:         parseCertainInt(fields[fieldMap["maxInodes"]]),
		allocInodes:       parseCertainInt(fields[fieldMap["allocInodes"]]),
		inodeSpaceMask:    parseCertainInt(fields[fieldMap["inodeSpaceMask"]]),
		snapID:            parseCertainInt(fields[fieldMap["snapId"]]),
		permChangeFlag:    fields[fieldMap["permChangeFlag"]],
		freeInodes:        parseCertainInt(fields[fieldMap["freeInodes"]]),
	}
}
