import React from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'

import { App } from './App'
import './main.css'
import { ApiProvider } from './lib/api/ApiProvider'
import { ApiClient, DEFAULT_API_ENDPOINT } from './lib/api/api-client'
import { DemoApiClient } from './lib/api/demo-api-client'

const apiClient = import.meta.env.VITE_DEMO_MODE
  ? new DemoApiClient()
  : new ApiClient(DEFAULT_API_ENDPOINT)

// biome-ignore lint/suspicious/noExplicitAny: This is a hack to expose things
const _w = window as any
_w.cupdate = {
  dump:
    import.meta.env.VITE_DEMO_MODE === 'true'
      ? () => console.error('Not available in demo mode')
      : () => new ApiClient(DEFAULT_API_ENDPOINT).dump(),
}

const root = document.getElementById('root')
if (root) {
  createRoot(root).render(
    <React.StrictMode>
      <BrowserRouter>
        <ApiProvider client={apiClient}>
          <App />
        </ApiProvider>
      </BrowserRouter>
    </React.StrictMode>
  )
}
