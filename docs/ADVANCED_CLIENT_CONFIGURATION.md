# Advanced Client Configuration

Руководство по расширенной настройке WireGuard клиентов.

## 📋 Содержание

- [Обзор](#обзор)
- [Создание клиента с полной настройкой](#создание-клиента-с-полной-настройкой)
- [Параметры](#параметры)
- [Обновление параметров](#обновление-параметров)
- [Примеры использования](#примеры-использования)

## 🎯 Обзор

Теперь вы можете полностью настроить каждый параметр WireGuard клиента:

### Сетевая конфигурация
- ✅ **IPv4 адрес** - custom или автоматический
- ✅ **IPv6 адрес** - опционально
- ✅ **AllowedIPs** - какой трафик маршрутизировать через VPN
- ✅ **DNS серверы** - custom DNS для клиента
- ✅ **MTU** - размер пакетов
- ✅ **Persistent Keepalive** - интервал keepalive

### Ключи безопасности
- ✅ **Private Key** - свой приватный ключ или автогенерация
- ✅ **Preshared Key** - дополнительный ключ для post-quantum защиты

### AmneziaWG параметры
- ✅ **Jc, Jmin, Jmax** - параметры junk packets
- ✅ **S1, S2** - размеры junk
- ✅ **H1, H2, H3, H4** - magic headers

## 📝 Создание клиента с полной настройкой

### Базовое создание (автоматические параметры)

```bash
POST /api/wireguard/client
Content-Type: application/json

{
  "name": "simple-client"
}
```

**Результат:**
- IPv4: автоматически (10.8.0.2, 10.8.0.3, ...)
- Ключи: автогенерированы
- Остальные параметры: из серверных настроек

### Полная кастомизация

```bash
POST /api/wireguard/client
Content-Type: application/json

{
  "name": "advanced-client",
  "expiredDate": "2025-12-31",
  
  // Сетевая конфигурация
  "address": "10.8.0.100",
  "address6": "fd42:42:42::100",
  "allowedIPs": "0.0.0.0/0, ::/0",
  "dns": "1.1.1.1, 1.0.0.1",
  "mtu": "1420",
  "persistentKeepalive": "25",
  
  // Собственные ключи (опционально)
  "privateKey": "YourPrivateKeyHere",
  "preSharedKey": "YourPresharedKeyHere",
  
  // AmneziaWG параметры
  "jc": "10",
  "jmin": "30",
  "jmax": "1200",
  "s1": "150",
  "s2": "150",
  "h1": "1234567891",
  "h2": "1234567892",
  "h3": "1234567893",
  "h4": "1234567894"
}
```

## 🔧 Параметры

### Network Configuration

#### `address` (string, optional)
IPv4 адрес клиента в VPN сети.

**По умолчанию:** Автоматически (10.8.0.2, 10.8.0.3, ...)  
**Пример:** `"10.8.0.100"`

```json
{
  "name": "client-with-static-ip",
  "address": "10.8.0.50"
}
```

#### `address6` (string, optional)
IPv6 адрес клиента (dual-stack).

**По умолчанию:** Отсутствует  
**Пример:** `"fd42:42:42::100"`

```json
{
  "name": "ipv6-client",
  "address": "10.8.0.100",
  "address6": "fd42:42:42::100"
}
```

**Конфигурация клиента будет содержать:**
```ini
[Interface]
Address = 10.8.0.100/24, fd42:42:42::100/64
```

#### `allowedIPs` (string, optional)
Какой трафик маршрутизировать через VPN.

**По умолчанию:** `0.0.0.0/0, ::/0` (весь трафик)  
**Примеры:**

**Только определенные сети:**
```json
{
  "name": "split-tunnel-client",
  "allowedIPs": "10.0.0.0/8, 192.168.0.0/16"
}
```

**Только IPv4:**
```json
{
  "name": "ipv4-only",
  "allowedIPs": "0.0.0.0/0"
}
```

#### `dns` (string, optional)
DNS серверы для клиента.

**По умолчанию:** Серверная настройка (обычно `1.1.1.1`)  
**Примеры:**

**Google DNS:**
```json
{
  "name": "google-dns-client",
  "dns": "8.8.8.8, 8.8.4.4"
}
```

**Cloudflare DNS:**
```json
{
  "name": "cloudflare-dns-client",
  "dns": "1.1.1.1, 1.0.0.1"
}
```

**Локальный DNS:**
```json
{
  "name": "local-dns-client",
  "dns": "192.168.1.1"
}
```

**Без DNS (пусто):**
```json
{
  "name": "no-dns-client",
  "dns": ""
}
```

#### `mtu` (string, optional)
Maximum Transmission Unit - размер пакетов.

**По умолчанию:** Серверная настройка или автоматически  
**Рекомендуется:** `1420` (стандарт WireGuard)  
**Диапазон:** обычно `1280` - `1500`

```json
{
  "name": "custom-mtu-client",
  "mtu": "1380"
}
```

**Когда изменять MTU:**
- ❌ Проблемы с соединением (timeouts, медленная скорость)
- ❌ PPPoE соединение (использовать 1412)
- ❌ Мобильные сети (использовать 1280)

#### `persistentKeepalive` (string, optional)
Интервал отправки keepalive пакетов (в секундах).

**По умолчанию:** `0` (отключено)  
**Рекомендуется:** `25` для NAT traversal

```json
{
  "name": "nat-client",
  "persistentKeepalive": "25"
}
```

**Когда использовать:**
- ✅ Клиент за NAT/файрволом
- ✅ Мобильные устройства
- ✅ Проблемы с разрывом соединения
- ❌ Стабильное соединение без NAT (экономит трафик)

### Security Keys

#### 🔑 Сценарии работы с ключами

API поддерживает 4 варианта создания клиентов:

| Что передать | publicKey | privateKey | Результат |
|--------------|-----------|------------|-----------|
| **Ничего** | Нет | Нет | ✅ Оба ключа автогенерируются |
| **Только privateKey** | Нет | ✅ Ваш | ✅ publicKey вычисляется автоматически |
| **Только publicKey** | ✅ Ваш | Нет | ✅ Road warrior (privateKey пустой) |
| **Оба ключа** | ✅ Ваш | ✅ Ваш | ✅ С валидацией соответствия |

#### `privateKey` (string, optional)
Приватный ключ клиента.

**По умолчанию:** Автогенерируется  
**Варианты использования:**

**Вариант 1: Только privateKey (publicKey вычислится автоматически)**
```json
{
  "name": "imported-client",
  "privateKey": "cGljYV9wcml2YXRlX2tleQ=="
}
```

**Вариант 2: Оба ключа (с валидацией соответствия)**
```json
{
  "name": "fully-imported-client",
  "privateKey": "cGljYV9wcml2YXRlX2tleQ==",
  "publicKey": "cGljYV9wdWJsaWNfa2V5"
}
```
⚠️ **API проверит, что publicKey соответствует privateKey!**

**Генерация ключей вручную:**
```bash
# Генерация приватного ключа
wg genkey > privatekey

# Получение публичного ключа из приватного
wg pubkey < privatekey > publickey
```

#### `publicKey` (string, optional)
Публичный ключ клиента.

**По умолчанию:** Автогенерируется (или вычисляется из privateKey)  
**Варианты использования:**

**Вариант 1: Только publicKey (road warrior setup)**
```json
{
  "name": "road-warrior-client",
  "publicKey": "cGljYV9wdWJsaWNfa2V5"
}
```
📝 **Результат:** privateKey будет пустым, клиент хранит ключ у себя

**Вариант 2: Оба ключа (с валидацией)**
```json
{
  "name": "validated-import",
  "privateKey": "cGljYV9wcml2YXRlX2tleQ==",
  "publicKey": "cGljYV9wdWJsaWNfa2V5"
}
```

#### `preSharedKey` (string, optional)
Preshared key для post-quantum защиты.

**По умолчанию:** Автогенерируется  
**Использование:** Дополнительная защита

```json
{
  "name": "quantum-safe-client",
  "preSharedKey": "cGljYV9wcmVzaGFyZWRfa2V5"
}
```

**Генерация PSK вручную:**
```bash
wg genpsk
```

### AmneziaWG Obfuscation

См. [AMNEZIAWG_PARAMETERS.md](./AMNEZIAWG_PARAMETERS.md) для подробной информации.

## 🔄 Обновление параметров

После создания клиента вы можете обновить любой параметр:

### Обновить IPv4 адрес

```bash
PUT /api/wireguard/client/:clientId/address
Content-Type: application/json

{
  "address": "10.8.0.200"
}
```

### Обновить IPv6 адрес

```bash
PUT /api/wireguard/client/:clientId/address6
Content-Type: application/json

{
  "address6": "fd42:42:42::200"
}
```

### Обновить AllowedIPs

```bash
PUT /api/wireguard/client/:clientId/allowedIPs
Content-Type: application/json

{
  "allowedIPs": "10.0.0.0/8, 192.168.0.0/16"
}
```

### Обновить DNS

```bash
PUT /api/wireguard/client/:clientId/dns
Content-Type: application/json

{
  "dns": "8.8.8.8, 8.8.4.4"
}
```

**Сбросить DNS (использовать серверный):**
```json
{
  "dns": ""
}
```

### Обновить MTU

```bash
PUT /api/wireguard/client/:clientId/mtu
Content-Type: application/json

{
  "mtu": "1380"
}
```

### Обновить Persistent Keepalive

```bash
PUT /api/wireguard/client/:clientId/keepalive
Content-Type: application/json

{
  "persistentKeepalive": "30"
}
```

## 💡 Примеры использования

### Пример 1: Мобильный клиент (за NAT)

```json
{
  "name": "mobile-client",
  "dns": "1.1.1.1, 1.0.0.1",
  "mtu": "1280",
  "persistentKeepalive": "25",
  "jc": "7",
  "s1": "100",
  "s2": "100"
}
```

**Почему:**
- MTU 1280 - для мобильных сетей
- Keepalive 25 - для NAT traversal
- AmneziaWG - обход DPI

### Пример 2: Split-tunnel клиент (только офисная сеть)

```json
{
  "name": "office-client",
  "allowedIPs": "10.0.0.0/8, 192.168.0.0/16",
  "dns": "192.168.1.1"
}
```

**Почему:**
- AllowedIPs - только внутренние сети
- DNS - локальный DNS офиса
- Остальной трафик идет напрямую

### Пример 3: Максимальная безопасность

```json
{
  "name": "secure-client",
  "dns": "1.1.1.1, 1.0.0.1",
  "preSharedKey": "custom-psk-here",
  "jc": "10",
  "jmin": "50",
  "jmax": "1500",
  "s1": "150",
  "s2": "150"
}
```

**Почему:**
- PSK - post-quantum защита
- AmneziaWG с максимальными параметрами
- Cloudflare DNS с шифрованием

### Пример 4: Dual-stack (IPv4 + IPv6)

```json
{
  "name": "dual-stack-client",
  "address": "10.8.0.100",
  "address6": "fd42:42:42::100",
  "allowedIPs": "0.0.0.0/0, ::/0",
  "dns": "2606:4700:4700::1111, 1.1.1.1"
}
```

**Почему:**
- Поддержка IPv6
- DNS работает по IPv6 и IPv4
- Весь трафик через VPN

### Пример 5: PPPoE соединение

```json
{
  "name": "pppoe-client",
  "mtu": "1412",
  "persistentKeepalive": "25"
}
```

**Почему:**
- MTU 1412 - для PPPoE (1492 - 80 overhead)
- Keepalive для стабильности

### Пример 6: Генерация ключей локально

```bash
# Генерируете ключи на своей машине
PRIVATE_KEY=$(wg genkey)
PUBLIC_KEY=$(echo "$PRIVATE_KEY" | wg pubkey)
PSK=$(wg genpsk)

# Создаете клиента через API
curl -X POST http://localhost:51821/api/wireguard/client \
  -H "Content-Type: application/json" \
  -H "Authorization: your_password" \
  -d '{
    "name": "secure-generated",
    "privateKey": "'$PRIVATE_KEY'",
    "publicKey": "'$PUBLIC_KEY'",
    "preSharedKey": "'$PSK'"
  }'
```

**Почему:**
- Ключи генерируются локально (не на сервере)
- Сервер проверяет соответствие privateKey и publicKey
- Полный контроль над процессом генерации

### Пример 7: Импорт существующего клиента

```json
{
  "name": "imported-from-other-server",
  "address": "10.8.0.150",
  "privateKey": "YourExistingPrivateKey==",
  "publicKey": "YourExistingPublicKey==",
  "preSharedKey": "YourExistingPSK==",
  "dns": "8.8.8.8"
}
```

**Почему:**
- Перенос с другого сервера
- Сохранение тех же ключей
- Валидация соответствия ключей
- Переопределение настроек

## 🔍 Просмотр конфигурации клиента

### Через API

```bash
GET /api/wireguard/client/:clientId
```

**Ответ включает все параметры:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "advanced-client",
  "address": "10.8.0.100",
  "address6": "fd42:42:42::100",
  "allowedIPs": "0.0.0.0/0, ::/0",
  "dns": "1.1.1.1, 1.0.0.1",
  "mtu": "1420",
  "persistentKeepalive": "25",
  ...
}
```

### Скачать конфигурацию

```bash
GET /api/wireguard/client/:clientId/configuration
```

**Результат:**
```ini
[Interface]
PrivateKey = <private-key>
Address = 10.8.0.100/24, fd42:42:42::100/64
DNS = 1.1.1.1, 1.0.0.1
MTU = 1420
Jc = 10
Jmin = 30
Jmax = 1200
S1 = 150
S2 = 150
H1 = 1234567891
H2 = 1234567892
H3 = 1234567893
H4 = 1234567894

[Peer]
PublicKey = <server-public-key>
PresharedKey = <preshared-key>
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 25
Endpoint = vpn.example.com:51820
```

## 🚨 Важные замечания

### Валидация ключей

⚠️ **При передаче обоих ключей (privateKey + publicKey):**

API автоматически проверяет, что publicKey соответствует privateKey!

```json
{
  "name": "test",
  "privateKey": "AAAA...",
  "publicKey": "BBBB..."  // Должен соответствовать!
}
```

❌ **Если ключи не соответствуют:**
```
Error: "provided publicKey does not match privateKey"
```

✅ **Проверка соответствия:**
```bash
# Вычислить publicKey из privateKey
echo "YOUR_PRIVATE_KEY" | wg pubkey

# Сравнить с вашим publicKey
```

### Конфликты IP адресов

❌ **Не используйте дублирующиеся IP:**
```json
// Если 10.8.0.100 уже занят:
{
  "name": "client2",
  "address": "10.8.0.100"  // ❌ ОШИБКА!
}
```

✅ **Проверяйте список клиентов перед созданием**

### Изменение ключей

⚠️ **Нельзя изменить ключи после создания!**

Если нужно изменить ключи:
1. Удалите клиента
2. Создайте нового с новыми ключами

### AllowedIPs и маршрутизация

**Full tunnel (весь трафик):**
```json
{
  "allowedIPs": "0.0.0.0/0, ::/0"
}
```

**Split tunnel (только VPN сеть):**
```json
{
  "allowedIPs": "10.8.0.0/24"
}
```

**Split tunnel (несколько сетей):**
```json
{
  "allowedIPs": "10.0.0.0/8, 192.168.0.0/16, 172.16.0.0/12"
}
```

### MTU проблемы

**Симптомы неправильного MTU:**
- ⚠️ Некоторые сайты не открываются
- ⚠️ SSH соединения зависают
- ⚠️ Медленная скорость передачи

**Решение:**
1. Попробуйте `1420` (стандарт)
2. Если не помогло: `1380`
3. Для PPPoE: `1412`
4. Для мобильных: `1280`

## 📚 Связанная документация

- [AmneziaWG Parameters](./AMNEZIAWG_PARAMETERS.md) - Параметры обфускации
- [Per-Client Parameters](./PER_CLIENT_PARAMETERS.md) - Индивидуальные настройки
- [API Reference](./API_REFERENCE.md) - Полный справочник API
- [Environment Variables](./ENVIRONMENT_VARIABLES.md) - Серверные настройки по умолчанию

## 🎓 Best Practices

### 1. Начните с настроек по умолчанию

```json
{
  "name": "test-client"
}
```

Если все работает - не меняйте!

### 2. Настраивайте только при необходимости

- ❌ Не меняйте MTU без причины
- ❌ Не включайте keepalive если не нужен
- ✅ Меняйте AllowedIPs для split-tunnel
- ✅ Меняйте DNS при необходимости

### 3. Документируйте изменения

Используйте понятные имена:
```json
{
  "name": "mobile-nat-keepalive25"
}
```

### 4. Тестируйте изменения

После изменения параметров:
1. Скачайте новую конфигурацию
2. Подключитесь
3. Проверьте работоспособность
4. При проблемах - верните к defaults

---

**Есть вопросы?** Откройте [issue на GitHub](https://github.com/little-secrets/amnezia-wg-easy/issues)!

