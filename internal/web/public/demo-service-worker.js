self.addEventListener('install', () => {
  console.log('Demo service worker installed')
})

self.addEventListener('activate', () => {
  console.log('Demo service worker activated')
})

const tagsResponse = `["replica set", "outdated", "deployment", "github", "up-to-date", "daemon set", "vulnerable", "minor", "job", "stateful set", "patch", "cron job"]`

/** @type {{matcher: RegExp, response: {body?: BodyInit, options?: ResponseInit}}[]} */
const data = [
  {
    matcher: /^GET \/api\/v1\/images.*$/,
    response: {
      body: imagesResponse,
      options: {
        status: 200,
        statusText: 'OK',
        headers: {},
      },
    },
  },
  {
    matcher: /^GET \/api\/v1\/tags$/,
    response: {
      body: tagsResponse,
      options: {
        status: 200,
        statusText: 'OK',
        headers: {},
      },
    },
  },
]

self.addEventListener('fetch', (event) => {
  const url = new URL(event.request.url)
  if (url.host === self.location.host && url.pathname.startsWith('/api/v1/')) {
    const key = `${event.request.method} ${url.pathname}${url.search}`
    for (const { matcher, response } of data) {
      if (matcher.test(key)) {
        event.respondWith(new Response(response.body, response.options))
        return
      }
    }
  }

  event.respondWith(
    fetch(event.request).catch(function () {
      return caches.match('/offline')
    })
  )
})
