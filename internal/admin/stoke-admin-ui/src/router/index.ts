/**
 * router/index.ts
 *
 * Automatic routes for `./src/pages/*.vue`
 * Router base is set from the current URL in main.ts (path up to and including /admin/).
 */

// Composables
import { createRouter, createWebHistory } from 'vue-router/auto'
import { routes } from 'vue-router/auto-routes'
import { setupLayouts } from 'virtual:generated-layouts'
import type { Router } from 'vue-router'

/**
 * Create the router with a history base (e.g. "/admin/" or "/auth/admin/").
 * Called from main.ts with base derived from window.location.pathname.
 */
export function createAppRouter(base: string): Router {
  const historyBase = base.endsWith('/') ? base : base + '/'
  return createRouter({
    routes: setupLayouts(routes),
    history: createWebHistory(historyBase),
  })
}

/** Fallback base when router is imported before bootstrap (e.g. tests): derive from current location. */
function getDefaultBase(): string {
  if (typeof window === 'undefined') return '/admin/'
  const match = window.location.pathname.match(/^(.+)\/admin\/?/)
  return match ? match[1] + '/admin/' : '/admin/'
}

const router = createRouter({
  routes: setupLayouts(routes),
  history: createWebHistory(getDefaultBase()),
})

export default router
