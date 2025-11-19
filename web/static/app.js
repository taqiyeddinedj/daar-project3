const API_BASE = '/api';

let currentPage = 1;
let currentQuery = '';
let totalPages = 1;
let currentBookId = null;
let readerFontSize = 16;
let searchType = 'keyword';
let searchResults = []; // Store full results with occurrences

document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('search-btn').addEventListener('click', () => performSearch(1));
    
    document.getElementById('search-input').addEventListener('keypress', (e) => {
        if (e.key === 'Enter') performSearch(1);
    });

    // Tab switching with placeholder update
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
            this.classList.add('active');
            searchType = this.dataset.type;
            
            // Update placeholder
            const placeholder = this.dataset.placeholder;
            document.getElementById('search-input').placeholder = placeholder;
        });
    });
});

function performSearch(page = 1) {
    const query = document.getElementById('search-input').value.trim();
    if (!query) {
        alert('Please enter a search query');
        return;
    }

    currentQuery = query;
    currentPage = page;

    showLoading();
    document.getElementById('results-section').classList.remove('hidden');

    fetch(`${API_BASE}/search?q=${encodeURIComponent(query)}&type=${searchType}&page=${page}`)
        .then(res => res.json())
        .then(data => {
            hideLoading();
            searchResults = data.results || []; // Store results with occurrences
            displayResults(data);
            updatePagination(data);
        })
        .catch(err => {
            hideLoading();
            console.error('Search error:', err);
            alert('Search failed: ' + err.message);
        });
}

function displayResults(data) {
    const grid = document.getElementById('books-grid');
    const count = document.getElementById('results-count');

    count.textContent = `Found ${data.total_count} books - Page ${data.page} of ${data.total_pages}`;

    if (!data.books || data.books.length === 0) {
        grid.innerHTML = '<p style="grid-column:1/-1;text-align:center;color:#808080;padding:2rem;">No books found</p>';
        return;
    }

    // Get occurrences for current page books
    const booksWithOccurrences = data.books.map(book => {
        const result = searchResults.find(r => r.book.id === book.id);
        return {
            ...book,
            occurrences: result ? result.occurrences : 0
        };
    });

    grid.innerHTML = booksWithOccurrences.map(book => `
        <div class="book-card" onclick="loadBookDetails(${book.id})">
            <div class="book-card-header">
                <div class="book-icon">📖</div>
                <div class="book-card-info">
                    <h4>${escapeHtml(book.title)}</h4>
                    <p class="book-author">by ${escapeHtml(book.author)}</p>
                    <div class="book-meta">
                        <span class="book-badge">ID: ${book.id}</span>
                        <span class="book-badge">${book.word_count.toLocaleString()} words</span>
                        ${book.occurrences > 0 ? `<span class="book-badge badge-matches">${book.occurrences} matches</span>` : ''}
                    </div>
                </div>
            </div>
        </div>
    `).join('');
}

function updatePagination(data) {
    const pagination = document.getElementById('pagination');
    totalPages = data.total_pages;

    if (totalPages <= 1) {
        pagination.classList.add('hidden');
        return;
    }

    pagination.classList.remove('hidden');

    let html = '';

    // Previous button
    html += `<button class="page-btn" onclick="performSearch(${currentPage - 1})" ${currentPage === 1 ? 'disabled' : ''}>
        <span class="page-arrow">←</span> Previous
    </button>`;

    // Page numbers
    const startPage = Math.max(1, currentPage - 2);
    const endPage = Math.min(totalPages, currentPage + 2);

    if (startPage > 1) {
        html += `<button class="page-num" onclick="performSearch(1)">1</button>`;
        if (startPage > 2) html += `<span class="page-dots">...</span>`;
    }

    for (let i = startPage; i <= endPage; i++) {
        html += `<button class="page-num ${i === currentPage ? 'active' : ''}" onclick="performSearch(${i})">${i}</button>`;
    }

    if (endPage < totalPages) {
        if (endPage < totalPages - 1) html += `<span class="page-dots">...</span>`;
        html += `<button class="page-num" onclick="performSearch(${totalPages})">${totalPages}</button>`;
    }

    // Next button
    html += `<button class="page-btn" onclick="performSearch(${currentPage + 1})" ${currentPage === totalPages ? 'disabled' : ''}>
        Next <span class="page-arrow">→</span>
    </button>`;

    pagination.innerHTML = html;
}

function loadBookDetails(bookId) {
    currentBookId = bookId;
    showBookView();
    showLoading();

    Promise.all([
        fetch(`${API_BASE}/book/${bookId}`).then(r => r.json()),
        fetch(`${API_BASE}/recommendations/${bookId}`).then(r => r.json())
    ])
    .then(([book, recommendations]) => {
        hideLoading();
        displayBookDetails(book);
        displayRecommendations(recommendations);
    })
    .catch(err => {
        hideLoading();
        console.error('Load book error:', err);
        alert('Failed to load book: ' + err.message);
    });
}

function displayBookDetails(book) {
    const details = document.getElementById('book-details');
    details.innerHTML = `
        <div class="book-header">
            <div class="book-cover">📚</div>
            <div class="book-info">
                <h2>${escapeHtml(book.title)}</h2>
                <div class="book-info-row">
                    <span class="book-info-label">Author:</span>
                    <span class="book-info-value">${escapeHtml(book.author)}</span>
                </div>
                <div class="book-info-row">
                    <span class="book-info-label">Book ID:</span>
                    <span class="book-info-value">${book.id}</span>
                </div>
                <div class="book-info-row">
                    <span class="book-info-label">Word Count:</span>
                    <span class="book-info-value">${book.word_count.toLocaleString()}</span>
                </div>
                <div class="book-actions">
                    <button class="btn-primary" onclick="loadBookContent(${book.id})">📖 Read Book</button>
                </div>
            </div>
        </div>
    `;
}

function displayRecommendations(recommendations) {
    const grid = document.getElementById('recommendations-grid');

    if (!recommendations || recommendations.length === 0) {
        grid.innerHTML = '<p style="grid-column:1/-1;color:#808080;padding:1rem;">No recommendations available</p>';
        return;
    }

    grid.innerHTML = recommendations.map(book => `
        <div class="book-card" onclick="loadBookDetails(${book.id})">
            <div class="book-card-header">
                <div class="book-icon">📖</div>
                <div class="book-card-info">
                    <h4>${escapeHtml(book.title)}</h4>
                    <p class="book-author">by ${escapeHtml(book.author)}</p>
                </div>
            </div>
        </div>
    `).join('');
}

function loadBookContent(bookId) {
    showReaderView();
    showLoading();

    fetch(`${API_BASE}/content/${bookId}`)
        .then(res => {
            if (!res.ok) {
                return res.json().then(err => {
                    throw new Error(err.error + ' - Check server logs for details');
                });
            }
            return res.json();
        })
        .then(data => {
            hideLoading();
            if (!data || !data.content) {
                throw new Error('No content received from server');
            }
            displayBookContent(data);
        })
        .catch(err => {
            hideLoading();
            console.error('Load content error:', err);
            alert('Failed to load book content: ' + err.message + '\n\nCheck the server terminal for file path details.');
            if (currentBookId) {
                loadBookDetails(currentBookId);
            }
        });
}

function displayBookContent(data) {
    const content = document.getElementById('reader-content');
    content.textContent = data.content;
    content.style.fontSize = `${readerFontSize}px`;
}

function closeReader() {
    if (currentBookId) {
        loadBookDetails(currentBookId);
    } else {
        showHome();
    }
}

function changeFontSize(delta) {
    readerFontSize = Math.max(12, Math.min(24, readerFontSize + delta));
    document.getElementById('reader-content').style.fontSize = `${readerFontSize}px`;
}

function toggleTheme() {
    document.querySelector('.reader-container').classList.toggle('light');
}

function showHome() {
    hideAllViews();
    document.getElementById('home-view').classList.add('active');
}

function showBookView() {
    hideAllViews();
    document.getElementById('book-view').classList.add('active');
}

function showReaderView() {
    hideAllViews();
    document.getElementById('reader-view').classList.add('active');
}

function hideAllViews() {
    document.querySelectorAll('.view').forEach(v => v.classList.remove('active'));
}

function showLoading() {
    document.getElementById('loading').classList.remove('hidden');
}

function hideLoading() {
    document.getElementById('loading').classList.add('hidden');
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}