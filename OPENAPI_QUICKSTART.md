# 🚀 OpenAPI Quick Start Guide

Быстрое руководство по использованию новой OpenAPI документации в AmneziaWG Easy.

## ✨ Что было добавлено

### 1. Интерактивная документация API (Swagger UI)

Доступна по адресу: **http://localhost:51821/api/docs**

**Возможности:**
- 📖 Просмотр всех API endpoints с описаниями
- 🧪 Тестирование запросов прямо в браузере
- 🔐 Встроенная аутентификация
- 📝 Схемы запросов и ответов с примерами
- 📋 Генерация кода на разных языках

### 2. OpenAPI 3.0 спецификация

Доступна по адресу: **http://localhost:51821/api/openapi.yaml**

**Использование:**
- Импорт в Postman/Insomnia
- Генерация клиентских библиотек
- Автоматизированное тестирование API
- Валидация запросов/ответов

### 3. Готовые примеры кода

**Python клиент:** `docs/examples/openapi_client.py`
```bash
pip install requests
python docs/examples/openapi_client.py
```

**JavaScript клиент:** `docs/examples/openapi_client.js`
```bash
npm install axios
node docs/examples/openapi_client.js
```

### 4. Подробные руководства

- **[OpenAPI Guide](./docs/OPENAPI_GUIDE.md)** - Полное руководство по использованию
- **[Examples README](./docs/examples/README.md)** - Примеры использования клиентов
- **[API Reference](./docs/API_REFERENCE.md)** - Обновленная справка API

## 🎯 Быстрый старт

### Шаг 1: Запустите приложение

```bash
docker compose up -d
```

### Шаг 2: Откройте Swagger UI

Перейдите в браузере:
```
http://localhost:51821/api/docs
```

### Шаг 3: Аутентифицируйтесь

1. Нажмите кнопку **"Authorize"** в правом верхнем углу
2. Выберите **headerAuth**
3. Введите ваш пароль
4. Нажмите **"Authorize"**, затем **"Close"**

### Шаг 4: Попробуйте API

1. Выберите endpoint, например `GET /api/wireguard/client`
2. Нажмите **"Try it out"**
3. Нажмите **"Execute"**
4. Посмотрите результат в секции **"Response"**

## 📚 Основные сценарии использования

### Тестирование API в браузере

1. Откройте http://localhost:51821/api/docs
2. Авторизуйтесь через кнопку "Authorize"
3. Выберите нужный endpoint
4. Нажмите "Try it out" → "Execute"

### Генерация клиента на Python

```bash
npm install -g @openapitools/openapi-generator-cli

openapi-generator-cli generate \
  -i http://localhost:51821/api/openapi.yaml \
  -g python \
  -o ./amnezia-client-python
```

### Генерация клиента на TypeScript

```bash
openapi-generator-cli generate \
  -i http://localhost:51821/api/openapi.yaml \
  -g typescript-axios \
  -o ./amnezia-client-ts
```

### Импорт в Postman

1. Откройте Postman
2. Нажмите **Import** → **Link**
3. Вставьте: `http://localhost:51821/api/openapi.yaml`
4. Нажмите **Continue** → **Import**

### Использование готовых примеров

**Python:**
```bash
cd docs/examples
# Отредактируйте openapi_client.py, установите пароль
python openapi_client.py
```

**JavaScript:**
```bash
cd docs/examples
# Отредактируйте openapi_client.js, установите пароль
node openapi_client.js
```

## 🔧 Примеры кода

### Создание клиента с кастомными параметрами AmneziaWG

**Python:**
```python
from openapi_client import AmneziaWGClient

client = AmneziaWGClient("http://localhost:51821")
client.login("your_password")

client.create_client(
    name="stealth-client",
    jc="10",
    jmin="50",
    jmax="1000",
    s1="150",
    s2="150",
    h1="1234567891",
    h2="1234567892",
    h3="1234567893",
    h4="1234567894"
)
```

**JavaScript:**
```javascript
const AmneziaWGClient = require('./openapi_client');

const client = new AmneziaWGClient('http://localhost:51821');
await client.login('your_password');

await client.createClient({
    name: 'stealth-client',
    jc: '10',
    jmin: '50',
    jmax: '1000',
    s1: '150',
    s2: '150',
    h1: '1234567891',
    h2: '1234567892',
    h3: '1234567893',
    h4: '1234567894'
});
```

**cURL:**
```bash
curl -X POST http://localhost:51821/api/wireguard/client \
  -H "Authorization: your_password" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "stealth-client",
    "jc": "10",
    "jmin": "50",
    "jmax": "1000",
    "s1": "150",
    "s2": "150",
    "h1": "1234567891",
    "h2": "1234567892",
    "h3": "1234567893",
    "h4": "1234567894"
  }'
```

### Массовое создание клиентов

**Python:**
```python
for i in range(1, 11):
    client.create_client(
        name=f"client-{i}",
        expired_date="2025-12-31"
    )
    print(f"Created client-{i}")
```

**JavaScript:**
```javascript
for (let i = 1; i <= 10; i++) {
    await client.createClient({
        name: `client-${i}`,
        expiredDate: '2025-12-31'
    });
    console.log(`Created client-${i}`);
}
```

## 📖 Структура документации

```
📁 Документация
├── 🌐 Swagger UI          → http://localhost:51821/api/docs
├── 📄 OpenAPI Spec        → http://localhost:51821/api/openapi.yaml
├── 📖 API Reference       → docs/API_REFERENCE.md
├── 📚 OpenAPI Guide       → docs/OPENAPI_GUIDE.md
├── 🐍 Python Client       → docs/examples/openapi_client.py
├── 🟨 JavaScript Client   → docs/examples/openapi_client.js
└── 📋 Examples Guide      → docs/examples/README.md
```

## 🔍 Категории API

API разделено на следующие категории:

1. **System** - Информация о системе, настройки UI
2. **Authentication** - Вход/выход, управление сессиями
3. **Clients** - CRUD операции с WireGuard клиентами
4. **Configuration** - Скачивание конфигураций и QR-кодов
5. **One-Time Links** - Временные ссылки для скачивания
6. **Backup** - Резервное копирование и восстановление
7. **Metrics** - Метрики Prometheus

## 🎓 Дополнительные ресурсы

### Внутренние

- [Environment Variables](./docs/ENVIRONMENT_VARIABLES.md)
- [AmneziaWG Parameters](./docs/AMNEZIAWG_PARAMETERS.md)
- [Per-Client Parameters](./docs/PER_CLIENT_PARAMETERS.md)
- [Architecture](./docs/ARCHITECTURE.md)

### Внешние

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [OpenAPI Generator](https://openapi-generator.tech/)
- [Postman Learning Center](https://learning.postman.com/)

## ❓ Часто задаваемые вопросы

### Как получить доступ к Swagger UI?

Откройте http://localhost:51821/api/docs в браузере.

### Как авторизоваться в Swagger UI?

Нажмите кнопку "Authorize" и введите пароль в поле Authorization.

### Где находится OpenAPI спецификация?

Доступна по адресу http://localhost:51821/api/openapi.yaml или в файле `docs/openapi.yaml`.

### Как сгенерировать клиент на моем языке?

Используйте OpenAPI Generator:
```bash
openapi-generator-cli generate \
  -i http://localhost:51821/api/openapi.yaml \
  -g YOUR_LANGUAGE \
  -o ./output-directory
```

Список поддерживаемых языков: https://openapi-generator.tech/docs/generators/

### Можно ли использовать API без Web UI?

Да! Установите `NO_WEB_UI=true` в переменных окружения. API останется доступным, включая Swagger UI.

### Как импортировать спецификацию в Postman?

1. Откройте Postman
2. Import → Link
3. Вставьте: `http://localhost:51821/api/openapi.yaml`
4. Import

## 🐛 Устранение проблем

### Swagger UI не загружается

**Проверьте:**
- Приложение запущено: `docker compose ps`
- Порт доступен: `curl http://localhost:51821/api/release`
- Файлы на месте: `ls -la www/swagger.html docs/openapi.yaml`

### Ошибка аутентификации

**Решение:**
1. Проверьте пароль в переменных окружения
2. Используйте `wgpw` для генерации нового хеша:
   ```bash
   docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw mypassword
   ```

### OpenAPI спецификация не загружается

**Проверьте:**
- Файл существует: `cat docs/openapi.yaml`
- Путь правильный в routes.go
- Перезапустите приложение: `docker compose restart`

## 🎉 Готово!

Теперь у вас есть полная интерактивная документация API! 

**Попробуйте:**
1. Откройте http://localhost:51821/api/docs
2. Создайте клиента через Swagger UI
3. Скачайте конфигурацию
4. Используйте готовые примеры кода

---

**Нужна помощь?** 
- 📖 Читайте [OpenAPI Guide](./docs/OPENAPI_GUIDE.md)
- 💬 Откройте [issue на GitHub](https://github.com/little-secrets/amnezia-wg-easy/issues)
- ⭐ Поставьте звезду проекту!

