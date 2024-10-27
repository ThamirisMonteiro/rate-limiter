# Rate Limiter em Go

## Objetivo
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

## Descrição
Este rate limiter controla o tráfego de requisições para um serviço web com base em dois critérios:
- **Endereço IP**: Limita o número de requisições recebidas de um endereço IP em um intervalo de tempo definido.
- **Token de Acesso**: Limita as requisições usando um token de acesso único, permitindo diferentes limites de expiração para cada token. O token deve ser enviado no header no formato: API_KEY: <TOKEN>

As configurações de limite do token de acesso têm prioridade sobre as do IP. Por exemplo, se o limite por IP é de 10 req/s e o de um token específico é de 100 req/s, o rate limiter deve usar o limite do token.

## Requisitos
1. Middleware injetável no servidor web para interceptar requisições.
2. Capacidade de configurar o número máximo de requisições permitidas por segundo.
3. Opção para definir o tempo de bloqueio de IP ou Token em caso de exceder o limite.
4. Configurações de limite definidas por variáveis de ambiente ou arquivo `.env` na raiz do projeto.
5. Configuração para limitação tanto por IP quanto por token de acesso.
6. Retorno apropriado quando o limite é excedido:
- **Código HTTP**: 429
- **Mensagem**: "you have reached the maximum number of requests or actions allowed within a certain time frame"
7. Armazenamento e consulta das informações de "limiter" em um banco de dados Redis, utilizando `docker-compose` para subir o Redis.
8. Estratégia para permitir trocar facilmente o Redis por outro mecanismo de persistência.
9. Lógica de limitação separada do middleware.

## Exemplos
- **Limitação por IP**: Se configurado para 5 req/s por IP, o IP `192.168.1.1` enviando 6 requisições em um segundo deve ter a sexta requisição bloqueada.
- **Limitação por Token**: Se um token `abc123` tem limite de 10 req/s, a décima primeira requisição no intervalo deve ser bloqueada.

Nos dois casos, as próximas requisições são permitidas apenas após o tempo de expiração definido. Exemplo: Se o tempo de expiração é 5 minutos, o IP/token pode realizar novas requisições somente após esse intervalo.

## Dicas
- Teste o rate limiter sob diferentes condições de carga para assegurar sua robustez e eficiência em alto tráfego.

## Entrega
- Código-fonte completo.
- Documentação explicando o funcionamento do rate limiter e como configurá-lo.
- Testes automatizados que demonstrem a eficácia do rate limiter em diferentes cenários.
- Utilização de `docker/docker-compose` para realizar testes da aplicação.
- O servidor web deve responder na porta `8080`.

## Checklist de Implementação
1. **Middleware**: Criação de um middleware para interceptar requisições.
2. **Lógica de Limitação**: Desenvolvimento da contagem de requisições para IP e tokens, priorizando tokens sobre IP.
3. **Persistência no Redis**: Uso de Redis para armazenar contagem e limites de requisições.
4. **Configurações**: Carregar as configurações via variáveis de ambiente ou `.env`.
4. **Estratégia de Persistência**: Definir uma interface para substituir facilmente o Redis por outro banco de dados.
4. **Resposta ao Exceder Limite**: Retornar erro `429` com a mensagem adequada.
4. **Docker-Compose**: Configuração de um ambiente Redis usando `docker-compose`. 
5. **Testes**: Implementação de testes automatizados para validação em diferentes condições.






