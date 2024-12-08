import { useEffect } from 'react'
import { NavLink, useNavigate, useSearchParams } from 'react-router-dom'

import { useImages, usePagination, useTags } from '../api'
import { ImageCard } from '../components/ImageCard'
import { Select } from '../components/Select'
import { TagSelect } from '../components/TagSelect'
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

  if (images.status !== 'resolved' || tags.status !== 'resolved') {
    return <></>
  }

  return (
    <div className="flex flex-col items-center pt-6 pb-10 px-2">
      <div className="grid grid-cols-2 md:grid-cols-4 divide-x dark:divide-[#333333]">
        <div className="py-2 px-4">
          <p className="text-sm">Outdated images</p>
          <p className="text-3xl font-semibold text-red-600">
            {images.value.summary.outdated}
          </p>
        </div>
        <div className="py-2 px-4">
          <p className="text-sm">Vulnerable images</p>
          <p className="text-3xl font-semibold text-red-600">
            {images.value.summary.vulnerable}
          </p>
        </div>
        <div className="py-2 px-4">
          <p className="text-sm">Queued images</p>
          <p className="text-3xl font-semibold">
            {images.value.summary.processing}
          </p>
        </div>
        <div className="py-2 px-4">
          <p className="text-sm">Total images</p>
          <p className="text-3xl font-semibold">
            {images.value.summary.images}
          </p>
        </div>
      </div>

      <hr className="my-6 w-3/4" />

      <div className="flex justify-between items-center w-full max-w-[800px]">
        <div className="flex gap-x-2">
          <Select
            value={sortProperty}
            onChange={(e) => setSortProperty(e.target.value)}
            defaultValue=""
          >
            <option value="" disabled hidden>
              Sort by
            </option>
            <option value="reference">Name</option>
          </Select>
          <Select
            value={sortOrder}
            onChange={(e) =>
              setSortOrder(e.target.value as 'asc' | 'desc' | undefined)
            }
            defaultValue=""
          >
            <option value="" disabled hidden>
              Sort order
            </option>
            <option value="asc">Ascending</option>
            <option value="desc">Descending</option>
          </Select>
          <TagSelect tags={tags.value} filter={filter} onChange={setFilter} />
        </div>
      </div>
      <div className="flex flex-col mt-2 gap-y-4 w-full max-w-[800px]">
        {images.value.images.map((x) => (
          <NavLink
            key={x.reference}
            to={`image?reference=${encodeURIComponent(x.reference)}`}
          >
            <ImageCard
              className="hover:shadow-md transition-shadow cursor-pointer"
              name={name(x.reference)}
              currentVersion={version(x.reference)}
              latestVersion={
                x.latestReference ? version(x.latestReference) : undefined
              }
              vulnerabilities={x.vulnerabilities.length}
              logo={x.image}
              description={x.description}
              tags={x.tags}
              // TODO:
              // updated={new Date(x.updated)}
            />
          </NavLink>
        ))}
      </div>

      {/* Pagination footer */}
      <div className="mt-4 flex flex-col md:flex-row items-center justify-center md:justify-between w-full max-w-[800px]">
        <p className="text-sm">
          Showing{' '}
          {Math.max(
            images.value.pagination.page * images.value.pagination.size,
            1
          )}
          -
          {images.value.pagination.page * images.value.pagination.size +
            images.value.images.length}{' '}
          of {images.value.pagination.total} entries
        </p>
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
                className={`m-1 w-6 h-6 text-center text-white dark:text-[#dddddd] leading-6 rounded ${current ? 'bg-blue-400 dark:bg-blue-800' : 'hover:bg-blue-400 hover:dark:bg-blue-800'}`}
              >
                <p>{label}</p>
              </NavLink>
            )
          )}
        </div>
      </div>
    </div>
  )
}
