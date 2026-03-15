# HabitFlow — Integration Checklist

Dokumen ini merangkum checklist integrasi end-to-end untuk memastikan semua komponen HabitFlow bekerja dengan benar secara bersama-sama.

---

## 1. Alur Utama (Critical Path)

### 1.1 Register → Login → Token

| # | Langkah | Expected | Status |
|---|---------|----------|--------|
| 1 | POST `/api/v1/auth/register` dengan name, email, password | 201, mendapat access_token + Set-Cookie refresh token (HttpOnly) | ☐ |
| 2 | POST `/api/v1/auth/login` dengan email + password yang sama | 200, mendapat token baru | ☐ |
| 3 | Gunakan access_token di header `Authorization: Bearer <token>` | Akses protected endpoint berhasil | ☐ |
| 4 | POST `/api/v1/auth/refresh` setelah token expired | 200, mendapat access_token baru, refresh dirotasi via cookie | ☐ |
| 5 | Akses protected endpoint tanpa token | 401 Unauthorized | ☐ |

### 1.2 Create Habit → Checkin → Streak → Report

| # | Langkah | Expected | Status |
|---|---------|----------|--------|
| 6 | POST `/api/v1/habits` — buat habit "Minum Air" | 201, habit tercreate + streak record initialized | ☐ |
| 7 | GET `/api/v1/habits` | 200, daftar habits include "Minum Air" | ☐ |
| 8 | POST `/api/v1/habits/:id/check` — checkin hari ini | 200, is_done=true, current_streak=1 | ☐ |
| 9 | POST `/api/v1/habits/:id/check` — checkin lagi (duplicate) | 409 Conflict, ErrAlreadyCheckedIn | ☐ |
| 10 | GET `/api/v1/habits/today` — status hari ini | 200, habit menunjukkan is_done_today=true | ☐ |
| 11 | DELETE `/api/v1/habits/:id/check` — undo checkin | 200, streak kembali ke 0 | ☐ |
| 12 | GET `/api/v1/reports/weekly` — laporan mingguan | 200, berisi score, insights, streaks | ☐ |
| 13 | GET `/api/v1/reports/score?days=7` — skor konsistensi | 200, overall_score sesuai komputasi | ☐ |
| 14 | GET `/api/v1/reports/insights` — pattern insights | 200, minimal 1 insight (encouragement jika baru) | ☐ |

### 1.3 Habit CRUD Lengkap

| # | Langkah | Expected | Status |
|---|---------|----------|--------|
| 15 | PUT `/api/v1/habits/:id` — update nama + notify_time | 200, data terupdate | ☐ |
| 16 | DELETE `/api/v1/habits/:id` — soft delete | 200, is_active=false | ☐ |
| 17 | GET `/api/v1/habits` — setelah soft delete | Habit yang dihapus tidak muncul | ☐ |

---

## 2. PWA & Frontend

| # | Item | Expected | Status |
|---|------|----------|--------|
| 18 | Buka `https://domain/` di Chrome | Halaman login tampil, dark mode | ☐ |
| 19 | Register akun baru lewat UI | Form validation bekerja, redirect ke dashboard | ☐ |
| 20 | Login lewat UI | access token tersimpan client-side, refresh token hanya di HttpOnly cookie, dashboard muncul | ☐ |
| 21 | Install PWA (Add to Home Screen) | Muncul opsi install, icon + name benar | ☐ |
| 22 | Buka PWA dari home screen | Standalone mode, splash screen, navigasi normal | ☐ |
| 23 | `manifest.json` valid | name, short_name, icons, start_url, display: standalone | ☐ |
| 24 | Service Worker registered | `sw.js` ter-register, caching statis bekerja | ☐ |
| 25 | Offline mode — buka halaman tanpa internet | Halaman cached tampil (minimal shell), data API graceful error | ☐ |

---

## 3. Push Notification

| # | Item | Expected | Status |
|---|------|----------|--------|
| 26 | GET `/api/v1/push/vapid-key` | 200, VAPID public key tersedia | ☐ |
| 27 | POST `/api/v1/push/subscribe` dengan subscription object | 201, subscription tersimpan | ☐ |
| 28 | Scheduler mengirim reminder sesuai notify_time | Push notification muncul di device | ☐ |
| 29 | POST `/api/v1/push/unsubscribe` | 200, subscription dihapus, notif berhenti | ☐ |

---

## 4. Security & Middleware

| # | Item | Expected | Status |
|---|------|----------|--------|
| 30 | CORS — request dari origin yang diizinkan | Header Access-Control-Allow-Origin benar | ☐ |
| 31 | CORS — preflight OPTIONS request | 204 untuk origin valid, 403 untuk origin tidak valid | ☐ |
| 32 | Rate limiting — burst request berlebihan | 429 Too Many Requests setelah limit tercapai | ☐ |
| 33 | JWT expired — akses dengan token expired | 401, message jelas tentang token expired | ☐ |
| 34 | JWT invalid — token malformed/signature salah | 401, ditolak | ☐ |
| 35 | User isolation — akses habit user lain | 404 (bukan 403, untuk privacy) | ☐ |
| 36 | Password hash — cek DB, password tidak plaintext | Tersimpan sebagai bcrypt hash | ☐ |
| 37 | SQL injection — input malicious di field name/email | Input di-escape/sanitize, tidak crash | ☐ |

---

## 5. Data Integrity

| # | Item | Expected | Status |
|---|------|----------|--------|
| 38 | Unique constraint — checkin 2x hari sama via DB | Hanya 1 record (habit_id + date unique) | ☐ |
| 39 | Streak calculation — checkin 3 hari berturut-turut | current_streak=3, longest_streak=3 | ☐ |
| 40 | Streak reset — skip 1 hari, lalu checkin | current_streak=1, longest_streak tetap | ☐ |
| 41 | Longest streak — tidak pernah turun | Setelah reset, longest_streak masih nilai tertinggi | ☐ |
| 42 | Score calculation — 7/7 hari done | Score = 100.0 | ☐ |
| 43 | Timezone WIB — checkin pk 23:59 WIB dan 00:01 WIB | Tanggal sesuai WIB, bukan UTC | ☐ |

---

## 6. Deployment Readiness

| # | Item | Expected | Status |
|---|------|----------|--------|
| 44 | `go build ./cmd/server` berhasil tanpa error | Binary executable terbuild | ☐ |
| 45 | `go test ./internal/services/... -cover` ≥ 80% | Coverage target tercapai untuk critical services | ☐ |
| 46 | Environment variables tersedia: `APP_ENV`, `JWT_SECRET`, `DATABASE_URL`, `CORS_ALLOWED_ORIGINS`, `MAX_BODY_BYTES`, `VAPID_*` | Server start tanpa panic dan hardening aktif | ☐ |
| 47 | Database migration otomatis (AutoMigrate) | Tabel User, RefreshToken, Habit, HabitLog, Streak, PushSubscription dibuat | ☐ |
| 48 | Static file serving — `/` serve `web/index.html` | SPA routing bekerja, fallback ke index.html | ☐ |
| 49 | Graceful error handling — endpoint salah | JSON error response, bukan HTML/stack trace | ☐ |
| 50 | Scheduler aktif — cron job reminder berjalan | Log menunjukkan scheduler start, job terdaftar | ☐ |

---

## Ringkasan

| Kategori | Total Item |
|----------|-----------|
| Alur Utama (Critical Path) | 17 |
| PWA & Frontend | 8 |
| Push Notification | 4 |
| Security & Middleware | 8 |
| Data Integrity | 6 |
| Deployment Readiness | 7 |
| **Total** | **50** |

> ✅ Tandai setiap item dengan ☑ setelah diverifikasi. Semua item harus ☑ sebelum deployment ke production.
