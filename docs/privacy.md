# Privacy

The GitHub Pages frontend does not include analytics.

Manifests are sent only to the backend URL configured by the user in the UI. If the default value is used, the browser calls:

http://localhost:25342

The backend stores temporary runtime artifacts under `WORK_DIR` on the backend host. It logs request metadata and scanner outcomes but avoids logging full manifest contents.

No secrets are required in the frontend. Optional local LLM configuration belongs in backend environment variables.
