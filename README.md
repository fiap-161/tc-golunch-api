# 🍔 GoLunch API

API desenvolvida em Go para gerenciamento de pedidos em uma lanchonete. A arquitetura da aplicação segue princípios da arquitetura hexagonal, com foco na separação entre os domínios.

### Link para o vídeo detalhando o projeto: https://www.youtube.com/watch?v=Il2WhYLpHsw

## 🧰 Tecnologias Utilizadas

- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin) – Framework HTTP
- [GORM](https://gorm.io/) – ORM para Go
- [Docker](https://www.docker.com/) – Containerização
- [PostgreSQL](https://www.postgresql.org/) – Banco de dados relacional

## 📦 Deploy com Kubernetes e Minikube

### ✅ Pré-requisitos

* [Minikube](https://minikube.sigs.k8s.io/docs/start/)
* [kubectl](https://kubernetes.io/docs/tasks/tools/)
* Habilitar o `metrics-server` do Minikube para utilizar HPA:

```bash
minikube addons enable metrics-server
```

---

### 🔪 Inicialização com Minikube

1. **Inicie o Minikube** (caso ainda não tenha iniciado):

```bash
minikube start
```

2. **Gere o Secret a partir do template**

Utilize o comando abaixo para criar o `secret.yaml` com os dados necessários (substitua pelos seus dados reais, se necessário):

IMPORTANTE
- Altere a variável WEBHOOK_URL para um link novo que deverá gerar aqui: https://webhook.site
- Também altere as variáveis do Mercado Pago para as descritas no documento PDF que foi enviado na entrega.
- Para gerar o QRCode (explicado no vídeo) pode-se utilizar esse site: https://www.qr-code-generator.com/

```bash
kubectl create secret generic app-secrets \
  --from-literal=DATABASE_URL=postgres://user:password@postgres:5432/dbname \
  --from-literal=POSTGRES_USER=user \
  --from-literal=POSTGRES_PASSWORD=password \
  --from-literal=POSTGRES_DB=dbname \
  --from-literal=SECRET_KEY=random_key \
  --from-literal=MERCADO_PAGO_ACCESS_TOKEN=APP_USR-8119906223498266-051516-0b1dc0cc2f9c6fb392955fb8e20dde55-2444053782 \
  --from-literal=MERCADO_PAGO_SELLER_APP_USER_ID=2444053782 \
  --from-literal=MERCADO_PAGO_EXTERNAL_POS_ID=DEFAULT \
  --dry-run=client -o yaml > secret.yaml
```

Aplique o secret:

```bash
kubectl apply -f secret.yaml
```

---

### 📂 Aplicando os manifestos do Kubernetes

Certifique-se de estar na pasta raiz onde os arquivos `YAML` estão localizados, nesse projeto, dentro da pasta `k8s`. Execute os comandos abaixo para aplicar os recursos:

```bash
kubectl apply -f configmap.yaml
kubectl apply -f postgres-service.yaml
kubectl apply -f postgres-statefulset.yaml
kubectl apply -f app-deployment.yaml
kubectl apply -f app-service.yaml
kubectl apply -f hpa.yaml
```

---

### 🌐 Acessando a aplicação

Exponha o serviço para acesso externo via Minikube:

```bash
minikube service go-web-service
```
IP  e porta da aplicação serão serão logados no terminal. 

---

### 🚰 Troubleshooting

* Verifique os pods:

```bash
kubectl get pods
```

* Verifique os logs da aplicação:

```bash
kubectl logs <nome-do-pod>
```

* Reinicie os recursos, se necessário:

```bash
kubectl delete -f <arquivo>.yaml
kubectl apply -f <arquivo>.yaml
```

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

## Testes

Os testes podem ser executados com o comando:
> go test ./... 

# Coleção Postman
### Pode ser encontrada no arquivo:

```FIAP TC1.json```

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
