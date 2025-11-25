# 🚗 Parking Lot Management System (Golang + MySQL)

A complete backend system built in **Go** for managing parking lots, parking slots, users, and daily parking analytics.

Supports **2 user roles**:
- **Manager**
- **Vehicle Owner (User)**

---

## ✨ Features

### 👨‍💼 Manager Features
- Create Parking Lots
- Create Parking Slots inside any lot
- View Parking Lot Status — which vehicle is parked in which slot
- Mark any slot:  
  - `active` (available)  
  - `engaged`  
  - `under_maintenance`
- View Daily Report:
  - Total vehicles parked today
  - Total parking duration
  - Total fee collected

---

### 🚘 User Features
- Choose a parking lot & park in **nearest available slot**
- Unpark vehicle and get:
  - `success` message
  - Parking fee = `10 × total hours rounded up`
    - Example: Parked 1h 5m → 10 × 2 = **20 Tk**

---
### 🚀 Running the App
- Install dependencies
```
go mod tidy
```

#### Run server
```
go run cmd/main.go
```

#### Server runs at:
```
http://localhost:8080
```
---
### ❤️ Author
- Safayet Shawn
- Golang Backend Developer
- GitHub: https://github.com/Safayet-Shawn
---
