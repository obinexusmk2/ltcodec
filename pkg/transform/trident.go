// Package transform implements the NSIGII isomorphic transformation and trident verification.
//
// Trident 3-channel verification architecture:
//	Channel 0 — TRANSMITTER  127.0.0.1  (1/3)  WRITE
//	Channel 1 — RECEIVER     127.0.0.2  (2/3)  READ
//	Channel 2 — VERIFIER     127.0.0.3  (3/3)  EXECUTE
//
// The Verifier uses Discriminant Flash: Δ = b² - 4ac
//	Δ > 0  → ORDER     (coherent signal)
//	Δ = 0  → CONSENSUS (flash point)
//	Δ < 0  → CHAOS     (repair needed)
//
// CRITICAL: All trident operations are READ-ONLY for verification.
// No mutation of payload data occurs here.
package transform

import (
	"math"
)

// TridentState represents the discriminant-derived state.
type TridentState int

const (
	StateOrder TridentState = iota // Δ > 0
	StateConsensus                  // Δ = 0
	StateChaos                      // Δ < 0
)

func (s TridentState) String() string {
	switch s {
	case StateOrder:
		return "ORDER"
	case StateConsensus:
		return "CONSENSUS"
	default:
		return "CHAOS"
	}
}

// RWXFlags mirrors Unix-style permissions.
const (
	RWXRead    uint8 = 0x04
	RWXWrite   uint8 = 0x02
	RWXExecute uint8 = 0x01
	RWXFull    uint8 = 0x07
)

// TridentResult carries verification outcome.
type TridentResult struct {
	State        TridentState
	RWXFlags     uint8
	WheelDeg     int
	Discriminant float64
	Verified     bool
	Polarity     byte
}

// RunTrident executes TRANSMIT → RECEIVE → VERIFY pipeline.
// All channels are READ-ONLY — data is never mutated.
func RunTrident(data []byte) TridentResult {
	// Channel 0: TRANSMITTER — pass-through
	transmitted := transmit(data)

	// Channel 1: RECEIVER — pass-through (NO mutation)
	received := receive(transmitted)

	// Channel 2: VERIFIER — discriminant check
	a, b, c := bipartiteConsensusParams(received)
	delta := b*b - 4*a*c

	var state TridentState
	var rwx uint8
	var wheelDeg int

	switch {
	case delta > 0:
		state = StateOrder
		rwx = RWXFull
		wheelDeg = 120
	case delta == 0:
		state = StateConsensus
		rwx = RWXFull
		wheelDeg = 240
	default:
		state = StateChaos
		rwx = RWXRead
		wheelDeg = 0
	}

	return TridentResult{
		State:        state,
		RWXFlags:     rwx,
		WheelDeg:     wheelDeg,
		Discriminant: delta,
		Verified:     state != StateChaos,
		Polarity:     PolaritySign(data),
	}
}

// DiscriminantState computes state classification only.
func DiscriminantState(data []byte) TridentState {
	a, b, c := bipartiteConsensusParams(data)
	delta := b*b - 4*a*c
	switch {
	case delta > 0:
		return StateOrder
	case delta == 0:
		return StateConsensus
	default:
		return StateChaos
	}
}

// Internal helpers — ALL READ-ONLY

func transmit(data []byte) []byte {
	out := make([]byte, len(data))
	copy(out, data)
	return out
}

// receive is READ-ONLY — pure copy, no bit manipulation
func receive(data []byte) []byte {
    out := make([]byte, len(data))
    copy(out, data)  // ← Pure copy, no mutation
    return out
}

// bipartiteConsensusParams maps bit density to discriminant B parameter.
// Returns A=1, B∈[0,4], C=1 where B = 4 * (setBits / totalBits)
func bipartiteConsensusParams(data []byte) (a, b, c float64) {
	if len(data) == 0 {
		return 1, 2, 1 // Neutral for empty
	}
	
	var setBits int
	for _, by := range data {
		for v := by; v != 0; v >>= 1 {
			setBits += int(v & 1)
		}
	}
	
	totalBits := len(data) * 8
	density := float64(setBits) / float64(totalBits)
	
	// Linear mapping: 0% ones → B=0 (CHAOS), 50% → B=2 (CONSENSUS), 100% → B=4 (ORDER)
	B := density * 4.0
	return 1.0, B, 1.0
}
