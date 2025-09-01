import { expect, test } from 'vitest'
import { type Purl, parsePurl } from './purl'

test.each([
  [
    'pkg:type/namespace/name@version?qualifiers#subpath',
    {
      scheme: 'pkg',
      type: 'type',
      namespace: 'namespace',
      name: 'name',
      version: 'version',
      qualifiers: { qualifiers: '' },
      subpath: 'subpath',
    },
  ],
  [
    'pkg:bitbucket/birkenfeld/pygments-main@244fd47e07d1014f0aed9c',
    {
      scheme: 'pkg',
      type: 'bitbucket',
      namespace: 'birkenfeld',
      name: 'pygments-main',
      version: '244fd47e07d1014f0aed9c',
    },
  ],
  [
    'pkg:deb/debian/curl@7.50.3-1?arch=i386&distro=jessie',
    {
      scheme: 'pkg',
      type: 'deb',
      namespace: 'debian',
      name: 'curl',
      version: '7.50.3-1',
      qualifiers: { arch: 'i386', distro: 'jessie' },
    },
  ],
  [
    'pkg:docker/cassandra@sha256:244fd47e07d1004f0aed9c',
    {
      scheme: 'pkg',
      type: 'docker',
      name: 'cassandra',
      version: 'sha256:244fd47e07d1004f0aed9c',
    },
  ],
  [
    'pkg:docker/customer/dockerimage@sha256:244fd47e07d1004f0aed9c?repository_url=gcr.io',
    {
      scheme: 'pkg',
      type: 'docker',
      namespace: 'customer',
      name: 'dockerimage',
      version: 'sha256:244fd47e07d1004f0aed9c',
      qualifiers: { repository_url: 'gcr.io' },
    },
  ],
  [
    'pkg:gem/jruby-launcher@1.1.2?platform=java',
    {
      scheme: 'pkg',
      type: 'gem',
      name: 'jruby-launcher',
      version: '1.1.2',
      qualifiers: { platform: 'java' },
    },
  ],
  [
    'pkg:gem/ruby-advisory-db-check@0.12.4',
    {
      scheme: 'pkg',
      type: 'gem',
      name: 'ruby-advisory-db-check',
      version: '0.12.4',
    },
  ],
  [
    'pkg:github/package-url/purl-spec@244fd47e07d1004f0aed9c',
    {
      scheme: 'pkg',
      type: 'github',
      namespace: 'package-url',
      name: 'purl-spec',
      version: '244fd47e07d1004f0aed9c',
    },
  ],
  [
    'pkg:golang/google.golang.org/genproto#googleapis/api/annotations',
    {
      scheme: 'pkg',
      type: 'golang',
      namespace: 'google.golang.org',
      name: 'genproto',
      subpath: 'googleapis/api/annotations',
    },
  ],
  [
    'pkg:maven/org.apache.xmlgraphics/batik-anim@1.9.1?packaging=sources',
    {
      scheme: 'pkg',
      type: 'maven',
      namespace: 'org.apache.xmlgraphics',
      name: 'batik-anim',
      version: '1.9.1',
      qualifiers: { packaging: 'sources' },
    },
  ],
  [
    'pkg:maven/org.apache.xmlgraphics/batik-anim@1.9.1?repository_url=repo.spring.io/release',
    {
      scheme: 'pkg',
      type: 'maven',
      namespace: 'org.apache.xmlgraphics',
      name: 'batik-anim',
      version: '1.9.1',
      qualifiers: { repository_url: 'repo.spring.io/release' },
    },
  ],
  [
    'pkg:npm/%40angular/animation@12.3.1',
    {
      scheme: 'pkg',
      type: 'npm',
      namespace: '@angular',
      name: 'animation',
      version: '12.3.1',
    },
  ],
  [
    'pkg:npm/foobar@12.3.1',
    {
      scheme: 'pkg',
      type: 'npm',
      name: 'foobar',
      version: '12.3.1',
    },
  ],
  [
    'pkg:nuget/EnterpriseLibrary.Common@6.0.1304',
    {
      scheme: 'pkg',
      type: 'nuget',
      name: 'EnterpriseLibrary.Common',
      version: '6.0.1304',
    },
  ],
  [
    'pkg:pypi/django@1.11.1',
    {
      scheme: 'pkg',
      type: 'pypi',
      name: 'django',
      version: '1.11.1',
    },
  ],
  [
    'pkg:rpm/fedora/curl@7.50.3-1.fc25?arch=i386&distro=fedora-25',
    {
      scheme: 'pkg',
      type: 'rpm',
      namespace: 'fedora',
      name: 'curl',
      version: '7.50.3-1.fc25',
      qualifiers: { arch: 'i386', distro: 'fedora-25' },
    },
  ],
  [
    'pkg:rpm/opensuse/curl@7.56.1-1.1.?arch=i386&distro=opensuse-tumbleweed',
    {
      scheme: 'pkg',
      type: 'rpm',
      namespace: 'opensuse',
      name: 'curl',
      version: '7.56.1-1.1.',
      qualifiers: { arch: 'i386', distro: 'opensuse-tumbleweed' },
    },
  ],
] as [string, Purl | null][])('parse purl %s', (purl, expected) => {
  expect(parsePurl(purl)).toEqual(expected)
})
