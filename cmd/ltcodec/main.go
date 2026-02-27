// ltcodec — Linkable Then Codec
// OBINexus Computing | github.com/obinexusmk2/ltcodec
//
// Constitutional pipeline:
//   riftlang.exe → .so.a → rift.exe → gosilang → ltcodec → nsigii
// Orchestration:
//   nlink → polybuild
//
// Usage:
//   ltcodec encode -input <file> [-type <payload_type>] [-output <file>]
//   ltcodec decode -input <file.ltf>
//   ltcodec state
//   ltcodec replay -time <RFC3339>
//   ltcodec replay -seq  <n>
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/obinexusmk2/ltcodec/pkg/codec"
	"github.com/obinexusmk2/ltcodec/pkg/gps"
)

const banner = `
╔══════════════════════════════════════════════╗
║  ltcodec v1.0.0 — OBINexus Computing        ║
║  Linkable Then Format (LTF) Codec            ║
║  Pipeline: riftlang → nlink → polybuild      ║
╚══════════════════════════════════════════════╝
`

func main() {
	fmt.Print(banner)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Default GPS resolver: IP geolocation with graceful fallback
	resolver := func() (gps.Coordinate, error) {
		log.Println("[ltcodec] Resolving location via IP geolocation…")
		coord, err := gps.FromIP()
		if err != nil {
			log.Printf("[ltcodec] IP geolocation failed: %v — using zero coordinate", err)
			return gps.Coordinate{Source: "unavailable"}, nil
		}
		log.Printf("[ltcodec] Location resolved: %.4f, %.4f [%s]",
			coord.Latitude, coord.Longitude, coord.Source)
		return coord, nil
	}

	c, err := codec.NewCodec(resolver)
	if err != nil {
		log.Fatalf("[ltcodec] Codec init failed: %v", err)
	}

	switch os.Args[1] {
	case "encode":
		runEncode(c)
	case "decode":
		runDecode(c)
	case "state":
		runState(c)
	case "replay":
		runReplay(c)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

// runEncode — LINK → THEN → EXECUTE: encode a file into an LTF packet
func runEncode(c *codec.Codec) {
	fs := flag.NewFlagSet("encode", flag.ExitOnError)
	inputFile  := fs.String("input",  "",    "input file to encode")
	payloadType := fs.String("type",  "raw", "payload type: raw | legal | nsigii")
	outputFile := fs.String("output", "",    "output .ltf file (default: <input>.ltf)")
	fs.Parse(os.Args[2:])

	if *inputFile == "" {
		log.Fatal("encode: -input required")
	}

	data, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("encode: read failed: %v", err)
	}

	packet, err := c.Encode(*payloadType, data)
	if err != nil {
		log.Fatalf("encode: %v", err)
	}

	out := *outputFile
	if out == "" {
		out = *inputFile + ".ltf"
	}

	packetJSON, err := packet.ToJSON()
	if err != nil {
		log.Fatalf("encode: serialise failed: %v", err)
	}
	if err := os.WriteFile(out, packetJSON, 0644); err != nil {
		log.Fatalf("encode: write failed: %v", err)
	}

	fmt.Printf("\n✓ Encoded → %s\n", out)
	fmt.Printf("  Spacetime fingerprint : %s\n", packet.State.Fingerprint)
	fmt.Printf("  Packet hash           : %s\n", packet.PacketHash)
	fmt.Printf("  Sequence              : %d\n", packet.State.Sequence)
	fmt.Printf("  Location              : %s\n", packet.State.Position.String())
	fmt.Printf("  Hardware (MAC)        : %s\n", packet.State.Hardware.MAC)
}

// runDecode — verify and inspect an LTF packet
func runDecode(c *codec.Codec) {
	fs := flag.NewFlagSet("decode", flag.ExitOnError)
	inputFile := fs.String("input", "", "input .ltf file")
	fs.Parse(os.Args[2:])

	if *inputFile == "" {
		log.Fatal("decode: -input required")
	}

	data, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("decode: read failed: %v", err)
	}

	fmt.Printf("\n✓ Packet loaded from %s\n", *inputFile)
	fmt.Printf("  Raw size : %d bytes\n", len(data))
	fmt.Printf("  Use 'state' or 'replay' to inspect the spacetime anchor.\n")
	_ = c // codec available for full decode in a future iteration
}

// runState — capture and print the current spacetime state
func runState(c *codec.Codec) {
	coord, _ := gps.FromIP()
	state, err := c.Session.Capture(coord)
	if err != nil {
		log.Fatalf("state: %v", err)
	}
	stateJSON, _ := state.ToJSON()
	fmt.Printf("\nCurrent spacetime state:\n%s\n", string(stateJSON))
}

// runReplay — replay a historical state by sequence or time
func runReplay(c *codec.Codec) {
	fs := flag.NewFlagSet("replay", flag.ExitOnError)
	timeStr := fs.String("time", "", "RFC3339 timestamp  e.g. 2026-02-27T07:10:00Z")
	seq     := fs.Uint64("seq",  0,  "monotonic sequence number to replay")
	fs.Parse(os.Args[2:])

	if *seq > 0 {
		state, ok := c.Session.Replay(uint64(*seq))
		if !ok {
			log.Fatalf("replay: sequence %d not found in session", *seq)
		}
		stateJSON, _ := state.ToJSON()
		fmt.Printf("\nReplayed state [seq=%d]:\n%s\n", *seq, string(stateJSON))
		return
	}

	if *timeStr != "" {
		t, err := time.Parse(time.RFC3339, *timeStr)
		if err != nil {
			log.Fatalf("replay: invalid time format — use RFC3339, e.g. 2026-02-27T07:10:00Z: %v", err)
		}
		state, ok := c.ReplayAt(t)
		if !ok {
			log.Fatal("replay: no states captured in this session yet")
		}
		stateJSON, _ := state.ToJSON()
		fmt.Printf("\nReplayed state [t=%s]:\n%s\n", *timeStr, string(stateJSON))
		return
	}

	log.Fatal("replay: specify -time <RFC3339> or -seq <n>")
}

func printUsage() {
	fmt.Print(`
Usage: ltcodec <command> [flags]

Commands:
  encode   -input <file> [-type <type>] [-output <file>]
  decode   -input <file.ltf>
  state                               Print current spacetime state
  replay   -time <RFC3339>            Replay state at given time
           -seq  <n>                  Replay state at sequence number

Payload types:
  raw      Raw binary data (default)
  legal    Legal document (Care Act / HRA claims)
  nsigii   NSIGII video codec output

Constitutional pipeline:
  riftlang.exe → .so.a → rift.exe → gosilang → ltcodec → nsigii
Orchestration:
  nlink → polybuild
`)
}
