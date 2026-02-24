package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Scenario struct {
	ID           string     `gorm:"type:text;primary_key" json:"id"`
	Title        string     `gorm:"type:varchar(255);not null" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	TotalCostHKD int        `gorm:"not null" json:"total_cost_hkd"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CostItems    []CostItem `gorm:"foreignKey:ScenarioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"cost_breakdown"`
	Payouts      []Payout   `gorm:"foreignKey:ScenarioID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"payouts"`
}

func (s *Scenario) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

type CostItem struct {
	ID         string `gorm:"type:text;primary_key" json:"id"`
	ScenarioID string `gorm:"type:text;not null;index" json:"-"`
	ItemName   string `gorm:"type:varchar(255);not null" json:"item_name"`
	AmountHKD  int    `gorm:"not null" json:"amount_hkd"`
}

func (c *CostItem) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

type Insurer struct {
	ID       string `gorm:"type:varchar(50);primary_key" json:"id"`
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	PlanName string `gorm:"type:varchar(255);not null" json:"plan_name"`
}

type Payout struct {
	ID                 string  `gorm:"type:text;primary_key" json:"id"`
	ScenarioID         string  `gorm:"type:text;not null;index" json:"-"`
	InsurerID          string  `gorm:"type:varchar(50);not null" json:"insurer_id"`
	Insurer            Insurer `gorm:"foreignKey:InsurerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	EstimatedPayoutHKD int     `gorm:"not null" json:"estimated_payout_hkd"`
	CoveragePercentage float64 `gorm:"type:real;not null" json:"coverage_percentage"`
	Analysis           string  `gorm:"type:text" json:"analysis"`
	IsRecommended      bool    `gorm:"default:false" json:"is_recommended"`
}

func (p *Payout) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}
