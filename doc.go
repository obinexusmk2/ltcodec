// Package ltcodec implements the Linkable Then Format (LTF) codec system.
//
// # Overview
//
// ltcodec is a spacetime-anchored encoding system developed by OBINexus Computing.
// Every payload encoded by ltcodec is bound to a non-reproducible identity:
//
//	Location (GPS/IP) + Hardware (MAC/UUID) + Time (delta-T) = unique state
//
// This state cannot be forged: it is captured at encoding time and embedded
// into every LTF packet as a cryptographic fingerprint.
//
// # Constitutional Pipeline
//
//	riftlang.exe → .so.a → rift.exe → gosilang → ltcodec → nsigii
//
// # Orchestration
//
//	nlink → polybuild
//
// # LTF Phases
//
// LINK — resolve hardware identity and GPS coordinate, produce a SpacetimeState.
//
// THEN — bind the payload to the state, compute a packet hash.
//
// EXECUTE — emit a signed, timestamped, location-anchored LTFPacket.
//
// # Packages
//
//   - [github.com/obinexusmk2/ltcodec/pkg/gps]       — GPS coordinate capture (real + IP fallback)
//   - [github.com/obinexusmk2/ltcodec/pkg/identity]  — hardware identity (MAC, UUID v6, dual-IP)
//   - [github.com/obinexusmk2/ltcodec/pkg/spacetime] — spacetime state and session replay
//   - [github.com/obinexusmk2/ltcodec/pkg/codec]     — LTF encode/decode/verify
//
// # Quick Start
//
//	c, err := codec.NewCodec(nil) // nil uses IP geolocation
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	packet, err := c.Encode("legal", data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Println(packet.State.Fingerprint) // SHA-256 spacetime proof
//	fmt.Println(packet.Verify())          // true
//
// # Relay-Replay
//
// Every state captured during a session is stored and can be replayed at any
// delta-T, providing accountability through time-ordered evidence:
//
//	state, ok := c.ReplayAt(time.Parse(time.RFC3339, "2026-02-27T17:22:38Z"))
//
// # Command-line Tool
//
// Install the CLI:
//
//	go install github.com/obinexusmk2/ltcodec/cmd/ltcodec@latest
//
// Encode a file:
//
//	ltcodec encode -input document.pdf -type legal -output document.ltf
//
// Inspect spacetime state:
//
//	ltcodec state
package ltcodec
