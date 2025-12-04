# MeetingAssist - Railway Deployment

## ✅ DEPLOYMENT BERHASIL!

Aplikasi sudah ter-deploy di Railway, tapi OCR dan Audio-to-Text menggunakan **Cloud API** untuk reliability.

## 🔑 API Keys yang Digunakan

### 1. OCR.space (Gratis)
- **Status**: ✅ Sudah terintegrasi
- **Limit**: 25,000 requests/bulan
- **API Key**: K87899142388957 (shared free key)
- **Upgrade**: Bisa daftar di https://ocr.space/ocrapi untuk key pribadi

### 2. Google Speech Recognition (Gratis)
- **Status**: ✅ Menggunakan library SpeechRecognition
- **Limit**: Unlimited (tapi ada rate limiting)
- **Bahasa**: Indonesia (id-ID)

## 📝 Cara Kerja di Production

### OCR (Document to Text):
1. User upload gambar/PDF
2. File di-preprocess (resize, grayscale, contrast)
3. Dikirim ke **OCR.space API**
4. Hasil text dikembalikan
5. Disimpan ke database

### Speech-to-Text (Audio):
1. User upload audio
2. Convert ke WAV format (FFmpeg)
3. Dikirim ke **Google Speech Recognition**
4. Hasil transkrip dikembalikan
5. Disimpan ke database

## 🚀 Testing

### Test OCR:
1. Buka https://[your-railway-url]/document/upload
2. Upload gambar dengan teks jelas
3. Tunggu processing (~5-10 detik)
4. Check di Riwayat → OCR Dokumen

### Test Audio:
1. Buka https://[your-railway-url]/audio/upload
2. Upload file audio (MP3/WAV)
3. Tunggu processing (~10-30 detik tergantung durasi)
4. Check di Riwayat → Transkrip Audio

## ⚠️ Limitations

### Free Tier Limits:
- **OCR.space**: 25,000 requests/bulan
- **Railway**: 500 execution hours/bulan
- **Google Speech**: Rate limited (biasanya cukup untuk personal use)

### Jika Limit Tercapai:
1. **OCR**: Daftar API key pribadi di ocr.space
2. **Audio**: Upgrade ke AssemblyAI atau Google Cloud Speech-to-Text
3. **Railway**: Upgrade ke paid plan

## 🔧 Troubleshooting

### Jika OCR Failed:
- Check Railway logs
- Pastikan gambar tidak terlalu besar (max 3MB)
- Pastikan gambar jelas dan kontras tinggi

### Jika Audio Failed:
- Check format audio (MP3, WAV, M4A supported)
- Pastikan audio tidak terlalu panjang (max 5 menit recommended)
- Check Railway logs untuk error detail

## 📊 Monitoring

Check Railway Dashboard untuk:
- Deployment logs
- Resource usage
- Error messages
- Request metrics

## 🎯 Next Steps

1. ✅ Test semua fitur
2. ✅ Monitor usage
3. 📝 Jika perlu upgrade, consider:
   - Railway Starter plan ($7/month)
   - OCR.space Pro ($9.99/month)
   - AssemblyAI ($0.00025/second)

## 💡 Tips

- Upload gambar dengan pencahayaan baik
- Pastikan teks jelas dan kontras
- Audio dengan noise minimal
- Gunakan format standar (JPG, PNG, MP3, WAV)
