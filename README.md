# MeetingAssist

Aplikasi web untuk transkripsi audio dan OCR dokumen.

## Fitur
- 🎤 Audio to Text (Bahasa Indonesia)
- 📄 OCR Dokumen (gambar/PDF)
- 📸 Camera Capture
- ✏️ Editor & Export (PDF/DOCX)
- 📱 Responsive Design

## Tech Stack
- Backend: Golang
- Database: PostgreSQL (Supabase)
- OCR: OCR.space API
- Speech-to-Text: AssemblyAI API
- Frontend: Bootstrap 5

## Railway Deployment

### 1. Push ke GitHub
```bash
git add .
git commit -m "Cloud-based deployment ready"
git push
```

### 2. Deploy di Railway
1. Login ke Railway.app
2. New Project → Deploy from GitHub
3. Pilih repository
4. Add environment variables (lihat .env.example)
5. Deploy!

### 3. Environment Variables
```
PORT=8082
ENV=production
DB_HOST=your-supabase-host
DB_PORT=6543
DB_USER=your-supabase-user
DB_PASSWORD=your-password
DB_NAME=postgres
DB_SSLMODE=disable
JWT_SECRET=your-secret
JWT_EXPIRATION=168h
FRONTEND_URL=https://your-app.up.railway.app
```

### 4. Generate Public URL
Di Railway Settings → Generate Domain

## API Keys (Gratis)
- **OCR.space**: K87899142388957 (25,000 requests/bulan)
- **AssemblyAI**: 4d2a2e6811d04ea99f0c5fa89090ddc6 ($50 credit)

Bisa daftar key pribadi untuk limit lebih besar.

## License
MIT

## Author
By [@ibnu.jz](https://instagram.com/ibnu.jz)
