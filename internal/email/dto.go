package email

// SendEmailRequest representa a requisição para envio de email
type SendEmailRequest struct {
	To      []string `json:"to" binding:"required,min=1,dive,email"`
	Cc      []string `json:"cc" binding:"omitempty,dive,email"`
	Bcc     []string `json:"bcc" binding:"omitempty,dive,email"`
	Subject string   `json:"subject" binding:"required,min=1,max=500"`
	Body    string   `json:"body" binding:"required,min=1"`
	IsHTML  bool     `json:"is_html"`
}

// SendTemplateEmailRequest representa a requisição para envio de email com template
type SendTemplateEmailRequest struct {
	To           []string               `json:"to" binding:"required,min=1,dive,email"`
	Cc           []string               `json:"cc" binding:"omitempty,dive,email"`
	Bcc          []string               `json:"bcc" binding:"omitempty,dive,email"`
	Subject      string                 `json:"subject" binding:"required,min=1,max=500"`
	TemplateName string                 `json:"template_name" binding:"required,oneof=default welcome notification"`
	TemplateData map[string]interface{} `json:"template_data"`
}

// EmailResponse representa a resposta do envio de email
type EmailResponse struct {
	Success   bool     `json:"success"`
	MessageID string   `json:"message_id,omitempty"`
	SentTo    []string `json:"sent_to"`
	Message   string   `json:"message"`
}
