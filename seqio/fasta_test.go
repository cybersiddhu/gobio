package seqio

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMultiFasta(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Error("Could not get current direcotry")
	}

	mfasta := filepath.Join(dir, "multi.fa")
	reader := NewFastaReader(mfasta)
	for i := 0; i <= 3; i++ {
		if !reader.HasEntry() {
			t.Error("Did not get expected iteration")
		}
		entry := reader.NextEntry()
		if entry == nil {
			t.Error("Did not get expected iteration")
		}

	}

}
