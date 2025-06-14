# Food Delivery App Server Side

<img src="https://go.dev/blog/gopher/gopher.png" alt="Gopher" width="100"/>

![Go](https://img.shields.io/badge/Go-1.24.4-blue?logo=go)


💫 In progress... To be written in GoLang

```txt
food-delivery-app/
├── 🔗 cmd/
│   └── 🛜 server/              # Entry point (main.go)
│
├── 🏢 infrastructure/          # Gin setup, routes, DB connect
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