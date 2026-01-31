package imoveis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiErrors "github.com/vahiiiid/go-rest-api-boilerplate/internal/errors"
)

// Handler defines HTTP handlers for imovel operations
type Handler struct {
	service       Service
	importService ImportService
}

// NewHandler creates a new imovel handler
func NewHandler(service Service, importService ImportService) *Handler {
	return &Handler{
		service:       service,
		importService: importService,
	}
}

// @Summary Import properties from external API
// @Description Import all published properties from dev-api-backend.pi8.com.br. Uses upsert logic - creates new properties and updates existing ones based on id_integracao mapping. Existing properties are detected via id_integracao field and updated with latest data. Attachments are deduplicated by URL.
// @Tags imoveis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Import completed with statistics (created, updated, failed counts)"
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/imoveis/import [post]
func (h *Handler) ImportProperties(c *gin.Context) {
	if err := h.importService.ImportPublishedProperties(c.Request.Context()); err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Import completed",
	})
}

// @Summary Get property by ID
// @Description Get a property by its ID
// @Tags imoveis
// @Accept json
// @Produce json
// @Param id path uint true "Property ID"
// @Success 200 {object} errors.Response{success=bool,data=ImovelResponse}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/imoveis/{id} [get]
func (h *Handler) GetImovel(c *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	imovel, err := h.service.GetImovel(c.Request.Context(), req.ID)
	if err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	if imovel == nil {
		_ = c.Error(apiErrors.NotFound("Property not found"))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(imovel))
}

// @Summary Create a new property
// @Description Create a new property
// @Tags imoveis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateImovelRequest true "Property creation request"
// @Success 201 {object} errors.Response{success=bool,data=ImovelResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 409 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/imoveis [post]
func (h *Handler) CreateImovel(c *gin.Context) {
	var req CreateImovelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	imovel, err := h.service.CreateImovel(c.Request.Context(), &req)
	if err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, apiErrors.Success(imovel))
}

// @Summary Update a property
// @Description Update an existing property
// @Tags imoveis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "Property ID"
// @Param request body UpdateImovelRequest true "Property update request"
// @Success 200 {object} errors.Response{success=bool,data=ImovelResponse}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/imoveis/{id} [put]
func (h *Handler) UpdateImovel(c *gin.Context) {
	var uriReq struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	var req UpdateImovelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	imovel, err := h.service.UpdateImovel(c.Request.Context(), uriReq.ID, &req)
	if err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(imovel))
}

// @Summary Delete a property
// @Description Soft delete a property
// @Tags imoveis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "Property ID"
// @Success 204 "No Content"
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/imoveis/{id} [delete]
func (h *Handler) DeleteImovel(c *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	if err := h.service.DeleteImovel(c.Request.Context(), req.ID); err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary List properties
// @Description Get paginated list of properties with filters
// @Tags imoveis
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param codigo query string false "Property code (partial match)"
// @Param tipo query string false "Property type (APARTAMENTO, CASA, COMERCIAL, SALA_COMERCIAL, TERRENO, GALPAO)"
// @Param objetivo query string false "Property objective (VENDER, ALUGAR)"
// @Param finalidade query string false "Property purpose (RESIDENCIAL, COMERCIAL, MISTO)"
// @Param status query string false "Property status (PUBLICADO, EM_EDICAO, ARQUIVADO)"
// @Param published query bool false "Published status"
// @Param min_preco query number false "Minimum price"
// @Param max_preco query number false "Maximum price"
// @Param min_metragem query number false "Minimum square meters"
// @Param max_metragem query number false "Maximum square meters"
// @Param rua query string false "Street name (partial match)"
// @Param cidade query string false "City name (partial match)"
// @Param bairro query string false "Neighborhood name (partial match)"
// @Param num_quartos query int false "Minimum number of bedrooms"
// @Param num_banheiros query int false "Minimum number of bathrooms"
// @Param num_garagens query int false "Minimum number of parking spaces"
// @Param empreendimento_id query uint false "Development ID"
// @Param sort query string false "Sort field (created_at, updated_at, preco, titulo, metragem)" default(created_at)
// @Param order query string false "Sort order (asc, desc)" default(desc)
// @Success 200 {object} errors.Response{success=bool,data=ImovelListResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/imoveis [get]
func (h *Handler) ListImoveis(c *gin.Context) {
	var query ImovelListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	result, err := h.service.ListImoveis(c.Request.Context(), &query)
	if err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(result))
}

// @Summary Add attachment to property
// @Description Add an image or document attachment to a property
// @Tags imoveis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "Property ID"
// @Param request body Anexo true "Attachment data"
// @Success 201 {object} map[string]interface{}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/imoveis/{id}/anexos [post]
func (h *Handler) AddAnexo(c *gin.Context) {
	var uriReq struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	var anexo Anexo
	if err := c.ShouldBindJSON(&anexo); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	if err := h.service.AddAnexo(c.Request.Context(), uriReq.ID, &anexo); err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Attachment added"})
}

// @Summary Get property attachments
// @Description Get all attachments for a property
// @Tags imoveis
// @Accept json
// @Produce json
// @Param id path uint true "Property ID"
// @Success 200 {object} errors.Response{success=bool,data=[]AnexoResponse}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/imoveis/{id}/anexos [get]
func (h *Handler) GetAnexos(c *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	anexos, err := h.service.GetAnexos(c.Request.Context(), req.ID)
	if err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(anexos))
}

// @Summary Add characteristics to property
// @Description Add multiple characteristics to a property
// @Tags imoveis
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "Property ID"
// @Param request body map[string][]uint true "Characteristics IDs"
// @Success 201 {object} map[string]interface{}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/imoveis/{id}/caracteristicas [post]
func (h *Handler) AddCaracteristicas(c *gin.Context) {
	var uriReq struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	var req struct {
		Caracteristicas []uint `json:"caracteristicas" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	if err := h.service.AddCaracteristicas(c.Request.Context(), uriReq.ID, req.Caracteristicas); err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Characteristics added"})
}

// @Summary Get property characteristics
// @Description Get all characteristics for a property
// @Tags imoveis
// @Accept json
// @Produce json
// @Param id path uint true "Property ID"
// @Success 200 {object} errors.Response{success=bool,data=[]CaracteristicaResponse}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/imoveis/{id}/caracteristicas [get]
func (h *Handler) GetCaracteristicas(c *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	caracteristicas, err := h.service.GetCaracteristicas(c.Request.Context(), req.ID)
	if err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(caracteristicas))
}
