# Changelog — github.com/obinexusmk2/ltcodec

All notable changes to this module are documented here.
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [v1.0.0] — 2026-02-27

### Added

**Core LTF pipeline — Linkable Then Format**

- `pkg/gps` — GPS coordinate type with Haversine distance, IP-geolocation fallback (`ip-api.com`), and manual coordinate override.
- `pkg/identity` — Hardware identity capture: physical MAC address, dual IPv4/IPv6, UUID v6 time-ordered identifier derived from MAC + nanosecond timestamp.
- `pkg/spacetime` — Spacetime state model binding WHO (hardware) + WHERE (GPS) + WHEN (delta-T) into a SHA-256 fingerprint. Session-based relay-replay mechanism for time-ordered state accountability.
- `pkg/codec` — `Codec` type implementing the LINK → THEN → EXECUTE pipeline. Produces `LTFPacket`: a self-verifying, location-anchored, hash-signed payload unit.
- `cmd/ltcodec` — Command-line interface: `encode`, `decode`, `state`, `replay` subcommands.

**Payload types supported**

- `raw` — arbitrary binary data
- `legal` — legal documents (Care Act 2014, Human Rights Act claims)
- `nsigii` — NSIGII video codec output

**Verified working**

- Location resolved: 53.8073, −2.7591 [IP geolocation]
- Care Act 2014 PDF encoded as `legal` LTF packet (`careact2014.nsigii`)
- Spacetime fingerprint: `900df61d3e10b796…`
- Packet hash: `bf04462c4ace99da…`

### Constitutional pipeline

```
riftlang.exe → .so.a → rift.exe → gosilang → ltcodec → nsigii
```

### Orchestration

```
nlink → polybuild
```

---

## Upcoming

- `v1.1.0` — NSIGII container writer integration (Trident channel architecture)
- `v1.2.0` — SMTP relay for LTF packets (legal evidence transmission)
- `v2.0.0` — CGo bridge for riftlang `.so.a` linkage
