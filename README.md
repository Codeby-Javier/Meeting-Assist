# 🎙️ MeetingAssist

> **Demo Live:** [https://meeting-assist-production.up.railway.app/](https://meeting-assist-production.up.railway.app/)

**MeetingAssist** is a smart web application designed to help you transcribe meeting recordings into text and extract text from documents (OCR) instantly. Built with modern Golang technology and Cloud AI APIs.

**MeetingAssist** adalah aplikasi web cerdas untuk membantu Anda mentranskripsi rekaman rapat menjadi teks dan mengambil teks dari dokumen (OCR) secara instan. Dibangun dengan teknologi Golang modern dan Cloud AI APIs.

---

## 🛠️ Technology Stack (Teknologi yang Digunakan)

### Core Backend
- **Language:** Go (Golang) v1.21+
  - *High performance, type-safe, and concurrent.*
- **Framework:** Gin Web Framework
  - *Fastest HTTP web framework for Go.*
- **Database:** PostgreSQL (via Supabase)
  - *Reliable relational database for storing users and history.*
- **ORM:** GORM
  - *Object Relational Mapping for easy database interactions.*

### Cloud AI Services (No Local Dependencies)
- **Audio-to-Text:** [AssemblyAI API](https://www.assemblyai.com/)
  - *Features: High accuracy for Indonesian/English, Automatic Punctuation, Speaker Diarization.*
- **OCR (Document-to-Text):** [OCR.space API](https://ocr.space/)
  - *Features: Supports PDF/JPG/PNG, auto-orientation, table recognition.*

### Frontend
- **Framework:** HTML5 + Bootstrap 5
  - *Responsive design for Mobile, Tablet, and Desktop.*
- **Styling:** Custom CSS + FontAwesome Icons
- **Scripting:** Vanilla JavaScript (ES6)

### Deployment
- **Platform:** Railway.app
- **Container:** Docker (Alpine Linux)
  - *Ultra-lightweight image (~50MB).*

---

## ✨ Features (Fitur Unggulan)

1.  **🎙️ Audio Transcription (Transkripsi Audio)**
    - Upload MP3, WAV, M4A files directly.
    - Automatic conversion to text (supports Bahasa Indonesia).
    - Status tracking (Uploaded -> Processing -> Completed).
    - **Cloud Upload:** Direct stream to cloud for stability.

2.  **📄 Document OCR (Scan Dokumen)**
    - Upload Images (JPG, PNG) or PDFs.
    - **Camera Capture:** Directly take photos of documents from mobile phone.
    - Extract text instantly for editing.

3.  **📝 Smart Editor & Export**
    - Built-in text editor to correct transcripts.
    - **Export:** Download results as PDF (`.pdf`) or Word (`.docx`).

4.  **📱 Fully Responsive**
    - Works perfectly on Smartphones, Tablets, and Laptops.
    - Modern Dashboard UI with Sidebar navigation.

---

## 🚀 Installation Guide (Panduan Instalasi)

### Prerequisites (Persyaratan)
Before you start, ensure you have:
- **Go 1.21** or newer installed.
- **Git** installed.
- **PostgreSQL Database** (recommended: [Supabase](https://supabase.com/) free tier).
- **AssemblyAI API Key** (Get free key at [assemblyai.com](https://www.assemblyai.com/)).

### Step 1: Clone Repository
```bash
git clone https://github.com/your-username/MeetingAssist.git
cd MeetingAssist
```

### Step 2: Configure Environment
Create a `.env` file in the root directory (or use System Environment Variables):

```env
# Server Config
PORT=8082
ENV=development

# Database (Supabase/PostgreSQL)
DB_HOST=your-db.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=postgres
DB_SSLMODE=disable

# Security (JWT)
JWT_SECRET=rahasia_super_aman_123
JWT_EXPIRATION=168

# Cloud APIs (REQUIRED)
# Get free key at assemblyai.com
ASSEMBLYAI_API_KEY=your-assemblyai-key
# OCR.space uses a shared free key by default, but you can set yours:
# OCR_API_KEY=your-ocr-key

# Limits
MAX_AUDIO_SIZE=104857600      # 100MB
MAX_DOCUMENT_SIZE=52428800    # 50MB
```

### Step 3: Run Locally
```bash
# Download dependencies
go mod tidy

# Run application
go run main.go
```
Access the app at: `http://localhost:8082`

---

## ☁️ Deployment Guide (Cara Deploy ke Railway)

This project is optimized for **Railway.app** (Free Trial / Starter Plan).

1.  **Fork/Push** this repository to your GitHub.
2.  Login to **Railway.app** and choose "Deploy from GitHub".
3.  Select this repository.
4.  Go to **Variables** tab and add the Environment Variables listed above (especially `DB_URL` and `ASSEMBLYAI_API_KEY`).
5.  Railway will automatically detect the `Dockerfile` and build the app.
6.  Go to **Settings** -> **Generate Domain** to get your public URL.

### Dockerfile Optimization
The project uses a **Multi-stage Dockerfile** to ensure a small footprint:
- **Build Stage:** Golang Image (compiles code).
- **Final Stage:** Alpine Linux (runs binary).
- *Result:* No heavy dependencies like Python or Tesseract needed on the server!

---

## 📝 License
MIT License - Free to use and modify.

---

**Developed with ❤️ using Golang**
