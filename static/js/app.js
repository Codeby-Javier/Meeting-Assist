const API_URL = '/api';

// Auth Helpers
function getToken() {
    return localStorage.getItem('token');
}

function setToken(token) {
    localStorage.setItem('token', token);
}

function removeToken() {
    localStorage.removeItem('token');
    window.location.href = '/login';
}

function checkAuth() {
    if (!getToken()) {
        window.location.href = '/login';
    }
}

// API Call Helper
async function apiCall(endpoint, method = 'GET', body = null, isFormData = false) {
    const headers = {};
    const token = getToken();
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }
    if (!isFormData) {
        headers['Content-Type'] = 'application/json';
    }

    const config = {
        method,
        headers,
    };

    if (body) {
        config.body = isFormData ? body : JSON.stringify(body);
    }

    try {
        const response = await fetch(`${API_URL}${endpoint}`, config);
        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || data.message || 'Terjadi kesalahan');
        }
        return data;
    } catch (error) {
        showToast(error.message, 'danger');
        throw error;
    }
}

// Toast Notification
function showToast(message, type = 'info') {
    const toastContainer = document.getElementById('toast-container');
    if (!toastContainer) return;

    const toast = document.createElement('div');
    toast.className = `toast align-items-center text-white bg-${type} border-0 show`;
    toast.role = 'alert';
    toast.innerHTML = `
        <div class="d-flex">
            <div class="toast-body">${message}</div>
            <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
        </div>
    `;

    toastContainer.appendChild(toast);
    setTimeout(() => toast.remove(), 3000);
}

// File Upload
function handleFileUpload(fileInputId, endpoint, type) {
    const fileInput = document.getElementById(fileInputId);
    if (!fileInput || !fileInput.files[0]) return;

    const formData = new FormData();
    formData.append(type, fileInput.files[0]);

    const progressBar = document.getElementById('progress-bar');
    if (progressBar) progressBar.style.width = '0%';

    // Show loading state
    if (progressBar) {
        progressBar.style.width = '50%';
        progressBar.classList.add('progress-bar-striped', 'progress-bar-animated');
    }

    apiCall(endpoint, 'POST', formData, true)
        .then(data => {
            if (progressBar) progressBar.style.width = '100%';
            // Redirect to history page with appropriate tab
            setTimeout(() => {
                if (type === 'audio') {
                    window.location.href = '/history#audio';
                } else if (type === 'document') {
                    window.location.href = '/history#document';
                }
            }, 500);
        })
        .catch(err => {
            console.error(err);
            if (progressBar) {
                progressBar.style.width = '0%';
                progressBar.classList.remove('progress-bar-striped', 'progress-bar-animated');
            }
        });
}

// Sidebar Toggle
document.addEventListener('DOMContentLoaded', () => {
    const toggleBtn = document.getElementById('sidebar-toggle');
    if (toggleBtn) {
        toggleBtn.addEventListener('click', () => {
            document.querySelector('.sidebar').classList.toggle('active');
            document.querySelector('.main-content').classList.toggle('active');
        });
    }
});
