# 🍔 GoLunch API

API desenvolvida em Go para gerenciamento de pedidos em uma lanchonete. A arquitetura da aplicação segue princípios da arquitetura hexagonal, com foco na separação entre domínio e infraestrutura.

## 🧰 Tecnologias Utilizadas

- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin) – Framework HTTP
- [GORM](https://gorm.io/) – ORM para Go
- [Docker](https://www.docker.com/) – Containerização
- [PostgreSQL](https://www.postgresql.org/) – Banco de dados relacional

## 🚀 Inicialização do Projeto Localmente

### Pré-requisitos

- Go 1.20+
- Docker e Docker Compose

### Passos

1. Clone o repositório:

```bash
git clone https://github.com/seu-usuario/lanchonete-api.git
cd lanchonete-api
```

2. Copie o arquivo de variáveis de ambiente:

```bash
cp .env.example .env
```

3. Suba os containers com Docker Compose:

```bash
docker-compose up --build
```

4. Acesse a aplicação:

A API estará disponível em `http://localhost:8080`.

## 📌 Endpoints

[//]: # (Adicionar aqui os endpoints conforme forem sendo desenvolvidos.)

<!-- Placeholder para endpoints -->

## 🧠 Modelagem do Sistema

### Event Storming (Miro)

[🔗 Link para o Miro](https://miro.com/app/board/uXjVI47kj_s=/?share_link_id=805239820203)

### Entidades (Diagrama Draw.io)

[🔗 Link para o Diagrama no Draw.io](https://drive.google.com/file/d/1JbteJHGAyQ__yRhp25sq0pfO-bhE2edP/view)

### Diagrama de Entidades

![Diagrama de Entidades](![alt text](image.png))

> ℹ️ O diagrama acima mostra as relações entre os usuários, pedidos, produtos e pagamentos dentro do sistema.

## 📂 Estrutura do Projeto

```
├── cmd/                # Arquivo principal de entrada da aplicação
├── internal/           # Domínio e casos de uso
│   ├── domain/         # Entidades e regras de negócio
│   ├── usecase/        # Casos de uso
│   └── infra/          # Adaptadores externos (DB, Web, etc)
├── pkg/                # Pacotes utilitários
├── api/                # Handlers HTTP
├── configs/            # Configurações (ex: env)
├── docs/               # Documentação e imagens
└── Dockerfile / docker-compose.yml
```

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
