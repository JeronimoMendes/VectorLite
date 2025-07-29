# VectorLite CLI Documentation

VectorLite provides a command-line interface for managing and running the vector database.

## Installation

Build the CLI from source:
```bash
make build
```

This creates the `vectorlite` binary in the `bin/` directory.

## Commands

### Root Command

```bash
vectorlite
```

**Description:** A CLI for VectorLite - Use this tool to manage your vector indexes.

**Flags:**
- `-t, --toggle`: Help message for toggle (boolean, default: false)
- `-h, --help`: Show help information

### serve

```bash
vectorlite serve [flags]
```

**Description:** Run a VectorLite instance

**Flags:**
- `--port int`: Port to run the server on (default: 9123)

**Example:**
```bash
# Start server on default port (9123)
./bin/vectorlite serve

# Start server on custom port
./bin/vectorlite serve --port 8080
```

### client

```bash
vectorlite client [flags]
```

**Description:** Interactive client for VectorLite server

**Flags:**
- `--server string`: VectorLite server URL (default: "http://localhost:9123")

**Example:**
```bash
# Connect to default server
./bin/vectorlite client

# Connect to custom server
./bin/vectorlite client --server http://localhost:8080
```

#### Interactive Commands

Once in the interactive client, you can use these commands:

- `help` - Show available commands
- `status` - Check server connection
- `add <vector> <metadata>` - Add vector entry
  - Example: `add [1.0,2.0,3.0] name=test,type=example`
- `query <vector> <k> <metric>` - Query similar vectors
  - Example: `query [1.0,2.0,3.0] 5 cosine`
  - Metrics: `cosine`, `dot_product`, `euclidean`
- `import <file>` - Import vectors from file (CSV format)
  - Example: `import vectors.csv`
  - Supports auto-detection of headers and batch processing
- `list` - List all entries
- `quit` or `exit` - Exit the client

## Usage Examples

### Start the Vector Database Server

```bash
# Build and start server
make build
./bin/vectorlite serve

# Or use the make target
make run
```

### Using the Interactive Client

```bash
# Start the client (server must be running)
./bin/vectorlite client

# Example session:
vectorlite> add [1.0,2.0,3.0] name=doc1,category=test
Entry added successfully

vectorlite> add [1.1,2.1,3.1] name=doc2,category=test  
Entry added successfully

vectorlite> query [1.0,2.0,3.0] 2 cosine
Found 2 similar entries:
1. ID: 1
   Vector: [1 2 3]
   Metadata: map[category:test name:doc1]

2. ID: 2
   Vector: [1.1 2.1 3.1]
   Metadata: map[category:test name:doc2]

vectorlite> list
Total entries: 2
1. ID: 1
   Vector: [1 2 3]
   Metadata: map[category:test name:doc1]

2. ID: 2
   Vector: [1.1 2.1 3.1]
   Metadata: map[category:test name:doc2]

vectorlite> quit
Goodbye!
```

### Bulk Import from CSV

```bash
# Example CSV format (vectors.csv):
# 1.0,2.0,3.0,name=doc1,category=test
# 1.1,2.1,3.1,name=doc2,category=test
# 0.9,1.9,2.9,name=doc3,category=example

vectorlite> import vectors.csv
Detected header row, skipping...
Processing 3 rows...
Importing 3 vectors in 1 batches...
Imported batch 1/1 (3 vectors)
Successfully imported vectors from vectors.csv
```

**CSV Format Options:**
- Pure numeric: `1.0,2.0,3.0` (gets default metadata)
- With metadata: `1.0,2.0,3.0,name=doc1,category=test`
- Mixed columns: numeric values become vector, non-numeric become metadata
- Headers are auto-detected and skipped

### Direct API Access

You can also interact with the REST API directly:

- **Add vectors:** `POST http://localhost:9123/entries`
- **Query vectors:** `POST http://localhost:9123/query`  
- **List entries:** `GET http://localhost:9123/entries`

## Development

- **Build:** `make build`
- **Run:** `make run` 
- **Test:** `make test`
- **Debug:** `make debug`
- **Clean:** `make clean`