function search() {
    const keyword = document.getElementById('searchInput').value;
    
    fetch(`/api/search?q=${encodeURIComponent(keyword)}`)
        .then(res => res.json())
        .then(data => displayResults(data.results))
        .catch(err => console.error(err));
}

function regexSearch() {
    const pattern = document.getElementById('regexInput').value;
    
    fetch(`/api/search/regex?pattern=${encodeURIComponent(pattern)}`)
        .then(res => res.json())
        .then(data => displayResults(data.results))
        .catch(err => console.error(err));
}

function displayResults(results) {
    const container = document.getElementById('results');
    
    if (!results || results.length === 0) {
        container.innerHTML = '<p>No results found</p>';
        return;
    }
    
    container.innerHTML = results.map(r => `
        <div class="result-item">
            <div class="result-title">${r.book.title}</div>
            <div class="result-author">by ${r.book.author}</div>
            <div class="result-stats">
                ${r.occurrences} occurrences | 
                Relevance: ${r.relevance.toFixed(2)} |
                ${r.book.word_count} words
            </div>
        </div>
    `).join('');
}

function toggleAdvanced() {
    const adv = document.getElementById('advancedSearch');
    adv.style.display = adv.style.display === 'none' ? 'block' : 'none';
}

// Search on Enter key
document.getElementById('searchInput').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') search();
});