# Email Service - Documentação

## Visão Geral

O serviço de email permite envio de emails através de SMTP com suporte a templates HTML personalizados.

## Configuração

### Variáveis de Ambiente

Configure as seguintes variáveis de ambiente (ou no arquivo `configs/config.yaml`):

```bash
# SMTP Server
EMAIL_HOST=smtp.gmail.com              # Servidor SMTP
EMAIL_PORT=587                         # Porta SMTP (587 para TLS, 465 para SSL)
EMAIL_USERNAME=seu-email@gmail.com     # Usuário SMTP
EMAIL_PASSWORD=sua-senha-app           # Senha do aplicativo SMTP
EMAIL_FROM=noreply@seudominio.com      # Email do remetente

# TLS/SSL Settings
EMAIL_USE_TLS=true                     # Habilitar TLS/SSL
EMAIL_USE_STARTTLS=true                # Usar STARTTLS (recomendado para porta 587)
```

### Exemplos de Configuração por Provedor

#### Gmail
```bash
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=seu-email@gmail.com
EMAIL_PASSWORD=sua-senha-de-app  # Gere em: https://myaccount.google.com/apppasswords
EMAIL_USE_TLS=true
EMAIL_USE_STARTTLS=true
```

#### Outlook/Hotmail
```bash
EMAIL_HOST=smtp-mail.outlook.com
EMAIL_PORT=587
EMAIL_USERNAME=seu-email@outlook.com
EMAIL_PASSWORD=sua-senha
EMAIL_USE_TLS=true
EMAIL_USE_STARTTLS=true
```

#### SendGrid
```bash
EMAIL_HOST=smtp.sendgrid.net
EMAIL_PORT=587
EMAIL_USERNAME=apikey
EMAIL_PASSWORD=SG.sua-api-key-aqui
EMAIL_USE_TLS=true
EMAIL_USE_STARTTLS=true
```

#### Mailgun
```bash
EMAIL_HOST=smtp.mailgun.org
EMAIL_PORT=587
EMAIL_USERNAME=postmaster@seu-dominio.mailgun.org
EMAIL_PASSWORD=sua-senha-smtp
EMAIL_USE_TLS=true
EMAIL_USE_STARTTLS=true
```

## Endpoints da API

### 1. Enviar Email Simples

**POST** `/api/v1/emails/send`

Envia um email simples com texto ou HTML personalizado.

**Headers:**
```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "to": ["destinatario@example.com"],
  "cc": ["copia@example.com"],          // Opcional
  "bcc": ["copia-oculta@example.com"],  // Opcional
  "subject": "Assunto do Email",
  "body": "Conteúdo do email",
  "is_html": false                       // true para HTML, false para texto
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "success": true,
    "sent_to": ["destinatario@example.com"],
    "message": "Email sent successfully"
  }
}
```

**Exemplo com cURL:**
```bash
curl -X POST http://localhost:8080/api/v1/emails/send \
  -H "Authorization: Bearer {seu-token}" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["destinatario@example.com"],
    "subject": "Teste de Email",
    "body": "Este é um email de teste",
    "is_html": false
  }'
```

### 2. Enviar Email com Template

**POST** `/api/v1/emails/send-template`

Envia um email usando um template HTML pré-definido.

**Headers:**
```
Authorization: Bearer {jwt_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "to": ["destinatario@example.com"],
  "subject": "Bem-vindo!",
  "template_name": "welcome",
  "template_data": {
    "UserName": "João Silva",
    "UserEmail": "joao@example.com",
    "ButtonText": "Acessar Plataforma",
    "ButtonURL": "https://seusite.com/login"
  }
}
```

**Templates Disponíveis:**
- `default` - Template padrão para mensagens gerais
- `welcome` - Template de boas-vindas
- `notification` - Template para notificações

**Response:**
```json
{
  "success": true,
  "data": {
    "success": true,
    "sent_to": ["destinatario@example.com"],
    "message": "Email sent successfully"
  }
}
```

## Templates HTML

### Template: `default`

Template genérico para mensagens.

**Variáveis disponíveis:**
- `AppName` - Nome da aplicação (automático)
- `Year` - Ano atual (automático)
- `Title` - Título do email
- `Message` - Mensagem principal
- `Body` - Conteúdo adicional (HTML)
- `ButtonText` - Texto do botão (opcional)
- `ButtonURL` - URL do botão (opcional)
- `UnsubscribeURL` - URL para cancelar inscrição (opcional)

**Exemplo de uso:**
```json
{
  "to": ["usuario@example.com"],
  "subject": "Atualização Importante",
  "template_name": "default",
  "template_data": {
    "Title": "Nova Funcionalidade",
    "Message": "Temos o prazer de anunciar...",
    "ButtonText": "Ver Detalhes",
    "ButtonURL": "https://seusite.com/novidades"
  }
}
```

### Template: `welcome`

Template para dar boas-vindas a novos usuários.

**Variáveis disponíveis:**
- `AppName` - Nome da aplicação (automático)
- `Year` - Ano atual (automático)
- `UserName` - Nome do usuário
- `UserEmail` - Email do usuário
- `Message` - Mensagem adicional (opcional)
- `ButtonText` - Texto do botão (opcional)
- `ButtonURL` - URL do botão (opcional)

**Exemplo de uso:**
```json
{
  "to": ["novo-usuario@example.com"],
  "subject": "Bem-vindo ao TRIIIO!",
  "template_name": "welcome",
  "template_data": {
    "UserName": "Maria Santos",
    "UserEmail": "maria@example.com",
    "ButtonText": "Começar Agora",
    "ButtonURL": "https://seusite.com/onboarding"
  }
}
```

### Template: `notification`

Template para notificações e alertas.

**Variáveis disponíveis:**
- `AppName` - Nome da aplicação (automático)
- `Year` - Ano atual (automático)
- `Type` - Tipo de notificação: `warning`, `success`, `error` (opcional)
- `Title` - Título da notificação
- `Message` - Mensagem principal
- `AlertMessage` - Mensagem de alerta destacada (opcional)
- `Body` - Conteúdo adicional (HTML)
- `Details` - Mapa de detalhes chave-valor (opcional)
- `ButtonText` - Texto do botão (opcional)
- `ButtonURL` - URL do botão (opcional)
- `Timestamp` - Data/hora da notificação (opcional)

**Exemplo de uso:**
```json
{
  "to": ["admin@example.com"],
  "subject": "Alerta de Sistema",
  "template_name": "notification",
  "template_data": {
    "Type": "warning",
    "Title": "Uso de Memória Elevado",
    "Message": "O servidor atingiu 85% de uso de memória",
    "AlertMessage": "Ação requerida em até 2 horas",
    "Details": {
      "Servidor": "API-01",
      "Memória Usada": "85%",
      "Última Verificação": "31/01/2026 10:30"
    },
    "ButtonText": "Ver Métricas",
    "ButtonURL": "https://monitoring.seusite.com"
  }
}
```

## Personalização de Templates

Os templates estão localizados em `internal/email/templates/`:

- `default.html` - Template padrão
- `welcome.html` - Template de boas-vindas
- `notification.html` - Template de notificação

Para adicionar um novo template:

1. Crie um arquivo `.html` em `internal/email/templates/`
2. Adicione o nome do template em `service.go`:
   ```go
   templateNames := []string{"default", "welcome", "notification", "novo-template"}
   ```
3. Atualize a validação em `dto.go`:
   ```go
   TemplateName string `json:"template_name" binding:"required,oneof=default welcome notification novo-template"`
   ```

## Erros Comuns

### 1. Autenticação SMTP Falhou

**Erro:**
```
failed to send email: authentication failed
```

**Solução:**
- Verifique `EMAIL_USERNAME` e `EMAIL_PASSWORD`
- Para Gmail, use uma senha de aplicativo
- Verifique se a autenticação de dois fatores está configurada

### 2. Conexão Recusada

**Erro:**
```
failed to create SMTP client: dial tcp: connection refused
```

**Solução:**
- Verifique `EMAIL_HOST` e `EMAIL_PORT`
- Confirme que o firewall não está bloqueando a porta
- Teste conectividade: `telnet smtp.gmail.com 587`

### 3. TLS/SSL Error

**Erro:**
```
x509: certificate signed by unknown authority
```

**Solução:**
- Verifique `EMAIL_USE_TLS` e `EMAIL_USE_STARTTLS`
- Para porta 587, use `EMAIL_USE_STARTTLS=true`
- Para porta 465, use `EMAIL_USE_TLS=true` e `EMAIL_USE_STARTTLS=false`

### 4. Template Não Encontrado

**Erro:**
```
Template 'xyz' not found
```

**Solução:**
- Verifique o nome do template no request
- Templates disponíveis: `default`, `welcome`, `notification`
- Nome é case-sensitive

## Testes

### Teste Manual com cURL

```bash
# 1. Login para obter token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"senha"}' \
  | jq -r '.data.access_token')

# 2. Enviar email simples
curl -X POST http://localhost:8080/api/v1/emails/send \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["teste@example.com"],
    "subject": "Teste",
    "body": "Este é um teste",
    "is_html": false
  }'

# 3. Enviar email com template
curl -X POST http://localhost:8080/api/v1/emails/send-template \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["teste@example.com"],
    "subject": "Bem-vindo!",
    "template_name": "welcome",
    "template_data": {
      "UserName": "Teste",
      "UserEmail": "teste@example.com"
    }
  }'
```

## Biblioteca Utilizada

Esta implementação usa [go-mail](https://github.com/wneessen/go-mail), uma biblioteca moderna e robusta para envio de emails em Go com:

- ✅ Suporte completo a SMTP/SMTPS
- ✅ TLS/SSL e STARTTLS
- ✅ Múltiplos métodos de autenticação
- ✅ Templates HTML
- ✅ Anexos (preparado para expansão futura)
- ✅ Altamente estável e testada
- ✅ Ativamente mantida

## Próximos Passos

Funcionalidades planejadas para futuras versões:

- [ ] Suporte a anexos de arquivos
- [ ] Fila de emails assíncrona (Redis/RabbitMQ)
- [ ] Log de emails enviados
- [ ] Retry automático em caso de falha
- [ ] Preview de templates antes de enviar
- [ ] Estatísticas de envio
- [ ] Webhooks para eventos (aberto, clicado, etc.)

## Suporte

Para issues e sugestões, visite: [GitHub Issues](https://github.com/vahiiiid/go-rest-api-boilerplate/issues)
