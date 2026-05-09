# Phase 3 Controls Audit

Date: 2026-05-10

| Control | Status before Phase 3 | Evidence | Decision |
| --- | --- | --- | --- |
| Backend URL input | Works partially | Value updates in memory and persists only when audit form submits. | Finish immediate persistence |
| Refresh backend tool status | Works fully | Calls TanStack query refetch. | Keep |
| Tool status chips | Works partially | Reflect status but are not controls and lack failure next-step guidance. | Keep with clearer text elsewhere |
| Sample buttons | Works partially | Load three samples only, while supported formats include lockfiles and pyproject. | Finish |
| File upload | Works partially | Single text file only; no batch, no state import. | Finish |
| Manifest content textarea | Works fully | User can paste/edit manifest text. | Keep |
| HTML paste textarea | Works partially | Requires user to know where HTML goes; no clipboard/URL helper. | Finish |
| Extract manifest | Works fully for supported pasted HTML | Calls `/api/v1/scrape`. | Keep |
| Error banner | Works partially | Shows error text; no dismissal or domain-specific action controls. | Keep |
| Run audit | Works fully for single manifest | Posts current content. | Keep |
| Cancel audit | Works partially | Aborts browser request; backend process may continue until request context/tool timeout. | Keep, document |
| Star on GitHub | Works fully | Link present and correct. | Keep |
| PayPal | Works fully | Link present and correct. | Keep |
| Debug view | Works partially | `?debug=1` works but is not discoverable. | Keep hidden, document |
| Report sections | Works partially | Summaries and tables display, but no export actions. | Finish |

No production UI stubs were found. The main problem is absent output/action controls after a report exists.
