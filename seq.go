package shutdown

import (
	"math/rand"
	"time"
)

const (
	Init = iota
	Before
	After
)

const (
	min = 1000000
	mid = 5000000
	max = 9999999
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (h *Handler) GenerateSeq(label string, pos int) int {
	if pos == Init {
		return randInt(mid, mid+1000)
	}

	var prev = min
	var serviceSeq = min+10000
	var next = max
	var after = false
	var finish = false
	h.connections.Range(func(k, v interface{}) bool {
		if c, ok := v.(container); ok {
			if finish {
				return false
			} else if c.label == label {
				serviceSeq = k.(int)
				after = true
			} else if after {
				next = k.(int)
				finish = true
			} else {
				prev = k.(int)
			}
			return true
		}
		return true
	})
	if pos == Before {
		return randInt(prev, serviceSeq)
	} else {
		return randInt(serviceSeq, next)
	}
}

func randInt(min, max int) int {
	return rand.Intn(max - min) + min
}