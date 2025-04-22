import React from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'

import { App } from './App'
import './main.css'
import { ApiProvider } from './lib/api/ApiProvider'
import { ApiClient, DEFAULT_API_ENDPOINT } from './lib/api/api-client'

const apiClient = new ApiClient(DEFAULT_API_ENDPOINT)

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
