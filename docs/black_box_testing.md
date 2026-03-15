# Black Box Testing — HabitFlow API

Tabel test case untuk pengujian API menggunakan metode **Black Box Testing** sesuai dengan outline skripsi Bab IV.

---

## 1. Auth Endpoints

### POST /api/v1/auth/register

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 1 | Register dengan data valid | `{"name":"John","email":"john@test.com","password":"password123"}` | 201, `success:true`, user data + access_token + refresh_token | ✅ |
| 2 | Register dengan email duplicate | Email yang sudah terdaftar | 409, `success:false`, `"email sudah terdaftar"` | ✅ |
| 3 | Register tanpa name | `{"email":"a@b.com","password":"12345678"}` | 400, `success:false`, `"Validasi gagal"` | ✅ |
| 4 | Register dengan password < 8 karakter | `{"name":"A","email":"a@b.com","password":"123"}` | 400, `success:false`, `"Validasi gagal"` | ✅ |
| 5 | Register dengan email format invalid | `{"name":"A","email":"bukan-email","password":"12345678"}` | 400, `success:false`, `"Validasi gagal"` | ✅ |
| 6 | Register dengan body kosong | `{}` | 400, `success:false`, `"Validasi gagal"` | ✅ |

### POST /api/v1/auth/login

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 7 | Login dengan kredensial valid | Email & password benar | 200, `success:true`, user data + tokens | ✅ |
| 8 | Login dengan password salah | Email benar, password salah | 401, `success:false`, `"email atau password salah"` | ✅ |
| 9 | Login dengan email tidak terdaftar | Email yang belum pernah register | 401, `success:false`, `"email atau password salah"` | ✅ |
| 10 | Login tanpa password | `{"email":"a@b.com"}` | 400, `success:false`, `"Validasi gagal"` | ✅ |

### POST /api/v1/auth/refresh

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 11 | Refresh dengan refresh token valid | `Authorization: Bearer <refresh_token>` | 200, `success:true`, access_token baru | ✅ |
| 12 | Refresh tanpa header Authorization | Tanpa header | 400, `success:false` | ✅ |
| 13 | Refresh dengan token expired/invalid | Token random | 401, `success:false` | ✅ |
| 14 | Refresh dengan access token (bukan refresh) | Access token di header | 401, `success:false`, `"token bukan refresh token"` | ✅ |

---

## 2. Habit Endpoints

### POST /api/v1/habits

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 15 | Buat habit dengan data valid | `{"name":"Minum Air","category":"health"}` | 201, `success:true`, habit data + streak 0 | ✅ |
| 16 | Buat habit dengan notify_time | `{"name":"Baca","category":"learning","notify_time":"07:00"}` | 201, habit data dengan notify_time | ✅ |
| 17 | Buat habit tanpa name | `{"category":"health"}` | 400, `success:false`, `"Validasi gagal"` | ✅ |
| 18 | Buat habit dengan notify_time format salah | `{"name":"X","notify_time":"25:00"}` | 400, `success:false`, format HH:MM | ✅ |
| 19 | Buat habit tanpa auth token | Tanpa Authorization header | 401, unauthorized | ✅ |

### GET /api/v1/habits

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 20 | Ambil semua habit user | Token valid | 200, array of habits | ✅ |
| 21 | Ambil habit ketika belum punya | Token valid, user baru | 200, array kosong `[]` | ✅ |

### GET /api/v1/habits/:id

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 22 | Ambil habit by ID valid | ID milik user | 200, habit data | ✅ |
| 23 | Ambil habit by ID milik user lain | ID milik user lain | 404, `success:false` | ✅ |
| 24 | Ambil habit dengan ID non-numerik | `/habits/abc` | 400, `"ID habit tidak valid"` | ✅ |

### PUT /api/v1/habits/:id

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 25 | Update nama habit | `{"name":"Nama Baru"}` | 200, habit terupdate | ✅ |
| 26 | Update notify_time format salah | `{"notify_time":"99:99"}` | 400, error format | ✅ |
| 27 | Update habit milik user lain | ID habit user lain | 404, `success:false` | ✅ |

### DELETE /api/v1/habits/:id

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 28 | Hapus habit (soft delete) | ID habit valid | 200, `"Habit berhasil dihapus"` | ✅ |
| 29 | Hapus habit yang sudah dihapus | ID habit is_active=false | 404, `success:false` | ✅ |
| 30 | Hapus habit milik user lain | ID habit user lain | 404, `success:false` | ✅ |

---

## 3. Check-in Endpoints

### POST /api/v1/habits/:id/check

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 31 | Check-in habit pertama kali | POST ke habit yang belum pernah checkin | 200, `is_done:true`, `current_streak:1` | ✅ |
| 32 | Check-in berturut-turut (streak lanjut) | POST setelah kemarin juga checkin | 200, streak bertambah 1 | ✅ |
| 33 | Check-in setelah gap (streak reset) | POST setelah 2+ hari tidak checkin | 200, `current_streak:1` | ✅ |
| 34 | Check-in habit yang sudah done hari ini | POST kedua kali di hari yang sama | 409, `"habit sudah dicheckin hari ini"` | ✅ |
| 35 | Check-in habit yang tidak ada | ID habit tidak valid | 404, `"Habit tidak ditemukan"` | ✅ |
| 36 | Check-in dengan catatan | `{"note":"2 liter"}` | 200, note tersimpan | ✅ |

### DELETE /api/v1/habits/:id/check

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 37 | Undo checkin hari ini | DELETE setelah checkin hari ini | 200, `"Checkin hari ini berhasil dibatalkan"` | ✅ |
| 38 | Undo tanpa ada checkin hari ini | DELETE tanpa pernah checkin | 404, `"tidak ada checkin hari ini untuk di-undo"` | ✅ |

### GET /api/v1/habits/today

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 39 | Ambil status hari ini | Token valid | 200, `date`, array `habits` dengan `is_done_today` | ✅ |
| 40 | Ambil status ketika belum ada habit | User baru tanpa habit | 200, `habits: []` | ✅ |

---

## 4. Report Endpoints

### GET /api/v1/reports/weekly

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 41 | Ambil laporan minggu ini | Token valid | 200, `period`, `start_date`, `end_date`, `score`, `streaks`, `insights` | ✅ |
| 42 | Ambil laporan periode custom | `?start=2026-03-01&end=2026-03-07` | 200, laporan untuk periode tersebut | ✅ |
| 43 | Ambil laporan user tanpa habit | User baru | 200, `total_habits:0`, `total_checkins:0` | ✅ |

### GET /api/v1/reports/score

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 44 | Ambil skor konsistensi 7 hari | Token valid (default) | 200, `overall_score`, `habit_scores[]` | ✅ |
| 45 | Ambil skor konsistensi 30 hari | `?days=30` | 200, `period:"last_30_days"` | ✅ |
| 46 | Ambil skor user tanpa habit | User baru | 200, `overall_score:0`, `habit_scores:[]` | ✅ |

### GET /api/v1/reports/insights

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 47 | Ambil insights user aktif | User dengan data checkin | 200, array insights berisi best_day/consistency/declining | ✅ |
| 48 | Ambil insights user baru | User tanpa data checkin | 200, insight `encouragement` | ✅ |

---

## 5. Push Notification Endpoints

### GET /api/v1/push/vapid-key

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 49 | Ambil VAPID public key | Token valid | 200, `vapid_public_key` | ✅ |

### POST /api/v1/push/subscribe

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 50 | Subscribe push notification | `{"endpoint":"https://...","keys":{"p256dh":"...","auth":"..."}}` | 201, subscription data | ✅ |
| 51 | Subscribe tanpa endpoint | `{}` | 400, `"Validasi gagal"` | ✅ |

### DELETE /api/v1/push/unsubscribe

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 52 | Unsubscribe push | `{"endpoint":"https://..."}` | 200, subscription dihapus | ✅ |
| 53 | Unsubscribe tanpa endpoint | `{}` | 400, `"endpoint diperlukan"` | ✅ |

---

## 6. Middleware & Security

| No | Skenario | Input | Expected Output | Status |
|----|----------|-------|-----------------|--------|
| 54 | Request tanpa token ke protected route | GET /api/v1/habits tanpa header | 401, `"Authorization header required"` | ✅ |
| 55 | Request dengan token expired | Token JWT yang sudah expired | 401, unauthorized | ✅ |
| 56 | CORS preflight request | OPTIONS /api/v1/habits | 200, CORS headers present | ✅ |
| 57 | Rate limit login (brute force) | 10+ login gagal berturut-turut | 429, rate limited | ✅ |
| 58 | API response format konsisten | Semua endpoint | Selalu `{"success":bool,"message":string,"data":...}` | ✅ |

---

**Total Test Case: 58**

> Dokumen ini digunakan sebagai lampiran di Bab IV skripsi bagian "Pengujian Sistem" menggunakan metode Black Box Testing.
