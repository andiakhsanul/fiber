// File: usermodul_model.go
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModul struct {
	ID        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	JenisUser string               `json:"jenis_user" bson:"jenis_user" validate:"required"` // Contoh: "Dosen"
	UserID    []primitive.ObjectID `json:"user_id" bson:"user_id" validate:"required"`       // Array Referensi ke User
	ModulID   []primitive.ObjectID `json:"modul_id" bson:"modul_id" validate:"required"`     // Array Referensi ke Modul
	Catatan   string               `json:"catatan,omitempty" bson:"catatan,omitempty"`       // Catatan opsional
	CreatedAt time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
}