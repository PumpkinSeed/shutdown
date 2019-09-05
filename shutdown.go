package shutdown

import (
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
)

type ServiceDescriptor interface {
	Stop() error
	Ping() error
}

type Handler struct {
	log         Log
	connections *sync.Map

	Healthcheck *Healthcheck
}

type container struct {
	label string
	conn  ServiceDescriptor
}

func NewHandler(l Log) *Handler {
	return &Handler{
		log:         l,
		connections: &sync.Map{},
	}
}

func (h *Handler) InitHealthcheck() {
	h.Healthcheck = NewHC(DefaultHealthcheckConfig())
}

func (h *Handler) HealthcheckWithExistingConnections() {
	h.connections.Range(func(_, connection interface{}) bool {
		if connectionv, ok := connection.(container); ok {
			h.Healthcheck.Add(connectionv.label, connectionv.conn)
		}
		return true
	})
	h.Healthcheck.Serve()
}

func (h *Handler) SetupHealthcheck(hc *Healthcheck) {
	h.Healthcheck = hc
}

func (h *Handler) Add(label string, labelPos string, pos int, stop ServiceDescriptor) {
	c := container{
		label: label,
		conn:  stop,
	}
	seq := h.GenerateSeq(labelPos, pos)
	h.connections.Store(seq, c)
}

func (h *Handler) Stop() error {
	var err error
	var keys []int
	h.connections.Range(func(seq, _ interface{}) bool {
		if seqv, ok := seq.(int); ok {
			keys = append(keys, seqv)
		}
		return true
	})
	sort.Ints(keys)

	for _, seq := range keys {
		if conn, ok := h.connections.Load(seq); ok {
			if connv, ok := conn.(container); ok {
				if connv.conn != nil {
					err = connv.conn.Stop()
					if err != nil {
						h.log.Errorf("Stop service error: %s", err.Error())
						os.Exit(1)
						return err
					}
					h.log.Infof("Stop service: %s, seq: %d", connv.label, seq)
				}
			}
		}
	}
	return err
}

func (h *Handler) debug() map[string]int {
	var result = make(map[string]int)
	h.connections.Range(func(seq, containerEntity interface{}) bool {
		var seqv int
		var containerEntityv container
		var ok bool
		if seqv, ok = seq.(int); !ok {
			return true
		}
		if containerEntityv, ok = containerEntity.(container); !ok {
			return true
		}
		result[containerEntityv.label] = seqv
		return true
	})

	return result
}

func (h *Handler) debugSeq(label string) int {
	var seq int
	h.connections.Range(func(k, v interface{}) bool {
		if c, ok := v.(container); ok {
			if c.label == label {
				seq = k.(int)
				return false
			}
		}
		return true
	})
	return seq
}

func GracefulExit(handler *Handler, wg *sync.WaitGroup) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case sig := <-c:
				handler.log.Infof("Got %s signal. Aborting...\n", sig)
				handler.Stop() // @TODO unhandeled error
				handler.log.Infof("Shutdown succesful.\n")
				os.Exit(0)
				return
			}
		}
	}()
}
