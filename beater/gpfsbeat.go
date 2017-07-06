package beater

import (
	"fmt"
	"os"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/hpcugent/gpfsbeat/config"
)

// Gpfsbeat generated structure
type Gpfsbeat struct {
	done   chan struct{}
	config config.Config
	client publisher.Client
}

// New Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	logp.Info("Gathering information from devices %q", config.Devices)
	logp.Info("Running every %d nanoseconds", config.Period)

	bt := &Gpfsbeat{
		done:   make(chan struct{}),
		config: config,
	}

	// make sure we get the devices, request them from mmlsfs is they are not provided explicitly
	if len(bt.config.Devices) == 1 && bt.config.Devices[0] == "all" {
		logp.Info("Requested information from 'all' devices. Gathering devices.")
		devices, err := bt.MmLsFs()
		if err != nil {
			logp.Err("Cannot get required devices information. Stopping.")
			os.Exit(-1)
		}
		bt.config.Devices = devices
		logp.Info("Renewed devices list: %s", bt.config.Devices)
	}

	return bt, nil
}

// Run does the actual things
func (bt *Gpfsbeat) Run(b *beat.Beat) error {
	logp.Info("gpfsbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		gpfsQuota, err := bt.MmRepQuota()
		logp.Info("retrieved quota information from mmrepquota")
		if err == nil {
			for _, q := range gpfsQuota {
				quota := q.ToMapStr()
				event := common.MapStr{
					"@timestamp": common.Time(time.Now()),
					"type":       b.Name,
					"counter":    counter,
					"quota":      quota,
				}
				if bt.config.UserIDField != "none" && q.Kind() == "USR" {
					event.Update(common.MapStr{bt.config.UserIDField: q.Entity()})
				}
				bt.client.PublishEvent(event)
			}
			logp.Info("mmrepquota events sent")
		} else {
			logp.Err("Could not retrieve mmrequota information")
		}

		mmdfinfos, err := bt.MmDf()
		logp.Info("Retrieved usage information from mmdf")
		if err == nil {
			for _, i := range mmdfinfos {
				info := i.ToMapStr()
				event := common.MapStr{
					"@timestamp": common.Time(time.Now()),
					"type":       b.Name,
					"counter":    counter,
					"mmdf":       info,
				}
				bt.client.PublishEvent(event)
			}
			logp.Info("mmdf events sent")
		} else {
			logp.Err("Could not retrieve mmdf information")
		}

		mmlsfilesetinfos, err := bt.MmLsFileset()
		logp.Info("Retrieved usage information from mmlsfileset")
		if err == nil {
			for _, i := range mmlsfilesetinfos {
				info := i.ToMapStr()
				event := common.MapStr{
					"@timestamp":  common.Time(time.Now()),
					"type":        b.Name,
					"counter":     counter,
					"mmlsfileset": info,
				}
				bt.client.PublishEvent(event)
			}
			logp.Info("mmlsfileset events sent")
		} else {
			logp.Err("Could not retrieve mmlsfileset information")
		}

		counter++
	}
}

// Stop shuts down the beat
func (bt *Gpfsbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
