package shutdown

import (
	"os"
	"time"
)

const (
	statusServing    int32 = 1
	statusNotServing int32 = 2
)

const (
	defaultInterval      = 1000
	defaultRetryAmount   = 100
	defaultCheckInterval = 3000
	discoveryName        = "discovery"
)

type HealthcheckConfig struct {
	RetryAmount   int `json:"retry_amount"`
	Interval      int `json:"interval"`       // in millisecond
	CheckInterval int `json:"check_interval"` // in millisecond
}

type Healthcheck struct {
	status int32
	log Log

	cfg *HealthcheckConfig
	hcs []*serviceHealthcheck
}

type Ping func() error
type Reconnect func() error

type HealthcheckDescriptor interface {
	Ping() error
	Reconnect() error
}

type serviceHealthcheck struct {
	Name   string
	status int32
	i      HealthcheckDescriptor
}

//type Response struct {
//	ID        string          `json:"id"`
//	Type      string          `json:"type"`
//	StartedAt time.Time       `json:"started_at"`
//	Status    string          `json:"status"`
//	Health    map[string]bool `json:"health"`
//}

func NewHC(config *HealthcheckConfig, l Log) *Healthcheck {
	// resp.Health = make(map[string]bool)

	return &Healthcheck{
		log: l,
		cfg:    config,
		status: statusServing,
	}
}

func (h *Healthcheck) Serve() {
	h.Policy()

	go func(log Log) {
		for {
			// check all hcs
			var errHCS = make(map[string]error)
			for _, hcs := range h.hcs {
				err := hcs.i.Ping()
				if err != nil {
					errHCS[hcs.Name] = err
					log.Debugf("%s ping caused error", hcs.Name)
					hcs.status = statusNotServing
					h.status = statusNotServing

					// retry hcs
					err = h.retry(hcs.Name, hcs.i, log)
				}
				if err != nil {
					errHCS[hcs.Name] = err
				}
			}

			// if all healthy set Response
			if len(errHCS) == 0 {
				h.status = statusServing
			} else {
				h.status = statusNotServing
			}

			time.Sleep(2 * time.Second)
		}
	}(h.log)
}

func (h *Healthcheck) Add(name string, hcs HealthcheckDescriptor) {
	if err := hcs.Ping(); err != nil {
		h.status = statusNotServing
	}

	h.hcs = append(h.hcs, &serviceHealthcheck{
		Name: name,
		i:    hcs,
	})
}

func (h *Healthcheck) CheckHCS() bool {
	for _, hcs := range h.hcs {
		if hcs.status == statusNotServing {
			return false
		}
	}
	return true
}

func (h *Healthcheck) Status() int32 {
	// var s = map[int32]string{
	// 	1: "serving",
	// 	2: "not_serving",
	// }

	// h.loggit.Debugf("Healthcheck happened, status: %s", s[h.status])
	return h.status
}

func (h *Healthcheck) Policy() *HealthcheckConfig {
	if h.cfg.RetryAmount == 0 {
		h.cfg.RetryAmount = defaultRetryAmount
	}

	if h.cfg.Interval == 0 {
		h.cfg.Interval = defaultInterval
	}

	if h.cfg.CheckInterval == 0 {
		h.cfg.CheckInterval = defaultCheckInterval
	}

	return h.cfg
}

func (h *Healthcheck) retry(name string, hcs HealthcheckDescriptor, log Log) error {
	var i int
	var retryErr error

	for i = 0; i < h.cfg.RetryAmount; i++ {
		retryErr = hcs.Reconnect()

		if retryErr == nil {
			log.Warnf("retry happened for %s with error: - %v -", name, retryErr)
			if h.CheckHCS() {
				h.status = statusServing
			}
			break
		}
		log.Debugf("%s still unavailable", name)
		time.Sleep(time.Duration(h.cfg.Interval) * time.Millisecond)
	}
	if retryErr != nil {
		log.Errorf("%s still unavailable os.Exit(1)", name)
		os.Exit(1)
	}
	return nil
}

func DefaultHealthcheckConfig() *HealthcheckConfig {
	return &HealthcheckConfig{
		RetryAmount: defaultRetryAmount,
		Interval: defaultInterval,
		CheckInterval: defaultCheckInterval,
	}
}