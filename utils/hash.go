package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// Password yang akan di-hash
const password = "12345"

func main() {
	// 1. Proses Hashing (Menghasilkan Hash yang Aman)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Fatal("Gagal menghash password:", err)
	}

	fmt.Println("Password Asli:", password)
	// Hasil hash akan berbeda setiap kali kode ini dijalankan karena penggunaan 'salt'
	fmt.Println("Hash yang Dihasilkan:", string(hashedPassword))
	fmt.Println("---")

	// 2. Proses Verifikasi (Saat Pengguna Mencoba Login)
	// Membandingkan password input dengan hash yang tersimpan
	isMatch := CheckPasswordHash(password, string(hashedPassword))

	fmt.Printf("Verifikasi (Password Benar) berhasil: %v\n", isMatch)

	// Uji coba dengan password yang salah
	isWrongMatch := CheckPasswordHash("password_salah", string(hashedPassword))
	fmt.Printf("Verifikasi (Password Salah) berhasil: %v\n", isWrongMatch)
}

// HashPassword mengenkripsi password menggunakan bcrypt
func HashPassword(password string) ([]byte, error) {
	// bcrypt.DefaultCost menggunakan cost (tingkat kelambatan) yang disarankan.
	// Cost yang lebih tinggi lebih aman, tetapi lebih lambat.
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CheckPasswordHash membandingkan password plaintext dengan hash yang tersimpan
func CheckPasswordHash(password string, hash string) bool {
	// Fungsi ini menangani proses pembandingan dan 'unsalting' secara otomatis
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}