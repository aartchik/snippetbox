# Snippetbox

https://snippetbox-aartchik.duckdns.org/ - ready to go link

`Snippetbox` is a small web application for storing, searching, and managing text snippets. The project is written in Go and uses server-side rendered HTML templates, MySQL for persistent storage, Redis for cache, and Docker Compose for local infrastructure.

The repository is positioned as a compact production-style learning project: it has a layered structure, database migrations, session-based user flows, integration tests, and a containerized development setup. The current codebase already includes the base for user accounts and protected actions, and it is designed to grow further with more explicit authentication and authorization features.

## What the project does

The application allows users to:

- create and store text snippets;
- open a snippet by direct URL;
- update and delete their own snippets;
- search snippets by title and content;
- sign up, log in, log out, and view account information;
- change password inside the account area;
- work with protected pages behind a session-based login flow.

Technical characteristics:

- Go HTTP server built on `net/http`;
- routing via `httprouter`;
- HTML rendering through Go templates;
- MySQL as the main database;
- Redis for snippet cache;
- sessions stored in MySQL via `scs`;
- CSRF protection and middleware-based request pipeline;
- Docker and Docker Compose for local startup.

## Project structure

```text
.
├── cmd/web                # HTTP server, routes, handlers, middleware
├── internal/models        # data access layer for snippets and users
├── internal/assert        # test helpers
├── internal/validator     # form and input validation helpers
├── migrations             # SQL migrations for schema and indexes
├── ui/html                # templates
├── ui/static              # CSS, JS, images
├── tls                    # local TLS certificates
├── Dockerfile
└── docker-compose.yml
```

## Stack

- Go
- MySQL 8
- Redis 7
- Docker Compose
- Go templates

## Run locally with Docker Compose

The simplest way to run the project is through Docker Compose. It starts MySQL, Redis, applies migrations, builds the application image, and runs the web server.

### 1. Clone the repository

```bash
git clone https://github.com/aartchik/snippetbox.git
cd snippetbox
```

### 2. Prepare environment variables

The compose file expects database variables:

```env
MYSQL_ROOT_PASSWORD=rootpass
MYSQL_DATABASE=snippetbox
MYSQL_USER=web
MYSQL_PASSWORD=pass
DB_PORT=3308
TEST_DB_PORT=3307
REDIS_PORT=6380
```

If you use a local `.env` file, `docker compose` will pick it up automatically.

### 3. Start the application

```bash
docker compose up --build
```

After startup the app will be available at:

- `http://localhost:4000`

The Compose setup includes:

- `db` - main MySQL database;
- `redis` - cache layer;
- `migrate` - one-shot migrations runner;
- `app` - Go web application.

### 4. Stop the environment

```bash
docker compose down
```

If you also want to remove volumes:

```bash
docker compose down -v
```

## Run without Docker

If you want to run the application directly on your machine, you need:

- Go installed;
- MySQL running locally on `localhost:3306`;
- Redis running locally, or the Compose Redis service exposed on `localhost:6380`;
- the schema applied from the `migrations/` directory.

Then start the server:

```bash
go run ./cmd/web \
  -dsn="web:pass@tcp(localhost:3306)/snippetbox?parseTime=true" \
  -redis-addr="localhost:6380" \
  -redis-password="" \
  -redis-db=0 \
  -tls=false
```

By default the application listens on `:4000`.

## Docker image

The repository also contains a multi-stage `Dockerfile`.

Build image:

```bash
docker build -t snippetbox .
```

Run container manually:

```bash
docker run --rm -p 4000:4000 snippetbox
```

In practice, `docker compose up --build` is the preferred workflow because it starts the dependent services too.

## Deploy on a VPS

The repository includes a production reverse-proxy profile. It expects a Linux VPS with Docker, a domain pointing to the VPS IP, and inbound ports `80` and `443` allowed by the firewall. Docker installation instructions are available in the [official Docker documentation](https://docs.docker.com/engine/install/).

On the server:

```bash
git clone https://github.com/aartchik/snippetbox.git
cd snippetbox
```

Create `.env` (it is ignored by Git) with strong, unique values:

```env
DOMAIN=app.example.com
MYSQL_ROOT_PASSWORD=replace-with-a-long-random-value
MYSQL_DATABASE=snippetbox
MYSQL_USER=snippetbox
MYSQL_PASSWORD=replace-with-another-long-random-value
SNIPPETBOX_DSN=snippetbox:replace-with-another-long-random-value@tcp(db:3306)/snippetbox?parseTime=true
APP_PORT=4000
```

Replace `app.example.com` with the real domain and use the same password in `MYSQL_PASSWORD` and `SNIPPETBOX_DSN`. Start the application and Caddy:

```bash
docker compose --profile prod up -d --build
docker compose ps
curl -f http://127.0.0.1:4000/ping
```

After DNS has propagated, open `https://app.example.com`. Caddy terminates HTTPS and proxies traffic to the app; MySQL, Redis, and the app port are bound to localhost only. The Caddy profile stores certificates in a persistent Docker volume.

Useful maintenance commands:

```bash
docker compose logs -f app caddy
docker compose --profile prod pull
docker compose --profile prod up -d --build
```

## Database migrations

Migrations live in the [`migrations/`](./migrations) directory.

They currently set up:

- `users` table;
- `sessions` table for server-side sessions;
- `snippets` table;
- index on snippet creation date;
- unique constraint and index for user email;
- full-text index for snippet search by title and content;
- user avatar URL column;
- snippet visibility level column.

Compose runs migrations automatically through the `migrate/migrate` container before the app starts.

Current migration files:

- `000001_create_snippetbox_table.*`
- `000002_snippetbox_create_index.*`
- `000003_search_title_content.*`
- `000004_avatars.*`
- `000005_visibility_snippets.*`

This means a fresh environment can be brought up from scratch without creating tables manually.

## Main features

### Snippets

- create snippet;
- view snippet;
- update snippet;
- delete snippet;
- list latest snippets;
- full-text search across title and content.

### Users and sessions

- sign up with unique email;
- log in and log out;
- persistent server-side session management;
- account page;
- password change flow;
- protected routes for authenticated users.

### Performance and UX

- Redis caching for snippet reads;
- server-rendered HTML pages;
- static assets for UI;
- middleware for security headers, CSRF, panic recovery, and request logging.

## Tests

The repository already contains tests for both HTTP handlers and the model layer.

Covered areas include:

- route and handler behavior;
- user model methods;
- snippet model methods;
- cache helpers;
- template-related behavior.

### Run tests locally

Start test infrastructure:

```bash
make test-env-up
```

Then run:

```bash
make test
```

The local test defaults are:

- MySQL: `test_web:pass@tcp(localhost:3307)/test_snippetbox`
- Redis: `localhost:6380`, DB `1`

You can override them with:

```bash
SNIPPETBOX_TEST_DSN="test_web:pass@tcp(localhost:3307)/test_snippetbox?parseTime=true&multiStatements=true&time_zone=%27%2B00%3A00%27" \
SNIPPETBOX_TEST_REDIS_ADDR="localhost:6380" \
SNIPPETBOX_TEST_REDIS_DB=1 \
go test ./...
```

### Run tests in Docker

There is also a dedicated test service in `docker-compose.yml`:

```bash
make test-compose
```

This is useful when you want a repeatable test environment without relying on locally installed Go tooling.

## Development notes

- the app entry point is [`cmd/web/main.go`](./cmd/web/main.go);
- routes are defined in [`cmd/web/routes.go`](./cmd/web/routes.go);
- data access is implemented in [`internal/models/`](./internal/models);
- templates are stored in [`ui/html/`](./ui/html);
- static files are served from [`ui/static/`](./ui/static).

## Roadmap

Planned next steps for the project:

- richer authentication and authorization rules;
- better access isolation for user-owned content;
- profile improvements and account settings;
- broader test coverage for protected flows and edge cases;
- deployment-oriented configuration cleanup.
