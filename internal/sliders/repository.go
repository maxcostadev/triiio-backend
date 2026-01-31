package sliders

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type txKey struct{}

// Repository defines slider repository interface
type Repository interface {
	Create(ctx context.Context, slider *Slider) error
	FindByID(ctx context.Context, id uint) (*Slider, error)
	FindByLocation(ctx context.Context, location string) (*Slider, error)
	Update(ctx context.Context, slider *Slider) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, perPage int) ([]Slider, int64, error)
	CreateItem(ctx context.Context, item *SliderItem) error
	FindItemByID(ctx context.Context, id uint) (*SliderItem, error)
	UpdateItem(ctx context.Context, item *SliderItem) error
	DeleteItem(ctx context.Context, id uint) error
	GetSliderItems(ctx context.Context, sliderID uint) ([]SliderItem, error)
	Transaction(ctx context.Context, fn func(context.Context) error) error
}

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new slider repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// getDB returns the DB from context if in transaction, otherwise returns the repository's DB
func (r *repository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return r.db
}

// Create creates a new slider in the database
func (r *repository) Create(ctx context.Context, slider *Slider) error {
	result := r.getDB(ctx).WithContext(ctx).Create(slider)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByID finds a slider by ID
func (r *repository) FindByID(ctx context.Context, id uint) (*Slider, error) {
	var slider Slider
	result := r.getDB(ctx).WithContext(ctx).Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC")
	}).First(&slider, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &slider, nil
}

// FindByLocation finds a slider by location
func (r *repository) FindByLocation(ctx context.Context, location string) (*Slider, error) {
	var slider Slider
	result := r.getDB(ctx).WithContext(ctx).Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC")
	}).Where("location = ?", location).First(&slider)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &slider, nil
}

// Update updates a slider in the database
func (r *repository) Update(ctx context.Context, slider *Slider) error {
	result := r.getDB(ctx).WithContext(ctx).Model(slider).Select("name", "type", "location", "updated_at").Save(slider)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete soft deletes a slider from the database
func (r *repository) Delete(ctx context.Context, id uint) error {
	result := r.getDB(ctx).WithContext(ctx).Delete(&Slider{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// List retrieves paginated list of sliders
func (r *repository) List(ctx context.Context, page, perPage int) ([]Slider, int64, error) {
	var sliders []Slider
	var total int64

	query := r.getDB(ctx).WithContext(ctx).Model(&Slider{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	if err := query.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC")
	}).Offset(offset).Limit(perPage).Find(&sliders).Error; err != nil {
		return nil, 0, err
	}

	return sliders, total, nil
}

// CreateItem creates a new slider item
func (r *repository) CreateItem(ctx context.Context, item *SliderItem) error {
	result := r.getDB(ctx).WithContext(ctx).Create(item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindItemByID finds a slider item by ID
func (r *repository) FindItemByID(ctx context.Context, id uint) (*SliderItem, error) {
	var item SliderItem
	result := r.getDB(ctx).WithContext(ctx).First(&item, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &item, nil
}

// UpdateItem updates a slider item
func (r *repository) UpdateItem(ctx context.Context, item *SliderItem) error {
	result := r.getDB(ctx).WithContext(ctx).Model(item).Select("image_url", "link_url", "content", "order", "tags", "titulo", "updated_at").Save(item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteItem soft deletes a slider item
func (r *repository) DeleteItem(ctx context.Context, id uint) error {
	result := r.getDB(ctx).WithContext(ctx).Delete(&SliderItem{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetSliderItems retrieves all items for a slider
func (r *repository) GetSliderItems(ctx context.Context, sliderID uint) ([]SliderItem, error) {
	var items []SliderItem
	result := r.getDB(ctx).WithContext(ctx).Where("slider_id = ?", sliderID).Order("\"order\" ASC").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

// Transaction executes a function within a database transaction
func (r *repository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Inject transaction into context
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}
