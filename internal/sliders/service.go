package sliders

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var (
	// ErrSliderNotFound is returned when slider is not found
	ErrSliderNotFound = errors.New("slider not found")
	// ErrSliderItemNotFound is returned when slider item is not found
	ErrSliderItemNotFound = errors.New("slider item not found")
	// ErrLocationExists is returned when location already exists
	ErrLocationExists = errors.New("location already exists")
	// ErrInvalidType is returned when slider type is invalid
	ErrInvalidType = errors.New("invalid slider type")
)

// Service defines slider service interface
type Service interface {
	CreateSlider(ctx context.Context, req *CreateSliderRequest) (*SliderResponse, error)
	GetSlider(ctx context.Context, id uint) (*SliderResponse, error)
	GetSliderByLocation(ctx context.Context, location string) (*SliderResponse, error)
	UpdateSlider(ctx context.Context, id uint, req *UpdateSliderRequest) (*SliderResponse, error)
	DeleteSlider(ctx context.Context, id uint) error
	ListSliders(ctx context.Context, page, perPage int) ([]SliderResponse, int64, error)
	AddSliderItem(ctx context.Context, sliderID uint, req *CreateSliderItemRequest) (*SliderItemResponse, error)
	GetSliderItem(ctx context.Context, itemID uint) (*SliderItemResponse, error)
	UpdateSliderItem(ctx context.Context, itemID uint, req *UpdateSliderItemRequest) (*SliderItemResponse, error)
	DeleteSliderItem(ctx context.Context, itemID uint) error
	GetSliderItems(ctx context.Context, sliderID uint) ([]SliderItemResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new slider service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateSlider creates a new slider
func (s *service) CreateSlider(ctx context.Context, req *CreateSliderRequest) (*SliderResponse, error) {
	if req.Type < 0 || req.Type > 2 {
		return nil, ErrInvalidType
	}

	existingSlider, err := s.repo.FindByLocation(ctx, req.Location)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing location: %w", err)
	}
	if existingSlider != nil {
		return nil, ErrLocationExists
	}

	slider := &Slider{
		Name:     req.Name,
		Type:     SliderType(req.Type),
		Location: req.Location,
	}

	err = s.repo.Transaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, slider); err != nil {
			return fmt.Errorf("failed to create slider: %w", err)
		}

		for _, itemReq := range req.Items {
			item := &SliderItem{
				SliderID: slider.ID,
				ImageURL: itemReq.ImageURL,
				LinkURL:  itemReq.LinkURL,
				Content:  itemReq.Content,
				Order:    itemReq.Order,
				Tags:     itemReq.Tags,
				Titulo:   itemReq.Titulo,
			}
			if err := s.repo.CreateItem(txCtx, item); err != nil {
				return fmt.Errorf("failed to create slider item: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	slider, err = s.repo.FindByID(ctx, slider.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload slider: %w", err)
	}
	if slider == nil {
		return nil, fmt.Errorf("failed to reload slider: slider not found after creation")
	}

	return s.sliderToResponse(slider), nil
}

// GetSlider retrieves a slider by ID
func (s *service) GetSlider(ctx context.Context, id uint) (*SliderResponse, error) {
	slider, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find slider: %w", err)
	}
	if slider == nil {
		return nil, ErrSliderNotFound
	}
	return s.sliderToResponse(slider), nil
}

// GetSliderByLocation retrieves a slider by location
func (s *service) GetSliderByLocation(ctx context.Context, location string) (*SliderResponse, error) {
	slider, err := s.repo.FindByLocation(ctx, location)
	if err != nil {
		return nil, fmt.Errorf("failed to find slider: %w", err)
	}
	if slider == nil {
		return nil, ErrSliderNotFound
	}
	return s.sliderToResponse(slider), nil
}

// UpdateSlider updates a slider
func (s *service) UpdateSlider(ctx context.Context, id uint, req *UpdateSliderRequest) (*SliderResponse, error) {
	slider, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find slider: %w", err)
	}
	if slider == nil {
		return nil, ErrSliderNotFound
	}

	if req.Name != "" {
		slider.Name = req.Name
	}
	if req.Type != nil {
		if *req.Type < 0 || *req.Type > 2 {
			return nil, ErrInvalidType
		}
		slider.Type = SliderType(*req.Type)
	}
	if req.Location != "" {
		if req.Location != slider.Location {
			existingSlider, err := s.repo.FindByLocation(ctx, req.Location)
			if err != nil {
				return nil, fmt.Errorf("failed to check existing location: %w", err)
			}
			if existingSlider != nil {
				return nil, ErrLocationExists
			}
		}
		slider.Location = req.Location
	}

	if err := s.repo.Update(ctx, slider); err != nil {
		return nil, fmt.Errorf("failed to update slider: %w", err)
	}

	slider, err = s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to reload slider: %w", err)
	}

	return s.sliderToResponse(slider), nil
}

// DeleteSlider deletes a slider
func (s *service) DeleteSlider(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSliderNotFound
		}
		return fmt.Errorf("failed to delete slider: %w", err)
	}
	return nil
}

// ListSliders retrieves paginated list of sliders
func (s *service) ListSliders(ctx context.Context, page, perPage int) ([]SliderResponse, int64, error) {
	if page < 1 {
		return nil, 0, fmt.Errorf("page must be >= 1")
	}
	if perPage < 1 {
		return nil, 0, fmt.Errorf("perPage must be >= 1")
	}
	if perPage > 100 {
		return nil, 0, fmt.Errorf("perPage must be <= 100")
	}

	sliders, total, err := s.repo.List(ctx, page, perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list sliders: %w", err)
	}

	responses := make([]SliderResponse, len(sliders))
	for i, slider := range sliders {
		responses[i] = *s.sliderToResponse(&slider)
	}

	return responses, total, nil
}

// AddSliderItem adds a new item to a slider
func (s *service) AddSliderItem(ctx context.Context, sliderID uint, req *CreateSliderItemRequest) (*SliderItemResponse, error) {
	slider, err := s.repo.FindByID(ctx, sliderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find slider: %w", err)
	}
	if slider == nil {
		return nil, ErrSliderNotFound
	}

	item := &SliderItem{
		SliderID: sliderID,
		ImageURL: req.ImageURL,
		LinkURL:  req.LinkURL,
		Content:  req.Content,
		Order:    req.Order,
		Tags:     req.Tags,
		Titulo:   req.Titulo,
	}

	if err := s.repo.CreateItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to create slider item: %w", err)
	}

	return s.itemToResponse(item), nil
}

// GetSliderItem retrieves a slider item by ID
func (s *service) GetSliderItem(ctx context.Context, itemID uint) (*SliderItemResponse, error) {
	item, err := s.repo.FindItemByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find slider item: %w", err)
	}
	if item == nil {
		return nil, ErrSliderItemNotFound
	}
	return s.itemToResponse(item), nil
}

// UpdateSliderItem updates a slider item
func (s *service) UpdateSliderItem(ctx context.Context, itemID uint, req *UpdateSliderItemRequest) (*SliderItemResponse, error) {
	item, err := s.repo.FindItemByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to find slider item: %w", err)
	}
	if item == nil {
		return nil, ErrSliderItemNotFound
	}

	if req.ImageURL != "" {
		item.ImageURL = req.ImageURL
	}
	if req.LinkURL != "" {
		item.LinkURL = req.LinkURL
	}
	if req.Content != "" {
		item.Content = req.Content
	}
	if req.Order != nil {
		item.Order = *req.Order
	}
	if req.Tags != nil {
		item.Tags = req.Tags
	}
	if req.Titulo != "" {
		item.Titulo = req.Titulo
	}

	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update slider item: %w", err)
	}

	return s.itemToResponse(item), nil
}

// DeleteSliderItem deletes a slider item
func (s *service) DeleteSliderItem(ctx context.Context, itemID uint) error {
	if err := s.repo.DeleteItem(ctx, itemID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSliderItemNotFound
		}
		return fmt.Errorf("failed to delete slider item: %w", err)
	}
	return nil
}

// GetSliderItems retrieves all items for a slider
func (s *service) GetSliderItems(ctx context.Context, sliderID uint) ([]SliderItemResponse, error) {
	slider, err := s.repo.FindByID(ctx, sliderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find slider: %w", err)
	}
	if slider == nil {
		return nil, ErrSliderNotFound
	}

	items, err := s.repo.GetSliderItems(ctx, sliderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get slider items: %w", err)
	}

	responses := make([]SliderItemResponse, len(items))
	for i, item := range items {
		responses[i] = *s.itemToResponse(&item)
	}

	return responses, nil
}

// Helper methods to convert models to responses

func (s *service) sliderToResponse(slider *Slider) *SliderResponse {
	items := make([]SliderItemResponse, len(slider.Items))
	for i, item := range slider.Items {
		items[i] = *s.itemToResponse(&item)
	}

	return &SliderResponse{
		ID:        slider.ID,
		Name:      slider.Name,
		Type:      int(slider.Type),
		Location:  slider.Location,
		Items:     items,
		CreatedAt: slider.CreatedAt,
		UpdatedAt: slider.UpdatedAt,
	}
}

func (s *service) itemToResponse(item *SliderItem) *SliderItemResponse {
	return &SliderItemResponse{
		ID:        item.ID,
		SliderID:  item.SliderID,
		ImageURL:  item.ImageURL,
		LinkURL:   item.LinkURL,
		Content:   item.Content,
		Order:     item.Order,
		Tags:      item.Tags,
		Titulo:    item.Titulo,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
