import type { WebPushSubscription } from './models'

// Polyfill until proper typing support exists
declare global {
  interface Uint8Array {
    toBase64?(options?: { alphabet?: 'base64url' }): string
  }
}

/** Creates a digest uniquely identifying a {WebPushSubscription} */
export async function webPushSubscriptionDigest(
  subscription: WebPushSubscription | PushSubscriptionJSON
): Promise<string> {
  // The content is the JSON-encoded subscription, with lexicographically sorted
  // properties
  const content = JSON.stringify({
    endpoint: subscription.endpoint,
    keys: { auth: subscription.keys?.auth, p256dh: subscription.keys?.p256dh },
  })

  const plaintext = new TextEncoder().encode(content)

  const digest = await crypto.subtle.digest('sha256', plaintext)

  const buffer = new Uint8Array(digest)
  if (!buffer.toBase64) {
    // I'm tired of implementing Base64-stuff in web browsers and now that the
    // new functionality is basically supported everywhere (except for Chrome),
    // let's just wait them out
    throw new Error('unsupported browser')
  }

  const string = buffer.toBase64({
    alphabet: 'base64url',
  })

  return `sha256-${string}`
}
