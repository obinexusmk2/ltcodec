package codec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/obinexusmk2/ltcodec/pkg/format"
	"github.com/obinexusmk2/ltcodec/pkg/transform"
)

// DecoderConfig holds options for the decoder subcommand.
type DecoderConfig struct {
	InputPath  string // .lt archive
	OutputPath string // recovered file (default: <original_name>)
	Verbose    bool
}

// Decode reads a .lt archive, reverses the isomorphic XOR transform,
// verifies trident integrity, and writes the original payload to disk.
//
// CLI: ltcodec decoder -input <file.lt> [-output <file>]
func Decode(cfg DecoderConfig) error {
	// ── Read .lt archive ──────────────────────────────────────────────────
	ltBytes, err := os.ReadFile(cfg.InputPath)
	if err != nil {
		return fmt.Errorf("decoder: read %q: %w", cfg.InputPath, err)
	}

	// ── Open and verify archive ───────────────────────────────────────────
	meta, payload, idx, err := format.Open(ltBytes)
	if err != nil {
		return fmt.Errorf("decoder: open archive: %w", err)
	}

	if cfg.Verbose {
		fmt.Printf("[DECODER] input:      %s (%d bytes)\n", cfg.InputPath, len(ltBytes))
		fmt.Printf("[DECODER] uuid:       %s\n", meta.UUID)
		fmt.Printf("[DECODER] original:   %s (%s)\n", meta.OriginalName, meta.ContentType)
		fmt.Printf("[DECODER] stateless:  %v\n", meta.Stateless)
		fmt.Printf("[DECODER] sections:   %d\n", len(idx))
	}

	// ── Trident verification (read-only integrity check) ─────────────────
	result := transform.RunTrident(payload)

	if cfg.Verbose {
		fmt.Printf("[DECODER] trident:    state=%s Δ=%.4f verified=%v\n",
			result.State, result.Discriminant, result.Verified)
	}

	// ── Reverse isomorphic XOR transform (self-inverse) ───────────────────
	key := transform.DeriveKey(meta.UUID)
	// The stored payload is the pure XOR-encoded data (no receive flip).
	// XOR is its own inverse, so Decode(Encode(raw, key), key) == raw.
	recovered := transform.Decode(payload, key)

	// ── Resolve output path ───────────────────────────────────────────────
	if cfg.OutputPath == "" {
		cfg.OutputPath = resolveOutputPath(cfg.InputPath, meta.OriginalName)
	}

	// ── Write recovered file ──────────────────────────────────────────────
	if err := os.WriteFile(cfg.OutputPath, recovered, 0644); err != nil {
		return fmt.Errorf("decoder: write %q: %w", cfg.OutputPath, err)
	}

	fmt.Printf("[DECODER] output:     %s (%d bytes)\n", cfg.OutputPath, len(recovered))
	fmt.Printf("[DECODER] state:      %s | polarity: %c\n",
		result.State, result.Polarity)

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────────────────────────────────────

// resolveOutputPath builds a sensible output path for the decoded file.
func resolveOutputPath(ltPath, originalName string) string {
	dir := filepath.Dir(ltPath)

	if originalName != "" && originalName != "." {
		return filepath.Join(dir, "decoded_"+originalName)
	}

	// Fallback: strip .lt extension
	base := filepath.Base(ltPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return filepath.Join(dir, name+"_decoded")
}
