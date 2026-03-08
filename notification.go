package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ============================================================================
// ТИПЫ УВЕДОМЛЕНИЙ
// ============================================================================

// NotificationType представляет тип уведомления
type NotificationType string

const (
	NotificationDailyQuests    NotificationType = "daily_quests"    // Ежедневные квесты
	NotificationFinalBattle    NotificationType = "final_battle"    // Финальная битва
	NotificationUnfinished     NotificationType = "unfinished"      // Незавершённые квесты
	NotificationDayStreak      NotificationType = "day_streak"      // Серия дней
	NotificationLevelUp        NotificationType = "level_up"        // Повышение уровня
	NotificationBossDefeated   NotificationType = "boss_defeated"   // Победа над боссом
	NotificationWelcomeBack    NotificationType = "welcome_back"    // Возвращение в игру
)

// NotificationSettings представляет настройки уведомлений игрока
type NotificationSettings struct {
	ChatID             int64 `json:"chat_id"`
	Enabled            bool  `json:"enabled"`              // Все уведомления включены
	DailyQuestsEnabled bool  `json:"daily_quests_enabled"` // Ежедневные квесты
	FinalBattleEnabled bool  `json:"final_battle_enabled"` // Финальная битва
	UnfinishedEnabled  bool  `json:"unfinished_enabled"`   // Незавершённые квесты
	QuestsHour         int   `json:"quests_hour"`          // Время напоминания о квестах (час)
	BattleHour         int   `json:"battle_hour"`          // Время напоминания о битве (час)
}

// DefaultNotificationSettings возвращает настройки по умолчанию
func DefaultNotificationSettings(chatID int64) *NotificationSettings {
	return &NotificationSettings{
		ChatID:             chatID,
		Enabled:            true,
		DailyQuestsEnabled: true,
		FinalBattleEnabled: true,
		UnfinishedEnabled:  true,
		QuestsHour:         9,  // 9:00 утра
		BattleHour:         20, // 20:00 вечера
	}
}

// ============================================================================
// ГЛОБАЛЬНЫЕ ПЕРЕМЕННЫЕ
// ============================================================================

// NotificationScheduler — планировщик уведомлений
var NotificationScheduler *NotificationManager

// NotificationManager управляет уведомлениями
type NotificationManager struct {
	ticker     *time.Ticker
	done       chan bool
	settings   map[int64]*NotificationSettings // Настройки по chat_id
}

// InitNotifications инициализирует менеджер уведомлений
func InitNotifications() {
	NotificationScheduler = &NotificationManager{
		ticker:   time.NewTicker(time.Minute * 30), // Проверяем каждые 30 минут
		done:     make(chan bool),
		settings: make(map[int64]*NotificationSettings),
	}

	// Запускаем горутину для проверки уведомлений
	go NotificationScheduler.run()

	log.Println("✅ Менеджер уведомлений запущен")
}

// StopNotifications останавливает менеджер уведомлений
func StopNotifications() {
	if NotificationScheduler != nil {
		NotificationScheduler.ticker.Stop()
		NotificationScheduler.done <- true
		log.Println("🛑 Менеджер уведомлений остановлен")
	}
}

// run запускает цикл проверки уведомлений
func (nm *NotificationManager) run() {
	for {
		select {
		case <-nm.ticker.C:
			nm.checkNotifications()
		case <-nm.done:
			return
		}
	}
}

// checkNotifications проверяет и отправляет уведомления
func (nm *NotificationManager) checkNotifications() {
	currentHour := time.Now().Hour()
	currentMinute := time.Now().Minute()

	// Проверяем каждый час в пределах 30-минутного окна
	if currentMinute >= 0 && currentMinute < 30 {
		// Уведомление о квестах (утро)
		if currentHour == 9 {
			nm.sendDailyQuestsReminders()
		}

		// Уведомление о финальной битве (вечер)
		if currentHour == 20 {
			nm.sendFinalBattleReminders()
		}

		// Уведомление о незавершённых квестах (поздний вечер)
		if currentHour == 22 {
			nm.sendUnfinishedQuestsReminders()
		}
	}
}

// ============================================================================
// ОТПРАВКА УВЕДОМЛЕНИЙ
// ============================================================================

// SendDailyQuestsReminders отправляет напоминания о квестах всем игрокам
func (nm *NotificationManager) sendDailyQuestsReminders() {
	log.Println("📋 Отправка напоминаний о квестах...")

	query := `SELECT chat_id, name FROM players WHERE game_active = 1`
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("❌ Ошибка получения игроков: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var chatID int64
		var name string
		if err := rows.Scan(&chatID, &name); err != nil {
			continue
		}

		// Проверяем настройки уведомлений
		settings := nm.getSettings(chatID)
		if !settings.Enabled || !settings.DailyQuestsEnabled {
			continue
		}

		// Отправляем уведомление
		if sendDailyQuestsNotification(chatID, name) {
			count++
		}
	}

	log.Printf("✅ Отправлено %d напоминаний о квестах", count)
}

// SendFinalBattleReminders отправляет напоминания о финальной битве
func (nm *NotificationManager) sendFinalBattleReminders() {
	log.Println("⚔️  Отправка напоминаний о финальной битве...")

	query := `SELECT chat_id, name FROM players WHERE game_active = 1`
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("❌ Ошибка получения игроков: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var chatID int64
		var name string
		if err := rows.Scan(&chatID, &name); err != nil {
			continue
		}

		// Проверяем настройки уведомлений
		settings := nm.getSettings(chatID)
		if !settings.Enabled || !settings.FinalBattleEnabled {
			continue
		}

		// Отправляем уведомление
		if sendFinalBattleNotification(chatID, name) {
			count++
		}
	}

	log.Printf("✅ Отправлено %d напоминаний о финальной битве", count)
}

// SendUnfinishedQuestsReminders отправляет напоминания о незавершённых квестах
func (nm *NotificationManager) sendUnfinishedQuestsReminders() {
	log.Println("⏳ Отправка напоминаний о незавершённых квестах...")

	query := `SELECT chat_id, name FROM players WHERE game_active = 1`
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("❌ Ошибка получения игроков: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var chatID int64
		var name string
		if err := rows.Scan(&chatID, &name); err != nil {
			continue
		}

		// Проверяем настройки уведомлений
		settings := nm.getSettings(chatID)
		if !settings.Enabled || !settings.UnfinishedEnabled {
			continue
		}

		// Проверяем, есть ли незавершённые квесты
		hasUnfinished := hasUnfinishedQuests(chatID)
		if !hasUnfinished {
			continue
		}

		// Отправляем уведомление
		if sendUnfinishedQuestsNotification(chatID, name) {
			count++
		}
	}

	log.Printf("✅ Отправлено %d напоминаний о незавершённых квестах", count)
}

// ============================================================================
// ФУНКЦИИ ОТПРАВКИ
// ============================================================================

// SendDailyQuestsNotification отправляет уведомление о квестах
func sendDailyQuestsNotification(chatID int64, name string) bool {
	text := fmt.Sprintf(`🌅 <b>ДОБРОЕ УТРО, %s!</b>

📋 <b>ЕЖЕДНЕВНЫЕ КВЕСТЫ</b>

Новый день — новые возможности!
Выполни 5 ежедневных квестов и получи очки навыков!

🎯 <b>Сегодня:</b>
• 30 минут Go
• Борец с искушениями
• Практика кода
• Утренний ритуал
• Цифровой детокс

💡 <b>Совет:</b>
Начни с изучения Go — это даст опыт и очки навыков!

🎮 <b>Начать день:</b>
/play`,
		name)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎮 Начать игру", "start_game"),
		),
	)

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("❌ Ошибка отправки уведомления: %v", err)
		return false
	}

	return true
}

// SendFinalBattleNotification отправляет уведомление о финальной битве
func sendFinalBattleNotification(chatID int64, name string) bool {
	text := fmt.Sprintf(`🌙 <b>ВЕЧЕР, %s!</b>

⚔️  <b>ВРЕМЯ ФИНАЛЬНОЙ БИТВЫ!</b>

День подходит к концу, но впереди главное испытание!
Сразись с боссом-искушением и докажи свою силу воли!

💪 <b>Подготовка:</b>
• Восстанови фокус (отдохни)
• Проверь силу воли
• Настройся на победу!

🏆 <b>Награда за победу:</b>
• +200 опыта
• Фокус и воля восстановлены
• Достижение "Победитель искушений"

🎮 <b>В бой:</b>
/play`,
		name)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ В бой!", "start_game"),
		),
	)

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("❌ Ошибка отправки уведомления: %v", err)
		return false
	}

	return true
}

// SendUnfinishedQuestsNotification отправляет уведомление о незавершённых квестах
func sendUnfinishedQuestsNotification(chatID int64, name string) bool {
	player, err := LoadPlayer(chatID)
	if err != nil || player == nil {
		return false
	}

	// Считаем незавершённые квесты
	unfinished := 0
	for _, quest := range player.Quests.Quests {
		if !quest.Completed {
			unfinished++
		}
	}

	if unfinished == 0 {
		return false
	}

	text := fmt.Sprintf(`⏰ <b>%s, ВРЕМЯ ПОДЖАЛО!</b>

📋 <b>НЕЗАВЕРШЁННЫЕ КВЕСТЫ</b>

До конца дня осталось мало времени!
Выполни квесты, чтобы получить очки навыков!

⚠️  <b>Осталось: %d квестов</b>

💡 <b>Совет:</b>
Даже если не успеешь всё — сделай максимум!
Каждый выполненный квест — это очки навыков!

🎮 <b>Выполнить квесты:</b>
/play`,
		name, unfinished)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎮 Выполнить квесты", "start_game"),
		),
	)

	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("❌ Ошибка отправки уведомления: %v", err)
		return false
	}

	return true
}

// SendLevelUpNotification отправляет уведомление о повышении уровня
func SendLevelUpNotification(chatID int64, name string, level int) bool {
	text := fmt.Sprintf(`🎉 <b>ПОЗДРАВЛЯЕМ, %s!</b>

🆙 <b>НОВЫЙ УРОВЕНЬ: %d!</b>

Твой прогресс растёт!
Продолжай в том же духе!

🎁 <b>Награда:</b>
• Очки навыков: +%d
• Фокус восстановлен: 100%%
• Сила воли восстановлена: 100%%

💪 <b>Так держать!</b>
Каждый уровень приближает тебя к цели!`,
		name, level, 2+(level/5))

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("❌ Ошибка отправки уведомления: %v", err)
		return false
	}

	return true
}

// SendWelcomeBackNotification отправляет уведомление о возвращении
func SendWelcomeBackNotification(chatID int64, name string, daysSinceLastLogin int) bool {
	var welcomeText string
	if daysSinceLastLogin == 1 {
		welcomeText = "С возвращением! Мы скучали!"
	} else if daysSinceLastLogin < 7 {
		welcomeText = fmt.Sprintf("Ты отсутствовал %d дней. Возвращайся в игру!", daysSinceLastLogin)
	} else {
		welcomeText = fmt.Sprintf("Давно не виделись! %d дней без Go — это слишком много!", daysSinceLastLogin)
	}

	text := fmt.Sprintf(`👋 <b>%s!</b>

%s

🎮 <b>FOCUSGO ждёт тебя!</b>

📋 <b>Что тебя ждёт:</b>
• Новые ежедневные квесты
• Возможность улучшить навыки
• Борьба с искушениями
• Путь к Go-Мастеру!

💡 <b>Совет:</b>
Начни с 30 минут Go — это даст опыт и очки навыков!

🎮 <b>Продолжить:</b>
/play`,
		name, welcomeText)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎮 Продолжить", "start_game"),
		),
	)

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("❌ Ошибка отправки уведомления: %v", err)
		return false
	}

	return true
}

// ============================================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// ============================================================================

// HasUnfinishedQuests проверяет, есть ли незавершённые квесты
func hasUnfinishedQuests(chatID int64) bool {
	player, err := LoadPlayer(chatID)
	if err != nil || player == nil {
		return false
	}

	for _, quest := range player.Quests.Quests {
		if !quest.Completed {
			return true
		}
	}

	return false
}

// GetSettings получает настройки уведомлений игрока
func (nm *NotificationManager) getSettings(chatID int64) *NotificationSettings {
	if settings, exists := nm.settings[chatID]; exists {
		return settings
	}

	// Загружаем из БД
	settings, err := LoadNotificationSettings(chatID)
	if err != nil || settings == nil {
		// Создаём по умолчанию
		settings = DefaultNotificationSettings(chatID)
	}

	nm.settings[chatID] = settings
	return settings
}

// SaveSettings сохраняет настройки уведомлений
func (nm *NotificationManager) saveSettings(settings *NotificationSettings) error {
	nm.settings[settings.ChatID] = settings
	return SaveNotificationSettings(settings)
}

// ============================================================================
// БАЗА ДАННЫХ
// ============================================================================

// LoadNotificationSettings загружает настройки уведомлений из БД
func LoadNotificationSettings(chatID int64) (*NotificationSettings, error) {
	query := `
		SELECT enabled, daily_quests_enabled, final_battle_enabled, 
		       unfinished_enabled, quests_hour, battle_hour
		FROM notification_settings
		WHERE chat_id = ?
	`

	row := DB.QueryRow(query, chatID)

	settings := &NotificationSettings{
		ChatID: chatID,
	}

	err := row.Scan(
		&settings.Enabled,
		&settings.DailyQuestsEnabled,
		&settings.FinalBattleEnabled,
		&settings.UnfinishedEnabled,
		&settings.QuestsHour,
		&settings.BattleHour,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return settings, nil
}

// SaveNotificationSettings сохраняет настройки уведомлений в БД
func SaveNotificationSettings(settings *NotificationSettings) error {
	query := `
		INSERT OR REPLACE INTO notification_settings 
		(chat_id, enabled, daily_quests_enabled, final_battle_enabled, 
		 unfinished_enabled, quests_hour, battle_hour)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := DB.Exec(query,
		settings.ChatID,
		settings.Enabled,
		settings.DailyQuestsEnabled,
		settings.FinalBattleEnabled,
		settings.UnfinishedEnabled,
		settings.QuestsHour,
		settings.BattleHour,
	)

	return err
}

// CreateNotificationSettingsTable создаёт таблицу настроек уведомлений
func CreateNotificationSettingsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS notification_settings (
			chat_id INTEGER PRIMARY KEY,
			enabled INTEGER DEFAULT 1,
			daily_quests_enabled INTEGER DEFAULT 1,
			final_battle_enabled INTEGER DEFAULT 1,
			unfinished_enabled INTEGER DEFAULT 1,
			quests_hour INTEGER DEFAULT 9,
			battle_hour INTEGER DEFAULT 20,
			FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE
		)
	`

	_, err := DB.Exec(query)
	return err
}
