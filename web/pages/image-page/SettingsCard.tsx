import type { JSX } from 'react'
import { FluentSettings16Regular } from '../../components/icons/fluent-settings-regular-16'
import { Card } from './Card'

export function SettingsCard(): JSX.Element {
  return (
    <Card
      persistenceKey="settings"
      tabs={[
        {
          icon: <FluentSettings16Regular />,
          label: 'Settings',
          content: (
            <p>
              Cupdate version:{' '}
              {import.meta.env.VITE_CUPDATE_VERSION || 'development build'}.
            </p>
          ),
        },
      ]}
    />
  )
}
