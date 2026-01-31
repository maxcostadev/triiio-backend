package email

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiErrors "github.com/vahiiiid/go-rest-api-boilerplate/internal/errors"
)

// Handler gerencia as requisições HTTP relacionadas a emails
type Handler struct {
	service Service
}

// NewHandler cria uma nova instância do handler de email
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// SendEmail envia um email simples
// @Summary Send email
// @Description Send a simple email to one or more recipients
// @Tags emails
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SendEmailRequest true "Email data"
// @Success 200 {object} errors.Response{success=bool,data=EmailResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 401 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/emails/send [post]
func (h *Handler) SendEmail(c *gin.Context) {
	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	result, err := h.service.SendEmail(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(result))
}

// SendTemplateEmail envia um email usando um template HTML
// @Summary Send template email
// @Description Send an email using a predefined HTML template
// @Tags emails
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SendTemplateEmailRequest true "Template email data"
// @Success 200 {object} errors.Response{success=bool,data=EmailResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 401 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/emails/send-template [post]
func (h *Handler) SendTemplateEmail(c *gin.Context) {
	var req SendTemplateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	result, err := h.service.SendTemplateEmail(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(result))
}
