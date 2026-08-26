package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// 200 OK — untuk respons data biasa
func ok(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// 200 OK — khusus respons daftar data dengan paginasi
func okList(c *fiber.Ctx, message string, data any, meta *Meta) error {
	return c.Status(fiber.StatusOK).JSON(WebResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// 201 Created — tambah data berhasil, wajib menyertakan Header Location
func created(c *fiber.Ctx, message string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// 204 No Content — operasi berhasil tapi tanpa body (umumnya untuk DELETE)
func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Gagal umum (400, 404, 409, 415, 500, dll.)
func fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(WebResponse{
		Success: false,
		Message: message,
	})
}

// 422 Unprocessable Entity — gagal validasi input form
func failValidation(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(WebResponse{
		Success: false,
		Message: "validasi gagal",
		Errors:  errs,
	})
}

// Whitelist kolom yang boleh di-sort (mencegah arbitrary input)
var allowedSort = map[string]bool{
	"id":         true,
	"nim":        true,
	"name":       true,
	"grade":      true,
	"created_at": true,
}

// Membaca dan membatasi query string dari URL
func parseListQuery(c *fiber.Ctx) ListQuery {
	q := ListQuery{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 10),
		Search: strings.TrimSpace(c.Query("search")),
		Sort:   c.Query("sort", "id"),
		Order:  strings.ToLower(c.Query("order", "asc")),
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	// Batas atas limit untuk mencegah memori jebol
	if q.Limit > 50 {
		q.Limit = 50
	}
	// Jika klien memasukkan kolom sort yang tidak terdaftar, kembalikan ke default 'id'
	if !allowedSort[q.Sort] {
		q.Sort = "id"
	}
	if q.Order != "desc" {
		q.Order = "asc"
	}

	if raw := c.Query("is_active"); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			q.IsActive = &v
		}
	}

	return q
}