// Переключение между формами
function showForm(formName) {
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const profileForm = document.getElementById('profile-form');
    const message = document.getElementById('message');

    // Скрываем сообщение
    message.className = 'message';

    // Скрываем все формы
    loginForm.classList.remove('active');
    registerForm.classList.remove('active');
    profileForm.style.display = 'none';

    if (formName === 'login') {
        loginForm.classList.add('active');
    } else if (formName === 'register') {
        registerForm.classList.add('active');
    } else if (formName === 'profile') {
        profileForm.style.display = 'block';
    }
}

// Показать профиль пользователя
function showProfile(user) {
    document.getElementById('profile-name').textContent = user.name || user.Name;
    document.getElementById('profile-email').textContent = user.email || user.Email;
    showForm('profile');
}

// Проверка токена при загрузке
async function checkAuth() {
    const token = localStorage.getItem('token');
    if (!token) return;

    try {
        const response = await fetch('/api/auth/me', {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const user = await response.json();
            showProfile(user);
        } else {
            localStorage.removeItem('token');
        }
    } catch (error) {
        console.error('Auth check error:', error);
    }
}

// Выход
function logout() {
    localStorage.removeItem('token');
    showForm('login');
    showMessage('Вы вышли из системы', 'success');
}

// Инициализация при загрузке
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();

    // Кнопка выхода
    const logoutBtn = document.getElementById('logoutBtn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', logout);
    }

    // Кнопка теста /me
    const testMeBtn = document.getElementById('testMeBtn');
    const meResult = document.getElementById('meResult');
    if (testMeBtn) {
        testMeBtn.addEventListener('click', async () => {
            const token = localStorage.getItem('token');
            try {
                const response = await fetch('/api/auth/me', {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                const data = await response.json();
                meResult.style.display = 'block';
                meResult.textContent = JSON.stringify(data, null, 2);
            } catch (error) {
                meResult.style.display = 'block';
                meResult.textContent = 'Error: ' + error.message;
            }
        });
    }
});

// Показать сообщение
function showMessage(text, type) {
    const message = document.getElementById('message');
    message.textContent = text;
    message.className = `message ${type}`;

    setTimeout(() => {
        message.className = 'message';
    }, 3000);
}

// Обработка формы входа
document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;

    try {
        const response = await fetch('/api/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        });

        const data = await response.json();

        if (response.ok) {
            showMessage('Успешный вход!', 'success');
            localStorage.setItem('token', data.token);
        } else {
            // Ошибка от сервера (401, 400 и т.д.)
            showMessage(data.message || 'Неверный email или пароль', 'error');
        }
    } catch (error) {
        // Действительная ошибка подключения
        showMessage('Ошибка: ' + (error.message || 'Не удалось подключиться к серверу'), 'error');
        console.error('Login error:', error);
    }
});

// Обработка формы регистрации
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;
    const confirm = document.getElementById('register-confirm').value;

    // Проверка совпадения паролей
    if (password !== confirm) {
        showMessage('Пароли не совпадают', 'error');
        return;
    }

    try {
        const response = await fetch('/api/auth/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, email, password }),
        });

        const data = await response.json();

        if (response.ok) {
            showMessage('Регистрация успешна!', 'success');
            // Сохраняем токен и переключаем на форму входа
            localStorage.setItem('token', data.token);
            setTimeout(() => showForm('login'), 1500);
        } else {
            showMessage(data.message || 'Ошибка регистрации', 'error');
        }
    } catch (error) {
        showMessage('Ошибка: ' + (error.message || 'Не удалось подключиться к серверу'), 'error');
        console.error('Register error:', error);
    }
});
