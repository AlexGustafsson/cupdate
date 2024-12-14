// TODO: Create a tool that extracts this data from a running instance by
// getting the image page and then recording all the responses

class Recorder {
  /** @param {string} basePath */
  constructor(basePath) {
    this.data = {}
    this.basePath = basePath
  }

  /** @returns {Promise<{body?: BodyInit, options?: RequestInit}>} */
  async recordResponse(path) {
    const url = new URL(`${this.basePath}${path}`)

    const res = await fetch(url)

    const body = await res.text()

    /** @type {ResponseInit} */
    const options = {
      status: res.status,
      statusText: res.statusText,
      headers: Object.fromEntries(res.headers.entries()),
    }

    delete options.headers?.['access-control-allow-origin']
    delete options.headers?.['date']
    delete options.headers?.['transfer-encoding']

    this.data[`GET ${url.pathname}${url.search}`] = {
      body,
      options,
    }

    return {
      body,
      options,
    }
  }
}

async function main() {
  const recorder = new Recorder('http://localhost:8080')

  const imagesPage = await recorder.recordResponse('/api/v1/images')
  const images = JSON.parse(imagesPage.body)

  await recorder.recordResponse('/api/v1/tags')

  for (const image of images.images) {
    await recorder.recordResponse(
      `/api/v1/image?reference=${encodeURIComponent(image.reference)}`
    )
    await recorder.recordResponse(
      `/api/v1/image/release-notes?reference=${encodeURIComponent(image.reference)}`
    )
    await recorder.recordResponse(
      `/api/v1/image/graph?reference=${encodeURIComponent(image.reference)}`
    )
    await recorder.recordResponse(
      `/api/v1/image/description?reference=${encodeURIComponent(image.reference)}`
    )
  }

  console.log(JSON.stringify(recorder.data, null, 2))
}

main()
