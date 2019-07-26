package shutdown

import (
	"fmt"
	"testing"
)

func TestGenerateSeqAfter(t *testing.T) {
	sh := NewHandler(&blankLog{})
	sh.Add("test", "", min+90000, &serviceWithStop{})
	seq := sh.GenerateSeq("test", After)
	fmt.Println(seq)
	if seq < min+90000 {
		t.Error("should be bigger")
	}
}

func TestGenerateSeqBefore(t *testing.T) {
	sh := NewHandler(&blankLog{})
	sh.Add("test", "", min+90000, &serviceWithStop{})
	seq := sh.GenerateSeq("test", Before)
	fmt.Println(seq)
	if seq > min+90000 {
		t.Error("should be smaller")
	}
}
