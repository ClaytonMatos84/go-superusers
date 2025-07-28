# Go Superusers API

API para gerenciamento e consulta de superusuários.

## Descrição
Este projeto implementa uma API RESTful em Go para manipulação de dados de superusuários, incluindo funcionalidades de listagem, filtragem e logs de ações.

## Estrutura do Projeto
- `cmd/api-server/`: Código principal do servidor HTTP
- `internal/model/`: Modelos de domínio e DTOs
- `internal/routers/`: Definição das rotas da API
- `internal/service/`: Serviços de negócio
- `pkg/`: Utilitários e helpers
- `deployments/`: Arquivos para deploy, Docker e scripts
- `docs/`: Documentação e exemplos de dados
- `test/postman/`: Coleções para testes com Postman

## Requisitos
- Go 1.21+
- Docker (opcional, para deploy)

## Como rodar localmente

> ### Usando Go
> ```bash
> git clone <repo-url>
> cd go-superusers
> go mod tidy
> go run ./cmd/api-server/main.go
> ```
> Acesse: http://localhost:8080/health

> ### Usando Makefile
> ```bash
> cd ./deployments
> make create
> ```
> Acesse: http://localhost:9797/health

> ### Usando Docker
> ```bash
> docker build -t go-superusers-api -f deployments/Dockerfile .
> docker run -p 8080:8080 go-superusers-api
> ```
> Acesse: http://localhost:8080/health

> ### Usando Docker Compose
> ```bash
> cd ./deployments
> docker-compose up --build
> ```
> Acesse: http://localhost:9797/health

## Testes
Utilize a coleção Postman disponível em `test/postman/` para testar os endpoints.

## Autor
Clayton Matos
