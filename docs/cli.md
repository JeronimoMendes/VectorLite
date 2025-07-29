# VectorLite CLI Documentation

VectorLite provides a command-line interface for managing and running the vector database with support for multiple databases and configurable search algorithms.

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

**General Commands:**
- `help` - Show available commands
- `status` - Check server connection
- `quit` or `exit` - Exit the client

**Database Management Commands:**
- `create-db <name> <algorithm>` - Create a new database
  - Example: `create-db mydb bruteforce`
  - Algorithms: `bruteforce`, `hnsw`
- `use-db <name>` - Select database to use for operations
  - Example: `use-db mydb`
- `list-dbs` - List all available databases

**Vector Operations Commands:** *(require database selection)*
- `add <vector> <metadata>` - Add vector entry to selected database
  - Example: `add [1.0,2.0,3.0] name=test,type=example`
- `query <vector> <k> <metric>` - Query similar vectors in selected database
  - Example: `query [1.0,2.0,3.0] 5 cosine`
  - Metrics: `cosine`, `dot_product`, `euclidean`
- `import <file>` - Import vectors from file (CSV format) to selected database
  - Example: `import vectors.csv`
  - Supports auto-detection of headers and batch processing
- `list` - List all entries in selected database

## Usage Examples

### Start the Vector Database Server

```bash
# Build and start server
make build
./bin/vectorlite serve

# Or use the make target
make run
```

### Using the Interactive Client with Multi-Database Support

```bash
# Start the client (server must be running)
./bin/vectorlite client

# Example session demonstrating multi-database workflow:
vectorlite> list-dbs
Available databases (0):

vectorlite> create-db documents bruteforce
Database 'documents' created successfully with algorithm 'bruteforce'

vectorlite> create-db images hnsw
Database 'images' created successfully with algorithm 'hnsw'

vectorlite> list-dbs
Available databases (2):
1. documents
2. images

vectorlite> use-db documents
Now using database: documents

vectorlite[documents]> add [1.0,2.0,3.0] name=doc1,category=text
Entry added successfully

vectorlite[documents]> add [1.1,2.1,3.1] name=doc2,category=text  
Entry added successfully

vectorlite[documents]> query [1.0,2.0,3.0] 2 cosine
Found 2 similar entries:
1. ID: 1
   Vector: [1 2 3]
   Metadata: map[category:text name:doc1]

2. ID: 2
   Vector: [1.1 2.1 3.1]
   Metadata: map[category:text name:doc2]

vectorlite[documents]> use-db images
Now using database: images

vectorlite[images]> add [0.5,1.5,2.5] name=img1,type=photo
Entry added successfully

vectorlite[images]> list
Total entries: 1
1. ID: 1
   Vector: [0.5 1.5 2.5]
   Metadata: map[name:img1 type:photo]

vectorlite[images]> quit
Goodbye!
```

### Bulk Import from CSV

```bash
# Example CSV format (vectors.csv):
# 1.0,2.0,3.0,name=doc1,category=test
# 1.1,2.1,3.1,name=doc2,category=test
# 0.9,1.9,2.9,name=doc3,category=example

# First select a database
vectorlite> use-db documents
Now using database: documents

vectorlite[documents]> import vectors.csv
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

**Database Management:**
- **Create database:** `POST http://localhost:9123/databases`
- **List databases:** `GET http://localhost:9123/databases`
- **Delete database:** `DELETE http://localhost:9123/databases/{name}`

**Vector Operations:** *(all require database parameter)*
- **Add vectors:** `POST http://localhost:9123/entries`
- **Query vectors:** `POST http://localhost:9123/query`  
- **List entries:** `GET http://localhost:9123/entries?database={name}`

#### API Examples

```bash
# Create a database
curl -X POST http://localhost:9123/databases \
  -H "Content-Type: application/json" \
  -d '{"name": "my_db", "algorithm": "bruteforce"}'

# Add vectors to database
curl -X POST http://localhost:9123/entries \
  -H "Content-Type: application/json" \
  -d '{
    "database": "my_db",
    "vectors": [[1.0, 2.0, 3.0], [2.0, 3.0, 4.0]],
    "metadatas": [{"name": "doc1"}, {"name": "doc2"}]
  }'

# Query vectors
curl -X POST http://localhost:9123/query \
  -H "Content-Type: application/json" \
  -d '{
    "database": "my_db",
    "vector": [1.1, 2.1, 3.1],
    "k": 5,
    "metric": "cosine"
  }'

# List entries from database
curl "http://localhost:9123/entries?database=my_db"

# List all databases
curl http://localhost:9123/databases
```

## Algorithm Selection

VectorLite supports multiple search algorithms that can be chosen when creating a database:

### Brute Force Algorithm
- **Best for:** Small to medium datasets (< 10K vectors)
- **Characteristics:** 
  - Exact search results
  - Linear time complexity O(n)
  - Low memory overhead
  - Simple and reliable

### HNSW (Hierarchical Navigable Small World)
- **Best for:** Large datasets requiring fast approximate search
- **Characteristics:**
  - Approximate search results (high accuracy)
  - Logarithmic time complexity O(log n) 
  - Higher memory usage
  - Excellent for high-dimensional vectors

### Choosing an Algorithm

```bash
# For small datasets or when exact results are required
vectorlite> create-db exact_search bruteforce

# For large datasets where speed is more important than perfect accuracy
vectorlite> create-db fast_search hnsw
```

**Recommendations:**
- Use `bruteforce` for datasets under 10,000 vectors or when you need exact results
- Use `hnsw` for larger datasets or when query speed is critical
- You can create multiple databases with different algorithms for different use cases

## Development

- **Build:** `make build`
- **Run:** `make run` 
- **Test:** `make test`
- **Debug:** `make debug`
- **Clean:** `make clean`