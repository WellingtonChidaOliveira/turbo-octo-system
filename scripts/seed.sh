# #Seed de Dados
# Inclua um script (scripts/seed.sh) que publique pelo menos 20 mensagens na fila raw-events via AWS CLI apontando para o LocalStack, incluindo:

# Mensagens válidas (vários developers e metric types)
# 2-3 mensagens inválidas (para mostrar validação e DLQ funcionando)
# 1-2 mensagens duplicadas (para mostrar idempotência no Aggregator)
# #