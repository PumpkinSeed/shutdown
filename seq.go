package shutdown

import (
	"math/rand"
	"time"
)

const (
	Before = iota
	After
)

const (
	min = 1000000
	max = 9999999
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (h *Handler) GenerateSeq(label string, pos int) int {
	var prev = min
	var next = min+10000
	var after = false
	h.connections.Range(func(k, v interface{}) bool {
		if c, ok := v.(container); ok {
			if c.label == label {
				prev = next
				next = k.(int)
				if pos == Before || after == true {
					return false
				} else if pos == After {
					after = true
					return true
				} else {
					return false
				}
			}
		}
		return true
	})
	if after == true {
		prev = next
		next = max
	}
	if next > max {
		next = max
	}
	return randInt(prev, next)
}

func randInt(min, max int) int {
	return rand.Intn(max - min) + min
}