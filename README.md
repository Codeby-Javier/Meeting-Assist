# 🎙️ MeetingAssist

> **Live Demo:** [https://meeting-assist-production.up.railway.app/](https://meeting-assist-production.up.railway.app/)

**MeetingAssist** adalah aplikasi web cerdas yang saya kembangkan untuk membantu efisiensi rapat. Aplikasi ini dapat mengubah rekaman suara menjadi teks (transkripsi) dan mengekstrak teks dari dokumen foto (OCR) secara otomatis dan akurat.

---

## ✨ Fitur Unggulan

1.  **🎙️ Audio ke Teks (Transkripsi Cerdas)**
    - Mendukung unggah file MP3, WAV, M4A.
    - Menggunakan engine AI canggih (AssemblyAI) untuk akurasi tinggi dalam Bahasa Indonesia.
    - Deteksi pembicara otomatis.

2.  **📷 Scan Dokumen (OCR)**
    - Ambil foto dokumen langsung dari HP atau unggah file.
    - Ekstrak teks dari gambar/PDF dalam hitungan detik.
    - Hasil bisa langsung diedit dan disalin.

3.  **📝 Editor & Ekspor**
    - Edit hasil transkripsi langsung di aplikasi.
    - Download hasil dalam format PDF atau Word (DOCX) untuk laporan rapat.

4.  **📱 Desain Responsif**
    - Tampilan modern dan ringan.
    - Dapat diakses dengan nyaman melalui HP, Tablet, maupun Laptop.

---

## 🛠️ Teknologi yang Digunakan

Project ini dibangun dengan **Golang** untuk performa backend yang maksimal dan stabil.

- **Backend:** Go (Golang) v1.21, Gin Framework
- **Database:** PostgreSQL (via Supabase)
- **Frontend:** HTML5, Bootstrap 5, Vanilla JS
- **Cloud Services:** AssemblyAI (Audio), OCR.space (Dokumen)
- **Infrastructure:** Docker & Railway

---

## 🚀 Cara Instalasi

Jika ingin menjalankan project ini di komputer lokal Anda:

1.  **Clone Repository**
    ```bash
    git clone https://github.com/ibnujz/MeetingAssist.git
    cd MeetingAssist
    ```

2.  **Setup Environment**
    Buat file `.env` dan sesuaikan konfigurasi database dan API Key.
    ```env
    PORT=8082
    DB_HOST=...
    ASSEMBLYAI_API_KEY=...
    ```

3.  **Jalankan Aplikasi**
    ```bash
    go mod tidy
    go run main.go
    ```
    Buka `http://localhost:8082` di browser.

---

## 👨‍💻 Author

Dikembangkan oleh **Ibnu JZ**.

[![Instagram](https://img.shields.io/badge/Instagram-%40ibnu.jz-E4405F?style=for-the-badge&logo=instagram&logoColor=white)](https://instagram.com/ibnu.jz)

---

**© 2025 MeetingAssist. All Rights Reserved.**
