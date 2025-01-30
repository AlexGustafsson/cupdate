import type { JSX } from 'react'

export function SettingsCard(): JSX.Element {
  return (
    <div className="rounded-lg bg-white dark:bg-[#1e1e1e] px-4 py-6 shadow">
      <p>
        Cupdate version:{' '}
        {import.meta.env.VITE_CUPDATE_VERSION || 'development build'}.
      </p>
    </div>
  )
}
