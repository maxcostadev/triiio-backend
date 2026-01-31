package sliders

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	apiErrors "github.com/vahiiiid/go-rest-api-boilerplate/internal/errors"
)

type Handler struct {
	service Service
}

// NewHandler creates a new slider handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// @Summary Create slider
// @Description Create a new slider with items
// @Tags sliders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSliderRequest true "Slider creation request"
// @Success 201 {object} errors.Response{success=bool,data=SliderResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 409 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders [post]
func (h *Handler) CreateSlider(c *gin.Context) {
	var req CreateSliderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	slider, err := h.service.CreateSlider(c.Request.Context(), &req)
	if err != nil {
		if err == ErrLocationExists {
			_ = c.Error(apiErrors.Conflict("Location already exists"))
			return
		}
		if err == ErrInvalidType {
			_ = c.Error(apiErrors.BadRequest("Invalid slider type"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, apiErrors.Success(slider))
}

// @Summary Get slider by ID
// @Description Retrieve a slider and its items by ID
// @Tags sliders
// @Accept json
// @Produce json
// @Param id path int true "Slider ID"
// @Success 200 {object} errors.Response{success=bool,data=SliderResponse}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders/{id} [get]
func (h *Handler) GetSlider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apiErrors.BadRequest("Invalid slider ID"))
		return
	}

	slider, err := h.service.GetSlider(c.Request.Context(), uint(id))
	if err != nil {
		if err == ErrSliderNotFound {
			_ = c.Error(apiErrors.NotFound("Slider not found"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(slider))
}

// @Summary Get slider by location
// @Description Retrieve a slider and its items by location
// @Tags sliders
// @Accept json
// @Produce json
// @Param location query string true "Slider location"
// @Success 200 {object} errors.Response{success=bool,data=SliderResponse}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders/location [get]
func (h *Handler) GetSliderByLocation(c *gin.Context) {
	location := c.Query("location")
	if location == "" {
		_ = c.Error(apiErrors.BadRequest("Location parameter is required"))
		return
	}

	slider, err := h.service.GetSliderByLocation(c.Request.Context(), location)
	if err != nil {
		if err == ErrSliderNotFound {
			_ = c.Error(apiErrors.NotFound("Slider not found"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(slider))
}

// @Summary Update slider
// @Description Update an existing slider
// @Tags sliders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Slider ID"
// @Param request body UpdateSliderRequest true "Slider update request"
// @Success 200 {object} errors.Response{success=bool,data=SliderResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 409 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders/{id} [put]
func (h *Handler) UpdateSlider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apiErrors.BadRequest("Invalid slider ID"))
		return
	}

	var req UpdateSliderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	slider, err := h.service.UpdateSlider(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err == ErrSliderNotFound {
			_ = c.Error(apiErrors.NotFound("Slider not found"))
			return
		}
		if err == ErrLocationExists {
			_ = c.Error(apiErrors.Conflict("Location already exists"))
			return
		}
		if err == ErrInvalidType {
			_ = c.Error(apiErrors.BadRequest("Invalid slider type"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(slider))
}

// @Summary Delete slider
// @Description Delete a slider and all its items
// @Tags sliders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Slider ID"
// @Success 204
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders/{id} [delete]
func (h *Handler) DeleteSlider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apiErrors.BadRequest("Invalid slider ID"))
		return
	}

	err = h.service.DeleteSlider(c.Request.Context(), uint(id))
	if err != nil {
		if err == ErrSliderNotFound {
			_ = c.Error(apiErrors.NotFound("Slider not found"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary List sliders
// @Description Retrieve paginated list of sliders
// @Tags sliders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {object} errors.Response{success=bool,data=[]SliderResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders [get]
func (h *Handler) ListSliders(c *gin.Context) {
	page := 1
	perPage := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if pp := c.Query("per_page"); pp != "" {
		if parsed, err := strconv.Atoi(pp); err == nil && parsed > 0 && parsed <= 100 {
			perPage = parsed
		}
	}

	sliders, total, err := h.service.ListSliders(c.Request.Context(), page, perPage)
	if err != nil {
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sliders,
		"pagination": gin.H{
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}

// @Summary Add slider item
// @Description Add a new item to an existing slider
// @Tags sliders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Slider ID"
// @Param request body CreateSliderItemRequest true "Slider item creation request"
// @Success 201 {object} errors.Response{success=bool,data=SliderItemResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders/{slider_id}/items [post]
func (h *Handler) AddSliderItem(c *gin.Context) {
	sliderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apiErrors.BadRequest("Invalid slider ID"))
		return
	}

	var req CreateSliderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	item, err := h.service.AddSliderItem(c.Request.Context(), uint(sliderID), &req)
	if err != nil {
		if err == ErrSliderNotFound {
			_ = c.Error(apiErrors.NotFound("Slider not found"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusCreated, apiErrors.Success(item))
}

// @Summary Get slider item
// @Description Retrieve a specific slider item by ID
// @Tags sliders
// @Accept json
// @Produce json
// @Param item_id path int true "Slider item ID"
// @Success 200 {object} errors.Response{success=bool,data=SliderItemResponse}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders/items/{item_id} [get]
func (h *Handler) GetSliderItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		_ = c.Error(apiErrors.BadRequest("Invalid item ID"))
		return
	}

	item, err := h.service.GetSliderItem(c.Request.Context(), uint(itemID))
	if err != nil {
		if err == ErrSliderItemNotFound {
			_ = c.Error(apiErrors.NotFound("Slider item not found"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(item))
}

// @Summary Update slider item
// @Description Update an existing slider item
// @Tags sliders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item_id path int true "Slider item ID"
// @Param request body UpdateSliderItemRequest true "Slider item update request"
// @Success 200 {object} errors.Response{success=bool,data=SliderItemResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders/items/{item_id} [put]
func (h *Handler) UpdateSliderItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		_ = c.Error(apiErrors.BadRequest("Invalid item ID"))
		return
	}

	var req UpdateSliderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apiErrors.FromGinValidation(err))
		return
	}

	item, err := h.service.UpdateSliderItem(c.Request.Context(), uint(itemID), &req)
	if err != nil {
		if err == ErrSliderItemNotFound {
			_ = c.Error(apiErrors.NotFound("Slider item not found"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(item))
}

// @Summary Delete slider item
// @Description Delete a slider item
// @Tags sliders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item_id path int true "Slider item ID"
// @Success 204
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders/items/{item_id} [delete]
func (h *Handler) DeleteSliderItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		_ = c.Error(apiErrors.BadRequest("Invalid item ID"))
		return
	}

	err = h.service.DeleteSliderItem(c.Request.Context(), uint(itemID))
	if err != nil {
		if err == ErrSliderItemNotFound {
			_ = c.Error(apiErrors.NotFound("Slider item not found"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary Get slider items
// @Description Retrieve all items for a specific slider
// @Tags sliders
// @Accept json
// @Produce json
// @Param id path int true "Slider ID"
// @Success 200 {object} errors.Response{success=bool,data=[]SliderItemResponse}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/sliders/{slider_id}/items [get]
func (h *Handler) GetSliderItems(c *gin.Context) {
	sliderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		_ = c.Error(apiErrors.BadRequest("Invalid slider ID"))
		return
	}

	items, err := h.service.GetSliderItems(c.Request.Context(), uint(sliderID))
	if err != nil {
		if err == ErrSliderNotFound {
			_ = c.Error(apiErrors.NotFound("Slider not found"))
			return
		}
		_ = c.Error(apiErrors.InternalServerError(err))
		return
	}

	c.JSON(http.StatusOK, apiErrors.Success(items))
}
