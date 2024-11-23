## Offi web extension for firefox and chrome

- `src/index.ts` contains entrypoint content script that calls functions for mutating page content.
- `src/api/` contains functions for fetching data from backend server.
- `src/options.ts, src/options_page.ts` contains logic for managing extension options.

### Building
You need to have `yarn` installed in order to build the extension.

```sh
yarn install

yarn lint

yarn build_f # firefox
yarn build_c # chrome
```
