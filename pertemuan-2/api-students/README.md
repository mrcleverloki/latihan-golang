# Student REST API (Pertemuan 2)

Praktikum Pemrograman Backend Lanjut — REST API & HTTP Deep Dive.

## Kontrak API

| Metode | Endpoint | Parameter / Query | Contoh Body Request | Status Response | Keterangan Respons |
|---|---|---|---|---|---|
| `GET` | `/api/v1/students` | `page`, `limit`, `search`, `sort`, `order`, `is_active` | - | `200 OK` | Mengambil daftar mahasiswa dengan paginasi & meta |
| `GET` | `/api/v1/students/:id` | `id` (path) | - | `200 OK`, `400 Bad Request`, `404 Not Found` | Mengambil detail satu data mahasiswa |
| `POST` | `/api/v1/students` | - | `{"nim":"18221001","name":"Ahmad","grade":3.85,"is_active":true}` | `201 Created`, `409 Conflict`, `415 Unsupported Media Type`, `422 Unprocessable Entity` | Menambah mahasiswa baru (disertai Header `Location`) |
| `PUT` | `/api/v1/students/:id` | `id` (path) | `{"nim":"18221001","name":"Ahmad Update","grade":3.90,"is_active":true}` | `200 OK`, `400 Bad Request`, `404 Not Found`, `409 Conflict`, `422 Unprocessable Entity` | Mengganti seluruh data mahasiswa |
| `PATCH`| `/api/v1/students/:id` | `id` (path) | `{"grade": 3.95}` | `200 OK`, `400 Bad Request`, `404 Not Found`, `409 Conflict`, `422 Unprocessable Entity` | Memperbarui sebagian field mahasiswa |
| `DELETE`| `/api/v1/students/:id`| `id` (path) | - | `204 No Content`, `400 Bad Request`, `404 Not Found` | Menghapus data mahasiswa |