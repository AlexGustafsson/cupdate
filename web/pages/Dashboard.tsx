import { type JSX, useEffect, useMemo, useState } from 'react'
import { Link, useNavigate, useSearchParams } from 'react-router-dom'
import { DemoWarning } from '../components/DemoWarning'
import { ImageCard } from '../components/ImageCard'
import { FluentAlignSpaceEvenlyVertical20Filled } from '../components/icons/fluent-align-space-evenly-vertical-20-filled'
import { FluentAlignSpaceEvenlyVertical20Regular } from '../components/icons/fluent-align-space-evenly-vertical-20-regular'
import { FluentArrowSync20Regular } from '../components/icons/fluent-arrow-sync-20-regular'
import { FluentGrid20Filled } from '../components/icons/fluent-grid-20-filled'
import { FluentGrid20Regular } from '../components/icons/fluent-grid-20-regular'
import { Select } from '../components/Select'
import { TagSelect } from '../components/TagSelect'
import { Toast } from '../components/Toast'
import { type Event, useEvents } from '../EventProvider'
import {
  useDebouncedEffect,
  useFilter,
  useLayout,
  useQuery,
  useSort,
} from '../hooks'
import {
  useImages,
  usePagination,
  useRefreshImages,
  useTags,
} from '../lib/api/ApiProvider'
import { formattedVersion, fullVersion, name } from '../oci'
import { DashboardSkeleton } from './dashboard-page/DashboardSkeleton'

export function Dashboard(): JSX.Element {
  const [filter, setFilter] = useFilter()

  const [sort, setSort, sortOrder, setSortOrder] = useSort()

  const [query, setQuery] = useQuery()

  const [queryInput, setQueryInput] = useState('')

  const [searchParams, _] = useSearchParams()
  const page = useMemo(() => {
    const value = searchParams.get('page')
    if (!value) {
      return undefined
    }

    const number = Number(value)
    // Page index starts at 1
    if (Number.isNaN(number) || number < 1) {
      searchParams.delete('page')
      return undefined
    }

    // Page index starts at 1
    return number - 1
  }, [searchParams])

  const [layout, setLayout] = useLayout()

  const navigate = useNavigate()

  useDebouncedEffect(() => {
    setQuery(queryInput)
  }, [queryInput])

  useEffect(() => {
    setQueryInput(query || '')
  }, [query])

  const imageSearchOptions = useMemo(
    () => ({
      tags: filter.tags,
      tagop: filter.operator,
      sort: sort as 'reference' | 'bump' | undefined,
      order: sortOrder,
      page: page,
      limit: 30,
      query: query,
    }),
    [filter, sort, sortOrder, page, query]
  )
  const [images, updateImages] = useImages(imageSearchOptions)
  const [isRefreshingImages, refreshImages] = useRefreshImages()

  const pages = usePagination(
    images.status === 'resolved' ? images.value : undefined,
    searchParams
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
    // Page index in the API starts at 1.
    if (images.value.pagination.page - 1 >= totalPages) {
      searchParams.delete('page')
      navigate(`/?${searchParams.toString()}`)
    }
  }, [images, navigate, searchParams])

  const [tags, updateTags] = useTags()

  // Go to the first page whenever the set of tags are changed
  // biome-ignore lint/correctness/useExhaustiveDependencies: Run every time tags is changed
  useEffect(() => {
    navigate(`/?${searchParams.toString()}`)
  }, [tags, navigate, searchParams])

  const [isUpdateAvailable, setIsUpdateAvailable] = useState(false)

  useEvents((e: Event) => {
    if (e.type === 'imageUpdated' || e.type === 'graphUpdated') {
      setIsUpdateAvailable(true)
    }
  }, [])

  if (images.status !== 'resolved' || tags.status !== 'resolved') {
    return <DashboardSkeleton />
  }

  return (
    <>
      <div className="fixed bottom-[env(safe-area-inset-bottom))] flex justify-center w-full sm:w-auto sm:right-0 p-4 z-100">
        {isUpdateAvailable && (
          <Toast
            title="New data available"
            body="One or more images have been updated. Update to show the latest data."
            secondaryAction="Dismiss"
            onSecondaryAction={() => setIsUpdateAvailable(false)}
            primaryAction="Update"
            onPrimaryAction={() => {
              setIsUpdateAvailable(false)
              updateImages()
              updateTags()
            }}
          />
        )}
      </div>
      <div className="flex flex-col items-center pt-2 pb-10 px-2">
        <DemoWarning />
        <div className="grid grid-cols-3 sm:grid-cols-5 mt-4">
          <Link
            to="/?tag=outdated"
            className="rounded-lg focus:bg-white hover:bg-white dark:focus:bg-[#1e1e1e] dark:hover:bg-[#1e1e1e] transition-colors"
            tabIndex={0}
          >
            <div className="py-2 px-4">
              <p className="text-sm">Outdated</p>
              <p
                className={`text-3xl font-semibold ${images.value.summary.outdated === 0 ? 'text-green-600' : 'text-red-600'}`}
              >
                {images.value.summary.outdated}
              </p>
            </div>
          </Link>
          <Link
            to="/?tag=vulnerability:critical&tag=vulnerability:high&tag=vulnerability:medium&tag=vulnerability:low&tag=vulnerability:unspecified&tagop=or"
            className="rounded-lg focus:bg-white hover:bg-white dark:focus:bg-[#1e1e1e] dark:hover:bg-[#1e1e1e] transition-colors"
            tabIndex={0}
          >
            <div className="py-2 px-4">
              <p className="text-sm">Vulnerable</p>
              <p
                className={`text-3xl font-semibold ${images.value.summary.vulnerable === 0 ? 'text-green-600' : 'text-red-600'}`}
              >
                {images.value.summary.vulnerable}
              </p>
            </div>
          </Link>
          <Link
            to="/?tag=status:failed"
            className="rounded-lg focus:bg-white hover:bg-white dark:focus:bg-[#1e1e1e] dark:hover:bg-[#1e1e1e] transition-colors"
            tabIndex={0}
          >
            <div className="py-2 px-4">
              <p className="text-sm">Failed</p>
              <p
                className={`text-3xl font-semibold ${images.value.summary.failed === 0 ? 'text-green-600' : 'text-red-600'}`}
              >
                {images.value.summary.failed}
              </p>
            </div>
          </Link>
          <div className="py-2 px-4">
            <p className="text-sm">Queued</p>
            <p className="text-3xl font-semibold">
              {images.value.summary.processing}
            </p>
          </div>
          <div className="py-2 px-4">
            <p className="text-sm">Total</p>
            <p className="text-3xl font-semibold">
              {images.value.summary.images}
            </p>
          </div>
        </div>

        <hr className="my-6 w-3/4" />

        {/* Filters / controls */}
        <div className="flex flex-col items-center w-full mt-2 gap-y-2 max-w-[800px]">
          {/* Search */}
          <input
            type="text"
            placeholder="Search"
            enterKeyHint="search"
            value={queryInput}
            onChange={(e) => setQueryInput(e.target.value)}
            onKeyUp={(e) =>
              e.key === 'Enter' ? e.currentTarget.blur() : undefined
            }
            className="bg-white dark:bg-[#1e1e1e] pl-3 pr-8 py-2 text-sm rounded-sm flex-grow shrink-0 w-full border border-[#e5e5e5] dark:border-[#333333]"
          />
          {/* Filters */}
          <div className="grid grid-cols-3 gap-x-2 w-full">
            <Select
              value={sort || ''}
              onChange={(e) => setSort(e.target.value)}
            >
              <option value="" disabled hidden>
                Sort by
              </option>
              <option value="bump">Bump</option>
              <option value="reference">Name</option>
            </Select>
            <Select
              value={sortOrder || ''}
              onChange={(e) =>
                setSortOrder(e.target.value as 'asc' | 'desc' | undefined)
              }
            >
              <option value="" disabled hidden>
                Sort order
              </option>
              <option value="asc">Ascending</option>
              <option value="desc">Descending</option>
            </Select>
            <TagSelect tags={tags.value} filter={filter} onChange={setFilter} />
          </div>
          {/* Toolbar */}
          <div className="flex flex-row items-center justify-end gap-x-2 w-full">
            <div className="grid grid-cols-2 divide-x divide-[#e5e5e5] dark:divide-[#333333] border border-[#e5e5e5] dark:border-[#333333] rounded-sm transition-colors focus:border-[#f0f0f0] dark:focus:border-[#333333] hover:border-[#f0f0f0] dark:hover:border-[#333333] shadow-xs focus:shadow-md bg-white dark:bg-[#1e1e1e] dark:hover:bg-[#262626] h-[38px]">
              <button
                type="button"
                title="Enable list view"
                className="pl-2 pr-1 cursor-pointer focus:bg-[#f5f5f5] dark:focus:bg-[#262626]"
                onClick={() => setLayout('list')}
                tabIndex={0}
              >
                {layout === 'list' ? (
                  <FluentAlignSpaceEvenlyVertical20Filled />
                ) : (
                  <FluentAlignSpaceEvenlyVertical20Regular />
                )}
              </button>
              <button
                type="button"
                title="Enable grid view"
                className="pl-1 pr-2 cursor-pointer focus:bg-[#f5f5f5] dark:focus:bg-[#262626]"
                onClick={() => setLayout('grid')}
                tabIndex={0}
              >
                {layout === 'grid' ? (
                  <FluentGrid20Filled />
                ) : (
                  <FluentGrid20Regular />
                )}
              </button>
            </div>
            <button
              type="button"
              title="Refresh images"
              className="border border-[#e5e5e5] dark:border-[#333333] rounded-sm transition-colors focus:border-[#f0f0f0] dark:focus:border-[#333333] hover:border-[#f0f0f0] dark:hover:border-[#333333] shadow-xs bg-white dark:bg-[#1e1e1e] dark:hover:bg-[#262626] h-[38px] px-2 cursor-pointer focus:bg-[#f5f5f5] dark:focus:bg-[#262626]"
              tabIndex={0}
              onClick={refreshImages}
              disabled={isRefreshingImages}
            >
              <FluentArrowSync20Regular
                className={isRefreshingImages ? 'animate-spin' : ''}
              />
            </button>
          </div>
        </div>

        {/* Images */}
        <div
          className={`mt-2 w-full ${layout === 'list' ? 'flex flex-col max-w-[800px] gap-y-4' : 'grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-1'}`}
        >
          {images.value.images.map((x) => (
            <Link
              key={x.reference}
              to={`image?reference=${encodeURIComponent(x.reference)}`}
              state={window.location.href}
              tabIndex={0}
              className="group/link"
            >
              <ImageCard
                className={`group/link-focus:shadow-md hover:shadow-md transition-shadow-sm cursor-pointer dark:transition-colors group-focus/link:bg-[#f5f5f5] dark:group-focus/link:bg-[#262626] dark:hover:bg-[#262626] ${layout === 'list' ? '' : 'h-[150px]'}`}
                reference={x.reference}
                name={name(x.reference)}
                currentVersion={formattedVersion(x.reference, x.annotations)}
                fullCurrentVersion={fullVersion(x.reference)}
                latestVersion={
                  x.latestReference
                    ? formattedVersion(x.latestReference, x.latestAnnotations)
                    : undefined
                }
                fullLatestVersion={
                  x.latestReference ? fullVersion(x.latestReference) : undefined
                }
                vulnerabilities={x.vulnerabilities}
                logo={x.image}
                description={x.description}
                tags={x.tags}
                compact={layout === 'grid'}
                // TODO:
                // updated={new Date(x.updated)}
              />
            </Link>
          ))}
        </div>

        {/* Pagination footer */}
        <div className="mt-4 flex flex-col md:flex-row items-center justify-center md:justify-between w-full max-w-[800px]">
          <p className="text-sm">
            Showing{' '}
            {Math.max(
              (images.value.pagination.page - 1) * images.value.pagination.size,
              1
            )}
            -
            {(images.value.pagination.page - 1) * images.value.pagination.size +
              images.value.images.length}{' '}
            of {images.value.pagination.total} entries
          </p>
          <div className="flex items-center justify-center text-sm">
            {pages.map((page) =>
              page.index === undefined ? (
                <p key={page.index} className="m-1 cursor-default">
                  {page.label}
                </p>
              ) : (
                <Link
                  key={page.index}
                  to={page.href}
                  tabIndex={0}
                  className={`m-1 w-6 h-6 text-center text-white dark:text-[#dddddd] leading-6 rounded-sm ${page.current ? 'bg-blue-400 dark:bg-blue-700' : 'bg-blue-200 dark:bg-blue-900 focus:bg-blue-400 hover:bg-blue-400  hover:dark:bg-blue-700 focus:dark:bg-blue-700'}`}
                >
                  <p>{page.label}</p>
                </Link>
              )
            )}
          </div>
        </div>
      </div>
    </>
  )
}
