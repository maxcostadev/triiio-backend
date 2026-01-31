package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"embed"
	"fmt"
	"html/template"
	"time"

	"github.com/vahiiiid/go-rest-api-boilerplate/internal/config"
	"github.com/vahiiiid/go-rest-api-boilerplate/internal/errors"
	mail "github.com/wneessen/go-mail"
)

//go:embed templates/*.html
var templatesFS embed.FS

// Service define a interface do serviço de email
type Service interface {
	SendEmail(ctx context.Context, req *SendEmailRequest) (*EmailResponse, error)
	SendTemplateEmail(ctx context.Context, req *SendTemplateEmailRequest) (*EmailResponse, error)
}

type service struct {
	cfg       *config.Config
	templates map[string]*template.Template
}

// NewService cria uma nova instância do serviço de email
func NewService(cfg *config.Config) (Service, error) {
	s := &service{
		cfg:       cfg,
		templates: make(map[string]*template.Template),
	}

	// Carrega os templates HTML
	if err := s.loadTemplates(); err != nil {
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}

	return s, nil
}

// loadTemplates carrega todos os templates HTML do embed.FS
func (s *service) loadTemplates() error {
	templateNames := []string{"default", "welcome", "notification"}

	for _, name := range templateNames {
		tmplPath := fmt.Sprintf("templates/%s.html", name)
		content, err := templatesFS.ReadFile(tmplPath)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", name, err)
		}

		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", name, err)
		}

		s.templates[name] = tmpl
	}

	return nil
}

// SendEmail envia um email simples
func (s *service) SendEmail(ctx context.Context, req *SendEmailRequest) (*EmailResponse, error) {
	// Validação das configurações de email
	if err := s.validateConfig(); err != nil {
		return nil, err
	}

	// Cria o cliente SMTP
	client, err := s.createSMTPClient()
	if err != nil {
		return nil, errors.InternalServerError(fmt.Errorf("failed to create SMTP client: %w", err))
	}
	defer client.Close()

	// Cria a mensagem
	msg := mail.NewMsg()

	// Define o remetente
	if err := msg.From(s.cfg.Email.From); err != nil {
		return nil, errors.InternalServerError(fmt.Errorf("failed to set from address: %w", err))
	}

	// Define os destinatários
	if err := msg.To(req.To...); err != nil {
		return nil, errors.BadRequest("Invalid 'to' addresses")
	}

	// Define CC se fornecido
	if len(req.Cc) > 0 {
		if err := msg.Cc(req.Cc...); err != nil {
			return nil, errors.BadRequest("Invalid 'cc' addresses")
		}
	}

	// Define BCC se fornecido
	if len(req.Bcc) > 0 {
		if err := msg.Bcc(req.Bcc...); err != nil {
			return nil, errors.BadRequest("Invalid 'bcc' addresses")
		}
	}

	// Define o assunto
	msg.Subject(req.Subject)

	// Define o corpo do email
	if req.IsHTML {
		msg.SetBodyString(mail.TypeTextHTML, req.Body)
	} else {
		msg.SetBodyString(mail.TypeTextPlain, req.Body)
	}

	// Envia o email
	if err := client.DialAndSend(msg); err != nil {
		return nil, errors.InternalServerError(fmt.Errorf("failed to send email: %w", err))
	}

	return &EmailResponse{
		Success: true,
		SentTo:  req.To,
		Message: "Email sent successfully",
	}, nil
}

// SendTemplateEmail envia um email usando um template HTML
func (s *service) SendTemplateEmail(ctx context.Context, req *SendTemplateEmailRequest) (*EmailResponse, error) {
	// Validação das configurações de email
	if err := s.validateConfig(); err != nil {
		return nil, err
	}

	// Verifica se o template existe
	tmpl, exists := s.templates[req.TemplateName]
	if !exists {
		return nil, errors.BadRequest(fmt.Sprintf("Template '%s' not found", req.TemplateName))
	}

	// Renderiza o template
	var body bytes.Buffer
	if req.TemplateData == nil {
		req.TemplateData = make(map[string]interface{})
	}

	// Adiciona dados padrão ao template
	req.TemplateData["Year"] = time.Now().Year()
	req.TemplateData["AppName"] = s.cfg.App.Name

	if err := tmpl.Execute(&body, req.TemplateData); err != nil {
		return nil, errors.InternalServerError(fmt.Errorf("failed to render template: %w", err))
	}

	// Cria a requisição de email com o corpo renderizado
	emailReq := &SendEmailRequest{
		To:      req.To,
		Cc:      req.Cc,
		Bcc:     req.Bcc,
		Subject: req.Subject,
		Body:    body.String(),
		IsHTML:  true,
	}

	return s.SendEmail(ctx, emailReq)
}

// createSMTPClient cria e configura o cliente SMTP
func (s *service) createSMTPClient() (*mail.Client, error) {
	options := []mail.Option{
		mail.WithPort(s.cfg.Email.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(s.cfg.Email.Username),
		mail.WithPassword(s.cfg.Email.Password),
		mail.WithTimeout(30 * time.Second),
	}

	// Configura TLS se habilitado
	if s.cfg.Email.UseTLS {
		tlsConfig := &tls.Config{
			ServerName:         s.cfg.Email.Host,
			InsecureSkipVerify: false,
		}
		options = append(options, mail.WithTLSConfig(tlsConfig))

		// Se StartTLS estiver habilitado, usa essa opção
		if s.cfg.Email.UseStartTLS {
			options = append(options, mail.WithTLSPolicy(mail.TLSMandatory))
		} else {
			options = append(options, mail.WithSSL())
		}
	}

	client, err := mail.NewClient(s.cfg.Email.Host, options...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// validateConfig valida as configurações de email
func (s *service) validateConfig() error {
	if s.cfg.Email.Host == "" {
		return errors.InternalServerError(fmt.Errorf("email host not configured"))
	}
	if s.cfg.Email.Port == 0 {
		return errors.InternalServerError(fmt.Errorf("email port not configured"))
	}
	if s.cfg.Email.From == "" {
		return errors.InternalServerError(fmt.Errorf("email from address not configured"))
	}
	if s.cfg.Email.Username == "" {
		return errors.InternalServerError(fmt.Errorf("email username not configured"))
	}
	if s.cfg.Email.Password == "" {
		return errors.InternalServerError(fmt.Errorf("email password not configured"))
	}
	return nil
}
