# 🚀 DEPLOYMENT GUIDE — FocusGo

**Инструкция по деплою на Ubuntu 24.04**

---

## 📋 Требования

- Ubuntu 24.04 (твой сервер)
- Go 1.21+
- Telegram Bot Token (от @BotFather)

---

## 🔧 Шаг 1: Подготовка сервера

### 1.1. Установи Go (если не установлен)

```bash
# Проверь версию
go version

# Если нужно установить:
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 1.2. Создай пользователя для бота (опционально)

```bash
sudo useradd -m -s /bin/bash focusgo
sudo su - focusgo
```

---

## 📦 Шаг 2: Загрузка бота на сервер

### Вариант 1: Клонирование репозитория

```bash
cd /home/gofer/godev/projects
git clone https://github.com/Folombas/focusgo.git
cd focusgo
```

### Вариант 2: Загрузка бинарника

```bash
# Локально (в WSL)
cd /home/gofer/godev/projects/focusgo
go build -o focusgo ./cmd/focusgo

# Копирование на сервер
scp focusgo user@server:/opt/focusgo/
scp .env user@server:/opt/focusgo/
```

---

## ⚙️ Шаг 3: Настройка

### 3.1. Настрой .env

```bash
cd /opt/focusgo
nano .env
```

**Содержимое .env:**
```bash
TELEGRAM_BOT_TOKEN=1234567890:AABBccDDeeFFggHHiiJJkkLLmmNNooP
DATABASE_PATH=focusgo.db
DEBUG=false
```

### 3.2. Установи права доступа

```bash
chmod +x focusgo
chmod 600 .env
```

---

## 🎯 Шаг 4: Создание systemd service

### 4.1. Создай файл сервиса

```bash
sudo nano /etc/systemd/system/focusgo.service
```

**Содержимое:**
```ini
[Unit]
Description=FocusGo Telegram Bot
After=network.target

[Service]
Type=simple
User=focusgo
WorkingDirectory=/opt/focusgo
ExecStart=/opt/focusgo/focusgo
Restart=always
RestartSec=10
Environment="PATH=/usr/local/go/bin:/usr/bin:/bin"

# Логирование
StandardOutput=journal
StandardError=journal
SyslogIdentifier=focusgo

[Install]
WantedBy=multi-user.target
```

### 4.2. Включи и запусти сервис

```bash
# Перезагрузи systemd
sudo systemctl daemon-reload

# Включи автозагрузку
sudo systemctl enable focusgo

# Запусти бота
sudo systemctl start focusgo

# Проверь статус
sudo systemctl status focusgo
```

---

## 📊 Шаг 5: Мониторинг

### Просмотр логов

```bash
# Последние 50 строк
sudo journalctl -u focusgo -n 50

# В реальном времени
sudo journalctl -u focusgo -f

# Логи за сегодня
sudo journalctl -u focusgo --since today
```

### Управление сервисом

```bash
# Остановить
sudo systemctl stop focusgo

# Перезапустить
sudo systemctl restart focusgo

# Перезагрузить конфигурацию
sudo systemctl daemon-reload
```

---

## 🔐 Шаг 6: Безопасность

### 6.1. Настрой фаервол

```bash
# Раз разреши только исходящие соединения
sudo ufw allow out 443/tcp  # HTTPS для Telegram API
sudo ufw allow out 53/udp   # DNS
```

### 6.2. Резервное копирование БД

```bash
# Создай скрипт бэкапа
sudo nano /opt/focusgo/backup.sh
```

**Содержимое backup.sh:**
```bash
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
cp /opt/focusgo/focusgo.db /opt/focusgo/backups/focusgo_backup_$DATE.db
# Удаляем бэкапы старше 7 дней
find /opt/focusgo/backups -name "*.db" -mtime +7 -delete
```

**Добавь в crontab:**
```bash
crontab -e
```

**Добавь строку:**
```bash
0 2 * * * /opt/focusgo/backup.sh
```

---

## 🐛 Шаг 7: Troubleshooting

### Бот не запускается

```bash
# Проверь логи
sudo journalctl -u focusgo -n 100

# Проверь .env
cat /opt/focusgo/.env

# Проверь права
ls -la /opt/focusgo/
```

### Ошибка "token not found"

```bash
# Проверь переменную окружения
sudo systemctl show focusgo | grep Environment

# Перезагрузи сервис
sudo systemctl daemon-reload
sudo systemctl restart focusgo
```

### Бот падает через несколько минут

```bash
# Проверь использование памяти
sudo systemctl status focusgo

# Проверь логи на панику
sudo journalctl -u focusgo | grep -i panic
```

---

## 📈 Шаг 8: Обновление

```bash
# Останови бота
sudo systemctl stop focusgo

# Обнови код
cd /opt/focusgo
git pull

# Перекомпилируй
go build -o focusgo ./cmd/focusgo

# Запусти
sudo systemctl start focusgo

# Проверь статус
sudo systemctl status focusgo
```

---

## ✅ Чек-лист деплоя

- [ ] Go установлен
- [ ] Репозиторий склонирован / бинарник загружен
- [ ] .env настроен с правильным токеном
- [ ] systemd service создан
- [ ] Сервис запущен и работает
- [ ] Логи пишутся в journalctl
- [ ] Бэкапы настроены
- [ ] Фаервол настроен

---

## 🎉 Готово!

Бот работает! Проверь в Telegram:
```
/start — начать игру
/help — справка
/leaderboard — таблица лидеров
/achievements — достижения
```

**Удачи с деплоем!** 🚀
