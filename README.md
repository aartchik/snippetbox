# Snippetbox

A simple web application for storing and viewing text snippets.

## Run locally

### 1. Clone the repository

```bash
git clone https://github.com/aartchik/snippetbox.git
cd snippetbox
```

### 2. Start infrastructure 

```bash
docker compose up 
```


The app will be available at:

```
http://localhost:4000
```

---

## Run tests

Make sure test database is running:

```bash
docker compose --profile test up -d test-db redis
```

Then run tests locally:

```bash
go test ./...
```
