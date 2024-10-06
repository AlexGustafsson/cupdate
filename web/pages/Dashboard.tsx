import { useCallback } from 'react'
import { NavLink, useSearchParams } from 'react-router-dom'

import { useImages, usePagination, useTags } from '../api'
import { Badge } from '../components/Badge'
import { FluentChevronRight24Regular } from '../components/icons/fluent-chevron-right-24-regular'
import { FluentArrowSortDown24Filled } from '../components/icons/fluent-sort-arrow-down-24-filled'
import { FluentArrowSortUp24Filled } from '../components/icons/fluent-sort-arrow-up-24-filled'
import { SimpleIconsOci } from '../components/icons/simple-icons-oci'
import { useFilter, useSort } from '../hooks'
import { name, version } from '../oci'

export function Dashboard(): JSX.Element {
  const [filter, setFilter] = useFilter()

  const [sortProperty, setSortProperty, sortOrder, setSortOrder] = useSort()

  const [searchParams, _] = useSearchParams()

  const images = useImages({
    tags: filter,
    sort: sortProperty,
    order: sortOrder,
    page: searchParams.get('page') ? Number(searchParams.get('page')) : 0,
    limit: 30,
  })

  const pages = usePagination(
    images.status === 'resolved' ? images.value : undefined
  )

  const tags = useTags()

  // // If the filter excludes all tags, default to show all tags again
  // useEffect(() => {
  //   if (filter.length === 0) {
  //     setFilter(tags.map((x) => x.label))
  //   }
  // }, [filter, setFilter])

  const toggleTag = useCallback(
    (tag: string) => {
      setFilter((previous) => {
        // If all are selected, only select the clicked tag
        // if (previous.length === tags.length) {
        //   return [tag]
        // }

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
      if (sortProperty === name) {
        setSortOrder((current) => (current === 'asc' ? 'desc' : 'asc'))
      } else {
        setSortOrder('desc')
        setSortProperty(name)
      }
    },
    [sortProperty, sortOrder, setSortProperty, setSortOrder]
  )

  if (images.status !== 'resolved' || tags.status !== 'resolved') {
    return <p>Loading</p>
  }

  return (
    <>
      <div className="flex flex-col items-center w-full py-[40px] px-[20px]">
        {/* Header with summary */}
        <div className="p-3 flex space-x-5">
          {
            <div className="p-5 w-32 h-32 bg-blue-100 rounded-lg">
              <p className="text-md font-medium">Images</p>
              <p className="text-3xl font-bold">
                {images.value.summary.images}
              </p>
            </div>
          }
          {
            <div className="p-5 w-32 h-32 bg-orange-100 rounded-lg">
              <p className="text-md font-medium">Outdated</p>
              <p className="text-3xl font-bold">
                {images.value.summary.outdated}
              </p>
            </div>
          }
          {images.value.summary.pods !== 0 && (
            <div className="p-5 w-32 h-32 bg-purple-100 rounded-lg">
              <p className="text-md font-medium">Pods</p>
              <p className="text-3xl font-bold">{images.value.summary.pods}</p>
            </div>
          )}
        </div>

        {/* Table card */}
        <div className="relative mt-6">
          <div className="rounded-lg bg-white px-4 py-2 shadow">
            <table>
              <thead>
                <tr>
                  <th
                    scope="col"
                    colSpan={2}
                    className="text-nowrap text-center cursor-pointer pr-[24px]"
                    onClick={() => toggleSort('reference')}
                  >
                    Image
                    <div className="inline-block relative py-[9px]">
                      {sortProperty === 'reference' && sortOrder === 'asc' && (
                        <FluentArrowSortUp24Filled className="absolute top-0" />
                      )}
                      {sortProperty === 'reference' && sortOrder === 'desc' && (
                        <FluentArrowSortDown24Filled className="absolute top-0" />
                      )}
                    </div>
                  </th>
                  <th scope="col" className="text-nowrap text-center pr-[24px]">
                    Version
                  </th>
                  <th scope="col" className="text-nowrap text-center pr-[24px]">
                    New version
                  </th>
                  <th scope="col" className="text-nowrap text-center">
                    Tags
                  </th>
                  <th scope="col"></th>
                </tr>
              </thead>
              <tbody>
                {images.value.images.map((image) => (
                  <tr key={image.reference} className="">
                    <td>
                      {image.image ? (
                        <img
                          className="w-10 h-10 rounded"
                          src={image.image}
                          referrerPolicy="no-referrer"
                        />
                      ) : (
                        <div className="w-10 h-10 rounded bg-blue-500 flex items-center justify-center">
                          <SimpleIconsOci className="text-white" />
                        </div>
                      )}
                    </td>
                    <td className="pr-[24px]">
                      <p>{name(image.reference)} </p>
                      {image.description && (
                        <p className="text-xs">{image.description}</p>
                      )}
                    </td>
                    <td
                      className={`text-end pr-[24px] ${image.reference === image.latestReference ? '' : 'text-red-400'}`}
                    >
                      {version(image.latestReference)}
                    </td>
                    <td
                      className={`text-end pr-[24px] ${image.reference === image.latestReference ? '' : 'text-green-400'}`}
                    >
                      {version(image.reference)}
                    </td>
                    <td className="flex flex-wrap">
                      {tags.value
                        .filter((tag) => image.tags.includes(tag.name))
                        .map((tag) => (
                          <Badge
                            key={tag.name}
                            label={tag.name}
                            color={tag.color}
                            className="cursor-pointer"
                            onClick={() => setFilter([tag.name])}
                          />
                        ))}
                    </td>
                    <td>
                      <NavLink
                        to={`/image?reference=${encodeURIComponent(image.reference)}`}
                      >
                        <FluentChevronRight24Regular />
                      </NavLink>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>

            {/* Pagination footer */}
            <div className="mt-4">
              <div className="flex items-center justify-center text-sm">
                {pages.map(({ index, label, current }) =>
                  index === undefined ? (
                    <p className="m-1 cursor-default">{label}</p>
                  ) : (
                    <NavLink
                      // TODO
                      to={`/?page=${index}`}
                      className={`m-1 w-6 h-6 text-center leading-6 rounded ${current ? 'bg-blue-400' : 'hover:bg-blue-400'}`}
                    >
                      <p>{label}</p>
                    </NavLink>
                  )
                )}
              </div>
              <p className="text-sm">
                Showing {images.value.images.length} of{' '}
                {images.value.pagination.total} entries
              </p>
            </div>
          </div>

          {/* Side menu with tag filters */}
          <div className="absolute left-full top-0 px-2">
            <div className="rounded-lg bg-white p-4 w-64 shadow">
              <p>Tags</p>
              <div className="flex flex-wrap mt-2">
                {tags.value.map((tag) => (
                  <Badge
                    key={tag.name}
                    label={tag.name}
                    color={tag.color}
                    disabled={!filter.includes(tag.name)}
                    onClick={() => toggleTag(tag.name)}
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
