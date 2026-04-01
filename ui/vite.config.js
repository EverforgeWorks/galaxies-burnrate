/**
 * @file vite.config.js
 * @description Configuration for the Vite development server and build process.
 * Specifically tailored for routing through a Cloudflare Tunnel securely with HMR.
 */

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

/**
 * Configures and returns the Vite server settings.
 * Configures the WebSocket for Hot Module Replacement to route correctly 
 * through the Cloudflare edge via WSS to prevent connection drops.
 *
 * @returns {import('vite').UserConfig} The Vite configuration object.
 */
export default defineConfig({
  plugins: [vue()],
  server: {
    host: '127.0.0.1',
    port: 5173,
    strictPort: true,
    hmr: {
      host: 'dev.playburnrate.com',
      protocol: 'wss',
      clientPort: 443
    }
  }
})