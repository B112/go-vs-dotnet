# Go vs .NET — Notes API Comparison

Two identical REST APIs with a web UI, built in Go and C#/.NET, to compare the languages side by side.

## Project structure

```
go-notes/                       dotnet-notes/
├── go.mod                      ├── DotnetNotes.csproj
├── models.go                   ├── Models.cs
├── store.go                    ├── NoteStore.cs
├── handlers.go                 ├── Program.cs
├── html.go (//go:embed)        └── wwwroot/
└── static/                         └── index.html
    └── index.html
```

## Prerequisites

**Go** — install from [go.dev/dl](https://go.dev/dl) or `brew install go` on Mac.

**.NET 10** — install from [dotnet.microsoft.com](https://dotnet.microsoft.com/download) or `brew install dotnet` on Mac.

## Running

```bash
# Go (port 5080, green badge)
cd go-notes
go run .

# .NET (port 5081, purple badge)
cd dotnet-notes
dotnet run
```

Open http://localhost:5080 (Go) or http://localhost:5081 (.NET) for the web UI.

## API endpoints

| Method | Path                   | Description                          |
|--------|------------------------|--------------------------------------|
| GET    | /api/notes             | List all notes                       |
| GET    | /api/notes?category=x  | Filter by category                   |
| GET    | /api/notes/{id}        | Get a single note                    |
| POST   | /api/notes             | Create a note                        |
| DELETE | /api/notes/{id}        | Delete a note                        |
| POST   | /api/notes/import      | Import notes from a JSON file        |
| GET    | /api/stats             | Get note count, categories, uptime   |

## curl examples

All examples use port 5080 (Go). Replace with 5081 for .NET.

### List all notes

```bash
curl -s http://localhost:5080/api/notes | jq
```

### Filter by category

```bash
curl -s "http://localhost:5080/api/notes?category=work" | jq
```

### Get a single note

```bash
curl -s http://localhost:5080/api/notes/1 | jq
```

### Create a note

```bash
curl -s -X POST http://localhost:5080/api/notes \
  -H "Content-Type: application/json" \
  -d '{"title": "Update fleet config", "content": "Push new env vars to Balena", "category": "work"}' | jq
```

### Create a note (minimal — only title required)

```bash
curl -s -X POST http://localhost:5080/api/notes \
  -H "Content-Type: application/json" \
  -d '{"title": "Quick reminder"}' | jq
```

### Delete a note

```bash
curl -s -X DELETE http://localhost:5080/api/notes/3
```

### Get stats

```bash
curl -s http://localhost:5080/api/stats | jq
```

### Import notes from a JSON file

First, create a test file:

```bash
cat > /tmp/notes.json << 'EOF'
[
  {"title": "Deploy OPX update", "content": "Push v2.3 to fleet", "category": "work"},
  {"title": "Book dentist", "content": "", "category": "personal"},
  {"title": "", "content": "This one will fail — no title"}
]
EOF
```

Then import it:

```bash
curl -s -X POST http://localhost:5080/api/notes/import \
  -H "Content-Type: application/json" \
  -d '{"filePath": "/tmp/notes.json"}' | jq
```

Expected response (2 imported, 1 failed validation):

```json
{
  "imported": 2,
  "errors": ["note 2: title required"]
}
```

### Error handling examples

```bash
# File not found
curl -s -X POST http://localhost:5080/api/notes/import \
  -H "Content-Type: application/json" \
  -d '{"filePath": "/tmp/does-not-exist.json"}' | jq

# Invalid JSON in file
echo "not valid json" > /tmp/bad.json
curl -s -X POST http://localhost:5080/api/notes/import \
  -H "Content-Type: application/json" \
  -d '{"filePath": "/tmp/bad.json"}' | jq

# Invalid request body
curl -s -X POST http://localhost:5080/api/notes \
  -H "Content-Type: application/json" \
  -d 'not json' | jq

# Note not found
curl -s http://localhost:5080/api/notes/999 | jq
```
