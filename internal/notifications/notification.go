package notifications

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"focusgo/internal/database"
)

// NotificationType представляет тип уведомления
type NotificationType string

const (
	NotificationDailyQuests  NotificationType = "daily_quests"
	NotificationFinalBattle  NotificationType = "final_battle"
	NotificationUnfinished   NotificationType = "unfinished"
)

// NotificationSettings представляет настройки уведомлений игрока
type NotificationSettings struct {
	ChatID             int64
	Enabled            bool
	DailyQuestsEnabled bool
	FinalBattleEnabled bool
	UnfinishedEnabled  bool
	QuestsHour         int
	BattleHour         int
}

// DefaultNotificationSettings возвращает настройки по умолчанию
func DefaultNotificationSettings(chatID int64) *NotificationSettings {
	return &NotificationSettings{
		ChatID:             chatID,
		Enabled:            true,
		DailyQuestsEnabled: true,
		FinalBattleEnabled: true,
		UnfinishedEnabled:  true,
		QuestsHour:         9,
		BattleHour:         20,
	}
}

// NotificationManager управляет уведомлениями
type NotificationManager struct {
	ticker   *time.Ticker
	done     chan bool
	settings map[int64]*NotificationSettings
	bot      *tgbotapi.BotAPI
}

// NotificationScheduler — глобальный планировщик
var NotificationScheduler *NotificationManager

// InitNotifications инициализирует менеджер уведомлений
func InitNotifications() {
	NotificationScheduler = &NotificationManager{
		ticker:   time.NewTicker(time.Minute * 30),
		done:     make(chan bool),
		settings: make(map[int64]*NotificationSettings),
	}
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

// SetBot устанавливает бота для отправки уведомлений
func SetBot(bot *tgbotapi.BotAPI) {
	if NotificationScheduler != nil {
		NotificationScheduler.bot = bot
	}
}

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

func (nm *NotificationManager) checkNotifications() {
	if nm.bot == nil {
		return
	}

	currentHour := time.Now().Hour()
	currentMinute := time.Now().Minute()

	if currentMinute >= 0 && currentMinute < 30 {
		if currentHour == 9 {
			nm.sendDailyQuestsReminders()
		}
		if currentHour == 20 {
			nm.sendFinalBattleReminders()
		}
		if currentHour == 22 {
			nm.sendUnfinishedQuestsReminders()
		}
	}
}

func (nm *NotificationManager) sendDailyQuestsReminders() {
	log.Println("📋 Отправка напоминаний о квестах...")
	nm.sendReminders(func(chatID int64, name string) bool {
		return sendDailyQuestsNotification(nm.bot, chatID, name)
	})
}

func (nm *NotificationManager) sendFinalBattleReminders() {
	log.Println("⚔️  Отправка напоминаний о финальной битве...")
	nm.sendReminders(func(chatID int64, name string) bool {
		return sendFinalBattleNotification(nm.bot, chatID, name)
	})
}

func (nm *NotificationManager) sendUnfinishedQuestsReminders() {
	log.Println("⏰ Отправка напоминаний о незавершённых квестах...")
	nm.sendReminders(func(chatID int64, name string) bool {
		return sendUnfinishedQuestsNotification(nm.bot, chatID, name)
	})
}

func (nm *NotificationManager) sendReminders(sendFunc func(int64, string) bool) {
	query := `SELECT chat_id, name FROM players WHERE game_active = 1`
	rows, err := database.DB.Query(query)
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

		settings := nm.getSettings(chatID)
		if !settings.Enabled {
			continue
		}

		if sendFunc(chatID, name) {
			count++
		}
	}
	log.Printf("✅ Отправлено %d уведомлений", count)
}

func sendDailyQuestsNotification(bot *tgbotapi.BotAPI, chatID int64, name string) bool {
	text := fmt.Sprintf(`🌅 <b>ДОБРОЕ УТРО, %s!</b>

📋 <b>ЕЖЕДНЕВНЫЕ КВЕСТЫ</b>

Новый день — новые возможности!
Выполни 5 ежедневных квестов и получи очки навыков!

🎮 Начать игру: /play`, name)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	return err == nil
}

func sendFinalBattleNotification(bot *tgbotapi.BotAPI, chatID int64, name string) bool {
	text := fmt.Sprintf(`🌙 <b>ВЕЧЕР, %s!</b>

⚔️  ВРЕМЯ ФИНАЛЬНОЙ БИТВЫ!

Сразись с боссом-искушением и докажи свою силу воли!

🎮 В бой: /play`, name)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	return err == nil
}

func sendUnfinishedQuestsNotification(bot *tgbotapi.BotAPI, chatID int64, name string) bool {
	// Упрощённая версия — просто отправляем уведомление
	text := fmt.Sprintf(`⏰ <b>%s, ВРЕМЯ ПОДЖАЛО!</b>

У тебя есть незавершённые квесты!

🎮 Выполнить: /play`, name)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	return err == nil
}

// GetSettings получает настройки уведомлений
func GetSettings(chatID int64) *NotificationSettings {
	if NotificationScheduler != nil {
		return NotificationScheduler.getSettings(chatID)
	}
	return DefaultNotificationSettings(chatID)
}

// SaveSettings сохраняет настройки уведомлений
func SaveSettings(settings *NotificationSettings) error {
	if NotificationScheduler != nil {
		return NotificationScheduler.saveSettings(settings)
	}
	return SaveNotificationSettings(settings)
}

func (nm *NotificationManager) getSettings(chatID int64) *NotificationSettings {
	if settings, exists := nm.settings[chatID]; exists {
		return settings
	}
	settings, _ := LoadNotificationSettings(chatID)
	if settings == nil {
		settings = DefaultNotificationSettings(chatID)
	}
	nm.settings[chatID] = settings
	return settings
}

func (nm *NotificationManager) saveSettings(settings *NotificationSettings) error {
	nm.settings[settings.ChatID] = settings
	return SaveNotificationSettings(settings)
}

// LoadNotificationSettings загружает настройки из БД
func LoadNotificationSettings(chatID int64) (*NotificationSettings, error) {
	query := `SELECT enabled, daily_quests_enabled, final_battle_enabled, unfinished_enabled, quests_hour, battle_hour FROM notification_settings WHERE chat_id = ?`
	row := database.DB.QueryRow(query, chatID)

	settings := &NotificationSettings{ChatID: chatID}
	err := row.Scan(&settings.Enabled, &settings.DailyQuestsEnabled, &settings.FinalBattleEnabled, &settings.UnfinishedEnabled, &settings.QuestsHour, &settings.BattleHour)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return settings, nil
}

// SaveNotificationSettings сохраняет настройки в БД
func SaveNotificationSettings(settings *NotificationSettings) error {
	query := `INSERT OR REPLACE INTO notification_settings (chat_id, enabled, daily_quests_enabled, final_battle_enabled, unfinished_enabled, quests_hour, battle_hour) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, settings.ChatID, settings.Enabled, settings.DailyQuestsEnabled, settings.FinalBattleEnabled, settings.UnfinishedEnabled, settings.QuestsHour, settings.BattleHour)
	return err
}
