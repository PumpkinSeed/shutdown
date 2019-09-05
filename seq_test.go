package shutdown

import (
	"testing"
)

func TestGenerateSeqAfter(t *testing.T) {
	sh := NewHandler(&blankLog{})
	sh.Add("test", "", Init, &serviceWithStop{})
	seqAfter := sh.GenerateSeq("test", After)
	seq := sh.debugSeq("test")
	if seq > seqAfter {
		t.Errorf("%d should be bigger %d", seqAfter, seq )
	}
}

func TestGenerateSeqBefore(t *testing.T) {
	sh := NewHandler(&blankLog{})
	sh.Add("test", "", Init, &serviceWithStop{})
	seqBefore := sh.GenerateSeq("test", Before)
	seq := sh.debugSeq("test")
	if seq < seqBefore {
		t.Errorf("%d should be smaller %d", seqBefore, seq )
	}
}

func TestGenerateSeqBeforeHeavy(t *testing.T) {
	sh := NewHandler(&blankLog{})
	sh.Add("test", "", Init, &serviceWithStop{})
	for i := 0; i<10000; i++ {
		seqBefore := sh.GenerateSeq("test", Before)
		seq := sh.debugSeq("test")
		if seq < seqBefore {
			t.Errorf("%d should be smaller %d", seqBefore, seq)
		}
	}
}

func TestGenerateSeqAfterHeavy(t *testing.T) {
	sh := NewHandler(&blankLog{})
	sh.Add("test", "", Init, &serviceWithStop{})
	for i := 0; i<10000; i++ {
		seqAfter := sh.GenerateSeq("test", After)
		seq := sh.debugSeq("test")
		if seq > seqAfter {
			t.Errorf("%d should be bigger %d", seqAfter, seq)
		}
	}
}
