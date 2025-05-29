# 🍔 GoLunch API

API desenvolvida em Go para gerenciamento de pedidos em uma lanchonete. A arquitetura da aplicação segue princípios da arquitetura hexagonal, com foco na separação entre os domínios.

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
git clone https://github.com/fiap-161/tech-challenge-fiap161.git
cd tech-challenge-fiap161
```

2. Certifique-se que o docker está em execução:
   
```bash
docker ps
```

3. Crie um arquivo com as variáveis de ambiente:

```bash
DATABASE_URL=postgres://pg:pg@postgres-db:5432/pg?sslmode=disable
POSTGRES_USER=pg
POSTGRES_PASSWORD=pg
POSTGRES_DB=pg
SECRET_KEY=random_key
UPLOAD_DIR=./uploads
PUBLIC_URL=http://localhost:8080 
```

4. Suba os containers com Docker Compose:

```bash
docker-compose up --build
```

5. Acesse a aplicação:

A API estará disponível em `http://localhost:8080`.

6. Troubleshoot:
   - Em caso de falhas para subir a aplicação é válido tentar derrubar os containers e volumes criados previamente
     
```bash
docker-compose down -v --remove-orphans
```

## 📌 Swagger
O link para a documentação do swagger está aqui: http://localhost:8080/swagger/index.html

## 🧠 Modelagem do Sistema

### Event Storming (Miro)

[🔗 Link para o Miro](https://miro.com/app/board/uXjVI47kj_s=/?share_link_id=805239820203)

### Entidades (Diagrama Draw.io)

[🔗 Link para o Diagrama no Draw.io](https://drive.google.com/file/d/1JbteJHGAyQ__yRhp25sq0pfO-bhE2edP/view)

### Diagrama de Entidades

![lanchonete-fase01drawio drawio](https://github.com/user-attachments/assets/7d68f06b-056a-4252-9608-f2bea084c8cc)


> ℹ️ O diagrama acima mostra as relações entre os usuários, pedidos, produtos e pagamentos dentro do sistema.

## 📂 Estrutura do Projeto
```
├── cmd/                    # Arquivo principal de entrada da aplicação
│   └── api/
│       └── main.go
├── internal/               # Domínio, regras de negócio e adaptadores
│   ├── http/               # Camada HTTP (middlewares compartilhados)
│   └── dominio/            # Um diretório para cada domínio
│       ├── adapters/       # Adaptadores (drivers/drivens)
│       │   ├── drivens/    # Infraestrutura externa (DB)
│       │   └── drivers/    # Interface com frameworks (HTTP)
│       ├── core/           # Núcleo do domínio do produto
│       │   ├── model/      # Modelos e entidades do domínio
│       │   └── ports/      # Interfaces (portas) para repository e services
│       └── services/       # Lógica de aplicação (casos de uso)
├── shared/                 # Componentes compartilhados entre domínios
├── uploads/                # Diretório para salvar imagens
├── docs/                   # Documentação swagger
├── .env                    # Arquivo de variáveis de ambiente
├── .env.example            # Exemplo de variáveis de ambiente
├── docker-compose.yml      # Orquestração com Docker
└──  Dockerfile              # Docker build da aplicação
```
## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
