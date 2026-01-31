package sliders

import "time"

type Slider struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	Name      string       `gorm:"not null" json:"name"`
	Type      SliderType   `gorm:"not null" json:"type"`
	Location  string       `gorm:"not null" json:"location"`
	Items     []SliderItem `gorm:"foreignKey:SliderID" json:"items"`
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

type SliderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SliderID  uint      `gorm:"not null" json:"slider_id"`
	ImageURL  string    `gorm:"not null" json:"image_url"`
	LinkURL   string    `gorm:"not null" json:"link_url"`
	Content   string    `gorm:"not null" json:"content"`
	Order     int       `gorm:"not null" json:"order"`
	Tags      []string  `gorm:"type:jsonb" json:"tags"`
	Titulo    string    `gorm:"not null" json:"titulo"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type SliderType int

const (
	SliderType_Slideshow SliderType = iota
	SliderType_Carousel
	SliderType_Static
)

func (st SliderType) String() string {
	switch st {
	case SliderType_Slideshow:
		return "slideshow"
	case SliderType_Carousel:
		return "carousel"
	case SliderType_Static:
		return "static"
	default:
		return "unknown"
	}
}

func (Slider) TableName() string {
	return "sliders"
}
