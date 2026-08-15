import type { JSX } from 'react/jsx-runtime'
import { FluentBug16Regular } from '../components/icons/fluent-bug-16-regular'
import { FluentPaintBrush16Regular } from '../components/icons/fluent-paint-brush-16-regular'
import { Card } from './image-page/Card'

export function SettingsPage(): JSX.Element {
  return (
    <div className="flex flex-col items-center w-full pt-2 pb-10 px-2">
      <main className="min-w-[200px] max-w-[800px] w-full box-border space-y-6 mt-6">
        <Card
          persistenceKey="theming"
          tabs={[
            {
              icon: <FluentPaintBrush16Regular />,
              label: 'Theming',
              content: <p>Theme: default</p>,
            },
          ]}
        />

        <Card
          persistenceKey="build"
          tabs={[
            {
              icon: <FluentBug16Regular />,
              label: 'Build',
              content: (
                <p>
                  Cupdate version:{' '}
                  {import.meta.env.VITE_CUPDATE_VERSION || 'development build'}.
                </p>
              ),
            },
          ]}
        />
      </main>
    </div>
  )
}
