# 🍔 GoLunch API

API desenvolvida em Go para gerenciamento de pedidos em uma lanchonete. A arquitetura da aplicação segue princípios da arquitetura hexagonal, com foco na separação entre os domínios.

### Link para o vídeo detalhando o projeto: https://www.youtube.com/watch?v=Il2WhYLpHsw

## 🧰 Tecnologias Utilizadas

- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin) – Framework HTTP
- [GORM](https://gorm.io/) – ORM para Go
- [Docker](https://www.docker.com/) – Containerização
- [PostgreSQL](https://www.postgresql.org/) – Banco de dados relacional

# 🍔 GoLunch API

  

API desenvolvida em Go para gerenciamento de pedidos em uma lanchonete. A arquitetura da aplicação segue princípios da arquitetura hexagonal, com foco na separação entre os domínios.

### Link para o vídeo detalhando o projeto: https://www.youtube.com/watch?v=Il2WhYLpHsw

## 🧰 Tecnologias Utilizadas

- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin) – Framework HTTP
- [GORM](https://gorm.io/) – ORM para Go
- [Docker](https://www.docker.com/) – Containerização
- [PostgreSQL](https://www.postgresql.org/) – Banco de dados relacional

# 🚀 Guia: Rodando o projeto no Kind

Este guia explica como instalar e executar o projeto localmente usando **kind** e **Kubernetes**, incluindo configuração do **Metrics Server**, criação de recursos, exposição da aplicação, geração de carga com **Fortio** e monitoramento com **HPA**.

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
kubectl apply -f app-uploads-pvc.yaml
```

## 6️⃣ Subir a aplicação e expor porta para uso local

### 6.1 Deployment e Service

```bash
kubectl apply -f app-deployment.yaml
kubectl apply -f app-service.yaml
```

### 6.2 Verificar pods

```bash
kubectl get pods -l app=go-web-api
```

## 7️⃣ Criar o HPA

```bash
kubectl apply -f hpa.yaml
kubectl get hpa go-web-api-hpa ## verify hpa status
kubectl describe hpa go-web-api-hpa # describe hpa info
```

## 8️⃣ Gerar carga com Fortio

```bash
kubectl apply -f fortio-stress-job.yaml
kubectl get jobs kubectl logs job/fortio-stress-job`
```

----------

## 9️⃣ Monitorar escalonamento em tempo real

Supondo que você não possua o **watch**, é possível rodar os comandos abaixo removendo o primeiro comando.

Em terminais separados:
```bash
watch kubectl get hpa go-web-api-hpa
watch kubectl get pods -l app=go-web-api
watch kubectl top pods -l app=go-web-api`
```
## 🔟 Acessar a aplicação localmente

###  Port-forward - Mapeamento de porta

```bash
kubectl port-forward svc/go-web-api-service 8080:8080
```

```bash
curl http://localhost:8080/ping
```

----------

## 1️⃣1️⃣ Limpeza

Caso queira fazer a deleção do cluster, basta rodar o seguinte comando:

```bash
kind delete cluster --name meu-cluster
```
