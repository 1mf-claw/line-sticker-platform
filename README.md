# LINE 貼圖製作平台 / LINE Sticker Platform

![CI](https://github.com/1mf-claw/line-sticker-platform/actions/workflows/ci-release.yml/badge.svg)
![Release](https://img.shields.io/github/v/release/1mf-claw/line-sticker-platform)
![License](https://img.shields.io/badge/license-MIT-green)

---

## 中文

本專案是「LINE 貼圖製作平台」的正式專案，提供從 **主題 → 草稿 → 圖像 → 去背 → 匯出** 的完整流程，並支援 **BYOK（自帶 API Key）** 與 **每個功能可選不同 provider/model**。

### 主要功能
- Go + Vue 全端架構
- BYOK（API Key 只存記憶體，不落地）
- 文字 / 圖像 / 去背可分別選擇不同 provider/model
- 真實去背流程（支援 Replicate）
- 產出 ZIP（含 main.png / tab.png / metadata / report）
- i18n 多語系（繁體 / 簡體 / 英文 / 日文 / 韓文）

### 本機開發
**後端**
```bash
cd backend
go run ./cmd/app
```

**前端**
```bash
cd frontend
npm install
npm run dev
```

### 測試方式
1. 開啟 `http://localhost:5173`
2. 填入 BYOK Key
3. 依流程建立草稿 → 生成貼圖 → 去背 → 匯出

### 部署說明（簡版）
**後端**
```bash
cd backend
GOOS=linux GOARCH=amd64 go build -o app ./cmd/app
./app
```

**前端**
```bash
cd frontend
npm install
npm run build
# 將 dist/ 部署到任意靜態主機 (Vercel / Netlify / Cloudflare Pages)
```

### 截圖
> TODO: 請補上 UI 截圖

### 版本紀錄
- `v0.1.x`：MVP（草稿 / 生成 / 去背 / 匯出）

---

## English

This project is the official **LINE Sticker Platform**, providing an end‑to‑end flow from **theme → drafts → images → background removal → export**. It supports **BYOK (Bring Your Own Key)** and allows **per‑task provider/model selection**.

### Features
- Full‑stack Go + Vue
- BYOK (API keys are kept in memory only)
- Per‑task provider/model selection for text/image/bg
- Real background removal (Replicate supported)
- ZIP export (includes main.png / tab.png / metadata / report)
- i18n (zh‑TW / zh‑CN / EN / JA / KO)

### Local Development
**Backend**
```bash
cd backend
go run ./cmd/app
```

**Frontend**
```bash
cd frontend
npm install
npm run dev
```

### Quick Test
1. Open `http://localhost:5173`
2. Enter your BYOK key
3. Follow the flow → generate → remove background → export

### Deployment (short)
**Backend**
```bash
cd backend
GOOS=linux GOARCH=amd64 go build -o app ./cmd/app
./app
```

**Frontend**
```bash
cd frontend
npm install
npm run build
# Deploy dist/ to any static host (Vercel / Netlify / Cloudflare Pages)
```

### Screenshots
> TODO: add UI screenshots

### Changelog
- `v0.1.x`: MVP (drafts / generate / remove BG / export)
