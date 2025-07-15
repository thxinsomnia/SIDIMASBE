# 🧾 SIDIMASBE – Backend untuk SIDIMAS  
*Backend system untuk SIDIMAS (Sistem Informasi Distribusi Makan Siang Gratis)*

Frontend live: 🌐 [sidimas.vercel.app](https://sidimas.vercel.app)

---

## 🎯 Deskripsi  
SIDIMAS membantu kelancaran distribusi makan siang gratis bagi siswa, mendukung

- pencatatan **stok bahan mentah**  
- penyusunan **menu dan katering harian**  
- distribusi secara **real‑time** via suplai backend dari dapur ke pemasok :contentReference[oaicite:3]{index=3}

---

## ⚙️ Fitur Utama (API)

- 🔐 Autentikasi ( admin, pemasok)
- 📦 Kelola stok bahan mentah & menu
- 📨 Permintaan & distribusi bahan antar pihak
- 📊 Laporan statistik distribusi & aktivitas sistem

---

## 🛠️ Stacks & Teknologi

| Komponen    | Teknologi               |
|-------------|-------------------------|
| Backend API | Golang (Gin Framework)  |
| Database    | Supabase (PostgreSQL)   |
| ORM         | GORM                    |
| Auth        | JWT                     |
| Hosting     | Vercel                  |
| Frontend    | Vue.js (di repo terpisah) :contentReference[oaicite:4]{index=4} |

---

## 🗂️ Struktur Proyek

```

SIDIMASBE/
├── config/        # DB, env, Supabase setup
├── controllers/   # API handlers
├── models/        # Struct & model DB
├── routes/        # Endpoint definitions
├── middleware/    # Auth, CORS, logger
└── main.go        # Entry-point

````

---

## 🚀 Instalasi Lokal

1. Clone repo  
   ```bash
   git clone https://github.com/thxinsomnia/SIDIMASBE.git
   cd SIDIMASBE

2. Buat `.env`:

   ```env
   SUPABASE_URL=https://...
   SUPABASE_KEY=your_supabase_key
   JWT_SECRET=secret_key
   ```

3. Pastikan schema database sudah tersedia di Supabase.

4. Jalankan project:

   ```bash
   go run main.go
   ```

   Akses server di `http://localhost:8080`

---

## 📬 Contoh Endpoint

| Method | Endpoint              | Deskripsi                                |
| ------ | --------------------- | ---------------------------------------- |
| POST   | `/login`              | Autentikasi pengguna                     |
| GET    | `/stok`               | Lihat stok bahan mentah                  |
| POST   | `/stok`               | Tambah stok bahan                        |
| POST   | `/request-distribusi` | Buat permintaan distribusi bahan         |
| GET    | `/laporan`            | Ambil ringkasan distribusi dan aktivitas |

---

## 🔗 Integrasi dengan Frontend

Cek **frontend live** di [sidimas.vercel.app](https://sidimas.vercel.app), yang terhubung langsung ke backend ini.

Dokumentasi lengkap bisa ditemukan di repo frontend (jika tersedia).

---

## 📄 Lisensi

Dirilis di bawah **MIT License**. Cek file LICENSE untuk detail.

---

## 👤 Developer

Made with 💙 by @thxinsomnia
📚 Mahasiswa D4 Teknik Informatika, ULBI
📌 Fokus: Web Dev, Backend, dan Sistem Informasi
