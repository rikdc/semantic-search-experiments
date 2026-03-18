# Semantic Search Experiments

Companion code for the [It's Just Vectors](https://claydon.co/series/its-just-vectors/) article series. Six parts working through the mechanics of embeddings by building real Go tools — vector math, anomaly detection, semantic search, clustering, and more.

Each part has a `start/` directory with stubbed-out functions and a `complete/` directory with the full implementation. Walkthroughs guide you through the exercises; unit tests let you verify your work before wiring anything to an API.

## Structure

```text
tutorial/
├── shared/              # Shared embedder package (OpenAI + Ollama clients)
├── part1/               # Vectors, cosine similarity, centroids, embed analyze
│   ├── start/           # Stubbed — implement the TODOs
│   ├── complete/        # Reference solution
│   ├── walkthrough.md   # Step-by-step guide
│   └── questions.md     # Checkpoint questions
├── part2/               # Anomaly detection with z-scores (coming soon)
├── part3/               # (coming soon)
├── part4/               # (coming soon)
├── part5/               # (coming soon)
└── part6/               # (coming soon)
```

## Getting started

```bash
cd tutorial/part1/start/
go mod tidy
go build -o embed .
```

Set up an embedding provider:

```bash
# OpenAI
export OPENAI_API_KEY="sk-your-key-here"
./embed analyze --provider openai --model text-embedding-3-small

# Or Ollama (local, free)
ollama pull qwen2.5:latest
export OLLAMA_HOST="http://localhost:11434"
./embed analyze --provider ollama --model qwen2.5:latest
```

## Prerequisites

- Go 1.24+
- An OpenAI API key, or [Ollama](https://ollama.ai) installed locally

## License

This work is licensed under [Creative Commons Attribution 4.0 International](https://creativecommons.org/licenses/by/4.0/). See [LICENSE](LICENSE) for details.
