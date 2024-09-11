import { useCallback, useEffect } from 'react'
import { NavLink } from 'react-router-dom'

import { Badge } from '../components/Badge'
import { FluentChevronRight24Regular } from '../components/icons/fluent-chevron-right-24-regular'
import { FluentArrowSortDown24Filled } from '../components/icons/fluent-sort-arrow-down-24-filled'
import { FluentArrowSortUp24Filled } from '../components/icons/fluent-sort-arrow-up-24-filled'
import { SimpleIconsRss } from '../components/icons/simple-icon-rss'
import { useFilter, useSort } from '../hooks'

interface Tag {
  label: string
  color: string
}

export function Dashboard(): JSX.Element {
  const tags: Tag[] = [
    { label: 'k8s', color: 'bg-blue-100' },
    { label: 'docker', color: 'bg-red-100' },
    { label: 'Pod', color: 'bg-orange-100' },
    { label: 'Job', color: 'bg-blue-100' },
    { label: 'ghcr', color: 'bg-blue-100' },
    { label: 'up-to-date', color: 'bg-green-100' },
    { label: 'outdated', color: 'bg-red-100' },
  ]

  const [filter, setFilter] = useFilter()

  const [sort, setSort] = useSort()

  const rowItems = [
    {
      image: 'home-assistant',
      imageUrl:
        'https://www.gravatar.com/avatar/461df105cc6cfcf386ebd5af85b925dc?s=120&r=g&d=404',
      current: '2024.4.4',
      new: '2024.8.3',
      tags: ['k8s', 'Pod', 'docker'],
    },
    {
      image: 'jacobalberty/unifi',
      imageUrl:
        'https://www.gravatar.com/avatar/461df105cc6cfcf386ebd5af85b925dc?s=120&r=g&d=404',
      current: 'v7',
      new: 'v8.4',
      tags: ['k8s', 'Pod', 'ghcr'],
    },
    {
      image: 'yooooomi/your_spotify_server',
      imageUrl:
        'https://www.gravatar.com/avatar/461df105cc6cfcf386ebd5af85b925dc?s=120&r=g&d=404',
      current: '1.11.0',
      new: '1.11.0',
      tags: ['k8s', 'Pod', 'up-to-date'],
    },
    {
      image: 'hashicorp/vault',
      imageUrl:
        'https://www.gravatar.com/avatar/461df105cc6cfcf386ebd5af85b925dc?s=120&r=g&d=404',
      current: '2024.4.4',
      new: '2024.8.3',
      tags: ['k8s', 'Pod', 'outdated'],
    },
  ]

  // If the filter excludes all tags, default to show all tags again
  useEffect(() => {
    if (filter.length === 0) {
      setFilter(tags.map((x) => x.label))
    }
  }, [filter, setFilter])

  const toggleTag = useCallback(
    (tag: string) => {
      setFilter((previous) => {
        // If all are selected, only select the clicked tag
        if (previous.length === tags.length) {
          return [tag]
        }

        // Filter the selection
        const selection = previous.includes(tag)
          ? previous.filter((x) => x !== tag)
          : [...previous, tag]

        return selection
      })
    },
    [filter, setFilter]
  )

  const toggleSort = useCallback(
    (name: string) => {
      setSort((current) => {
        if (current === `${name}_asc`) {
          return `${name}_desc`
        } else if (current === `${name}_desc`) {
          return `${name}_asc`
        } else {
          // Default
          return `${name}_desc`
        }
      })
    },
    [setSort]
  )

  return (
    <>
      <div className="flex flex-col items-center w-full py-[40px] px-[20px]">
        <div className="p-3 flex space-x-5">
          <div className="p-5 w-32 h-32 bg-blue-100 rounded-lg">
            <p className="text-md font-medium">Images</p>
            <p className="text-3xl font-bold">15</p>
          </div>
          <div className="p-5 w-32 h-32 bg-orange-100 rounded-lg">
            <p className="text-md font-medium">Outdated</p>
            <p className="text-3xl font-bold">15</p>
          </div>
          <div className="p-5 w-32 h-32 bg-purple-100 rounded-lg">
            <p className="text-md font-medium">Pods</p>
            <p className="text-3xl font-bold">71</p>
          </div>
        </div>
        <div className="relative mt-6">
          <div className="rounded-lg bg-white px-4 py-2 shadow">
            <table>
              <thead>
                <tr>
                  <th
                    scope="col"
                    colSpan={2}
                    className="text-nowrap text-center cursor-pointer pr-[24px]"
                    onClick={() => toggleSort('image')}
                  >
                    Image
                    <div className="inline-block relative py-[9px]">
                      {sort === 'image_asc' && (
                        <FluentArrowSortUp24Filled className="absolute top-0" />
                      )}
                      {sort === 'image_desc' && (
                        <FluentArrowSortDown24Filled className="absolute top-0" />
                      )}
                    </div>
                  </th>
                  <th
                    scope="col"
                    className="text-nowrap text-center cursor-pointer pr-[24px]"
                    onClick={() => toggleSort('version')}
                  >
                    Version
                    <div className="inline-block relative py-[9px]">
                      {sort === 'version_asc' && (
                        <FluentArrowSortUp24Filled className="absolute top-0" />
                      )}
                      {sort === 'version_desc' && (
                        <FluentArrowSortDown24Filled className="absolute top-0" />
                      )}
                    </div>
                  </th>
                  <th
                    scope="col"
                    className="text-nowrap text-center cursor-pointer pr-[24px]"
                    onClick={() => toggleSort('new_version')}
                  >
                    New version
                    <div className="inline-block relative py-[9px]">
                      {sort === 'new_version_asc' && (
                        <FluentArrowSortUp24Filled className="absolute top-0" />
                      )}
                      {sort === 'new_version_desc' && (
                        <FluentArrowSortDown24Filled className="absolute top-0" />
                      )}
                    </div>
                  </th>
                  <th scope="col" className="text-nowrap text-center">
                    Tags
                  </th>
                  <th scope="col"></th>
                </tr>
              </thead>
              <tbody>
                {rowItems.map((x) => (
                  <tr key={x.image} className="">
                    <td>
                      <img className="w-10 rounded" src={x.imageUrl} />
                    </td>
                    <td className="pr-[24px]">{x.image}</td>
                    <td className="text-end pr-[24px]">{x.current}</td>
                    <td className="text-end pr-[24px]">{x.new}</td>
                    <td className="flex flex-wrap">
                      {x.tags.map((x) => (
                        <Badge
                          key={x}
                          label={x}
                          color={
                            tags.find((y) => y.label === x)?.color ||
                            'bg-blue-100'
                          }
                          className="cursor-pointer"
                          onClick={() => setFilter([x])}
                        />
                      ))}
                    </td>
                    <td>
                      <NavLink
                        to={`/image?name=${x.image}&version=${x.current}`}
                      >
                        <FluentChevronRight24Regular />
                      </NavLink>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
            <div className="mt-4">
              <p className="text-sm">Showing 50 of 50 entries</p>
            </div>
          </div>
          <div className="absolute left-full top-0 px-2">
            <div className="rounded-lg bg-white p-4 w-64 shadow">
              <p>Tags</p>
              <div className="flex flex-wrap mt-2">
                {tags.map((x) => (
                  <Badge
                    key={x.label}
                    label={x.label}
                    color={x.color}
                    disabled={!filter.includes(x.label)}
                    onClick={() => toggleTag(x.label)}
                    className="cursor-pointer"
                  />
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}
