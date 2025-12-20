import { expect, test } from 'vitest'
import { fullVersion, name, version } from './oci'

test.each([
  ['mongo', 'mongo', 'latest', 'latest'],
  ['mongo:4', 'mongo', '4', '4'],
  ['ghcr.io/mongo/mongo', 'ghcr.io/mongo/mongo', 'latest', 'latest'],
  ['ghcr.io/mongo/mongo:4', 'ghcr.io/mongo/mongo', '4', '4'],
  [
    'mongo@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
    'mongo',
    'sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
    'sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
  ],
  [
    'ghcr.io/mongo/mongo@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
    'ghcr.io/mongo/mongo',
    'sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
    'sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
  ],
  [
    'ghcr.io/mongo/mongo:6.0.0@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
    'ghcr.io/mongo/mongo',
    '6.0.0',
    '6.0.0@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
  ],
  [
    '[1234:5678:9012:1::1]:8080/test/image:1.0.0',
    '[1234:5678:9012:1::1]:8080/test/image',
    '1.0.0',
    '1.0.0',
  ],
  [
    '192.168.1.2:8080/test/image:1.0.0',
    '192.168.1.2:8080/test/image',
    '1.0.0',
    '1.0.0',
  ],
  [
    '[1234:5678:9012:1::1]/test/image:1.0.0',
    '[1234:5678:9012:1::1]/test/image',
    '1.0.0',
    '1.0.0',
  ],
  ['192.168.1.2/test/image:1.0.0', '192.168.1.2/test/image', '1.0.0', '1.0.0'],
])('oci name and version from %s', (reference, expectedName, expectedVersion, expectedFullVersion) => {
  expect(name(reference)).toBe(expectedName)
  expect(version(reference)).toBe(expectedVersion)
  expect(fullVersion(reference)).toBe(expectedFullVersion)
})
