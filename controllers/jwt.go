package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Fungsi encode Base64 URL
func encodeBase64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

// Fungsi untuk signature menggunakan HMAC-SHA256
func createSignature(message, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return encodeBase64URL(h.Sum(nil))
}

// Fungsi generate JWT
func generateJWT(username string) string {
	// Header JWT
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, _ := json.Marshal(header)
	headerEncoded := encodeBase64URL(headerJSON)

	// Payload JWT dengan waktu kedaluwarsa
	payload := map[string]interface{}{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 jam
	}
	payloadJSON, _ := json.Marshal(payload)
	payloadEncoded := encodeBase64URL(payloadJSON)

	// Signature
	secret := "mySecretKey"
	signature := createSignature(headerEncoded+"."+payloadEncoded, secret)

	// Token JWT
	token := fmt.Sprintf("%s.%s.%s", headerEncoded, payloadEncoded, signature)
	return token
}

// Middleware untuk validasi JWT
func JWTMiddleware(c *fiber.Ctx) error {
	// Mengambil token dari header Authorization
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if token == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "no token provided"})
	}

	// Memisahkan bagian JWT menjadi header, payload, dan signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token format"})
	}

	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	// Memverifikasi signature
	secret := "mySecretKey"
	expectedSignature := createSignature(header+"."+payload, secret)

	if signature != expectedSignature {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	// Decode payload dan cek expired time
	payloadDecoded, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token payload"})
	}

	var payloadData map[string]interface{}
	if err := json.Unmarshal(payloadDecoded, &payloadData); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token payload"})
	}

	// Cek apakah token sudah expired
	if exp, ok := payloadData["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "token expired"})
		}
	}

	// Lanjutkan ke handler berikutnya jika token valid
	return c.Next()
}
