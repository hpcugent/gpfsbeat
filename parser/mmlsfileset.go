package parser

import (
	"strings"
	"time"

	"github.com/elastic/beats/v7/libbeat/common"
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

	creationTime, err := time.Parse("Mon Jan 2 15%3A04%3A05 2006", fields[fieldMap["created"]])
	if err != nil {
		panic(err)
	}

	var parentID int64
	parentID = -1
	if fields[fieldMap["parentId"]] == "-" {
		parentID = parseCertainInt(fields[fieldMap["parentId"]])
	}

	return &MmLsFilesetInfo{
		version:           parseCertainInt(fields[fieldMap["version"]]),
		filesystemName:    fields[fieldMap["filesystemName"]],
		filesetName:       fields[fieldMap["filesetName"]],
		ID:                parseCertainInt(fields[fieldMap["id"]]),
		rootInode:         parseCertainInt(fields[fieldMap["rootInode"]]),
		status:            fields[fieldMap["status"]],
		path:              strings.Replace(fields[fieldMap["path"]], "%2F", "/", -1),
		parentID:          parentID,
		created:           creationTime,
		comment:           fields[fieldMap["comment"]],
		filesetMode:       fields[fieldMap["filesetMode"]],
		inodeSpace:        parseCertainInt(fields[fieldMap["inodeSpace"]]),
		isInodeSpaceOwner: fields[fieldMap["isInodeSpaceOwner"]] == "1",
		maxInodes:         parseCertainInt(fields[fieldMap["maxInodes"]]),
		allocInodes:       parseCertainInt(fields[fieldMap["allocInodes"]]),
		inodeSpaceMask:    parseCertainInt(fields[fieldMap["inodeSpaceMask"]]),
		snapID:            parseCertainInt(fields[fieldMap["snapId"]]),
		permChangeFlag:    fields[fieldMap["permChangeFlag"]],
		freeInodes:        parseCertainInt(fields[fieldMap["freeInodes"]]),
	}
}

// ParseMmLsFileset converts the output lines to the desired format
func ParseMmLsFileset(device string, output string) ([]MmLsFilesetInfo, error) {
	var prefixFieldlocation = 0
	var identifierFieldLocation = 1
	var headerFieldLocation = 2

	fs, _ := parseGpfsYOutput(prefixFieldlocation, identifierFieldLocation, headerFieldLocation, "mmlsfileset", output, parseMmLsFilesetCallback)

	var filesetInfos = make([]MmLsFilesetInfo, 0, len(fs))
	for _, f := range fs {
		filesetInfos = append(filesetInfos, *(f.(*MmLsFilesetInfo)))
	}

	return filesetInfos, nil
}
