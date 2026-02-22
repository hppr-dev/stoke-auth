/**
 * main.ts
 *
 * Bootstraps the app with router base derived from the current URL so it works
 * at /admin/ or behind a proxy (e.g. /auth/admin/) without fetching config.
 */

// Plugins
import { registerPlugins } from '@/plugins'

// Components
import App from './App.vue'

// Composables
import { createApp } from 'vue'
import { createAppRouter } from './router'

/** Router history base: path up to and including /admin/ (e.g. /admin/ or /auth/admin/). */
function getRouterBase(): string {
  if (typeof window === 'undefined') return '/admin/'
  const pathname = window.location.pathname
  const match = pathname.match(/^(.+)\/admin\/?/)
  return match ? match[1] + '/admin/' : '/admin/'
}

const router = createAppRouter(getRouterBase())
const app = createApp(App)
registerPlugins(app, router)
app.mount('#app')
