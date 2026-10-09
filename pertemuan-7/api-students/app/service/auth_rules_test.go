package service_test

import (
	"testing"

	"api-students/app/service"
)

func TestCheckPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected string
	}{
		{
			name:     "Password kurang dari 8 karakter",
			password: "pass1",
			expected: "minimal 8 karakter",
		},
		{
			name:     "Password hanya berisi huruf",
			password: "hanyahurufsaja",
			expected: "harus memuat huruf dan angka",
		},
		{
			name:     "Password hanya berisi angka",
			password: "1234567890",
			expected: "harus memuat huruf dan angka",
		},
		{
			name:     "Password umum yang dilarang",
			password: "password123",
			expected: "password terlalu umum",
		},
		{
			name:     "Password valid memenuhi kriteria",
			password: "kombinasiAman99",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.CheckPasswordStrength(tt.password)
			if got != tt.expected {
				t.Errorf("CheckPasswordStrength(%q) = %q, ingin %q", tt.password, got, tt.expected)
			}
		})
	}
}