# RigidSearch

RigidSearch is a lightweight search engine backend written in Go. It implements document indexing, persistent index storage, and ranked search using TF-IDF and BM25.

The project is intended as a systems-oriented backend project exploring how search engines represent documents, build inverted indexes, rank results, and serve search APIs.

## Features

- Inverted index with postings lists
- Document indexing through REST APIs
- TF-IDF based search
- BM25 based search
- Stop-word removal
- Porter stemming
- Persistent on-disk index storage
- Thread-safe indexing and querying
- Basic document retrieval and deletion APIs

## Tech Stack

- Go
- chi router
- REST APIs
- Inverted indexes
- TF-IDF
- BM25
- Porter stemming

## How It Works

RigidSearch stores each indexed document on disk and maintains an in-memory search index containing:

- a term dictionary
- postings lists for each indexed term
- term frequency and document frequency statistics
- document metadata
- deleted document tracking

At startup, the index is loaded from disk. During shutdown, the updated index is serialized back to disk.

Search queries are normalized using the same preprocessing pipeline as documents: tokenization, lowercasing, stop-word filtering, and stemming. Results are ranked using either TF-IDF or BM25.

## API

### Health Check

```bash
GET /
```

Returns:

```text
Operational
```

### Index a Document

```bash
POST /index
Content-Type: application/json
```

Request body:

```json
{
  "name": "example-document",
  "text": "This is the document text to be indexed."
}
```

Example:

```bash
curl -X POST http://localhost:8000/index \
  -H "Content-Type: application/json" \
  -d '{"name":"example-document","text":"Go is a productive language for backend systems."}'
```

Response:

```json
{
  "document_id": 0
}
```

### Search Documents

```bash
GET /search?query=<query>&num_results=<n>&method=<method>
```

Supported methods:

- `tf_idf`
- `bm_25`

If no method is provided, `tf_idf` is used by default.

Example:

```bash
curl "http://localhost:8000/search?query=backend%20systems&num_results=5&method=bm_25"
```

Response:

```json
[
  {
    "id": 0,
    "name": "example-document",
    "score": 0.42
  }
]
```

### Get a Document

```bash
GET /documents/{documentId}
```

Example:

```bash
curl http://localhost:8000/documents/0
```

Response:

```json
{
  "name": "example-document",
  "text": "Go is a productive language for backend systems."
}
```

### Delete a Document

```bash
DELETE /documents/{documentId}
```

Example:

```bash
curl -X DELETE http://localhost:8000/documents/0
```

Response:

```text
success
```

## Running Locally

### 1. Clone the repository

```bash
git clone https://github.com/Ciryandil/rigidsearch.git
cd rigidsearch
```

### 2. Create required storage directory

```bash
mkdir -p doc_storage
```

### 3. Configure environment variables

Create a `.env` file:

```env
INDEX_FILE=rigidsearch_index.json
STORAGE_LOC=./doc_storage
```

### 4. Run the server

```bash
go run main.go
```

The server starts on:

```text
http://localhost:8000
```

## Example Workflow

Index a few documents:

```bash
curl -X POST http://localhost:8000/index \
  -H "Content-Type: application/json" \
  -d '{"name":"go-backend","text":"Go is commonly used for backend systems and infrastructure services."}'

curl -X POST http://localhost:8000/index \
  -H "Content-Type: application/json" \
  -d '{"name":"search-engines","text":"Search engines use inverted indexes and ranking algorithms to retrieve documents."}'
```

Search using TF-IDF:

```bash
curl "http://localhost:8000/search?query=backend%20systems&num_results=5&method=tf_idf"
```

Search using BM25:

```bash
curl "http://localhost:8000/search?query=inverted%20index&num_results=5&method=bm_25"
```

## Project Structure

```text
.
├── constants/       # Runtime configuration
├── data_models/     # Request, document, and result models
├── heap/            # Generic heap implementation for ranking results
├── indexing/        # Document indexing and index persistence
├── router/          # HTTP routes and API handlers
├── search/          # TF-IDF and BM25 ranking
├── stemming/        # Porter stemming implementation
├── stop_words/      # Stop-word list
├── string_utils/    # Text normalization helpers
└── main.go          # Server startup and graceful shutdown
```

## Current Limitations

RigidSearch is a learning-oriented search backend and is not intended for production use yet. Some current limitations include:

- index persistence is file-based
- no distributed indexing or replication
- limited query parsing
- no phrase search
- no pagination
- no authentication
- deleted documents are tracked lazily in the index

## Future Improvements

Possible future directions:

- phrase and boolean queries
- pagination and result highlighting
- compact binary index format
- incremental index snapshots
- better document update support
- benchmark suite for indexing and query latency
- distributed indexing across shards
- Docker support
