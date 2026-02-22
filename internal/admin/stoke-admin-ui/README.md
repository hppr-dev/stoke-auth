# Stoke Admin Console

This is the UI code for the stoke admin console.
It is embedded and served directly from the stoke executable.
The admin console supports login with username and password (local) and with configured OIDC providers; which options are shown is determined by the `/api/available_providers` endpoint.
When a change is denied by policy (e.g. protected user, group or claim), the UI shows an error message.

npm 22 or later is required.
Run `npm install && npm run build --emptyOutDir` to build UI assets to be embedded into the stoke go executable.

Run `npm run dev` to start the development server when developing the admin UI.
This allows the use of hot reload without needing to rebuild the go executable.
