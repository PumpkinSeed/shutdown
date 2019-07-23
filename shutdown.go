package shutdown

import (
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
)

type Stop interface {
	Stop() error
}

type Handler struct {
	log Log
	connections *sync.Map
	hc          *HC
}

type container struct {
	label string
	conn  Stop
}

func NewHandler(l Log, hc *HC) *Handler {
	return &Handler{
		log:      l,
		connections: &sync.Map{},
		hc:          hc,
	}
}

func (h *Handler) Add(label string, seq int, stop Stop) {
	c := container{
		label: label,
		conn:  stop,
	}
	h.connections.Store(seq, c)
}

func (h *Handler) Stop() error {
	h.hc.status = statusNotServing

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
				handler.Stop()
				handler.log.Infof("Shutdown succesful.\n")
				os.Exit(0)
				return
			}
		}
	}()
}
