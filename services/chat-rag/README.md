# Chat-RAG 🚀

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://golang.org/doc/go1.24) [![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE) [![Docker](https://img.shields.io/badge/docker-available-blue.svg)](Dockerfile) [![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)

[English](#english) | [中文](./README.zh-CN.md)

</div>

## 🎯 Overview

Chat-RAG is a high-performance, enterprise-grade chat service that combines Large Language Models (LLM) with Retrieval-Augmented Generation (RAG) capabilities. It provides intelligent context processing, tool integration, and streaming responses for modern AI applications.

### Key Features

- **🧠 Intelligent Context Processing**: Advanced prompt engineering with context compression and filtering
- **🔧 Tool Integration**: Seamless integration with semantic search, code definition lookup, and knowledge base queries
- **⚡ Streaming Support**: Real-time streaming responses with Server-Sent Events (SSE)
- **🛡️ Enterprise Security**: JWT-based authentication and request validation
- **📊 Comprehensive Monitoring**: Built-in metrics and logging with Prometheus support
- **🔄 Multi-Modal Support**: Support for various LLM models and function calling
- **🚀 High Performance**: Optimized for low-latency responses and high throughput
 - **🤖 Semantic Router (migrated from ai-llm-router)**: Optional auto model selection via semantic classification; emits `x-select-llm` and `x-user-input` response headers

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │───▶│  Chat Handler   │───▶│  Prompt Engine  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Authentication│    │  LLM Client     │    │  Tool Executor  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Metrics       │    │  Redis Cache    │    │  Search Tools   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24.2 or higher
- Redis 6.0+ (optional, for caching)
- Docker (optional, for containerized deployment)

### Installation

```bash
# Clone the repository
git clone https://github.com/zgsm-ai/chat-rag.git
cd chat-rag

# Install dependencies
make deps

# Build the application
make build

# Run with default configuration
make run
```

### Docker Deployment

```bash
# Build Docker image
make docker-build

# Run container
make docker-run
```

## ⚙️ Configuration

The service is configured via YAML files. See [`etc/chat-api.yaml`](etc/chat-api.yaml) for the default configuration:

```yaml
# Server
Host: 0.0.0.0
Port: 8080

# LLM upstream (single endpoint; model is specified in the request body)
LLM:
  Endpoint: "http://localhost:8000/v1/chat/completions"
  # Optional: models that support function-calling
  FuncCallingModels: ["gpt-4o-mini", "o4-mini"]

# Context compression
ContextCompressConfig:
  EnableCompress: true
  TokenThreshold: 5000
  SummaryModel: "deepseek-v3"
  SummaryModelTokenThreshold: 4000
  RecentUserMsgUsedNums: 4

# Tool backends (RAG)
Tools:
  SemanticSearch:
    SearchEndpoint: "http://localhost:8002/codebase-indexer/api/v1/semantics"
    ApiReadyEndpoint: "http://localhost:8002/healthz"
    TopK: 5
    ScoreThreshold: 0.3
  DefinitionSearch:
    SearchEndpoint: "http://localhost:8002/codebase-indexer/api/v1/definitions"
    ApiReadyEndpoint: "http://localhost:8002/healthz"
  ReferenceSearch:
    SearchEndpoint: "http://localhost:8002/codebase-indexer/api/v1/references"
    ApiReadyEndpoint: "http://localhost:8002/healthz"
  KnowledgeSearch:
    SearchEndpoint: "http://localhost:8003/knowledge/api/v1/search"
    ApiReadyEndpoint: "http://localhost:8003/healthz"
    TopK: 5
    ScoreThreshold: 0.3

# Logging and classification
Log:
  LogFilePath: "logs/chat-rag.log"
  LokiEndpoint: "http://localhost:3100/loki/api/v1/push"
  LogScanIntervalSec: 60
  ClassifyModel: "deepseek-v3"
  EnableClassification: true

# Redis (optional)
Redis:
  Addr: "127.0.0.1:6379"
  Password: ""
  DB: 0

# Semantic Router (migrated from ai-llm-router). Triggered when request body model == "auto".
router:
  enabled: true
  strategy: semantic
  semantic:
    analyzer:
      model: gpt-4o-mini
      timeoutMs: 3000
      # endpoint and apiToken can override global LLM only for analyzer
      # endpoint: "http://higress-gateway.costrict.svc.cluster.local/v1/chat/completions"
      # apiToken: "<your-token>"
      # Optional advanced fields:
      # totalTimeoutMs: 5000
      # maxInputBytes: 8192
      # promptTemplate: ""   # custom classification prompt; default is built-in
      # analysisLabels: ["simple_request", "planning_request", "code_modification"]
      # dynamicMetrics:
      #   enabled: false
      #   redisPrefix: "ai_router:metrics:"
      #   metrics: ["error_rate", "p99", "circuit"]
    inputExtraction:
      protocol: openai
      userJoinSep: "\n\n"
      stripCodeFences: true
      codeFenceRegex: ""
      maxUserMessages: 100
      maxHistoryBytes: 4096
    routing:
      candidates:
        - modelName: "gpt-4o-mini"
          enabled: true
          scores:
            simple_request: 10
            planning_request: 5
            code_modification: 3
        - modelName: "o4-mini"
          enabled: true
          scores:
            simple_request: 4
            planning_request: 8
            code_modification: 6
      minScore: 0
      tieBreakOrder: ["o4-mini", "gpt-4o-mini"]
      fallbackModelName: "gpt-4o-mini"
    ruleEngine:
      enabled: false
      inlineRules: []
      bodyPrefix: "body."
      headerPrefix: "header."
```

#### Configuration details (highlights)

- LLM
  - Endpoint: Single Chat Completions endpoint. Final model is carried by request body `model`.
  - FuncCallingModels: Models supporting function-calling to enable tools.
- ContextCompressConfig
  - EnableCompress: Whether to compress long prompts.
  - TokenThreshold: Trigger threshold for compression (input tokens).
  - SummaryModel / SummaryModelTokenThreshold: Model and threshold used for summarization.
  - RecentUserMsgUsedNums: Number of recent user messages considered for compression.
- Tools (RAG)
  - Each search block provides HTTP endpoints. TopK/ScoreThreshold control recall count and quality.
- Log
  - LogFilePath: Local log file persisted before background upload to Loki.
  - LokiEndpoint: Loki push endpoint.
  - LogScanIntervalSec: Scan/upload interval in seconds.
  - ClassifyModel / EnableClassification: Optional LLM-based log categorization.
- Redis: Optional; used by tools, router dynamic metrics, and transient statuses.
- router (Semantic Router)
  - enabled/strategy: Enable semantic router; current strategy is `semantic`.
  - semantic.analyzer: Classification model/timeouts; can override endpoint/apiToken for analyzer-only calls; uses a separate non-streaming client in auto mode; optional custom prompt/labels; optional dynamic metrics via Redis.
  - semantic.inputExtraction: Controls extraction of current user input and bounded history; supports stripping code fences.
  - semantic.routing: Candidate model score table; tie-break via `tieBreakOrder`; fallback via `fallbackModelName`.
  - semantic.ruleEngine: Optional rule engine to pre-filter candidates (disabled by default).

## 📡 API Endpoints

### Chat Completion (non-streaming)

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [
      {"role": "user", "content": "What is the weather like today?"}
    ],
    "stream": false
  }'
```

### Enable Semantic Router (auto selection)

Set request body `model` to `auto` and enable `router.enabled: true` in config:

```bash
curl -i -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "auto",
    "messages": [
      {"role": "user", "content": "Give me a detailed refactor plan with code examples"}
    ],
    "stream": false
  }'
```

Response headers:
- `x-select-llm`: selected downstream model name
- `x-user-input`: extracted user input for classification (sanitized and base64-encoded)

### Streaming Response

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o-mini",
    "messages": [
      {"role": "user", "content": "Write a Python function"}
    ],
    "stream": true
  }'
```

### Metrics

Prometheus metrics are exposed at `/metrics`. See `METRICS.md` for full metric names and labels.

## 🔧 Development

### Project Structure

```
chat-rag/
├── internal/
│   ├── handler/          # HTTP handlers
│   ├── logic/           # Business logic
│   ├── client/          # External service clients
│   ├── router/          # Semantic router (strategy + factory)
│   ├── promptflow/      # Prompt processing pipeline
│   ├── functions/       # Tool execution engine
│   └── config/          # Configuration management
├── etc/                 # Configuration files
├── test/               # Test files
└── deploy/             # Deployment configurations
```

### Available Commands

```bash
make help              # Show available commands
make build            # Build the application
make test             # Run tests
make fmt              # Format code
make vet              # Vet code
make docker-build     # Build Docker image
make dev              # Run development server with auto-reload
```

### Testing

```bash
# Run all tests
make test

# Run specific test
go test -v ./internal/logic/

# Run with coverage
go test -cover ./...
```

## 🔍 Advanced Features

### Context Compression

Intelligent context compression to handle long conversations:

```yaml
ContextCompressConfig:
  EnableCompress: true
  TokenThreshold: 5000
  SummaryModel: "deepseek-v3"
  SummaryModelTokenThreshold: 4000
  RecentUserMsgUsedNums: 4
```

### Tool Integration

Support for multiple search and analysis tools:

- **Semantic Search**: Vector-based code and document search
- **Definition Search**: Code definition lookup
- **Reference Search**: Code reference analysis
- **Knowledge Search**: Document knowledge base queries

### Semantic Router (migrated from ai-llm-router)

When `router.enabled: true` and request body `model` is `auto`, the service selects the best downstream model automatically:

1. Input extraction: extract current user input and limited history per `router.semantic.inputExtraction` (can strip code fences)
2. Semantic classification: call `router.semantic.analyzer.model` to get a label (default: simple_request / planning_request / code_modification)
3. Candidate scoring: score `routing.candidates` by label; support `minScore` and optional dynamic metrics
4. Tie-break & fallback: break ties via `tieBreakOrder`; fallback to `fallbackModelName` on errors or low scores
5. Observability: write `x-select-llm` and `x-user-input` to HTTP response headers

### Agent-Based Processing

Configurable agent matching for specialized tasks:

```yaml
AgentsMatch:
  - AgentName: "strict"
    MatchKey: "a strict strategic workflow controller"
  - AgentName: "code"
    MatchKey: "a highly skilled software engineer"
```

## 📊 Monitoring & Observability

### Metrics

The service exposes Prometheus metrics at `/metrics` endpoint (see `METRICS.md` for full metric names and labels):

- Request count and latency
- Token usage statistics
- Tool execution metrics
- Error rates and types

Routing observability response headers:
- `x-select-llm`: selected model name
- `x-user-input`: base64 of extracted user input used for classification

### Logging

Structured logging with Zap logger:

- Request/response logging
- Error tracking
- Performance metrics
- Debug information

## 🔒 Security

- JWT-based authentication
- Request validation and sanitization
- Rate limiting support
- Secure header handling

## 🚢 Deployment

### Production Deployment

```bash
# Build for production
CGO_ENABLED=0 GOOS=linux go build -o chat-rag .

# Run with production config
./chat-rag -f etc/prod.yaml
```

### Kubernetes Deployment

See [`deploy/`](deploy/) directory for Kubernetes manifests and Helm charts.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- Create an issue in the GitHub repository
- Contact the maintainers

---

<div align="center">
  <b>⭐ If this project helps you, please give us a star!</b>
</div>