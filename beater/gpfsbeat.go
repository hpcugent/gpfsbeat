package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/itkovian/gpfsbeat/config"
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
	logp.Info("Checking quota for device %q", config.Device)
	logp.Info("Running every %d nanoseconds", config.Period)

	bt := &Gpfsbeat{
		done:   make(chan struct{}),
		config: config,
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

		gpfsQuota, err := bt.MmRepQuota() // TODO: get this for each device
		logp.Warn("retrieved quota information from mmrepquota")
		if err != nil {
			panic("Could not get quota information")
		}

		for _, q := range gpfsQuota {
			quota := bt.GetQuotaEvent(&q)
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       b.Name,
				"counter":    counter,
				"quota":      quota,
			}
			bt.client.PublishEvent(event)
		}
		logp.Info("Events sent")
		counter++
	}
}

// Stop shuts down the beat
func (bt *Gpfsbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
