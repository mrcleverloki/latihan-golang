package service

import (
	"strings"
)

// NormalizeTags membersihkan slice tag: huruf kecil, trim spasi, hapus duplikat dan string kosong.
// Menjamin return berupa slice kosong []string{} jika tidak ada tag (bukan nil).
func NormalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}

	seen := make(map[string]bool)
	result := make([]string, 0, len(tags))

	for _, t := range tags {
		cleaned := strings.ToLower(strings.TrimSpace(t))
		if cleaned == "" {
			continue
		}
		if !seen[cleaned] {
			seen[cleaned] = true
			result = append(result, cleaned)
		}
	}

	return result
}