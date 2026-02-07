# LINE 貼圖製作平台

![CI](https://github.com/1mf-claw/line-sticker-platform/actions/workflows/ci-release.yml/badge.svg)
![Release](https://img.shields.io/github/v/release/1mf-claw/line-sticker-platform)
![License](https://img.shields.io/badge/license-MIT-green)

[English Version → README.en.md](README.en.md)

---

本專案是「LINE 貼圖製作平台」的正式專案，提供從 **主題 → 草稿 → 圖像 → 去背 → 匯出** 的完整流程，並支援 **BYOK（自帶 API Key）** 與 **每個功能可選不同 provider/model**。

## 主要功能
- Go + Vue 全端架構
- BYOK（API Key 只存記憶體，不落地）
- 文字 / 圖像 / 去背可分別選擇不同 provider/model
- 真實去背流程（支援 Replicate）
- 產出 ZIP（含 main.png / tab.png / metadata / report）
- i18n 多語系（繁體 / 簡體 / 英文 / 日文 / 韓文）

## 本機開發
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

## 測試方式
1. 開啟 `http://localhost:5173`
2. 填入 BYOK Key
3. 依流程建立草稿 → 生成貼圖 → 去背 → 匯出

## 部署說明（簡版）
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

## 截圖
![UI Overview](screenshots/overview.svg)

## 版本紀錄
- `v0.1.x`：MVP（草稿 / 生成 / 去背 / 匯出）
