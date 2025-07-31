# 🍔 GoLunch API

API desenvolvida em Go para gerenciamento de pedidos em uma lanchonete. A arquitetura da aplicação segue princípios da arquitetura hexagonal, com foco na separação entre os domínios.

### Link para o vídeo detalhando o projeto: https://www.youtube.com/watch?v=Il2WhYLpHsw

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
- Ter uma conta de testes no Mercado Pago (serão enviadas credenciais de teste no arquivo da entrega, utilize-as para logar no app do Mercado Pago)

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
cp .env.example .env
```
IMPORTANTE
- Altere a variável WEBHOOK_URL para um link novo que deverá gerar aqui: https://webhook.site
- Também altere as variáveis do Mercado Pago para as descritas no documento PDF que foi enviado na entrega.
- Para gerar o QRCode (explicado no vídeo) pode-se utilizar esse site: https://www.qr-code-generator.com/

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

## ☸️ Inicialização do Projeto no Kubernetes (Minikube)

### Pré-requisitos

- Minikube instalado
- kubectl instalado e configurado
- Docker instalado

### Passos

1. **Iniciar o minikube:**
```bash
minikube start
```

2. **Verificar se o cluster está funcionando:**
```bash
kubectl cluster-info
```

3. **Construir e carregar a imagem Docker no minikube:**
```bash
# Configurar docker para usar o registry do minikube
eval $(minikube docker-env)

# Construir a imagem
docker build -t golunch-app:latest .
```

4. **Configurar secrets (IMPORTANTE):**
```bash
# Editar o arquivo k8s/secrets.yaml com suas credenciais
# Encode os valores em base64:
echo -n "sua-secret-key-jwt" | base64
echo -n "sua-senha-database" | base64
echo -n "seu-token-mercadopago" | base64
```

5. **Aplicar os recursos do Kubernetes:**
```bash
# IMPORTANTE: Aplicar em ordem específica para evitar erros

# 1. Primeiro criar o namespace
kubectl apply -f k8s/namespace.yaml

# 2. Aguardar o namespace estar pronto
kubectl wait --for=condition=Active namespace/tech-challenge-fiap161 --timeout=60s

# 3. Aplicar configurações
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml

# 4. Aplicar deployments e services
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml
kubectl apply -f k8s/app-deployment.yaml
kubectl apply -f k8s/app-service.yaml

# 5. Aplicar HPA por último
kubectl apply -f k8s/hpa.yaml

# ALTERNATIVA: Use kustomize (se disponível)
kubectl apply -k k8s/
```

6. **Verificar se os pods estão rodando:**
```bash
kubectl get pods -n tech-challenge-fiap161
```

7. **Acessar a aplicação:**
```bash
# Obter URL de acesso
minikube service golunch-app-service -n tech-challenge-fiap161 --url
```

### Troubleshoot Kubernetes

```bash
# Ver logs da aplicação
kubectl logs -f deployment/golunch-app-deployment -n tech-challenge-fiap161

# Ver logs do banco de dados
kubectl logs -f deployment/postgres-deployment -n tech-challenge-fiap161

# Verificar status do HPA
kubectl get hpa -n tech-challenge-fiap161

# Acessar o banco diretamente
kubectl exec -it deployment/postgres-deployment -n tech-challenge-fiap161 -- psql -U golunch_user -d golunch
```

## 📌 Swagger
O link para a documentação do swagger está aqui: http://localhost:8080/swagger/index.html (Docker) ou utilize a URL fornecida pelo minikube service

## 🧠 Modelagem do Sistema

### Event Storming (Miro)

[🔗 Link para o Miro](https://miro.com/app/board/uXjVI47kj_s=/?share_link_id=805239820203)

### Entidades (Diagrama Draw.io)

[🔗 Link para o Diagrama no Draw.io](https://drive.google.com/file/d/1JbteJHGAyQ__yRhp25sq0pfO-bhE2edP/view)

### Diagrama de Entidades

![image](https://github.com/user-attachments/assets/aac0e29d-3546-4cda-ac6b-a7c78a867dec)



> ℹ️ O diagrama acima mostra as relações entre os usuários, pedidos, produtos e pagamentos dentro do sistema.

## 📂 Estrutura do Projeto
```
├── cmd/                    # Arquivo principal de entrada da aplicação
│   └── api/
│       └── main.go
├── internal/               # Domínio, regras de negócio e adaptadores
│   ├── http/               # Camada HTTP (middlewares compartilhados)
│   ├── shared/             # Componentes compartilhados entre domínios
│   └── dominio/            # Um diretório para cada domínio
│       ├── adapters/       # Adaptadores (drivers/drivens)
│       │   ├── drivens/    # Infraestrutura externa (DB)
│       │   └── drivers/    # Interface com frameworks (HTTP)
│       ├── core/           # Núcleo do domínio do produto
│       │   ├── model/      # Modelos e entidades do domínio
│       │   └── ports/      # Interfaces (portas) para repository e services
│       └── services/       # Lógica de aplicação (casos de uso)
├── uploads/                # Diretório para salvar imagens
├── docs/                   # Documentação swagger
├── .env                    # Arquivo de variáveis de ambiente
├── .env.example            # Exemplo de variáveis de ambiente
├── docker-compose.yml      # Orquestração com Docker
└──  Dockerfile              # Docker build da aplicação
```

## 🛑 Como Parar Todos os Processos

### Para Minikube + Kubernetes

1. **Parar aplicação específica (recomendado):**
```bash
# Parar todos os recursos do namespace
kubectl delete all --all -n tech-challenge-fiap161

# Ou parar usando os arquivos de configuração
kubectl delete -f k8s/
# Ou com kustomize
kubectl delete -k k8s/
```

2. **Parar minikube completamente:**
```bash
# Parar o cluster minikube
minikube stop

# Deletar completamente o cluster (remove tudo)
minikube delete
```

### Para Docker Compose (se também estiver rodando)

```bash
# Parar e remover containers
docker-compose down -v --remove-orphans
```

### Verificar se tudo parou

```bash
# Verificar pods (deve estar vazio)
kubectl get pods -n tech-challenge-fiap161

# Verificar status do minikube
minikube status

# Verificar containers docker
docker ps
```

## Testes

Os testes podem ser executados com o comando:
> go test ./... 

# Coleção Postman
### Pode ser encontrada no arquivo:

```FIAP TC1.json```

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
