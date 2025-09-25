
# 🍔 Repositório da aplicação + manifestos k8s - GoLunch API

API desenvolvida em Go para gerenciamento de pedidos em uma lanchonete. A arquitetura da aplicação segue princípios da arquitetura hexagonal, com foco na separação entre os domínios.

### 🎥 Link para o vídeo detalhando o projeto: TO DO
  

## 🏛️ [Link Excalidraw - Arquitetura k8s + Fluxos funcionais](https://excalidraw.com/#room=19187e25c8f502969730,UYsX9MelEMWQAT8VN4Marg)


## 🔄 Fluxo de Trabalho (CI/CD)

Para garantir que a infraestrutura seja criada/atualizada corretamente via **GitHub Actions**, siga os passos abaixo:

1. Atualizar Secrets da Organização
  Antes de rodar o pipeline, verifique se as **secrets da organização** estão configuradas em:  
  [Configurações de Secrets](https://github.com/fiap-161/tc-golunch-infra/settings/secrets/actions)
  
  As secrets necessárias são:
  - `AWS_ACCESS_KEY_ID`
  - `AWS_SECRET_ACCESS_KEY`
  - `AWS_SESSION_TOKEN`
  - `DATABASE_USER`
  - `DATABASE_PASSWORD`
  - `MERCADO_PAGO_ACCESS_TOKEN`
  - `MERCADO_PAGO_SELLER_APP_USER_ID` 
  - `SECRET_KEY`

  Essas credenciais são utilizadas no `cd.yaml` para autenticar na AWS e pela app para se conectar ao banco de dados.

2. Criar uma Branch a partir da `main`, dar push nas alterações
3. Abrir um Pull Request.