package beater

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"

	"github.com/hpcugent/gpfsbeat/parser"

	"github.com/elastic/beats/libbeat/logp"
)

var debugf = logp.MakeDebug("gpfs")
var mmrepquotaTimeOut = 5 * 60 * 1000 * time.Millisecond
var mmlsfsTimeout = 1 * 60 * 1000 * time.Millisecond
var mmdfTimeout = 1 * 60 * 1000 * time.Millisecond

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

	devices, err := parser.ParseMmLsFs(out.String())
	if err != nil {
		var nope []string
		return nope, errors.New("mmlsfs info could not be parsed")
	}

	return devices, nil
}

// MmRepQuota is a wrapper around the mmrepquota command
func (bt *Gpfsbeat) MmRepQuota() ([]parser.QuotaInfo, error) {

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
			var nope []parser.QuotaInfo
			return nope, errors.New("mmrepquota failed")
		}

		var qs []parser.QuotaInfo
		qs, err = parser.ParseMmRepQuota(out.String())
		if err != nil {
			var nope []parser.QuotaInfo
			return nope, errors.New("mmrepquota info could not be parsed")
		}
		quotas = append(quotas, qs...)
	}
	return quotas, nil
}

// MmDf is a wrapper around the mmdf command
func (bt *Gpfsbeat) MmDf() ([]parser.ParseResult, error) {

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
			var nope []parser.ParseResult
			return nope, errors.New("mmdf failed")
		}

		var qs []parser.ParseResult
		qs, err = parser.ParseMmDf(device, out.String())
		if err != nil {
			var nope []parser.ParseResult
			return nope, errors.New("mmdf info could not be parsed")
		}
		mmdfinfos = append(mmdfinfos, qs...)
	}
	return mmdfinfos, nil
}
