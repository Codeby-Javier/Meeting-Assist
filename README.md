# MeetingAssist

Aplikasi web untuk transkripsi audio dan OCR dokumen dengan fitur edit dan export.

## Fitur

- 🎤 **Audio to Text**: Upload audio dan dapatkan transkrip otomatis (Bahasa Indonesia)
- 📄 **OCR Dokumen**: Upload gambar/PDF atau ambil foto langsung untuk extract text
- 📸 **Camera Capture**: Foto dokumen langsung dari HP untuk instant OCR
- ✏️ **Editor**: Edit hasil transkrip dan OCR
- 📥 **Export**: Download hasil dalam format PDF atau DOCX
- 📱 **Responsive**: Optimal di Desktop, Tablet, dan Mobile
- 🔐 **Authentication**: Sistem login dan register dengan JWT

## Tech Stack

- **Backend**: Golang (Gin Framework)
- **Database**: PostgreSQL (Supabase) dengan fallback SQLite
- **OCR**: Tesseract
- **Speech-to-Text**: Google Speech Recognition (Python)
- **PDF Processing**: PyMuPDF
- **Frontend**: Bootstrap 5, Vanilla JavaScript

## Prerequisites

- Go 1.21+
- Python 3.8+
- Tesseract OCR
- FFmpeg

## Installation

### 1. Clone Repository

```bash
git clone <repository-url>
cd MeetingAssist
```

### 2. Install Go Dependencies

```bash
go mod download
```

### 3. Install Python Dependencies

```bash
pip install SpeechRecognition PyMuPDF Pillow
```

### 4. Setup Environment Variables

Copy `.env.example` to `.env` and configure:

```env
PORT=8082
ENV=development

# Database (Supabase or leave as sqlite for local)
DB_HOST=your-supabase-host
DB_PORT=6543
DB_USER=postgres.xxx
DB_PASSWORD=your-password
DB_NAME=postgres
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=168h

# File Limits
MAX_AUDIO_SIZE=104857600
MAX_DOCUMENT_SIZE=52428800
MAX_STORAGE_PER_USER=1073741824

# Tesseract Path (Windows example)
TESSERACT_PATH=C:\Program Files\Tesseract-OCR\tesseract.exe

# Frontend URL
FRONTEND_URL=http://localhost:8082
```

### 5. Run Application

```bash
go run main.go
```

Application will be available at `http://localhost:8082`

## Deployment to Render

### 1. Create `render.yaml`

```yaml
services:
  - type: web
    name: meetingassist
    env: go
    buildCommand: go build -o main .
    startCommand: ./main
    envVars:
      - key: PORT
        value: 8082
      - key: ENV
        value: production
      # Add other env vars from your .env
```

### 2. Push to GitHub

```bash
git init
git add .
git commit -m "Initial commit"
git remote add origin <your-repo-url>
git push -u origin main
```

### 3. Deploy on Render

1. Go to [Render Dashboard](https://dashboard.render.com/)
2. Click "New +" → "Web Service"
3. Connect your GitHub repository
4. Render will auto-detect Go and deploy
5. Add environment variables in Render dashboard

## Project Structure

```
MeetingAssist/
├── config/          # Configuration files
├── handlers/        # HTTP request handlers
├── middleware/      # Auth, CORS, Logger middleware
├── models/          # Database models
├── routes/          # Route definitions
├── scripts/         # Python scripts for OCR & transcription
├── services/        # Business logic services
├── static/          # CSS, JS, images
├── templates/       # HTML templates
├── uploads/         # Uploaded files storage
├── utils/           # Utility functions
├── workers/         # Background workers for processing
├── main.go          # Application entry point
├── go.mod           # Go dependencies
└── .env             # Environment variables (not in git)
```

## Usage

### Upload Audio

1. Go to "Unggah Audio"
2. Drag & drop or select audio file
3. Wait for processing
4. View transcript in "Riwayat" → "Transkrip Audio"

### Upload Document

1. Go to "Unggah Dokumen"
2. Choose method:
   - Drag & drop file
   - Select from device
   - **Ambil Foto** (camera capture on mobile)
3. Wait for OCR processing
4. View result in "Riwayat" → "OCR Dokumen"

### Edit & Export

1. Click "Lihat/Edit" on completed item
2. Edit text as needed
3. Export as PDF or DOCX

## API Endpoints

### Authentication
- `POST /register` - Register new user
- `POST /login` - Login user

### Audio
- `POST /audio/upload` - Upload audio file
- `GET /audio` - Get user's audio list
- `GET /audio/:id` - Get audio detail
- `PUT /audio/:id` - Update transcript

### Document
- `POST /document/upload` - Upload document
- `GET /document` - Get user's documents
- `GET /document/:id` - Get document detail
- `PUT /document/:id` - Update OCR text

### Export
- `POST /export/pdf/:type/:id` - Export as PDF
- `POST /export/docx/:type/:id` - Export as DOCX

## License

MIT

## Author

By [@ibnu.jz](https://instagram.com/ibnu.jz)
