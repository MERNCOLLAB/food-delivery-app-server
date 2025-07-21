# Food Delivery App Server Side

<img src="https://go.dev/blog/gopher/gopher.png" alt="Gopher" width="100"/>

![Go](https://img.shields.io/badge/Go-1.24.5-blue?logo=go)


## 🚀 Getting Started

Follow these steps to set up and run the server locally:

1. **Clone the repository**
   ```bash
   git clone https://github.com/MERNCOLLAB/food-delivery-app-server
   cd food-delivery-app-server
   ```
2. **Install dependencies**
   ```bash
   go mod download
   ```
3. **Add your environment variables**
   - Copy the example file (if available) or create a `.env` file in the project root
   
   - Fill in the required environment variables (DB connection, secrets, etc)

4. **Run the server with Air (live reload)**

	Command for running is air (air.toml is already configured for this project)
   ```bash
   air
   ```
   - If you don't have Air installed, install it with:
     ```bash
     go install github.com/cosmtrek/air@latest
     ```


## 📂 Project Folder Structure Guide

```txt
food-delivery-app/
├── 🔗 cmd/
│   └── 🛜 server/              # Entry point (main.go)
│
├── 🏢 infrastructure/          # Gin setup, routes, DB connect, Redis connect
│
├── 🌐 internal/                # Features: auth, user, order (handlers, services, repos, DTOs)
│
├── 💾 models/                  # App-wide structs (User, Order, etc.)
│
├── ⚙️ config/                  # Environment loading, config helpers
│
├── 🔐 middleware/              # Auth, role guard, file upload
│
├── 📦 pkg/                     # Shared utilities and helpers
│
├── ✈️ .air.toml                # Live reload config
├── 📖 go.mod
└── 📝 README.md

```

## Delivery Transaction API Endpoints

![Delivery_Transaction](delivery_transaction.png)