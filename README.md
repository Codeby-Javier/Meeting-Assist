# 🎙️ MeetingAssist - Smart AI Transcription & OCR

> **Live Demo:** [https://meeting-assist-production.up.railway.app/](https://meeting-assist-production.up.railway.app/)

[![Author](https://img.shields.io/badge/Created_by-Ibnu_JZ-blue?style=for-the-badge&logo=instagram)](https://instagram.com/ibnu.jz)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ibnujz/MeetingAssist?style=for-the-badge)](https://github.com/ibnujz/MeetingAssist)
[![Deployment](https://img.shields.io/badge/Deploy-Railway-success?style=for-the-badge&logo=railway)](https://railway.app)

---

## 🇺🇸 English Version

**MeetingAssist** is a powerful **Open Source Web Application** designed to streamline your meeting productivity. It leverages advanced **Artificial Intelligence (AI)** to convert audio recordings into text (**Speech-to-Text**) and extract text from scanned documents (**OCR**) instantly.

Built with **Golang (Go)** for high performance and **Docker** for easy deployment, MeetingAssist is the perfect solution for students, professionals, and developers looking for a **free meeting minutes generator** and **document digitizer**.

### 🌟 Key Features
*   **🎙️ Smart Audio Transcription:** Upload MP3/WAV/M4A files and convert speech to text automatically using **AssemblyAI**. Specialized support for **Bahasa Indonesia** and English.
*   **📷 Document OCR Scanner:** Extract text from Images (JPG, PNG) and PDFs using **OCR.space**. Ideal for digitizing printed reports or handwritten notes.
*   **📝 Integrated Editor:** Review and edit your transcripts directly within the app before exporting.
*   **📱 Mobile Responsive:** Fully optimized for Smartphones, Tablets, and Desktop browsers.
*   **☁️ Cloud-Native:** Works seamlessly on **Railway**, Docker, or local environments without heavy dependencies.

### 🛠️ Tech Stack
*   **Backend:** Go (Golang) v1.21, Gin Framework (REST API)
*   **Database:** PostgreSQL (Supabase)
*   **Frontend:** HTML5, Bootstrap 5, Vanilla JavaScript
*   **AI Services:** AssemblyAI (Audio), OCR.space (Image/PDF)
*   **Infrastructure:** Docker, Railway

### 🚀 Installation Guide

#### Prerequisites
*   Go 1.21+
*   Git
*   PostgreSQL Database
*   API Keys (AssemblyAI & OCR.space - Free Tier)

#### Step-by-Step
1.  **Clone the Repository**
    ```bash
    git clone https://github.com/ibnujz/MeetingAssist.git
    cd MeetingAssist
    ```

2.  **Configuration (.env)**
    Create a `.env` file based on the example:
    ```env
    PORT=8080
    DB_HOST=your-supabase-url.co
    DB_PASSWORD=your-db-password
    ASSEMBLYAI_API_KEY=your-api-key
    ```

3.  **Run the App**
    ```bash
    go mod tidy
    go run main.go
    ```
    Visit `http://localhost:8080`

---
<br>

## 🇮🇩 Versi Bahasa Indonesia

**MeetingAssist** adalah **Aplikasi Web Transkripsi & OCR** cerdas yang dirancang untuk membantu Anda membuat notulen rapat dengan cepat. Menggunakan teknologi **AI (Kecerdasan Buatan)** terkini, aplikasi ini dapat mengubah rekaman suara menjadi teks dan memindai dokumen fisik menjadi teks digital yang bisa diedit.

Proyek ini sangat cocok untuk notulis, mahasiswa, dan profesional yang membutuhkan **aplikasi transkripsi bahasa Indonesia** gratis dan akurat.

### 🌟 Fitur Unggulan
*   **🎙️ Ubah Suara jadi Teks (Transkripsi):** Unggah rekaman rapat (MP3, WAV, M4A) dan dapatkan teksnya secara otomatis. Sangat akurat untuk **Bahasa Indonesia**.
*   **📷 Scan Foto ke Teks (OCR):** Ambil foto dokumen atau unggah PDF, lalu ekstrak tulisannya menjadi teks digital yang bisa dicopy-paste.
*   **📝 Editor Teks Bawaan:** Edit hasil transkrip langsung di browser sebelum disimpan atau dibagikan.
*   **📱 Akses Dimana Saja:** Tampilan modern yang ringan dan responsif di HP Android/iOS maupun Laptop.
*   **⚡ Cepat & Ringan:** Dibangun dengan **Golang**, aplikasi ini berjalan sangat cepat tanpa memberatkan server.

### 🛠️ Teknologi yang Digunakan
*   **Bahasa Pemrograman:** Go (Golang) - *Performa tinggi & stabil.*
*   **Framework Web:** Gin Gonic - *Framework HTTP tercepat untuk Go.*
*   **Database:** PostgreSQL - *Penyimpanan data aman & terstruktur.*
*   **Layanan AI:** AssemblyAI (untuk Audio), OCR.space (untuk Gambar).

### 🚀 Cara Instalasi (Lokal)

1.  **Download Source Code**
    Buka terminal dan jalankan:
    ```bash
    git clone https://github.com/ibnujz/MeetingAssist.git
    ```

2.  **Atur Konfigurasi**
    Buat file `.env` dan masukkan kredensial database serta API Key Anda (Dapatkan gratis di web AssemblyAI).

3.  **Jalankan Aplikasi**
    ```bash
    go mod tidy  # Download library
    go run main.go # Jalankan server
    ```
    Buka browser dan akses `http://localhost:8080`.

---

## 👨‍💻 Author & Social

Project ini dikembangkan oleh **Ibnu JZ**.
Jangan lupa follow Instagram saya untuk update proyek menarik lainnya!

<a href="https://instagram.com/ibnu.jz" target="_blank">
  <img src="https://img.shields.io/badge/Instagram-%40ibnu.jz-E4405F?style=for-the-badge&logo=instagram&logoColor=white" alt="Follow @ibnu.jz" height="40">
</a>

If you find this project useful, please give it a ⭐ **Star** on GitHub!

---
**© 2025 MeetingAssist by Ibnu JZ.**
