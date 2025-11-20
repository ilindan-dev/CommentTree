const API_URL = '/v1/comments';

// State
let state = {
    view: 'roots', // 'roots', 'search', 'thread'
    page: 1,
    searchQuery: '',
    currentThreadId: null
};

// DOM Elements
const container = document.getElementById('commentsContainer');
const pagination = document.getElementById('paginationControls');
const pageIndicator = document.getElementById('pageIndicator');
const rootForm = document.getElementById('rootFormContainer');

// --- API Calls ---

async function fetchComments(params = {}) {
    const url = new URL(window.location.origin + API_URL);
    Object.keys(params).forEach(key => url.searchParams.append(key, params[key]));

    try {
        const res = await fetch(url);
        if (!res.ok) throw new Error(await res.text());
        return await res.json();
    } catch (e) {
        alert('Ошибка загрузки: ' + e.message);
        return [];
    }
}

async function postComment(text, parentId = null) {
    try {
        const res = await fetch(API_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ text, parent_id: parentId ? parentId : "" })
        });
        if (!res.ok) throw new Error(await res.text());
        return await res.json();
    } catch (e) {
        alert('Ошибка отправки: ' + e.message);
    }
}

async function deleteComment(id) {
    if (!confirm('Вы уверены?')) return;
    try {
        const res = await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
        if (!res.ok) throw new Error(await res.text());
        loadCurrentView();
    } catch (e) {
        alert('Ошибка удаления: ' + e.message);
    }
}

// --- Rendering ---

function renderList(comments, isTree = false) {
    container.innerHTML = '';

    if (!comments || comments.length === 0) {
        container.innerHTML = '<p style="text-align:center; color: gray;">Здесь пока пусто.</p>';
        return;
    }

    if (isTree) {
        const tree = buildTreeFromFlat(comments);
        const ul = document.createElement('ul');
        ul.className = 'comment-list';
        tree.forEach(node => ul.appendChild(createCommentElement(node, true)));
        container.appendChild(ul);
    } else {
        comments.forEach(c => {
            const div = createCommentElement(c, false);
            container.appendChild(div);
        });
    }
}

function createCommentElement(comment, isTreeMode) {
    const isDeleted = comment.is_deleted;
    const li = document.createElement(isTreeMode ? 'li' : 'div');
    li.className = 'comment-item';

    const date = new Date(comment.created_at).toLocaleString('ru-RU');

    li.innerHTML = `
        <div class="comment-meta">
            <span>ID: ${comment.id.slice(0, 8)}...</span>
            <span>${date}</span>
        </div>
        <div class="comment-body ${isDeleted ? 'deleted' : ''}">${isDeleted ? '[Комментарий удален]' : escapeHtml(comment.text)}</div>
        
        <div class="comment-actions">
            ${!isTreeMode ? `<button class="btn-link" onclick="viewThread('${comment.id}')">💬 Смотреть ветку</button>` : ''}
            ${!isDeleted ? `<button class="btn-link" onclick="showReplyForm('${comment.id}')">↩️ Ответить</button>` : ''}
            ${!isDeleted ? `<button class="danger" onclick="deleteComment('${comment.id}')">Удалить</button>` : ''}
        </div>
        <div id="reply-form-${comment.id}"></div>
    `;

    if (comment.children && comment.children.length > 0) {
        const ul = document.createElement('ul');
        ul.className = 'nested-comments';
        comment.children.forEach(child => {
            ul.appendChild(createCommentElement(child, true));
        });
        li.appendChild(ul);
    }

    return li;
}

function buildTreeFromFlat(flatList) {
    const map = {};
    const roots = [];

    flatList.forEach(c => {
        c.children = [];
        map[c.id] = c;
    });

    flatList.forEach(c => {
        if (c.parent_id && map[c.parent_id]) {
            map[c.parent_id].children.push(c);
        } else {
            roots.push(c);
        }
    });

    return roots;
}

// --- Actions ---

window.loadCurrentView = () => {
    if (state.view === 'roots') loadRoots();
    else if (state.view === 'search') loadSearch();
    else if (state.view === 'thread') loadThread(state.currentThreadId);
};

async function loadRoots() {
    state.view = 'roots';
    state.currentThreadId = null;
    document.getElementById('listTitle').innerText = 'Последние комментарии';
    rootForm.classList.remove('hidden');
    pagination.classList.remove('hidden');
    pageIndicator.innerText = `Страница ${state.page}`;

    const comments = await fetchComments({ page: state.page, limit: 10, sort: 'desc' });
    renderList(comments, false);
}

async function loadSearch() {
    if (!state.searchQuery) return loadRoots();
    state.view = 'search';
    document.getElementById('listTitle').innerText = `Результаты поиска: "${state.searchQuery}"`;
    rootForm.classList.add('hidden');
    pagination.classList.add('hidden');

    const comments = await fetchComments({ query: state.searchQuery, limit: 50 });
    renderList(comments, false);
}

window.viewThread = async (rootId) => {
    state.view = 'thread';
    state.currentThreadId = rootId;
    document.getElementById('listTitle').innerText = 'Просмотр ветки';
    rootForm.classList.add('hidden');
    pagination.classList.add('hidden');

    const comments = await fetchComments({ parent: rootId, sort: 'asc' });
    renderList(comments, true);
};

window.showReplyForm = (parentId) => {
    const container = document.getElementById(`reply-form-${parentId}`);
    if (container.innerHTML !== '') {
        container.innerHTML = '';
        return;
    }

    container.innerHTML = `
        <div class="reply-form">
            <textarea id="reply-text-${parentId}" rows="2" style="width:95%" placeholder="Ваш ответ..."></textarea>
            <div style="margin-top:5px">
                <button onclick="submitReply('${parentId}')">Отправить</button>
                <button class="secondary" onclick="document.getElementById('reply-form-${parentId}').innerHTML=''">Отмена</button>
            </div>
        </div>
    `;
};

window.submitReply = async (parentId) => {
    const text = document.getElementById(`reply-text-${parentId}`).value;
    if (!text) return alert('Введите текст');

    await postComment(text, parentId);
    if (state.view === 'thread') loadCurrentView();
    else alert('Ответ отправлен! Перейдите в ветку, чтобы увидеть его.');

    document.getElementById(`reply-form-${parentId}`).innerHTML = '';
};

// --- Events ---

document.getElementById('btnPostRoot').onclick = async () => {
    const text = document.getElementById('rootCommentText').value;
    if (!text) return alert('Введите текст');
    await postComment(text);
    document.getElementById('rootCommentText').value = '';
    if (state.view === 'roots') loadRoots();
};

document.getElementById('btnSearch').onclick = () => {
    state.searchQuery = document.getElementById('searchInput').value;
    loadSearch();
};

document.getElementById('btnHome').onclick = () => {
    document.getElementById('searchInput').value = '';
    state.page = 1;
    loadRoots();
};

document.getElementById('btnNext').onclick = () => {
    state.page++;
    loadRoots();
};
document.getElementById('btnPrev').onclick = () => {
    if (state.page > 1) {
        state.page--;
        loadRoots();
    }
};

function escapeHtml(text) {
    if (!text) return "";
    return text
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

// Init
loadRoots();