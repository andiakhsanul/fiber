package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User struct represents a user in the MongoDB database
type User struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                      // ID unik dari MongoDB
	Username     string             `json:"username" bson:"username" validate:"required"`           // Nama pengguna
	NmUser       string             `json:"nm_user" bson:"nm_user" validate:"required"`             // Nama lengkap pengguna
	Password     string             `json:"password" bson:"pass" validate:"required"`               // Password yang di-hash
	Email        string             `json:"email" bson:"email" validate:"required,email"`           // Email pengguna
	Role         string             `json:"role" bson:"role" validate:"required"`                   // Peran pengguna, misalnya civitas
	CreatedAt    primitive.DateTime `json:"created_at" bson:"created_at,omitempty"`                 // Tanggal pembuatan akun
	JenisKelamin int                `json:"jenis_kelamin" bson:"jenis_kelamin" validate:"required"` // 1 untuk laki-laki, 2 untuk perempuan
	Photo        string             `json:"photo,omitempty" bson:"photo,omitempty"`                 // Path atau URL gambar profil
	Phone        string             `json:"phone" bson:"phone" validate:"required"`                 // Nomor telepon pengguna
	Token        string             `json:"token,omitempty" bson:"token,omitempty"`                 // Token autentikasi (opsional)
	JenisUser    string             `json:"jenis_user" bson:"jenis_user" validate:"required"`       // Jenis pengguna, misalnya Mahasiswa
	Pass_2       string             `json:"pass_2,omitempty" bson:"pass_2,omitempty"`               // Field tambahan (opsional)
}
