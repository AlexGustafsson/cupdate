import { useParams, useSearchParams } from 'react-router-dom'

export function ImagePage(): JSX.Element {
  const [params, _] = useSearchParams()

  const imageName = params.get('name')
  const imageVersion = params.get('version')

  return (
    <>
      <p>
        {imageName}:{imageVersion}
      </p>
    </>
  )
}
