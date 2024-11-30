import { useCallback, useEffect } from 'react'
import { NavLink, useNavigate, useSearchParams } from 'react-router-dom'

import { useImages, usePagination, useTags } from '../api'
import { Badge } from '../components/Badge'
import { InfoTooltip } from '../components/InfoTooltip'
import { FluentChevronRight24Regular } from '../components/icons/fluent-chevron-right-24-regular'
import { FluentShieldError16Filled } from '../components/icons/fluent-shield-error-16-filled'
import { FluentArrowSortDown24Filled } from '../components/icons/fluent-sort-arrow-down-24-filled'
import { FluentArrowSortUp24Filled } from '../components/icons/fluent-sort-arrow-up-24-filled'
import { SimpleIconsOci } from '../components/icons/simple-icons-oci'
import { useFilter, useSort } from '../hooks'
import { name, version } from '../oci'

export function Dashboard(): JSX.Element {
  const [filter, setFilter] = useFilter()

  const [sortProperty, setSortProperty, sortOrder, setSortOrder] = useSort()

  const [searchParams, _] = useSearchParams()

  const navigate = useNavigate()

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

  // Go to the first page if the current page exceeds the available number of
  // pages
  useEffect(() => {
    if (images.status !== 'resolved') {
      return
    }

    const totalPages = Math.ceil(
      images.value.pagination.total / images.value.pagination.size
    )
    if (images.value.pagination.page >= totalPages) {
      searchParams.delete('page')
      navigate('/?' + searchParams.toString())
    }
  }, [images])

  const tags = useTags()

  // Go to the first page whenever the set of tags are changed
  useEffect(() => {
    navigate('/?' + searchParams.toString())
  }, [tags])

  const toggleTag = useCallback(
    (tag: string) => {
      setFilter((previous) => {
        // If all are selected, only select the clicked tag
        if (
          tags.status === 'resolved' &&
          previous.length === tags.value.length
        ) {
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
    return <></>
  }

  return (
    <>
      <div className="flex flex-col items-center w-full py-[40px] lg:px-[20px]">
        {/* Header with summary */}
        <div className="p-3 flex space-x-5">
          {
            <div className="p-5 w-32 h-32 bg-blue-100 dark:bg-blue-800 rounded-lg">
              <p className="text-md font-medium">Images</p>
              <p className="text-3xl font-bold">
                {images.value.summary.images}
              </p>
            </div>
          }
          {
            <div
              className="p-5 w-32 h-32 bg-red-100 dark:bg-red-800 rounded-lg cursor-pointer"
              onClick={() => toggleTag('outdated')}
            >
              <p className="text-md font-medium">Outdated</p>
              <p className="text-3xl font-bold">
                {images.value.summary.outdated}
              </p>
            </div>
          }
          {images.value.summary.processing !== 0 && (
            <div className="p-5 w-32 h-32 bg-purple-100 dark:bg-purple-800 rounded-lg">
              <p className="text-md font-medium">Processing</p>
              <p className="text-3xl font-bold">
                {images.value.summary.processing}
              </p>
            </div>
          )}
        </div>

        <main className="relative mt-6">
          {/* Side menu with tag filters */}
          <div className="lg:absolute left-full h-full">
            <div className="sticky top-[80px] lg:ml-4 mb-2 w-full lg:w-64 rounded-lg bg-white dark:bg-[#121212] p-4 w-64 shadow">
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

          {/* Table card. Not a table due to responsiveness and sticky header limitations */}
          <div className="rounded-lg bg-white dark:bg-[#121212] px-1 lg:px-6 py-2 shadow w-full max-w-[1200px]">
            <div className="relative break-words grid items-center dashboard-table gap-y-2 md:gap-y-4">
              {/* Header row */}
              <>
                {/* Image column */}
                <div
                  className="sticky top-[63px] bg-white dark:bg-[#121212] col-span-2 border-b-[1px] border-[#ebebeb] dark:border-[#262626] text-nowrap text-sm lg:text-base font-bold text-center cursor-pointer py-2"
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
                </div>

                {/* Version column */}
                <div className="sticky top-[63px] bg-white dark:bg-[#121212] border-b-[1px] border-[#ebebeb] dark:border-[#262626] text-nowrap text-sm lg:text-base font-bold text-center py-2">
                  Current
                </div>

                {/* Latest version column */}
                <div className="sticky top-[63px] bg-white dark:bg-[#121212] border-b-[1px] border-[#ebebeb] dark:border-[#262626] text-nowrap text-sm lg:text-base font-bold text-center py-2">
                  Latest
                </div>

                {/* Tags column */}
                <div className="sticky col-span-2 top-[64px] bg-white dark:bg-[#121212] border-b-[1px] border-[#ebebeb] dark:border-[#262626] text-nowrap text-sm lg:text-base font-bold text-center py-2">
                  Tags
                </div>
              </>

              {/* Data rows */}
              {images.value.images.map((image) => (
                <>
                  {/* Image column */}
                  <div>
                    {image.image ? (
                      <img
                        className="w-10 h-10 rounded"
                        src={image.image}
                        referrerPolicy="no-referrer"
                      />
                    ) : (
                      <div className="w-10 h-10 rounded bg-blue-500 dark:dark:bg-blue-800 flex items-center justify-center">
                        <SimpleIconsOci className="text-white" />
                      </div>
                    )}
                  </div>

                  {/* Details column */}
                  <div className="px-2">
                    <p className="text-xs md:text-base font-semibold md:font-normal">
                      {name(image.reference)}{' '}
                    </p>
                    {image.description && (
                      <p className="text-xs">{image.description}</p>
                    )}
                  </div>

                  {/* Current version column */}
                  <div
                    className={`text-end text-xs lg:text-base text-nowrap px-1 ${image.latestReference && image.reference !== image.latestReference ? 'text-red-400' : ''}`}
                  >
                    <>
                      {version(image.reference)}
                      {image.vulnerabilities.length > 0 && (
                        <InfoTooltip icon={<FluentShieldError16Filled />}>
                          {image.vulnerabilities.length} vulnerabilities
                          reported.
                        </InfoTooltip>
                      )}
                    </>
                  </div>

                  {/* Latest version column */}
                  <div
                    className={`text-end text-xs lg:text-base text-nowrap px-1 ${image.latestReference && image.reference !== image.latestReference ? 'text-green-400' : ''}`}
                  >
                    {image.latestReference ? (
                      version(image.latestReference)
                    ) : (
                      <>
                        unknown{' '}
                        <InfoTooltip>
                          The latest version cannot be identified. This could be
                          due to the image not being available, the registry not
                          being supported, missing authentication or a temporary
                          issue.
                        </InfoTooltip>
                      </>
                    )}
                  </div>

                  {/* Tags column */}
                  <div className="flex flex-wrap px-2">
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
                  </div>

                  {/* Chevron column */}
                  <div className="self-center">
                    <NavLink
                      to={`/image?reference=${encodeURIComponent(image.reference)}`}
                    >
                      <FluentChevronRight24Regular />
                    </NavLink>
                  </div>
                  <hr className="col-span-6 border-b-[1px] border-[#ebebeb] dark:border-[#262626]" />
                </>
              ))}
            </div>

            {/* Pagination footer */}
            <div className="mt-4">
              <div className="flex items-center justify-center text-sm">
                {pages.map(({ index, label, current }) =>
                  index === undefined ? (
                    <p key={index} className="m-1 cursor-default">
                      {label}
                    </p>
                  ) : (
                    <NavLink
                      // TODO
                      key={index}
                      to={`/?page=${index}`}
                      className={`m-1 w-6 h-6 text-center leading-6 rounded ${current ? 'bg-blue-400 dark:bg-blue-800' : 'hover:bg-blue-400 hover:dark:bg-blue-800'}`}
                    >
                      <p>{label}</p>
                    </NavLink>
                  )
                )}
              </div>
              <p className="text-sm">
                Showing{' '}
                {images.value.pagination.page * images.value.pagination.size}-
                {images.value.pagination.page * images.value.pagination.size +
                  images.value.images.length}{' '}
                of {images.value.pagination.total} entries
              </p>
            </div>
          </div>
        </main>
      </div>
    </>
  )
}
