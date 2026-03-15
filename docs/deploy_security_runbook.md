# HabitFlow Deployment Security Runbook

Dokumen ini dipakai untuk final check sebelum dan sesudah deploy backend HabitFlow.

## 1) Pre-Deploy Config

Pastikan environment variable berikut tersedia di environment deployment:

- APP_ENV=production
- PORT (contoh: 8080)
- DATABASE_URL (PostgreSQL/SQLite sesuai target)
- JWT_SECRET (secret kuat, minimal 32 karakter, kombinasi huruf besar-kecil dan angka)
- CORS_ALLOWED_ORIGINS (comma-separated, contoh: https://your-app.github.io)
- MAX_BODY_BYTES (contoh: 1048576)
- VAPID_PUBLIC_KEY
- VAPID_PRIVATE_KEY
- VAPID_SUBJECT

Referensi template: .env.example

## 2) Build & Test Gate

Jalankan sebelum deploy:

```powershell
go test ./... -count=1
go build ./cmd/server
```

Harus lolos 100%.

## 3) Startup Validation

Server wajib gagal start bila APP_ENV=production dan JWT_SECRET lemah/default.

Check log startup:

- Valid config: server start normal
- Invalid config: exit dengan pesan "Invalid configuration"

## 4) Post-Deploy Health Checks

Ganti BASE_URL dengan URL backend production.

```powershell
$base='https://your-backend.example.com'
Invoke-RestMethod -Method Get -Uri "$base/api/v1/health"
Invoke-RestMethod -Method Get -Uri "$base/api/v1/health/security"
```

Expected:

- `/health` -> `success=true`, `data.status=ok`
- `/health/security` -> `success=true`, ada `checks.*` dan `warnings` kosong/terkontrol

## 5) Security Header Checks

```powershell
$base='https://your-backend.example.com'
$r = Invoke-WebRequest -Method Get -Uri "$base/api/v1/health"
$r.Headers['X-Content-Type-Options']
$r.Headers['X-Frame-Options']
$r.Headers['Referrer-Policy']
$r.Headers['Permissions-Policy']
$r.Headers['Content-Security-Policy']
```

Expected minimal:

- X-Content-Type-Options: nosniff
- X-Frame-Options: DENY
- Referrer-Policy: strict-origin-when-cross-origin
- Permissions-Policy: tersedia
- Content-Security-Policy: tersedia

## 6) CORS Verification

Origin valid harus lolos preflight, origin asing harus ditolak.

```powershell
$base='https://your-backend.example.com'
$allowed='https://your-app.github.io'
$blocked='https://evil.example'

# Allowed origin
Invoke-WebRequest -Method Options -Uri "$base/api/v1/auth/login" -Headers @{
  Origin=$allowed
  'Access-Control-Request-Method'='POST'
  'Access-Control-Request-Headers'='Content-Type,Authorization'
}

# Blocked origin
Invoke-WebRequest -Method Options -Uri "$base/api/v1/auth/login" -Headers @{
  Origin=$blocked
  'Access-Control-Request-Method'='POST'
  'Access-Control-Request-Headers'='Content-Type,Authorization'
}
```

Expected:

- Allowed: status 204
- Blocked: status 403

## 7) Auth Flow Smoke Test (Cookie-based refresh)

1. Register/login -> response JSON hanya berisi access token + user.
2. Pastikan refresh token tidak muncul di body JSON.
3. Pastikan Set-Cookie refresh token HttpOnly ada.
4. Panggil `/api/v1/auth/refresh` -> access token baru berhasil.
5. Panggil `/api/v1/auth/logout` -> refresh token direvoke.
6. Coba refresh lagi -> harus gagal (401).

## 8) Push Security Smoke Test

- `GET /api/v1/push/vapid-key` hanya public key.
- `POST /api/v1/push/subscribe` endpoint non-https harus ditolak (400).
- Lebih dari 5 subscription/user harus ditolak (429).

## 9) Incident Checklist (Jika Ada Gejala Security)

1. Rotasi `JWT_SECRET`.
2. Rotasi VAPID keys bila dicurigai bocor.
3. Re-deploy backend + force logout seluruh user (revoke refresh token table).
4. Audit log akses endpoint auth, push, dan refresh.
5. Jalankan ulang `go test ./... -count=1` dan health/security checks.
