<div align="center">

# ✅ HabitFlow

**Bangun kebiasaan tanpa drama.**

Habit tracker minimalis — satu tap check-off, streak tracking, consistency score, push notification.

[![Live App](https://img.shields.io/badge/Live-App-2ec2b3?style=for-the-badge)](https://habit-tracker-production-3749.up.railway.app/)
[![Landing Page](https://img.shields.io/badge/Landing-Page-0b1215?style=for-the-badge)](https://zhafirash04.github.io/habit-tracker)

</div>

---

## Tentang

HabitFlow adalah Progressive Web App untuk tracking kebiasaan harian. Dibuat karena tidak ada habit tracker yang cukup simpel untuk dipakai konsisten — form panjang, terlalu banyak klik, akhirnya tidak dipakai setelah seminggu.

**Filosofi desain:**
- Satu tap untuk check-in, tanpa form panjang
- Feedback yang jelas lewat streak & consistency score
- Notifikasi yang bisa dikustomisasi per habit, bukan generik

## Fitur

| Fitur | Deskripsi |
|-------|-----------|
| ✅ **Check-off Satu Tap** | Tandai habit selesai dengan satu tap. Tanpa friction. |
| 🔥 **Streak Tracking** | Lihat berapa hari konsisten berturut-turut. |
| 📊 **Consistency Score** | Skor berbasis completion rate — transparan, bisa diverifikasi. |
| 🔔 **Push Notification** | Pengingat per habit sesuai jadwal yang kamu tentukan. |
| 📱 **PWA** | Install di HP seperti native app. Offline-capable. |
| 📈 **Weekly Report** | Laporan mingguan otomatis dengan insight pola kebiasaan. |
| 🔐 **Auth** | Register/login dengan JWT + refresh token rotation. |

## Tech Stack

| Layer | Teknologi |
|-------|-----------|
| **Backend** | Go (Gin framework) |
| **Database** | PostgreSQL (Supabase) |
| **Frontend** | Vanilla JS + Tailwind CSS (SPA, client-side routing) |
| **Auth** | JWT access + refresh token, bcrypt password hashing |
| **Push** | Web Push API + VAPID |
| **Deploy** | Railway (backend) + GitHub Pages (landing page) |

## Struktur Project

```
habitflow/
├── cmd/
│   ├── server/           # Entry point aplikasi
│   └── generate-vapid/   # Tool generate VAPID keys
├── internal/
│   ├── config/           # Konfigurasi & validasi env
│   ├── database/         # Koneksi PostgreSQL via GORM
│   ├── handlers/         # HTTP handler (auth, habit, checkin, report, push)
│   ├── middleware/        # Auth, CORS, rate limit, security headers
│   ├── models/           # GORM model (User, Habit, HabitLog, Streak, dll)
│   ├── router/           # Route setup
│   ├── scheduler/        # Cron job (reminder, weekly report)
│   ├── services/         # Business logic
│   └── tests/            # Integration & security tests
├── web/
│   ├── css/              # Tailwind CSS (local build)
│   ├── js/               # SPA app (app.js, api.js, icons.js)
│   ├── icons/            # PWA icons
│   ├── index.html        # Single page entry
│   ├── manifest.json     # PWA manifest
│   └── sw.js             # Service Worker
├── docs/
│   └── index.html        # Landing page (GitHub Pages)
├── Dockerfile            # Multi-stage Docker build
├── railway.toml          # Railway deployment config
└── .env.example          # Template environment variables
```

## Setup Lokal

### Prerequisites

- Go 1.25+
- PostgreSQL 15+ (atau Docker)
- Node.js 18+ (untuk Tailwind CSS build)

### 1. Clone & Install

```bash
git clone https://github.com/zhafirash04/habit-tracker.git
cd habit-tracker
go mod download
npm install
```

### 2. Setup Database

```bash
# Via Docker
docker run --name habitflow-pg \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=habitflow \
  -p 5432:5432 -d postgres:15
```

### 3. Environment Variables

```bash
cp .env.example .env
# Edit .env — minimal:
# DATABASE_URL=postgresql://postgres:postgres@localhost:5432/habitflow
# JWT_SECRET=(min 32 karakter acak)
```

### 4. Generate VAPID Keys

```bash
go run cmd/generate-vapid/main.go
# Copy output ke .env
```

### 5. Build CSS & Run

```bash
npm run build:css
go run cmd/server/main.go
# 🚀 Server running at http://localhost:8080
```

## API Endpoints

### Auth
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `POST` | `/api/v1/auth/register` | Registrasi user baru |
| `POST` | `/api/v1/auth/login` | Login, return JWT |
| `POST` | `/api/v1/auth/refresh` | Refresh access token |
| `POST` | `/api/v1/auth/logout` | Logout (revoke token) |

### Habits *(Protected)*
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/api/v1/habits` | List semua habit user |
| `POST` | `/api/v1/habits` | Buat habit baru |
| `GET` | `/api/v1/habits/:id` | Detail habit |
| `PUT` | `/api/v1/habits/:id` | Update habit |
| `DELETE` | `/api/v1/habits/:id` | Hapus habit |
| `POST` | `/api/v1/habits/:id/check` | Check-in hari ini |
| `DELETE` | `/api/v1/habits/:id/check` | Undo check-in |
| `GET` | `/api/v1/habits/today` | Status check-in hari ini |

### Reports *(Protected)*
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/api/v1/reports/weekly` | Laporan mingguan |
| `GET` | `/api/v1/reports/daily` | Laporan harian |
| `GET` | `/api/v1/reports/score` | Consistency score |
| `GET` | `/api/v1/reports/insights` | Pattern insights |

### Push Notifications *(Protected)*
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/api/v1/push/vapid-key` | Get VAPID public key |
| `POST` | `/api/v1/push/subscribe` | Subscribe push |
| `DELETE` | `/api/v1/push/unsubscribe` | Unsubscribe push |

### Health
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/api/v1/health` | Liveness check |
| `GET` | `/api/v1/health/security` | Security self-check |

## Testing

```bash
go test ./... -v -count=1
```

Test menggunakan SQLite in-memory untuk kecepatan. Coverage mencakup:
- Handler tests (auth, habit, checkin, report, push)
- Service tests (streak, score, insight)
- Security tests (JWT, rate limiting, injection)

## Deployment

### Railway (Backend)
- Push ke `main` → auto-deploy via Dockerfile
- Set environment variables di Railway dashboard
- Health check: `/api/v1/health`

### GitHub Pages (Landing Page)
- Source: branch `main`, folder `/docs`
- URL: https://zhafirash04.github.io/habit-tracker

## License

Open source — built for personal use.

---

<div align="center">
  <sub>Dibuat dengan frustrasi dan terlalu banyak waktu luang.</sub>
</div>
