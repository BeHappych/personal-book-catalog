const API_BASE = '/api';

let currentFilters = {
    title: '',
    author: '',
    status: '',
    genre: ''
};

// Функция загрузки книг с фильтрами
async function loadBooks(filters = {}) {
    try {
        // Обновляем текущие фильтры
        if (filters) {
            currentFilters = { ...currentFilters, ...filters };
        }
        
        // Строим URL с параметрами запроса
        const params = new URLSearchParams();
        
        if (currentFilters.title) params.append('title', currentFilters.title);
        if (currentFilters.author) params.append('author', currentFilters.author);
        if (currentFilters.status) params.append('status', currentFilters.status);
        if (currentFilters.genre) params.append('genre', currentFilters.genre);
        
        const queryString = params.toString();
        const url = queryString ? `${API_BASE}/books?${queryString}` : `${API_BASE}/books`;
                
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }
        
        const books = await response.json();
        console.log(`Загружено ${books.length} книг`);
        
        updateStats(books);
        renderBooks(books);
        
    } catch (error) {
        console.error('Ошибка загрузки книг:', error);
        showError('Не удалось загрузить список книг: ' + error.message);
    }
}

// Функция синхронизации currentFilters с полями ввода
function syncFiltersFromInputs() {
    currentFilters = {
        title: document.getElementById('searchTitle').value.trim(),
        author: document.getElementById('searchAuthor').value.trim(),
        status: document.getElementById('filterStatus').value,
        genre: document.getElementById('filterGenre')?.value.trim() || ''
    };
    console.log('Фильтры синхронизированы:', currentFilters);
}

// Функция обновления отображения активных фильтров
function updateActiveFiltersDisplay() {
    const activeFiltersDiv = document.getElementById('activeFilters');
    const filterTagsDiv = document.getElementById('filterTags');
    
    if (!activeFiltersDiv || !filterTagsDiv) {
        console.warn('Элементы для отображения фильтров не найдены');
        return;
    }
    
    const filtersToCheck = { ...currentFilters };
    
    syncFiltersFromInputs();
    
    console.log('Обновление отображения. Текущие фильтры:', filtersToCheck);
    
    const activeFilters = [];
    
    ['title', 'author', 'genre', 'status'].forEach(key => {
        const value = filtersToCheck[key];
        if (value && value.toString().trim() !== '') {
            let label = '';
            
            switch(key) {
                case 'title':
                    label = `Название: ${value}`;
                    break;
                case 'author':
                    label = `Автор: ${value}`;
                    break;
                case 'genre':
                    label = `Жанр: ${value}`;
                    break;
                case 'status':
                    const statusText = value === 'available' ? 'Доступна' : 'Выдана';
                    label = `Статус: ${statusText}`;
                    break;
            }
            
            activeFilters.push({
                label: label,
                key: key,
                value: value
            });
        }
    });
    
    console.log('Найдено активных фильтров:', activeFilters.length);
    
    if (activeFilters.length > 0) {
        activeFiltersDiv.style.display = 'block';
        
        filterTagsDiv.innerHTML = activeFilters.map(filter => `
            <span class="badge bg-info text-dark d-inline-flex align-items-center me-1 mb-1">
                ${filter.label}
                <button type="button" 
                        class="btn-close btn-close-white ms-1" 
                        style="font-size: 0.6rem;"
                        onclick="removeFilter('${filter.key}')"
                        aria-label="Удалить фильр ${filter.label}">
                </button>
            </span>
        `).join('');
        
    } else {
        activeFiltersDiv.style.display = 'none';
        filterTagsDiv.innerHTML = '';
    }
}

// Функция удаления одного фильтра
function removeFilter(filterKey) {
    console.log('Удаление фильтра:', filterKey);
    
    let fieldToClear = '';
    
    switch(filterKey) {
        case 'title':
            fieldToClear = 'searchTitle';
            currentFilters.title = '';
            break;
        case 'author':
            fieldToClear = 'searchAuthor';
            currentFilters.author = '';
            break;
        case 'genre':
            fieldToClear = 'filterGenre';
            currentFilters.genre = '';
            break;
        case 'status':
            fieldToClear = 'filterStatus';
            currentFilters.status = '';
            break;
        default:
            console.warn('Неизвестный фильтр:', filterKey);
            return;
    }
    
    const field = document.getElementById(fieldToClear);
    if (field) {
        field.value = ''; 
    }
    
    console.log('Текущие фильтры после удаления:', currentFilters);
    
    updateActiveFiltersDisplay();
    
    setTimeout(() => {
        loadBooks();
    }, 100);
}

// Функция применения фильтров
function applyFilters() {
    syncFiltersFromInputs();
    
    console.log('Применяем фильтры:', currentFilters);
    
    loadBooks();
    
    updateActiveFiltersDisplay();
}

// Функция сброса фильтров
function resetFilters() {
    console.log('Сброс всех фильтров');
    
    document.getElementById('searchTitle').value = '';
    document.getElementById('searchAuthor').value = '';
    document.getElementById('filterStatus').value = '';
    
    const filterGenre = document.getElementById('filterGenre');
    if (filterGenre) {
        filterGenre.value = '';
    }
    
    currentFilters = {
        title: '',
        author: '',
        status: '',
        genre: ''
    };
    
    console.log('Фильтры после сброса:', currentFilters);
    
    updateActiveFiltersDisplay();
    
    loadBooks();
}

let searchTimeout;
// Функция автоматического поиска с задержкой
function setupSearchInputs() {
    const searchInputs = ['searchTitle', 'searchAuthor'];
    
    searchInputs.forEach(inputId => {
        const input = document.getElementById(inputId);
        if (input) {
            input.addEventListener('input', function() {
                clearTimeout(searchTimeout);
                searchTimeout = setTimeout(() => {
                    applyFilters();
                }, 500); 
            });
        }
    });
    
    const selectInputs = ['filterStatus', 'filterGenre'];
    selectInputs.forEach(selectId => {
        const select = document.getElementById(selectId);
        if (select) {
            select.addEventListener('change', applyFilters);
        }
    });
}

// Функция отображения книг
function renderBooks(books) {
    const container = document.getElementById('booksList');
    
    if (!books || books.length === 0) {
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-info">
                    <i class="bi bi-info-circle me-2"></i>
                    В библиотеке пока нет книг. Добавьте первую!
                </div>
            </div>
        `;
        return;
    }
    
    container.innerHTML = books.map(book => `
        <div class="col-md-6 col-lg-4">
            <div class="card book-card h-100">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-2">
                        <h5 class="card-title mb-0">${escapeHtml(book.title)}</h5>
                        <span class="badge ${book.status === 'available' ? 'bg-success' : 'bg-warning'}">
                            ${book.status === 'available' ? 'Доступна' : 'Выдана'}
                        </span>
                    </div>
                    <p class="card-text text-muted">
                        <i class="bi bi-person me-1"></i>${escapeHtml(book.author)}
                    </p>
                    ${book.genre ? `<p class="card-text"><small class="text-muted">${escapeHtml(book.genre)}</small></p>` : ''}
                    
                    <div class="mt-3">
                        <p class="card-text mb-1">
                            <i class="bi bi-geo-alt me-1"></i>
                            Местоположение: ${escapeHtml(book.room)}, шк. ${book.cabinet}, полка ${book.shelf}
                        </p>
                        ${book.lent_to ? `
                            <p class="card-text mb-1">
                                <i class="bi bi-person-check me-1"></i>
                                Выдана: ${escapeHtml(book.lent_to)}
                            </p>
                        ` : ''}
                        ${book.status === 'lent' && book.lent_date ? `
                            <p class="card-text">
                                <i class="bi bi-calendar me-1"></i>
                                Дата выдачи: ${new Date(book.lent_date).toLocaleDateString('ru-RU', {timeZone: 'UTC'})}
                            </p>
                        ` : ''}
                    </div>
                    
                    <div class="mt-3 action-buttons">
                        <div class="btn-group btn-group-sm w-100">
                            <button class="btn btn-outline-primary" onclick="editBook(${book.id})">
                                <i class="bi bi-pencil"></i> Редактировать
                            </button>
                            <button class="btn btn-outline-danger" onclick="deleteBook(${book.id})">
                                <i class="bi bi-trash"></i> Удалить
                            </button>
                            ${book.status === 'available' ? `
                                <button class="btn btn-outline-warning" onclick="lendBook(${book.id})">
                                    <i class="bi bi-arrow-right-circle"></i> Выдать
                                </button>
                            ` : `
                                <button class="btn btn-outline-success" onclick="returnBook(${book.id})">
                                    <i class="bi bi-arrow-left-circle"></i> Вернуть
                                </button>
                            `}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `).join('');
}

// Функция добавления книги
async function addBook() {
    const form = document.getElementById('addBookForm');
    const formData = new FormData(form);
    
    const bookData = {
        title: formData.get('title'),
        author: formData.get('author'),
        genre: formData.get('genre') || '',
        description: formData.get('description') || '', 
        room: formData.get('room'),
        cabinet: parseInt(formData.get('cabinet')),
        shelf: parseInt(formData.get('shelf')),
        row: 1,
        status: 'available'
    };
    
    try {
        const response = await fetch(`${API_BASE}/books`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(bookData)
        });
        
        if (response.ok) {
            const modal = bootstrap.Modal.getInstance(document.getElementById('addBookModal'));
            modal.hide();
            
            form.reset();
            
            loadBooks();
            
            showSuccess('Книга успешно добавлена!');
        } else {
            const error = await response.json();
            showError(error.error || 'Ошибка при добавлении книги');
        }
    } catch (error) {
        console.error('Ошибка:', error);
        showError('Не удалось добавить книгу');
    }
}

// Функция удаления книги
async function deleteBook(id) {
    if (!confirm('Вы уверены, что хотите удалить эту книгу?')) {
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/books/${id}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            loadBooks();
            showSuccess('Книга успешно удалена');
        } else {
            const error = await response.json();
            showError(error.error || 'Ошибка при удалении книги');
        }
    } catch (error) {
        console.error('Ошибка:', error);
        showError('Не удалось удалить книгу');
    }
}

// Функция обновления статистики
function updateStats(books) {
    if (!Array.isArray(books)) {
        console.error('books не является массивом:', books);
        books = [];
    }
    
    const total = books.length;
    const available = books.filter(b => b.status === 'available').length;
    const lent = total - available;
    
    document.getElementById('totalBooks').textContent = total;
    document.getElementById('availableBooks').textContent = available;
    document.getElementById('lentBooks').textContent = lent;
    
    const hasActiveFilters = currentFilters.title || currentFilters.author || 
                            currentFilters.genre || currentFilters.status;
    
    if (hasActiveFilters) {
        document.getElementById('totalBooks').parentElement.innerHTML = 
            `Найдено книг: <span id="totalBooks" class="badge bg-primary">${total}</span>`;
    } else {
        document.getElementById('totalBooks').parentElement.innerHTML = 
            `Всего книг: <span id="totalBooks" class="badge bg-primary">${total}</span>`;
    }
}

// Функция выдачи книги
async function lendBook(id) {
    const person = prompt('Кому выдать книгу?', 'Иванов Иван');
    if (!person) return;
    
    try {
        const response = await fetch(`${API_BASE}/books/${id}/lend`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ lent_to: person })
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка при выдаче книги');
        }
        
        loadBooks();  
        showSuccess('Книга успешно выдана!');
    } catch (error) {
        showError(error.message);
    }
}

// Функция возвращения книги
async function returnBook(id) {
    if (!confirm('Подтвердите возврат книги')) {
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/books/${id}/return`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Ошибка при возврате книги');
        }
        
        loadBooks();  
        showSuccess('Книга успешно возвращена!');
    } catch (error) {
        showError(error.message);
    }
}

// Функция редактирования книги
async function editBook(bookId) {
    console.log('Редактирование книги ID:', bookId);
    
    try {
        const response = await fetch(`${API_BASE}/books/${bookId}`);
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Не удалось загрузить данные книги');
        }
        
        const book = await response.json();
        console.log('Данные книги:', book);
        
        document.getElementById('editBookId').value = book.id;
        document.getElementById('editTitle').value = book.title;
        document.getElementById('editAuthor').value = book.author;
        document.getElementById('editGenre').value = book.genre || '';
        document.getElementById('editDescription').value = book.description || '';
        document.getElementById('editRoom').value = book.room;
        document.getElementById('editCabinet').value = book.cabinet;
        document.getElementById('editShelf').value = book.shelf;
        document.getElementById('editStatus').value = book.status || 'available';
        document.getElementById('editLentTo').value = book.lent_to || '';
        
        toggleLentInfoSection(book.status);
        
        const editModal = new bootstrap.Modal(document.getElementById('editBookModal'));
        editModal.show();
        
    } catch (error) {
        console.error('Ошибка при загрузке книги:', error);
        showError('Не удалось загрузить данные книги: ' + error.message);
    }
}

// Функция переключения видимости поля "Выдана кому"
function toggleLentInfoSection(status) {
    const lentInfoSection = document.getElementById('lentInfoSection');
    if (status === 'lent') {
        lentInfoSection.style.display = 'block';
    } else {
        lentInfoSection.style.display = 'none';
    }
}

// Функция слушатель изменения статуса в форме редактирования
document.getElementById('editStatus').addEventListener('change', function() {
    toggleLentInfoSection(this.value);
});

// Функция сохранить изменения книги
async function saveBookChanges() {
    const bookId = document.getElementById('editBookId').value;
    
    if (!bookId) {
        showError('ID книги не найден');
        return;
    }
    
    const bookData = {
        title: document.getElementById('editTitle').value.trim(),
        author: document.getElementById('editAuthor').value.trim(),
        genre: document.getElementById('editGenre').value.trim() || undefined,
        description: document.getElementById('editDescription').value.trim() || undefined,
        room: document.getElementById('editRoom').value.trim(),
        cabinet: parseInt(document.getElementById('editCabinet').value) || 1,
        shelf: parseInt(document.getElementById('editShelf').value) || 1,
        status: document.getElementById('editStatus').value
    };
    
    if (bookData.status === 'lent') {
        bookData.lent_to = document.getElementById('editLentTo').value.trim() || '';
    } else {
        bookData.lent_to = '';
    }
    
    if (!bookData.title || !bookData.author || !bookData.room) {
        showError('Заполните обязательные поля: название, автор, комната');
        return;
    }
    
    console.log('Сохранение книги ID:', bookId, 'Данные:', bookData);
    
    try {
        const response = await fetch(`${API_BASE}/books/${bookId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(bookData)
        });
        
        const responseText = await response.text();
        console.log('Ответ сервера:', responseText);
        
        if (!response.ok) {
            let errorMsg = 'Ошибка при обновлении книги';
            try {
                const errorData = JSON.parse(responseText);
                errorMsg = errorData.error || errorMsg;
            } catch {
                errorMsg = `${errorMsg}: ${responseText}`;
            }
            throw new Error(errorMsg);
        }
        
        const editModal = bootstrap.Modal.getInstance(document.getElementById('editBookModal'));
        editModal.hide();
        
        loadBooks();
        
        showSuccess('Книга успешно обновлена!');
        
    } catch (error) {
        console.error('Ошибка при сохранении:', error);
        showError('Не удалось обновить книгу: ' + error.message);
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function showError(message) {
    alert(`Ошибка: ${message}`);
}

function showSuccess(message) {
    alert(`✅ ${message}`);
}

window.loadBooks = loadBooks;
window.addBook = addBook;
window.deleteBook = deleteBook;
window.editBook = editBook;
window.lendBook = lendBook;
window.returnBook = returnBook;

document.addEventListener('DOMContentLoaded', function() {
    console.log('Инициализация приложения...');
    
    setupSearchInputs();
    
    loadBooks();
    
    updateActiveFiltersDisplay();
    
    const searchInputs = ['searchTitle', 'searchAuthor', 'filterGenre'];
    searchInputs.forEach(inputId => {
        const input = document.getElementById(inputId);
        if (input) {
            input.addEventListener('keypress', function(e) {
                if (e.key === 'Enter') {
                    applyFilters();
                }
            });
        }
    });
});