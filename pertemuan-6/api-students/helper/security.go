package helper

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// bcryptCost menentukan berapa kali proses hashing diulang secara internal.
// Cost 12 adalah standar yang aman: lambat untuk ditebak, cepat untuk sekali login.
const bcryptCost = 12

// dummyHash dipakai ketika username tidak ditemukan, agar waktu tanggap
// login tetap mirip dengan kasus password salah (mencegah timing attack).
var dummyHash = []byte("$2a$12$abcdefghijklmnopqrstuuLKa3Bt1TCmU/6zvhZ8x4nq1yBiuGvS")

// HashPassword mengubah password menjadi hash bcrypt yang aman.
func HashPassword(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword membandingkan password teks biasa dengan hash-nya.
func VerifyPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// VerifyDummyPassword sengaja membuang waktu komputasi yang sama,
// dipakai ketika user tidak ditemukan.
func VerifyDummyPassword(plain string) {
	_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(plain))
}

// RandomToken menghasilkan string acak aman secara kriptografis untuk refresh token.
func RandomToken(numBytes int) (string, error) {
	buf := make([]byte, numBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// SHA256Hex dipakai untuk menyimpan refresh token dalam bentuk hash ke database.
func SHA256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}