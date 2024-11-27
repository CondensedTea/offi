import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  eslint.configs.recommended,
  tseslint.configs.strict,
  {
    files: ['src/*.ts'],
    rules: {
      'object-curly-spacing': ['error', 'always'],
      'no-async-promise-executor': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
    },
  }
);
