
# 🍔 GoLunch API - GRUPO 275

  

API desenvolvida em Go para gerenciamento de pedidos em uma lanchonete. A arquitetura da aplicação segue princípios da arquitetura hexagonal, com foco na separação entre os domínios.

  

### 🎥 Link para o vídeo detalhando o projeto: https://www.youtube.com/watch?v=ujhWQOJ0Jo0

  

## 🧰 Tecnologias Utilizadas


- [Go](https://golang.org/)

- [Gin](https://github.com/gin-gonic/gin) – Framework HTTP

- [GORM](https://gorm.io/) – ORM para Go

- [Docker](https://www.docker.com/) – Containerização

- [PostgreSQL](https://www.postgresql.org/) – Banco de dados relacional

- [Kubernetes](https://kubernetes.io/pt-br/) - Orquestrador de containers

- [Fortio](https://fortio.org/) - Teste de estresse
  

## 🏛️ [Link Excalidraw - Arquitetura k8s + Fluxos funcionais](https://excalidraw.com/#room=19187e25c8f502969730,UYsX9MelEMWQAT8VN4Marg)

  

### Arquitetura Kubernetes

<img width="805" height="765" alt="image" src="https://github.com/user-attachments/assets/d04c4f4c-a54f-4f0b-9fce-01235d12ad92" />


### Fluxo de criação de pedido

<img width="1062" height="602" alt="image" src="https://github.com/user-attachments/assets/be42b3db-19d1-4939-a212-c48a230717de" />

### Fluxo de pagamento

<img width="1151" height="402" alt="image" src="https://github.com/user-attachments/assets/0f0f6963-e28e-452c-9d75-83df268814e1" />

### Fluxo de atualização de pedido

<img width="643" height="384" alt="image" src="https://github.com/user-attachments/assets/bddefb7e-a2c0-4a9c-91c5-7e348997ac45" />

## Desenho da arquitetura

  

## 🏗️ Arquitetura Limpa (Clean Architecture)

  

Este projeto implementa os princípios da **Arquitetura Limpa** (Clean Architecture), organizando o código em camadas bem definidas para garantir separação de responsabilidades, testabilidade e manutenibilidade.

  

### Estrutura das Camadas

  

#### **Entities (Entidades)**

-  **Localização**: `internal/{domain}/entity/`

-  **Responsabilidade**: Contém as regras de negócio fundamentais e estruturas de dados principais

-  **Exemplo**: `internal/product/entity/product.go` - Define a estrutura do produto e suas validações básicas

  

#### **Use Cases (Casos de Uso)**

-  **Localização**: `internal/{domain}/usecases/`

-  **Responsabilidade**: Contém a lógica de negócio específica da aplicação

-  **Exemplo**: `internal/product/usecases/usecases.go` - Implementa operações como criar, atualizar, buscar produtos

  

#### **Gateways (Portões/Interfaces)**

-  **Localização**: `internal/{domain}/gateway/`

-  **Responsabilidade**: Interfaces que abstraem o acesso a dados externos

-  **Exemplo**: `internal/product/gateway/gateway.go` - Abstrai operações de persistência de dados

  

#### **Controllers (Controladores)**

-  **Localização**: `internal/{domain}/controller/`

-  **Responsabilidade**: Coordena a interação entre as camadas, criando gateways e executando casos de uso

-  **Exemplo**: `internal/product/controller/controller.go` - Orquestra operações de produtos

  

#### **Handlers (Manipuladores Web)**

-  **Localização**: `internal/{domain}/handler/`

-  **Responsabilidade**: Gerencia requisições HTTP, validações de entrada e respostas

-  **Exemplo**: `internal/product/handler/handler.go` - Endpoints REST para produtos

  

#### **External/Infrastructure (Infraestrutura Externa)**

-  **Localização**: `internal/{domain}/external/`

-  **Responsabilidade**: Implementações concretas de interfaces externas (banco de dados, APIs, etc.)

-  **Exemplo**: `internal/product/external/datasource/` - Implementação com GORM para PostgreSQL

  
  

## 📁 Estrutura de Diretórios

  

```

.

├── cmd/

│ └── api/

│ └── main.go

├── conf/

│ └── environment/

│ └── default.yml

├── database/

│ ├── database.go

│ └── postgre.go

├── docs/

├── internal/

│ ├── admin/ # Domínio de administração

│ │ ├── controller/

│ │ ├── dto/

│ │ ├── entity/

│ │ ├── external/datasource/

│ │ ├── gateway/

│ │ ├── handler/

│ │ ├── usecases/

│ │ └── utils/

│ └── shared/ # Código compartilhado

│ ├── entity/

│ ├── errors/

│ └── helper/

├── k8s/ # Manifestos Kubernetes

│ ├── app-deployment.yaml

│ ├── app-service.yaml

│ └── app-uploads-pvc.yaml

│ ├── configmap.yaml

│ └── fortio-stress-job.yaml

│ ├── hpa.yaml

│ ├── postgre-statefulset.yaml

│ ├── postgre-service.yaml

│ └── secrets.yaml

├── uploads/

├── docker-compose.yml # Configuração Docker Compose

├── Dockerfile # Imagem Docker

├── go.mod # Dependências Go

├── go.sum # Checksums das dependências

└── Makefile

```

  

# 🚀 Guia: Rodando o projeto no Kind

  

Este guia explica como instalar e executar o projeto localmente usando **kind** e **Kubernetes**, incluindo configuração do **Metrics Server**, criação de recursos, exposição da aplicação, geração de carga com **Fortio** e monitoramento com **HPA**.

  

## ⚠️ IMPORTANTE

  

### 📊 **Configurações de Recursos: Teste vs Produção**

  

Os recursos estão **intencionalmente baixos** para demonstrar o **HPA (Horizontal Pod Autoscaler)** em ação. Isso permite ver facilmente o escalonamento automático durante os testes de carga.

  

#### 🧪 **Configuração Atual (Ideal para Testes de HPA)**

```yaml

# Configuração otimizada para demonstrar escalabilidade

resources:

requests:

cpu: "0.2"  # 200m - Baixo para triggerar HPA rapidamente

memory: 70Mi  # Baixo para demonstrar limitações

limits:

cpu: "0.3"  # 300m - Limite baixo força escalabilidade

memory: 70Mi  # Força o HPA a criar novos pods

```

  

**Vantagens desta configuração:**

- ✅ HPA escala rapidamente durante teste de carga

- ✅ Demonstra claramente os benefícios do auto-scaling

- ✅ Simula ambiente com recursos limitados

  

#### 🚀 **Configuração para Produção (Opcional)**

Se quiser usar em produção, ajuste os recursos:

  

**Arquivo**: `k8s/app-deployment.yaml`

```yaml

resources:

requests:

cpu: "500m"

memory: "256Mi"

limits:

cpu: "1000m"

memory: "512Mi"

```

  

**Arquivo**: `k8s/hpa.yaml`

```yaml

spec:

minReplicas: 2

maxReplicas: 10

metrics:

-  type: Resource

resource:

name: cpu

target:

type: Utilization

averageUtilization: 70

```

  

---

  

Antes de executar o projeto, certifique-se de configurar as seguintes variáveis:

  

### 1. WEBHOOK_URL

Altere a variável `WEBHOOK_URL` para um link novo que deverá gerar aqui: https://webhook.site

  

**Arquivo a editar**: `k8s/configmap.yaml`

```yaml

data:

WEBHOOK_URL: "https://webhook.site/SEU-NOVO-LINK-AQUI"

```

  

### 2. Variáveis do Mercado Pago

Altere as variáveis do Mercado Pago para as descritas no documento PDF que foi enviado na entrega.

  

**Arquivo a editar**: `k8s/secrets.yaml`

```yaml

stringData:

MERCADO_PAGO_ACCESS_TOKEN: "SEU_ACCESS_TOKEN_AQUI"

MERCADO_PAGO_SELLER_APP_USER_ID: "SEU_USER_ID_AQUI"

MERCADO_PAGO_EXTERNAL_POS_ID: "SEU_POS_ID_AQUI"

```

  

### 3. Credenciais do Banco de Dados (Opcional)

Se desejar alterar as credenciais padrão do PostgreSQL:

  

**Arquivo a editar**: `k8s/secrets.yaml`

```yaml

stringData:

DATABASE_URL: "postgres://seu_usuario:sua_senha@postgres:5432/seu_banco"

POSTGRES_USER: "seu_usuario"

POSTGRES_PASSWORD: "sua_senha"

POSTGRES_DB: "seu_banco"

```

  

---

  

## 📦 Pré-requisitos

  
  

- [kind](https://kind.sigs.k8s.io/) instalado

- [kubectl](https://kubernetes.io/docs/tasks/tools/) instalado e configurado

  

- Manifestos YAML disponíveis:

  

-  `secrets.yaml`

-  `configmap.yaml`

-  `postgre-statefulset.yaml`

-  `postgre-service.yaml`

-  `app-uploads-pvc.yaml`

-  `app-deployment.yaml`

-  `app-service.yaml`

-  `hpa.yaml`

-  `fortio-stress-job.yaml`

  

---

  

## 1️⃣ Criar o cluster kind

  
  

```bash

kind  create  cluster  --name  meu-cluster

  

kubectl  get  nodes

```

  

---

  

## 2️⃣ Instalar o Metrics Server

  

Necessário para o HPA baseado em CPU/memória.

  

```bash

kubectl  apply  -f  https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

  
  

kubectl  patch  deployment  metrics-server  -n  kube-system  --type='json'  -p='[

  

{

  

"op": "add",

  

"path": "/spec/template/spec/containers/0/args/-",

  

"value": "--kubelet-insecure-tls"

  

}

  

]'

```

  

Verificar instalação

  

```bash

kubectl  get  pods  -n  kube-system  |  grep  metrics-server

  

kubectl  top  nodes

```

  

## 3️⃣ Criar Secrets e ConfigMap

  

```bash

kubectl  apply  -f  secrets.yaml

  

kubectl  apply  -f  configmap.yaml

```

## 4️⃣ Subir o PostgreSQL

  

```bash

kubectl  apply  -f  postgre-statefulset.yaml

  

kubectl  apply  -f  postgre-service.yaml

```

  

Verificar a instalação:

  

```bash

kubectl  get  pods  -l  app=postgres

```

  

## 5️⃣ Criar volume de upload

  

```bash

kubectl  apply  -f  app-uploads-pvc.yaml

```

  

## 6️⃣ Subir a aplicação e expor porta para uso local

  

### 6.1 Deployment e Service

  

```bash

kubectl  apply  -f  app-deployment.yaml

kubectl  apply  -f  app-service.yaml

```

  

### 6.2 Verificar pods

  

```bash

kubectl  get  pods  -l  app=go-web-api

```

  

## 7️⃣ Criar o HPA

  

```bash

kubectl  apply  -f  hpa.yaml

kubectl  get  hpa  go-web-api-hpa  ## verify hpa status

kubectl  describe  hpa  go-web-api-hpa  # describe hpa info

```

  

## 8️⃣ Gerar carga com Fortio

  

```bash

kubectl  apply  -f  fortio-stress-job.yaml

kubectl  get  jobs

kubectl  logs  job/fortio-stress-job

```

  

----------

  

## 9️⃣ Monitorar escalonamento em tempo real

  

### 📺 **Para o Vídeo - Comandos Essenciais**

  

Supondo que você não possua o **watch**, é possível rodar os comandos abaixo removendo o primeiro comando.

  

**Em terminais separados** (recomendado para demonstração):

```bash

# Terminal 1: Monitorar HPA (mostra CPU%, target, replicas)

watch  kubectl  get  hpa  go-web-api-hpa

  

# Terminal 2: Monitorar pods (mostra pods sendo criados/removidos)

watch  kubectl  get  pods  -l  app=go-web-api

  

# Terminal 3: Monitorar recursos dos pods (mostra uso real de CPU/memória)

watch  kubectl  top  pods  -l  app=go-web-api

```

  

### 🎯 **O que observar durante o teste:**

  

1.  **Antes do teste de carga**:

- HPA mostra baixo uso de CPU (< 40%)

- Apenas 1-3 pods rodando

  

2.  **Durante o teste de carga**:

- CPU sobe rapidamente para 100%+

- HPA começa a escalar (TARGETS aumenta)

- Novos pods aparecem com status `Pending` → `Running`

  

3.  **Após o teste**:

- CPU diminui gradualmente

- HPA escala para baixo (com delay de 5 segundos configurado)

## 🔟 Acessar a aplicação localmente

  

### Port-forward - Mapeamento de porta

  

```bash

kubectl  port-forward  svc/go-web-api-service  8080:8080

```

  

```bash

curl  http://localhost:8080/ping

```

  

----------

  

## 1️⃣1️⃣ Limpeza

  

Caso queira fazer a deleção do cluster, basta rodar o seguinte comando:

  

```bash

kind  delete  cluster  --name  meu-cluster

```
