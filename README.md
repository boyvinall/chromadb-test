# chromadb-test

Explorations with [chromadb](https://github.com/chroma-core/chroma),
built fairly quickly with Cursor.

## Setup

```bash
make start   # start ChromaDB
make build   # build the binary
```

## Usage

```bash
# Add sample documents
./chromadb-test add

# Query for similar documents
./chromadb-test query

# Use a different embedding function
./chromadb-test -e gemini add
./chromadb-test -e gemini query
```

See <https://go-client.chromadb.dev/embeddings/> for info on different embeddings.
The following are currently supported:

- `default` – uses [all-MiniLM-L6-v2](https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2) on the [ONNX](https://onnx.ai/) runtime
- `gemini` – get an API key at <https://ai.google.dev/gemini-api/docs/api-key> and set `GEMINI_API_KEY` in your environment
- `hash` – a simple consistent hash, not particularly recommended
