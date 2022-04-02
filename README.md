## Offi - plugin and api-server

Pairs etf2l match pages and logs.tf logs.

Build with:
- Go with fiber framework
- Redis
- Typescript

## How to run
### Backend:
```shell
cd backend/
docker-compose --build -d up
```
### Plugin:
```shell
cd plugin/
npm run dev-build
```
import plugin directory in your browser.