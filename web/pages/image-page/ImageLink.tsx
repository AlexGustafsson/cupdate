import type { JSX, ReactNode } from 'react'
import { FluentBookQuestionMark24Filled } from '../../components/icons/fluent-book-question-mark-24-filled'
import { FluentLink24Filled } from '../../components/icons/fluent-link-24-filled'
import { Quay } from '../../components/icons/quay'
import { SimpleIconsDocker } from '../../components/icons/simple-icons-docker'
import { SimpleIconsGit } from '../../components/icons/simple-icons-git'
import { SimpleIconsGithub } from '../../components/icons/simple-icons-github'
import { SimpleIconsGitlab } from '../../components/icons/simple-icons-gitlab'

const titles: Record<string, string | undefined> = {
  github: "Project's page on GitHub",
  ghcr: "Project's package page on GHCR",
  gitlab: "Project's page on GitLab",
  docker: "Project's page on Docker Hub",
  quay: "Project's page on Quay.io",
  git: "Project's git page",
  docs: "Project's documentation",
  'oci-registry': "Project's OCI registry",
}

export function ImageLink({
  type,
  url,
}: {
  type: string
  url: string
}): JSX.Element {
  const { hostname } = new URL(url)

  const title = titles[type] || url
  let icon: ReactNode
  switch (type) {
    case 'github':
    case 'ghcr':
      icon = (
        <SimpleIconsGithub className="text-black dark:dark:text-[#dddddd]" />
      )
      break
    case 'gitlab':
      icon = <SimpleIconsGitlab className="text-orange-500" />
      break
    case 'docker':
      icon = <SimpleIconsDocker className="text-blue-500" />
      break
    case 'quay':
      icon = <Quay className="text-blue-700" />
      break
    case 'git':
      icon = <SimpleIconsGit className="text-orange-500" />
      break
    case 'svc':
      switch (hostname) {
        case 'github.com':
          return <ImageLink type="github" url={url} />
        case 'gitlab.com':
          return <ImageLink type="gitlab" url={url} />
        default:
          return <ImageLink type="git" url={url} />
      }
    case 'docs':
      icon = <FluentBookQuestionMark24Filled />
      break
    default:
      icon = <FluentLink24Filled />
  }

  return (
    <a
      title={title}
      href={url}
      target="_blank"
      rel="noreferrer"
      tabIndex={0}
      className="hover:opacity-90 focus:opacity-90"
    >
      {icon}
    </a>
  )
}
