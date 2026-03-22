# Nexus ROMs

![Status](https://img.shields.io/badge/Status-Beta-blue)
![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen)
![Coverage](https://img.shields.io/badge/Coverage-74.8%25-green)
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)
![Wails](https://img.shields.io/badge/Wails-v2-red)
![Svelte](https://img.shields.io/badge/Frontend-Svelte-ff3e00?logo=svelte)

Nexus ROMs is a high-performance desktop application built with **Wails** and **Svelte**, designed to organize, identify, and rename your retro gaming collection with precision and ease.

## 🌟 Key Features

- **🚀 Concurrent Orchestration**: Multi-threaded scanning engine (configurable worker counts) for blazing fast identification.
- **🔍 Multi-Provider Scraping**: Integrated support for world-class metadata providers:
  - [ScreenScraper](https://www.screenscraper.fr/)
  - [TheGamesDB](https://thegamesdb.net/)
- **🧩 Smart Templates**: Rename your collection using dynamic tokens like `{Name}`, `{Region}`, `{Platform}`, `{Year}`, and `{Developer}`.
- **💾 Session Recovery**: Robust session management allows you to pause, resume, or restart scanning operations without losing progress.
- **💎 Premium UI**: A reactive, dark-mode dashboard providing real-time logs, live identification tables, and progress statistics.
- **🔒 Security First**: API credentials and personal data are securely stored and handled.
- **📦 Wide Compression Support**: Native handling of `.zip`, `.rar`, and `.7z` archives.
- **🛡️ Smart Filtering**: Automatic detection and skipping of massive romsets to focus on high-quality single titles.
- **📉 API Quota Tracking**: Persistent monthly tracking of ScreenScraper requests to stay within tier limits.

## 🚀 Getting Started

### Prerequisites

To compile Nexus ROMs from source, you will need:

- **Go** (1.24 or higher)
- **Node.js** (v20+) & **NPM**
- **Wails CLI**: Install via `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **C Compiler**: Required for SQLite (e.g., GCC on Linux, MingW on Windows).

### 🔑 API Configuration

To use the identification features, you will need accounts from these providers:

1. **ScreenScraper**: Register at [screenscraper.fr](https://www.screenscraper.fr/).
2. **TheGamesDB**: Get an API key at [thegamesdb.net](https://thegamesdb.net/).

Once you have your credentials, enter them in the **Settings** tab within the application.

### 🏗️ Installation & Development

1. Clone the repository:

   ```bash
   git clone https://github.com/ricardoaevt/nexus-roms.git
   cd nexus-roms
   ```

2. Install dependencies:

   ```bash
   wails doctor
   ```

3. Run in development mode (Live Reload):

   ```bash
   wails dev
   ```

### 🔨 Building for Production

Generate a production-ready binary for your current OS:

```bash
wails build
```

The compiled binary will be located in the `build/bin` directory.

## 📦 Automated Releases

This repository is configured with **GitHub Actions** to automatically build and package the application for **Windows, Linux, and macOS**. 

**To trigger a new release:**

1. Update the version in `wails.json`.
2. Create and push a tag:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

The binaries will appear in the [Releases](https://github.com/ricardoaevt/nexus-roms/releases) section once the workflow finishes.

## 🛠️ Tech Stack

- **Backend**: Go 1.24+
- **Frontend**: Svelte
- **Database**: SQLite (Encrypted)
- **Framework**: [Wails v2](https://wails.io/)

## 🔄 Workflow & Smart Logic

Nexus ROMs goes beyond simple renaming by implementing advanced decision-making logic:

### 1. Smart Archive Filtering
The system automatically evaluates the contents of compressed files (`.zip`, `.rar`, `.7z`):
- **Scenario**: If an archive contains multiple unrelated ROMs (e.g., a "100-in-1" pack), it is flagged as a *romset* and skipped to prevent cluttering.
- **Scenario**: If multiple files share a similar base name (e.g., "Game Disc 1" and "Game Disc 2"), the system recognizes them as parts of a single title and proceeds with identification.

### 2. Intelligent Collision Handling
To prevent data loss and maintain backup integrity:
- **Scenario**: When a target filename already exists, the system doesn't overwrite it. Instead, it moves the source file to a dedicated `duplicados/` folder.
- **Scenario**: The original path structure is recreated inside `duplicados/` to ensure you know exactly where each file originated.

### 3. API Quota Management
For providers with daily or monthly limits like ScreenScraper:
- **Scenario**: The system tracks the current month and the number of requests made.
- **Scenario**: At the start of a new month, the counter automatically resets to ensure accurate tracking against your subscription tier.

### 4. Robust Error Reporting
At the end of every batch operation, a detailed summary is displayed:
- **Scenario**: Files locked by the OS or permission issues are gracefully logged and presented in a final report, allowing you to troubleshoot without stopping the entire process.

## 🤝 Contributing

1. Fork the Project.
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`).
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the Branch (`git push origin feature/AmazingFeature`).

## 🧪 Testing & Quality

Nexus ROMs is committed to high software quality. Our backend suite ensures critical logic remains robust.

| Package | Coverage | Status |
| :--- | :--- | :--- |
| `internal/renamer` | 90.9% | ✅ Robust |
| `internal/orchestrator` | 76.8% | ✅ Tested |
| `internal/scraper` | 70.6% | ✅ Tested |
| `internal/db` | 81.7% | ✅ Robust |
| `internal/crypto` | 75.9% | ✅ Secure |

To run the tests:

```bash
make test # Runs all unit tests
```

To view coverage:

```bash
make coverage # Generates and opens HTML report
```

---
*Developed with focus on speed, aesthetics, and user experience for the Retro Gaming Community.*
