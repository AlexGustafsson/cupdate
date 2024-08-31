import { useCallback, useState } from 'react'

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

  const [filteredTags, setFilteredTags] = useState<string[]>(
    tags.map((x) => x.label)
  )

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

  const toggleTag = useCallback(
    (tag: string) => {
      setFilteredTags((previous) => {
        // If all are selected, only select the clicked tag
        if (previous.length === tags.length) {
          return [tag]
        }

        // Filter the selection
        const selection = previous.includes(tag)
          ? previous.filter((x) => x !== tag)
          : [...previous, tag]

        // If the filter excludes all tags, default to show all tags again
        if (selection.length === 0) {
          return tags.map((x) => x.label)
        }

        // Use the filtered selection
        return selection
      })
    },
    [filteredTags, setFilteredTags]
  )

  return (
    <div className="flex flex-col items-center w-full">
      <div className="rounded-lg bg-white p-3 flex space-x-5">
        <div className="p-5 w-24 h-24 bg-blue-100 rounded-lg">
          <p className="text-sm">Images</p>
          <p className="text-xl font-medium">15</p>
        </div>
        <div className="p-5 w-24 h-24 bg-orange-100 rounded-lg">
          <p className="text-sm">Outdated</p>
          <p className="text-xl font-medium">15</p>
        </div>
        <div className="p-5 w-24 h-24 bg-purple-100 rounded-lg">
          <p className="text-sm">Pods</p>
          <p className="text-xl font-medium">71</p>
        </div>
      </div>
      <div className="relative mt-6">
        <div className="rounded-lg bg-white px-4 py-2 max-w-screen-sm">
          <table>
            <thead>
              <tr>
                <th scope="col" colSpan={2} className="text-nowrap text-center">
                  Image
                </th>
                <th scope="col" className="text-nowrap text-center">
                  Version
                </th>
                <th scope="col" className="text-nowrap text-center">
                  New version
                </th>
                <th scope="col" className="text-nowrap text-center">
                  Tags
                </th>
              </tr>
            </thead>
            <tbody>
              {rowItems.map((x) => (
                <tr key={x.image} className="">
                  <td>
                    <img className="w-10 rounded" src={x.imageUrl} />
                  </td>
                  <td>{x.image}</td>
                  <td className="text-end">{x.current}</td>
                  <td className="text-end">{x.new}</td>
                  <td className="flex flex-wrap">
                    {x.tags.map((x) => (
                      <span
                        key={x}
                        className="rounded-full bg-red-100 px-2 py-1 text-xs text-nowrap m-1 cursor-pointer"
                        onClick={() => toggleTag(x)}
                      >
                        {x}
                      </span>
                    ))}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <div className="mt-2">
            <p className="text-sm">Showing 50 of 50 entries</p>
          </div>
        </div>
        <div className="absolute left-full top-0 px-2">
          <div className="rounded-lg bg-white p-4 w-64">
            <p>Tags</p>
            <div className="flex flex-wrap mt-2">
              {tags.map((x) => (
                <span
                  key={x.label}
                  className={`rounded-full px-2 py-1 text-xs text-nowrap cursor-pointer ${x.color} ${filteredTags.includes(x.label) ? '' : 'opacity-50'} m-1`}
                  onClick={() => toggleTag(x.label)}
                >
                  {x.label}
                </span>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
