export default {
  mode: 'jit',
  content: ['./web/**/*.{html,tsx}'],
  darkMode: 'media',
  theme: {
    extend: {
      boxShadow: {
        around: '0 0 2px rgba(0, 0, 0, 0.22), 0 4px 8px rgba(0, 0, 0, 0.28)',
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
