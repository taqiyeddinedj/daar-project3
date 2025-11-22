# Digital Library - Book Search Engine

A fast search engine for digital libraries with 1664+ books from Project Gutenberg.

**Live Demo:** [library.taqiyeddine.tech](https://library.taqiyeddine.tech/)

## Features

- **Fast Search**: Find books in less than 1ms
- **Regex Search**: Advanced pattern matching
- **Smart Ranking**: PageRank algorithm
- **Recommendations**: Similar books based on Jaccard similarity
- **Modern UI**: Clean interface inspired by Z-Library

## Tech Stack

- **Backend**: Go + Gin framework
- **Frontend**: HTML, CSS, JavaScript
- **Algorithms**: Index inversé, PageRank, Jaccard similarity

---

## Quick Start

### Prerequisites

```bash
# Install Go 1.21+
sudo apt install golang-go
```

### Automated Setup (Recommended)

```bash
# Make scripts executable
chmod +x setup.sh setup.sh

# Run setup (download books+ builds index + graph + binary)
./setup.sh



# Run the Server

**Option A: Run directly (development)**

```bash
go run cmd/server/main.go

# Server runs on: http://localhost:8080
```

**Option B: Build binary (production)**

```bash
# Build the binary
go build -o book-server cmd/server/main.go

# Run it
./book-server
```

**Option C: Run as daemon (Linux systemd)**

```bash
# Create systemd service file
sudo nano /etc/systemd/system/book-library.service
```

Add this content:

```ini
[Unit]
Description=Book Library Search Engine
After=network.target

[Service]
Type=simple
User=YOUR_USERNAME
WorkingDirectory=/path/to/your/project
ExecStart=/path/to/your/project/book-server
Restart=always

[Install]
WantedBy=multi-user.target
```

Start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable book-library
sudo systemctl start book-library
sudo systemctl status book-library
```

---

## Project Structure

```
.
├── cmd/
│   ├── indexer/     # Build inverted index
│   ├── graph/       # Build Jaccard graph
│   └── server/      # Web server
├── pkg/
│   ├── indexer/     # Index data structures
│   ├── search/      # Search algorithms
│   ├── graph/       # Jaccard graph
│   └── ranking/     # PageRank algorithm
├── web/
│   ├── templates/   # HTML files
│   └── static/      # CSS + JavaScript
└── data/
    ├── books/       # Text files (not in repo)
    ├── index.json   # Inverted index
    └── jaccard_graph.json  # Similarity graph
```

---

## How It Works

### 1. Inverted Index

Maps each word to books containing it:

```
"love" → {book_1: 45 times, book_2: 12 times, ...}
"whale" → {book_3: 89 times, book_5: 3 times, ...}
```

**Why?** Allows search in O(1) time instead of scanning all books.

### 2. Jaccard Similarity

Measures how similar two books are:

```
Jaccard(A, B) = (common words) / (total unique words)
```

Example:
- Book A: {love, hate, war, peace}
- Book B: {love, war, hope, joy}
- Common: {love, war} = 2
- Total unique: 6
- Similarity: 2/6 = 0.33

**Why?** Books sharing many words are probably about similar topics.

### 3. PageRank

Ranks books by "importance" in the graph:

```
Important book = Many similar books point to it
```

Same algorithm Google uses for web pages!

**Why?** Better results show "central" books first.

---

## API Endpoints

```
GET  /                           # Home page
GET  /api/search?q=love          # Simple search
GET  /api/search?q=wha.*&type=regex  # Regex search
GET  /api/book/:id               # Book details
GET  /api/recommendations/:id    # Similar books
GET  /api/content/:id            # Book content
```

---

## Performance

Tested on 1664 books:

| Operation | Time |
|-----------|------|
| Simple search | < 1ms |
| Regex search | 150-500ms |
| Recommendations | < 0.01ms |
| Index build | ~4 min |
| Graph build | ~18 min |

---

## Configuration

Edit these values in `cmd/server/main.go`:

```go
// Port
r.Run(":8081")  // Change port here

// Jaccard threshold
threshold := 0.1  // Books with >10% similarity

// PageRank iterations
iterations := 20  // More = more accurate
```

---

## Deployment

### With Nginx (reverse proxy)

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:8081;
        proxy_set_header Host $host;
    }
}
```

### With SSL (Let's Encrypt)

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com
```

---

## Future Improvements

- [ ] Stemming (running → run)
- [ ] TF-IDF ranking
- [ ] Filters (language, author, genre)
- [ ] Snippets with highlighted keywords
- [ ] User accounts
- [ ] Reading history
- [ ] Mobile app

---

## Credits

- Books from [Project Gutenberg](https://www.gutenberg.org/)
- UI inspired by Z-Library
- Built for DAAR course at Sorbonne Université

---

## Author

**Taqiyeddine DJOUANI**
- Website: [taqiyeddine.tech](https://taqiyeddine.tech)
- Project: [library.taqiyeddine.tech](https://library.taqiyeddine.tech/)
