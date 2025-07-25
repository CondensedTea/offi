### Offi backend

Offi backend. Is used for caching match/log data and storing information about players' recruitment posts.

Redis is used for caching, PostgreSQL is used for persistent storage.

---



How to run:
1. Install Docker.
2. Run development stand:
```bash
docker-compose up
```
Main API is available at `http://localhost:8080`.

Tracing UI is available at `http://localhost:16686`.
