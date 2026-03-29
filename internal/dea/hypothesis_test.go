package dea

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHypothesisNodeDefaultStatus(t *testing.T) {
	h := &HypothesisNode{
		ID:   "H-001",
		Kind: KindRequiredField,
	}
	assert.Equal(t, HypothesisStatus(""), h.Status)
}

func TestHypothesisNodeResolve(t *testing.T) {
	h := &HypothesisNode{ID: "H-001", Kind: KindRequiredField}
	ev := &Evidence{ActualStatus: 400}
	h.Resolve(ev, true)
	assert.Equal(t, StatusConfirmed, h.Status)
	assert.Equal(t, ev, h.Evidence)
}

func TestHypothesisNodeRefute(t *testing.T) {
	h := &HypothesisNode{ID: "H-001", Kind: KindRequiredField}
	ev := &Evidence{ActualStatus: 201}
	h.Resolve(ev, false)
	assert.Equal(t, StatusRefuted, h.Status)
}
