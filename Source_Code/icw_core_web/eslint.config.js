import stylistic from '@stylistic/eslint-plugin';
import perfectionist from 'eslint-plugin-perfectionist';
import simpleImportSort from 'eslint-plugin-simple-import-sort';
import tseslint from 'typescript-eslint';

const styleRules = {
  '@stylistic/array-bracket-spacing': ['error', 'never'],
  '@stylistic/arrow-parens': ['error', 'always'],
  '@stylistic/brace-style': ['error', '1tbs', { allowSingleLine: true }],
  '@stylistic/comma-dangle': ['error', 'always-multiline'],
  '@stylistic/comma-spacing': ['error', { before: false, after: true }],
  '@stylistic/eol-last': ['error', 'always'],
  '@stylistic/indent': ['error', 2, { SwitchCase: 1 }],
  '@stylistic/jsx-curly-brace-presence': ['error', { children: 'never', props: 'never' }],
  '@stylistic/jsx-quotes': ['error', 'prefer-double'],
  '@stylistic/jsx-self-closing-comp': 'error',
  '@stylistic/key-spacing': ['error', { beforeColon: false, afterColon: true }],
  '@stylistic/keyword-spacing': ['error', { before: true, after: true }],
  '@stylistic/member-delimiter-style': [
    'error',
    {
      multiline: { delimiter: 'semi', requireLast: true },
      singleline: { delimiter: 'semi', requireLast: false },
    },
  ],
  '@stylistic/no-multiple-empty-lines': ['error', { max: 1, maxBOF: 0, maxEOF: 0 }],
  '@stylistic/no-trailing-spaces': 'error',
  '@stylistic/object-curly-spacing': ['error', 'always'],
  '@stylistic/quotes': ['error', 'single', { avoidEscape: true }],
  '@stylistic/semi': ['error', 'always'],
  '@stylistic/space-before-blocks': ['error', 'always'],
  '@stylistic/space-before-function-paren': ['error', { anonymous: 'always', named: 'never', asyncArrow: 'always' }],
  '@stylistic/space-infix-ops': 'error',
  '@stylistic/type-annotation-spacing': 'error',
  'no-restricted-syntax': [
    'error',
    {
      selector: 'Program > :not(ImportDeclaration) ~ ImportDeclaration',
      message: 'All imports must be placed at the top of the file.',
    },
    {
      selector: String.raw`Literal[value=/\b(?:rgb|rgba|hsl|hsla)\s*\(/i]`,
      message: 'Color literals must use #RRGGBB or #RRGGBBAA uppercase hex notation.',
    },
    {
      selector: String.raw`TemplateElement[value.raw=/\b(?:rgb|rgba|hsl|hsla)\s*\(/i]`,
      message: 'Color literals must use #RRGGBB or #RRGGBBAA uppercase hex notation.',
    },
    {
      selector: String.raw`Literal[value=/#[0-9A-Fa-f]*[a-f][0-9A-Fa-f]*(?![0-9A-Fa-f])/]`,
      message: 'Hex color literals must be uppercase.',
    },
    {
      selector: String.raw`TemplateElement[value.raw=/#[0-9A-Fa-f]*[a-f][0-9A-Fa-f]*(?![0-9A-Fa-f])/]`,
      message: 'Hex color literals must be uppercase.',
    },
    {
      selector: String.raw`Literal[value=/#(?:[0-9A-Fa-f]{1,5}|[0-9A-Fa-f]{7}|[0-9A-Fa-f]{9,})(?![0-9A-Fa-f])/]`,
      message: 'Hex color literals must use #RRGGBB or #RRGGBBAA notation.',
    },
    {
      selector: String.raw`TemplateElement[value.raw=/#(?:[0-9A-Fa-f]{1,5}|[0-9A-Fa-f]{7}|[0-9A-Fa-f]{9,})(?![0-9A-Fa-f])/]`,
      message: 'Hex color literals must use #RRGGBB or #RRGGBBAA notation.',
    },
  ],
  'perfectionist/sort-jsx-props': ['error', { ignoreCase: true, order: 'asc', type: 'alphabetical' }],
  'simple-import-sort/exports': 'error',
  'simple-import-sort/imports': 'error',
};

const qualityRules = {
  '@typescript-eslint/consistent-type-definitions': ['error', 'interface'],
  '@typescript-eslint/consistent-type-imports': [
    'error',
    { fixStyle: 'separate-type-imports', prefer: 'type-imports' },
  ],
  '@typescript-eslint/explicit-function-return-type': [
    'error',
    {
      allowExpressions: true,
      allowHigherOrderFunctions: false,
      allowTypedFunctionExpressions: true,
    },
  ],
  '@typescript-eslint/explicit-module-boundary-types': 'error',
  '@typescript-eslint/naming-convention': [
    'error',
    {
      selector: 'default',
      format: ['camelCase'],
      leadingUnderscore: 'allow',
      trailingUnderscore: 'allow',
    },
    {
      selector: 'variable',
      format: ['camelCase', 'PascalCase', 'UPPER_CASE'],
      leadingUnderscore: 'allow',
      trailingUnderscore: 'allow',
    },
    {
      selector: 'function',
      format: ['camelCase', 'PascalCase'],
    },
    {
      selector: 'import',
      format: ['camelCase', 'PascalCase'],
    },
    {
      selector: 'parameter',
      format: ['camelCase'],
      leadingUnderscore: 'allow',
    },
    {
      selector: 'typeLike',
      format: ['PascalCase'],
    },
    {
      selector: ['property', 'objectLiteralProperty', 'typeProperty'],
      format: ['camelCase', 'PascalCase', 'snake_case'],
      leadingUnderscore: 'allow',
      trailingUnderscore: 'allow',
    },
  ],
  '@typescript-eslint/no-floating-promises': ['error', { ignoreVoid: true, ignoreIIFE: false }],
  '@typescript-eslint/no-confusing-void-expression': [
    'error',
    { ignoreArrowShorthand: false, ignoreVoidOperator: false },
  ],
  '@typescript-eslint/no-explicit-any': 'error',
  '@typescript-eslint/no-import-type-side-effects': 'error',
  '@typescript-eslint/no-non-null-assertion': 'error',
  '@typescript-eslint/no-unnecessary-condition': 'error',
  '@typescript-eslint/prefer-nullish-coalescing': 'error',
  '@typescript-eslint/prefer-optional-chain': 'error',
  '@typescript-eslint/switch-exhaustiveness-check': 'error',
  complexity: ['warn', 12],
  curly: ['error', 'all'],
  eqeqeq: ['error', 'always'],
  'max-depth': ['warn', 4],
  'max-lines': ['warn', { max: 500, skipBlankLines: true, skipComments: true }],
  'max-lines-per-function': ['warn', { IIFEs: true, max: 120, skipBlankLines: true, skipComments: true }],
  'no-console': 'error',
  'no-debugger': 'error',
  'no-alert': 'error',
  'no-duplicate-imports': ['error', { allowSeparateTypeImports: true }],
  'no-else-return': 'error',
  'no-restricted-imports': [
    'error',
    {
      patterns: ['../*/*/*'],
    },
  ],
  'no-useless-concat': 'error',
  'no-useless-return': 'error',
  'no-var': 'error',
  'object-shorthand': ['error', 'always'],
  'prefer-destructuring': ['warn', { array: false, object: true }],
  'prefer-const': 'error',
  'prefer-template': 'error',
};

export default tseslint.config(
  {
    ignores: ['dist', 'node_modules'],
  },
  ...tseslint.configs.strictTypeChecked.map((config) => ({
    ...config,
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      ...config.languageOptions,
      parserOptions: {
        ...config.languageOptions?.parserOptions,
        project: ['./tsconfig.app.json', './tsconfig.node.json'],
        tsconfigRootDir: import.meta.dirname,
      },
    },
    plugins: {
      ...config.plugins,
      '@stylistic': stylistic,
      perfectionist,
      'simple-import-sort': simpleImportSort,
    },
    rules: {
      ...config.rules,
      ...qualityRules,
      ...styleRules,
    },
  })),
);
