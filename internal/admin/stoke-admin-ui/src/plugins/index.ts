/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 * Router is created in main from API base and passed in.
 */

// Plugins
import vuetify from './vuetify'
import pinia from '../stores'

// Types
import type { App } from 'vue'
import type { Router } from 'vue-router'

export function registerPlugins (app: App, router: Router) {
  app
    .use(vuetify)
    .use(router)
    .use(pinia)
}
