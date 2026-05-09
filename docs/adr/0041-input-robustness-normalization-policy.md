# 0041 - Input Robustness And Normalization Policy

## Status

Accepted

## Context

Users paste manifests with BOMs, CRLF newlines, NBSPs, smart quotes, truncated content, HTML page chrome, and lockfile formats.

## Decision

Normalize UTF-8 BOM, CRLF, NBSP, and common smart quotes at the API boundary. Keep original byte length and normalized content hash in report provenance. Detect truncated JSON and unsupported structured formats as recoverable domain errors.

## Consequences

Parsers receive consistent text. Error messages can explain whether the app failed to understand the input or the input is incomplete.

## Alternatives Considered

Letting every parser normalize independently was rejected because it scatters policy and makes failures inconsistent.
