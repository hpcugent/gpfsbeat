package beater

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"

	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/itkovian/gpfsbeat/parser"
)

var debugf = logp.MakeDebug("gpfs")

var mmrepquotaTimeOut = 5 * 60 * 1000 * time.Millisecond
var mmlsfsTimeout = 1 * 60 * 1000 * time.Millisecond
var mmdfTimeout = 5 * 60 * 1000 * time.Millisecond
var mmlsfilesetTimeout = 5 * 60 * 1000 * time.Millisecond

// MmLsFs returns an array of the devices known to the GPFS cluster
func (bt *gpfsbeat) MmLsFs() ([]string, error) {
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

	devices, err := parser.ParseMmLsFs(out.String())
	if err != nil {
		var nope []string
		return nope, errors.New("mmlsfs info could not be parsed")
	}

	return devices, nil
}

// MmRepQuota is a wrapper around the mmrepquota command
func (bt *gpfsbeat) MmRepQuota() ([]parser.QuotaInfo, error) {
	var quotas []parser.QuotaInfo

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
			return nil, errors.New("mmrepquota failed")
		}

		var qs []parser.QuotaInfo
		qs, err = parser.ParseMmRepQuota(out.String())
		if err != nil {
			return nil, errors.New("mmrepquota info could not be parsed")
		}
		quotas = append(quotas, qs...)
	}
	return quotas, nil
}

// MmDf is a wrapper around the mmdf command
func (bt *gpfsbeat) MmDf() ([]parser.ParseResult, error) {

	var mmdfinfos []parser.ParseResult

	for _, device := range bt.config.Devices {
		logp.Info("Running mmdf for device %s", device)

		ctx, cancel := context.WithTimeout(context.Background(), mmdfTimeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, bt.config.MMDfCommand, device, "-Y")
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			logp.Err("Command mmdf did not run correctly for device %s! Aborting. Error: %s", device, err)
			return nil, errors.New("mmdf failed")
		}

		var qs []parser.ParseResult
		qs, err = parser.ParseMmDf(device, out.String())
		if err != nil {
			return nil, errors.New("mmdf info could not be parsed")
		}
		mmdfinfos = append(mmdfinfos, qs...)
	}
	return mmdfinfos, nil
}

// MmLsFileset is a wrapper around the mmlsfileset command
func (bt *gpfsbeat) MmLsFileset() ([]parser.MmLsFilesetInfo, error) {

	var mmlsfilesetinfos []parser.MmLsFilesetInfo

	for _, device := range bt.config.Devices {

		logp.Info("Running mmlsfileset for device %s", device)

		ctx, cancel := context.WithTimeout(context.Background(), mmlsfilesetTimeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, bt.config.MMLsFilesetCommand, device, "-L", "-Y")
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			logp.Err("Command mmlsfileset did not runn correctly for device %s! Error: %s", device, err)
			return nil, err
		}

		var fs []parser.MmLsFilesetInfo
		fs, err = parser.ParseMmLsFileset(device, out.String())
		if err != nil {
			return nil, errors.New("mmlsfileset info could not be parsed")
		}
		mmlsfilesetinfos = append(mmlsfilesetinfos, fs...)
	}

	return mmlsfilesetinfos, nil
}
