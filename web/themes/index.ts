export interface Theme {
  name: string
  displayName: string
  description: string
  load?: () => Promise<{ default: typeof import('*.css') }>
}

const themes: Theme[] = [
  {
    name: 'system',
    displayName: 'System',
    description: "Cupdate's default auto theme",
  },
  {
    name: 'light',
    displayName: 'Light',
    description: "Cupdate's default light theme",
  },
  {
    name: 'dark',
    displayName: 'Dark',
    description: "Cupdate's default dark theme",
  },
  // External
  {
    name: 'gruvbox',
    displayName: 'Gruvbox',
    description: 'Gruvbox auto',
    load: () => import('./gruvbox.css'),
  },
  {
    name: 'dracula',
    displayName: 'Dracula',
    description: 'Dracula dark',
    load: () => import('./dracula.css'),
  },
  {
    name: 'catppuccin',
    displayName: 'Catppuccin',
    description: 'Catppuccin pastel auto (latte/mocha)',
    load: () => import('./catppuccin.css'),
  },
]

export default themes
