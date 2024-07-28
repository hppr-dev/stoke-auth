# Stoke Admin Console

This is the UI code for the stoke admin console.
It is embeded and served directly from the stoke executable.
npm 22 or later is required.
Run `npm install && npm run build --emptyOutDir` to build UI assets to be embedded into the stoke go executable.

Run `npm run dev` to start the development server when developing the admin UI.
This allows the use of hot reload without needing to rebuild the go executable.
