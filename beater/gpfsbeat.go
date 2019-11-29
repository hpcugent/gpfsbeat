package beater

import (
	"fmt"
	"os"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/itkovian/gpfsbeat/config"
)

// gpfsbeat configuration.
type gpfsbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of gpfsbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	logp.Info("Gathering information from devices %q", c.Devices)
	logp.Info("Running every %d nanoseconds", c.Period)

	bt := &gpfsbeat{
		done:   make(chan struct{}),
		config: c,
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

// Run starts gpfsbeat.
func (bt *gpfsbeat) Run(b *beat.Beat) error {
	logp.Info("gpfsbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

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
				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"type":    b.Info.Name,
						"counter": counter,
						"quota":   quota,
					},
				}
				bt.client.Publish(event)
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
				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"type":    b.Info.Name,
						"counter": counter,
						"mmdf":    info,
					},
				}
				bt.client.Publish(event)
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
				event := beat.Event{
					Timestamp: time.Now(),
					Fields: common.MapStr{
						"type":        b.Info.Name,
						"counter":     counter,
						"mmlsfileset": info,
					},
				}
				bt.client.Publish(event)
			}
			logp.Info("mmlsfileset events sent")
		} else {
			logp.Err("Could not retrieve mmlsfileset information")
		}

		counter++
	}
}

// Stop stops gpfsbeat.
func (bt *gpfsbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
