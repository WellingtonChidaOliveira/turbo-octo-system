# turbo-octo-system
Case Técnico — Desenvolvedor Pleno (Golang)

## Processor

O `processor` consome mensagens da fila `raw-events`, valida o payload, enriquece o evento com `processed_at` e `processor_id`, publica em `processed-events` e deleta a mensagem original somente depois do publish bem-sucedido.

### Configuração

Variáveis usadas pelo serviço:

```text
PROCESSOR_ID=processor-1
QUEUE_URL=http://localhost:4566/
RAW_QUEUE_URL=http://localhost:4566/000000000000/raw-events
PROCESSED_QUEUE_URL=http://localhost:4566/000000000000/processed-events
REGION=us-east-1
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
WORKER_COUNT=5
```

Dentro do Docker Compose, `QUEUE_URL`, `RAW_QUEUE_URL` e `PROCESSED_QUEUE_URL` devem apontar para `http://localstack:4566`, não para `localhost`.

### Testes

```bash
cd services/processor
GOCACHE=/tmp/go-build go test ./...
```

### Seed de Mensagens

Com o LocalStack rodando:

```bash
./scripts/seed.sh
```

O script publica mensagens válidas, inválidas e duplicadas na fila `raw-events`.

### Verificar Fila Processada

```bash
aws --endpoint-url=http://localhost:4566 sqs receive-message \
  --queue-url http://localhost:4566/000000000000/processed-events \
  --max-number-of-messages 10
```

### Garantias Atuais

- `WORKER_COUNT <= 0` é normalizado para `1`.
- O channel de jobs é fechado pelo consumer no shutdown, evitando envio em channel fechado.
- O worker pool aguarda os workers terminarem.
- Logs são emitidos em JSON via `log/slog` e incluem `event_id` quando o payload pode ser lido.
- Falhas de receive e publish usam backoff exponencial.
