# Phase 2 State Taxonomy

## Backend Audit States

- `accepted`: request passed boundary validation.
- `normalizing`: input bytes are decoded and normalized.
- `classified`: manifest shape and ecosystem inferred.
- `parsed`: first-party parser produced dependency evidence.
- `scanner_partial`: at least one scanner returned evidence and another failed or timed out.
- `scanner_complete`: all available scanners completed.
- `scored`: risk score and confidence calculated.
- `recoverable_error`: user can fix input without losing work.
- `fatal_error`: backend cannot continue, but response preserves actionable detail.
- `cancelled`: request context was cancelled and no partial report is committed.

## Frontend States

- `idle-empty`: no report yet.
- `editing`: user has content or HTML paste in the form.
- `checking-tools`: backend tool status is loading.
- `running`: audit request in flight.
- `cancelling`: abort requested; UI keeps previous report.
- `report-ready`: latest report is coherent and inspectable.
- `recoverable-error`: input remains editable and error includes next step.
- `backend-unreachable`: user can edit backend URL or keep working with last report.

Every state has an exit: edit input, retry, cancel, change backend URL, or keep prior report.
