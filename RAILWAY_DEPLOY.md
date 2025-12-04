# Railway Deployment Guide

## Deploy MeetingAssist ke Railway.app

### Prerequisites
- GitHub account
- Repository MeetingAssist sudah di-push ke GitHub

### Step-by-Step

#### 1. Sign Up Railway
1. Buka https://railway.app/
2. Klik "Login with GitHub"
3. Authorize Railway

#### 2. Create New Project
1. Klik "New Project"
2. Pilih "Deploy from GitHub repo"
3. Pilih repository "MeetingAssist"
4. Railway akan auto-detect Go dan mulai build

#### 3. Configure Environment Variables
Klik tab "Variables" dan tambahkan:

```
PORT=8082
ENV=production
DB_HOST=aws-0-ap-southeast-1.pooler.supabase.com
DB_PORT=6543
DB_USER=postgres.bosrsajnvtwjfgxqbput
DB_PASSWORD=[your-supabase-password]
DB_NAME=postgres
DB_SSLMODE=disable
JWT_SECRET=[generate-random-string]
JWT_EXPIRATION=168h
MAX_AUDIO_SIZE=104857600
MAX_DOCUMENT_SIZE=52428800
MAX_STORAGE_PER_USER=1073741824
TESSERACT_PATH=tesseract
FRONTEND_URL=https://meetingassist.up.railway.app
```

#### 4. Deploy
Railway akan otomatis build dan deploy.

#### 5. Get URL
Setelah deploy sukses, klik "Settings" → "Generate Domain" untuk mendapatkan public URL.

### Troubleshooting

**Jika build gagal:**
1. Check logs di Railway dashboard
2. Pastikan `nixpacks.toml` ada di root folder
3. Pastikan semua environment variables sudah di-set

**Jika app crash:**
1. Check logs untuk error message
2. Pastikan Supabase credentials benar
3. Test koneksi database

### Free Tier Limits
- 500 execution hours/month
- $5 credit/month
- Cukup untuk 1 app running 24/7

### Alternative: Fly.io

Jika Railway tidak work, coba Fly.io:

```bash
# Install Fly CLI
iwr https://fly.io/install.ps1 -useb | iex

# Login
fly auth login

# Launch app
fly launch

# Deploy
fly deploy
```

## Support

Jika ada masalah, check:
- Railway logs
- Supabase connection
- Environment variables
