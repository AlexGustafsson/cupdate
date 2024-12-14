import React from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'

import { App } from './App'
import './main.css'

const root = document.getElementById('root')
if (root) {
  if (import.meta.env['VITE_DEMO_MODE'] === 'true') {
    console.log('Starting Cupdate in demo mode')
    navigator.serviceWorker
      .getRegistrations()
      .then((registrations) =>
        Promise.all(registrations.map((x) => x.unregister()))
      )
      .then(() =>
        navigator.serviceWorker.register('/demo-service-worker.js', {
          scope: '/',
        })
      )
      .then(() => {
        createRoot(root).render(
          <React.StrictMode>
            <BrowserRouter>
              <App />
            </BrowserRouter>
          </React.StrictMode>
        )
      })
      .catch((error) => {
        createRoot(root).render(
          <React.StrictMode>
            <p>Service workers are required for the demo to work.</p>
            <pre>{error}</pre>
          </React.StrictMode>
        )
      })
  } else {
    createRoot(root).render(
      <React.StrictMode>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </React.StrictMode>
    )
  }
}
