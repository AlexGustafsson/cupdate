import { fixupConfigRules, fixupPluginRules } from '@eslint/compat'
import { FlatCompat } from '@eslint/eslintrc'
import js from '@eslint/js'
import tsParser from '@typescript-eslint/parser'
import eslintConfigPrettier from 'eslint-config-prettier'
import json from 'eslint-plugin-json'
import reactHooks from 'eslint-plugin-react-hooks'
import unusedImports from 'eslint-plugin-unused-imports'
import globals from 'globals'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import tseslint from 'typescript-eslint'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all,
})

// TODO: Once https://github.com/facebook/react/pull/30774 is closed,
// we can remove all compat things and fixup plugin rules.
// With a little bit of luck, alexg will remember what to do, but might need to
// be poked.

export default [
  // Ignore options
  {
    ignores: [
      'target/**',
      '**/.vscode',
      '**/__generated__/',
      '**/node_modules/',
      '**/build/',
      '**/coverage/',
      '**/dist/',
      '**/bin/',
      '**/*.log*',
      '**/yarn.lock',
      '**/.yarn/',
      '**/.pnp.js',
    ],
  },

  // Plugin options

  // eslint:recommended
  js.configs.recommended,
  // plugin/@typescript-eslint/recommended
  ...tseslint.configs.recommended,
  {
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: { ...globals.node, ...globals.browser },
      parser: tsParser,
    },
  },
  // plugin:json/recommended
  {
    files: ['**/*.json'],
    ...json.configs['recommended'],
  },
  // plugin:react-hooks/recommended
  // SEE: https://github.com/facebook/react/pull/30774
  ...fixupConfigRules(compat.extends('plugin:react-hooks/recommended')),
  // Note: prettier needs to always be last in the list, it disables eslint style rules.
  eslintConfigPrettier,

  // Plugins

  // unused-imports
  {
    plugins: {
      'unused-imports': unusedImports,
    },
    rules: {
      '@typescript-eslint/no-unused-vars': ['off'],
      'unused-imports/no-unused-imports': ['error'],

      'unused-imports/no-unused-vars': [
        'error',
        {
          vars: 'all',
          varsIgnorePattern: '^_',
          args: 'after-used',
          argsIgnorePattern: '^_',
        },
      ],
    },
  },
  // react-hooks
  // SEE: https://github.com/facebook/react/pull/30774
  {
    plugins: {
      'react-hooks': fixupPluginRules(reactHooks),
    },
    rules: {
      curly: ['error'],
      'spaced-comment': ['error', 'always'],

      'no-restricted-syntax': [
        'error',

        {
          selector: 'TSTypeAliasDeclaration[id.name=Props]',
          message:
            'Type declarations for props should be named to reflect what it is used for. Use\nProps as a suffix.\n',
        },
        {
          selector: 'TSInterfaceDeclaration[id.name=Props]',
          message:
            'Interface declarations for props should be named to reflect what it is used for.\nUse Props as a suffix.\n',
        },
      ],

      // Use unused-imports instead of the built-in support
      '@typescript-eslint/no-unused-vars': ['off'],
      'unused-imports/no-unused-imports': ['error'],

      'unused-imports/no-unused-vars': [
        'error',
        {
          vars: 'all',
          varsIgnorePattern: '^_',
          args: 'after-used',
          argsIgnorePattern: '^_',
        },
      ],
    },
  },
]
