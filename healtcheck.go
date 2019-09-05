package shutdown

import (
	"errors"
	"fmt"
	"time"
)

const (
	StatusServing    int32 = iota
	StatusNotServing
)

const (
	defaultInterval      = 1000
	defaultRetryAmount   = 100
	defaultCheckInterval = 3000
)

type HealthcheckConfig struct {
	RetryAmount   int `json:"retry_amount"`
	Interval      int `json:"interval"`       // in millisecond
	CheckInterval int `json:"check_interval"` // in millisecond
}

type Healthcheck struct {
	status int32
	errorCaused string

	cfg *HealthcheckConfig
	hcs []*serviceHealthcheck
}

type Ping func() error

type HealthcheckDescriptor interface {
	Ping() error
}

type serviceHealthcheck struct {
	name   string
	i      HealthcheckDescriptor
}

func NewHC(config *HealthcheckConfig) *Healthcheck {
	return &Healthcheck{
		cfg:    config,
		status: StatusServing,
	}
}

func (h *Healthcheck) Serve() {
	h.Policy()

	go func() {
		for {
			var err error = nil
			for _, hcs := range h.hcs {
				err = hcs.i.Ping()
				if err != nil {
					h.status = StatusNotServing
					h.errorCaused = fmt.Sprintf("%s -> %s", hcs.name, err.Error())

					break
				}
			}
			if err == nil {
				h.status = StatusServing
				h.errorCaused = ""
			}

			time.Sleep(2 * time.Second)
		}
	}()
}

func (h *Healthcheck) Add(name string, hcs HealthcheckDescriptor) {
	if err := hcs.Ping(); err != nil {
		h.status = StatusNotServing
	}

	h.hcs = append(h.hcs, &serviceHealthcheck{
		name: name,
		i:    hcs,
	})
}

func (h *Healthcheck) Status() (int32, error) {
	if h.errorCaused != "" {
		return h.status, errors.New(h.errorCaused)
	}
	return h.status, nil
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

func DefaultHealthcheckConfig() *HealthcheckConfig {
	return &HealthcheckConfig{
		RetryAmount: defaultRetryAmount,
		Interval: defaultInterval,
		CheckInterval: defaultCheckInterval,
	}
}