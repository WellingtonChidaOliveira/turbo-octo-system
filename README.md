# turbo-octo-system

Case Tecnico - Desenvolvedor Pleno (Golang)

Pipeline de metricas de desenvolvedores com dois servicos Go, SQS e DynamoDB rodando em LocalStack.

```text
raw-events -> processor -> processed-events -> aggregator -> DynamoDB -> API REST
```

## Como Rodar

Suba toda a aplicacao:

```bash
docker compose up --build
```

Se o LocalStack estiver com estado antigo e o init de recursos falhar, limpe o volume:

```bash
docker compose down -v
docker compose up --build
```

Popule a fila `raw-events`:

```bash
./scripts/seed.sh
```

Consulte a API do aggregator:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/metrics/dev-001
curl http://localhost:8080/metrics/dev-001/summary
```

## Servicos

### Processor

Consome mensagens da fila `raw-events`, valida o payload, enriquece o evento com `processed_at` e `processor_id`, publica em `processed-events` e deleta a mensagem original somente depois do publish bem-sucedido.

Principais garantias:

- Worker pool configuravel por `WORKER_COUNT`.
- Retry com backoff exponencial e jitter no receive/publish.
- Eventos invalidos nao sao deletados e seguem para DLQ apos as tentativas configuradas na fila.
- Logs estruturados em JSON com `event_id` quando disponivel.
- Shutdown com drenagem dos workers.

### Aggregator

Consome mensagens da fila `processed-events`, persiste o evento individual em `events`, atualiza o resumo incremental em `developer_summary` e expoe API REST com Fiber.

Principais garantias:

  - Idempotencia por `event_id` usando `TransactWriteItems` e `attribute_not_exists(event_id)`.
  - Mensagem duplicada e tratada como sucesso idempotente e deletada da fila.
  - `last_activity` preserva o maior timestamp recebido, mesmo com eventos fora de ordem.
- Consulta de eventos por `developer_id` via GSI `developer_id-index`.
- Health check valida SQS e DynamoDB.
- Shutdown encerra API, consumer e workers de forma coordenada.

## API

### `GET /health`

Retorna `200` quando SQS e DynamoDB estao acessiveis.

```json
{
  "status": "ok"
}
```

Retorna `503` quando alguma dependencia falha.

### `GET /metrics/:developer_id`

Retorna todos os eventos processados de um desenvolvedor.

```json
[
  {
    "event_id": "550e8400-e29b-41d4-a716-446655440000",
    "developer_id": "dev-001",
    "metric_type": "commits",
    "value": 12,
    "repository": "org/api",
    "timestamp": "2026-04-15T10:30:00Z",
    "processed_at": "2026-04-15T10:30:05Z",
    "processor_id": "processor-1"
  }
]
```

### `GET /metrics/:developer_id/summary`

Retorna o resumo agregado do desenvolvedor.

```json
{
  "developer_id": "dev-001",
  "total_commits": 12,
  "total_pull_requests": 3,
  "avg_review_time_minutes": 45,
  "events_processed": 3,
  "last_activity": "2026-04-15T10:32:00Z"
}
```

## Configuracao

Variaveis principais dos servicos:

```text
AWS_ENDPOINT_URL=http://localhost:4566
RAW_QUEUE_URL=http://localhost:4566/000000000000/raw-events
PROCESSED_QUEUE_URL=http://localhost:4566/000000000000/processed-events
EVENTS_TABLE_NAME=events
DEVELOPER_SUMMARY_TABLE_NAME=developer_summary
API_PORT=8080
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
WORKER_COUNT=5
JOB_BUFFER_SIZE=10
SHUTDOWN_TIMEOUT=30s
SQS_MAX_NUMBER_OF_MESSAGES=10
SQS_WAIT_TIME_SECONDS=20
SQS_VISIBILITY_TIMEOUT_SECONDS=30
RECEIVE_BACKOFF_INITIAL=1s
RECEIVE_BACKOFF_MAX=30s
RECEIVE_BACKOFF_JITTER=250ms
```

Dentro do Docker Compose, os endpoints apontam para `http://localstack:4566`.
Aliases legados ainda aceitos temporariamente: `QUEUE_URL` para `AWS_ENDPOINT_URL` e `REGION` para `AWS_REGION`.

## Testes

```bash
cd services/processor
GOCACHE=/tmp/go-build go test ./...

cd ../aggregator
GOCACHE=/tmp/go-build go test ./...
```

## Tradeoffs

- API implementada com Fiber para simplificar roteamento e testes HTTP.
- Sem OpenTelemetry nesta versao; seria o proximo passo para tracing distribuido.
- `last_activity` usa comparacao lexicografica de timestamps RFC3339. Os eventos gerados pelo Processor usam o formato UTC esperado para essa comparacao.
- Health check valida operacoes leves de SQS e DynamoDB, suficiente para demonstracao local.
