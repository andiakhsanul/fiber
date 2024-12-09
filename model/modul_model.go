package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Modul struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	NmModul   string             `json:"nm_modul" bson:"nm_modul" validate:"required"`
	KetModul  string             `json:"ket_modul" bson:"ket_modul"`
	IsAktif   bool               `json:"is_aktif" bson:"is_aktif"`
	Alamat    string             `json:"alamat" bson:"alamat"`
	Urutan    int                `json:"urutan" bson:"urutan"`
	Gbr_Icon  string             `json:"gbr_icon" bson:"gbr_icon"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
