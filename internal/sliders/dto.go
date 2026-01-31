package sliders

import "time"

// CreateSliderRequest represents slider creation request
type CreateSliderRequest struct {
	Name     string                    `json:"name" binding:"required,min=1,max=200"`
	Type     int                       `json:"type" binding:"required,min=0,max=2"`
	Location string                    `json:"location" binding:"required,min=1,max=255"`
	Items    []CreateSliderItemRequest `json:"items" binding:"dive"`
}

// UpdateSliderRequest represents slider update request
type UpdateSliderRequest struct {
	Name     string `json:"name" binding:"omitempty,min=1,max=200"`
	Type     *int   `json:"type" binding:"omitempty,min=0,max=2"`
	Location string `json:"location" binding:"omitempty,min=1,max=255"`
}

// CreateSliderItemRequest represents slider item creation request
type CreateSliderItemRequest struct {
	ImageURL string   `json:"image_url" binding:"required,min=1,max=2048"`
	LinkURL  string   `json:"link_url" binding:"omitempty,max=2048"`
	Content  string   `json:"content" binding:"omitempty,max=1000"`
	Order    int      `json:"order" binding:"required,min=0"`
	Tags     []string `json:"tags" binding:"omitempty,dive,max=100"`
	Titulo   string   `json:"titulo" binding:"omitempty,max=255"`
}

// UpdateSliderItemRequest represents slider item update request
type UpdateSliderItemRequest struct {
	ImageURL string   `json:"image_url" binding:"omitempty,min=1,max=2048"`
	LinkURL  string   `json:"link_url" binding:"omitempty,max=2048"`
	Content  string   `json:"content" binding:"omitempty,max=1000"`
	Order    *int     `json:"order" binding:"omitempty,min=0"`
	Tags     []string `json:"tags" binding:"omitempty,dive,max=100"`
	Titulo   string   `json:"titulo" binding:"omitempty,max=255"`
}

// SliderResponse represents slider response
type SliderResponse struct {
	ID        uint                 `json:"id"`
	Name      string               `json:"name"`
	Type      int                  `json:"type"`
	Location  string               `json:"location"`
	Items     []SliderItemResponse `json:"items"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

// SliderItemResponse represents slider item response
type SliderItemResponse struct {
	ID        uint      `json:"id"`
	SliderID  uint      `json:"slider_id"`
	ImageURL  string    `json:"image_url"`
	LinkURL   string    `json:"link_url"`
	Content   string    `json:"content"`
	Order     int       `json:"order"`
	Tags      []string  `json:"tags"`
	Titulo    string    `json:"titulo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
