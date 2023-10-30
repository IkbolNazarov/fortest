package models

import (
	"time"

	"github.com/lib/pq"
)

type ApartmentsResponse struct {
	Apartments []Apartment `json:"appartments"`
}

type Apartment struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	ObjectID uint    `gorm:"foreignKey" json:"object_id"`
	Type     string  `json:"type"`
	Area     float64 `json:"area"`
	Floor    int     `json:"floor"`
	Block    string  `json:"block"`
	Price    float64 `json:"price"`
	Layout   string  `json:"layout"`
}

type ApartmentType struct {
	Type      string `json:"type"`
	MinSquare string `json:"min_square"`
	MinPrice  string `json:"min_price"`
}

func AppartmentsTable() string {
	tableName := "apartments"
	return tableName
}

type Developer struct {
	ID                   uint     `gorm:"primaryKey" json:"id"`
	Name                 string   `gorm:"not null" json:"name"`
	Address              string   `json:"address"`
	Description          string   `json:"description"`
	Logo                 string   `json:"logo"`
	Website              string   `json:"website"`
	Facebook             string   `json:"facebook"`
	Instagram            string   `json:"instagram"`
	Whatsapp             string   `json:"whatsapp"`
	ResidentialComplexes int      `json:"residential_complexes"`
	LocationLon          string   `json:"location_lon"`
	LocationLat          string   `json:"location_lat"`
	Objects              []Object `gorm:"foreignKey:DeveloperID"`
}

type DeveloperWithObjectsResponse struct {
	DeveloperInfo Developer `json:"developer"`
	Objects       []Object  `json:"objects"`
}

type DeveloperObjectsSummary struct {
	CreationDate time.Time `json:"created_at"`
	Address      string    `json:"address"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
}

type DeveloperObjectsResponse struct {
	DeveloperID uint                      `json:"developer_id"`
	Objects     []DeveloperObjectsSummary `json:"objects"`
}

func DevelopersTable() string {
	tableName := "developers"
	return tableName
}

type Object struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Description  string         `json:"description"`
	Address      string         `json:"address"`
	Latitude     string         `json:"latitude"`
	Longitude    string         `json:"longitude"`
	Logo         string         `json:"logo"`
	Phone        string         `json:"phone"`
	MinPrice     float64        `json:"min_price"`
	MaxPrice     float64        `json:"max_price"`
	MinArea      float64        `json:"min_area"`
	MaxArea      float64        `json:"max_area"`
	Status       string         `json:"status"`
	Class        string         `json:"class"`
	Photos       pq.StringArray `gorm:"type:text[]" json:"photos"`
	YoutubeLink  string         `json:"youtube_link"`
	Amenities    pq.Int32Array  `gorm:"type:integer[]" json:"amenities"`
	DeveloperID  uint           `json:"developer_id"`
	Developer    Developer      `gorm:"foreignKey:DeveloperID"`
	CreationDate time.Time      `gorm:"column:created_at" json:"created_at"`
	Apartments   []Apartment    `gorm:"foreignKey:ObjectID"`
}

func ObjectsTable() string {
	tableName := "objects"
	return tableName
}

type ObjectResponse struct {
	Status    string          `json:"status"`
	Name      string          `json:"name"`
	Address   string          `json:"address"`
	Photos    string          `json:"photos"`
	PhoneNumb string          `json:"phone_numb"`
	Types     []ApartmentType `json:"types"`
}

type FinalResponse struct {
	Objects []ObjectResponse `json:"objects"`
}

type ObjectFilter struct {
	Name     string  `json:"name"`
	MinPrice float64 `json:"minPrice"`
	MaxPrice float64 `json:"maxPrice"`
	Status   string  `json:"status"`
}

type Pagination struct {
	Limit      int         `json:"limit,omitempty;query:limit"`
	Page       int         `json:"page,omitempty;query:page"`
	TotalPages int64       `json:"total_pages"`
	Records    interface{} `json:"records"`
}
