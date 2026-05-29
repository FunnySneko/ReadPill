# ReadPill
User-driven book review web app with suggestion algorithms.

**Backend:** Go + PostgreSQL    
**Frontend:** React     
**API:** OpenAPI     

## Features
- Add books to database
- Write and read reviews
- Advanced search with filters
- Personalized book suggestions

## Dependencies
- Go
- PostgreSQL
- migrate
- make
- npm

## Installation

### Install system dependencies
**Ubuntu/Debian**
```bash
sudo apt install golang-go postgresql postgresql-client make nodejs npm
```
**Fedora**
```bash
sudo dnf install golang postgresql postgresql-server make nodejs npm
```
**Arch Linux**
```bash
sudo pacman -S go postgresql make nodejs npm
```
**macOS**
```bash
brew install go postgresql node
```
### Install migrate CLI tool

See: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation

### Clone the repository

```bash
git clone https://github.com/FunnySneko/ReadPill.git
cd ReadPill
```

### Start backend

```bash
cd server
make setup
go run cmd/main.go
```

### Start frontend

```bash
cd client
npm install
npm start
```