# Rate Limiter

### Visão Geral

Este Rate Limiter é projetado para controlar o tráfego e impedir o abuso de endpoints limitando o número de requisições permitidas. Ele funciona com dois identificadores:
- **IP**: bloqueio baseado no endereço IP.
- **Token de Acesso**: bloqueio baseado em um token específico de acesso.

### Requisitos

- [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados
- [Go](https://golang.org/dl/) versão 1.19 ou superior (opcional, caso deseje rodar o projeto fora do contêiner)

### Configuração

O sistema de Rate Limiting é configurado com variáveis de ambiente que você pode definir no `docker-compose.yml`. Aqui estão as variáveis principais e seus propósitos:

- `REDIS_ADDR`: endereço do Redis, geralmente `redis:6379` quando rodando em contêiner.
- `REQ_LIMIT`: número máximo de requisições permitidas antes de bloquear.
- `BLOCK_TIME_IP`: tempo de bloqueio em segundos para o limite de IP.
- `BLOCK_TIME_TOKEN`: tempo de bloqueio em segundos para o limite de token.

### Execução

Para rodar o Rate Limiter e suas dependências, siga os passos abaixo:

1. Clone o repositório:
   ```bash
   git clone https://github.com/seu-usuario/rate-limiter.git
   cd rate-limiter
   ```
2. Execute o Docker Compose para iniciar o serviço:
   ```bash
   docker compose up --build
   ```
Isso inicia o serviço rate-limiter na porta 8080 e um contêiner do Redis.
3. Acesse o serviço em:
   ```bash
    http://localhost:8080
   ```
### Testes

Para executar os testes, o projeto inclui um contêiner dedicado que inicializa o Redis e executa automaticamente os testes.
Os testes cobrem a lógica do Rate Limiter e verificam os métodos principais de controle de requisições.

