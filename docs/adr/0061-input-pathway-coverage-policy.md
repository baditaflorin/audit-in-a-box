# 0061 - Input Pathway Coverage Policy

## Status

Accepted

## Context

Real users bring manifests as files, pasted text, GitHub URLs, rendered HTML, clipboard contents, and sometimes multiple files.

## Decision

Support text manifest paste, single upload, drag-drop, multi-file upload, clipboard read, GitHub/raw URL fetch when browser CORS allows it, all supported samples, imported state, hash state, and autosaved drafts.

Image/OCR parsing and folder import stay out of scope. URL fetch failures must explain CORS and tell users to paste rendered HTML or raw file contents.

## Consequences

The UI becomes useful for normal repo workflows without adding backend engine behavior.

## Alternatives Considered

A backend URL-fetch proxy was rejected for Phase 3 because it expands server trust and SSRF concerns.
