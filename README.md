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
# Add documents from a YAML file
./chromadb-test add --file documents.yaml

# Add or update documents (upsert)
./chromadb-test add --upsert --file documents.yaml

# List documents in the collection
./chromadb-test list

# Query for similar documents
./chromadb-test query --query "neural networks and AI"
./chromadb-test query -q "machine learning" --results 5

# Delete a document by ID
./chromadb-test delete <document-id>

# Use a different embedding function
./chromadb-test -e gemini add --file documents.yaml
./chromadb-test -e openai --openai-base-url https://api.openai.com/v1 add --file documents.yaml
```

See <https://go-client.chromadb.dev/embeddings/> for info on different embeddings.
The following are currently supported:

* `default`
    * uses [all-MiniLM-L6-v2](https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2) on the [ONNX](https://onnx.ai/) runtime

* `gemini`
    * get an API key at <https://ai.google.dev/gemini-api/docs/api-key> and set `GEMINI_API_KEY` in your environment

* `openai`
    * requires `--openai-api-key` and `--openai-base-url` flags or env vars
    * Use `-m` to specify the embedding model

* `hash` – a simple consistent hash, not particularly recommended

## Future

* See [[ENH] Implement Chroma Cloud Search API with rank expressions](https://github.com/amikos-tech/chroma-go/pull/291)