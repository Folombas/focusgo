package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB — глобальный экземпляр базы данных
var DB *sql.DB

// InitDB инициализирует базу данных и создаёт таблицы
func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Проверка подключения
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	log.Println("✅ Подключение к базе данных установлено")

	// Применяем миграции
	if err = InitMigrations(); err != nil {
		return fmt.Errorf("ошибка применения миграций: %w", err)
	}

	return nil
}

// CloseDB закрывает подключение к базе данных
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// createTables создаёт таблицы если они не существуют
func createTables() error {
	queries := []string{
		// Таблица игроков
		`CREATE TABLE IF NOT EXISTS players (
			chat_id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			level INTEGER DEFAULT 1,
			experience INTEGER DEFAULT 0,
			go_knowledge INTEGER DEFAULT 40,
			focus INTEGER DEFAULT 70,
			willpower INTEGER DEFAULT 65,
			money INTEGER DEFAULT 500,
			dopamine INTEGER DEFAULT 200,
			play_time INTEGER DEFAULT 0,
			days_played INTEGER DEFAULT 1,
			current_day INTEGER DEFAULT 1,
			hour INTEGER DEFAULT 8,
			game_active INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Таблица навыков
		`CREATE TABLE IF NOT EXISTS skills (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			skill_id TEXT NOT NULL,
			level INTEGER DEFAULT 0,
			unlocked INTEGER DEFAULT 0,
			FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE,
			UNIQUE(chat_id, skill_id)
		)`,

		// Таблица квестов
		`CREATE TABLE IF NOT EXISTS quests (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			quest_id TEXT NOT NULL,
			progress INTEGER DEFAULT 0,
			completed INTEGER DEFAULT 0,
			deadline DATE,
			FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE,
			UNIQUE(chat_id, quest_id)
		)`,

		// Таблица достижений
		`CREATE TABLE IF NOT EXISTS achievements (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			achievement TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE
		)`,

		// Таблица искушений (преодолённых)
		`CREATE TABLE IF NOT EXISTS temptations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			temptation_name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE
		)`,

		// Таблица игровых сессий (история дней)
		`CREATE TABLE IF NOT EXISTS game_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			chat_id INTEGER NOT NULL,
			day_number INTEGER NOT NULL,
			score INTEGER DEFAULT 0,
			boss_defeated INTEGER DEFAULT 0,
			quests_completed INTEGER DEFAULT 0,
			play_time INTEGER DEFAULT 0,
			completed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE
		)`,

		// Таблица серии дней (day streak)
		`CREATE TABLE IF NOT EXISTS day_streaks (
			chat_id INTEGER PRIMARY KEY,
			current_streak INTEGER DEFAULT 0,
			best_streak INTEGER DEFAULT 0,
			last_quest_date DATE,
			total_quests_completed INTEGER DEFAULT 0,
			FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE
		)`,

		// Индексы для ускорения
		`CREATE INDEX IF NOT EXISTS idx_skills_chat ON skills(chat_id)`,
		`CREATE INDEX IF NOT EXISTS idx_quests_chat ON quests(chat_id)`,
		`CREATE INDEX IF NOT EXISTS idx_achievements_chat ON achievements(chat_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_chat ON game_sessions(chat_id)`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			return fmt.Errorf("ошибка выполнения запроса: %w", err)
		}
	}

	log.Println("✅ Таблицы базы данных созданы")
	return nil
}

// SavePlayer сохраняет игрока в базу данных
func SavePlayer(player *Player) error {
	query := `
		INSERT OR REPLACE INTO players 
		(chat_id, name, level, experience, go_knowledge, focus, willpower, 
		 money, dopamine, play_time, days_played, current_day, hour, game_active, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := DB.Exec(query,
		player.ChatID,
		player.Name,
		player.Level,
		player.Experience,
		player.GoKnowledge,
		player.Focus,
		player.Willpower,
		player.Money,
		player.Dopamine,
		player.PlayTime,
		player.DaysPlayed,
		player.CurrentDay,
		player.Hour,
		boolToInt(player.GameActive),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("ошибка сохранения игрока: %w", err)
	}

	// Сохраняем навыки
	if err = saveSkills(player); err != nil {
		return err
	}

	// Сохраняем квесты
	if err = saveQuests(player); err != nil {
		return err
	}

	// Сохраняем достижения
	if err = saveAchievements(player); err != nil {
		return err
	}

	log.Printf("💾 Игрок %s (chat_id: %d) сохранён", player.Name, player.ChatID)
	return nil
}

// LoadPlayer загружает игрока из базы данных
func LoadPlayer(chatID int64) (*Player, error) {
	query := `
		SELECT chat_id, name, level, experience, go_knowledge, focus, willpower,
		       money, dopamine, play_time, days_played, current_day, hour, game_active
		FROM players
		WHERE chat_id = ?
	`

	row := DB.QueryRow(query, chatID)

	player := &Player{}
	var gameActive int

	err := row.Scan(
		&player.ChatID,
		&player.Name,
		&player.Level,
		&player.Experience,
		&player.GoKnowledge,
		&player.Focus,
		&player.Willpower,
		&player.Money,
		&player.Dopamine,
		&player.PlayTime,
		&player.DaysPlayed,
		&player.CurrentDay,
		&player.Hour,
		&gameActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Игрок не найден
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки игрока: %w", err)
	}

	player.GameActive = intToBool(gameActive)

	// Инициализируем системы
	player.SkillTree = NewSkillTree()
	player.Quests = NewQuestSystem()

	// Загружаем навыки
	if err = loadSkills(player); err != nil {
		return nil, err
	}

	// Загружаем квесты
	if err = loadQuests(player); err != nil {
		return nil, err
	}

	// Загружаем достижения
	if err = loadAchievements(player); err != nil {
		return nil, err
	}

	// Загружаем серию дней
	if err = loadDayStreak(player); err != nil {
		return nil, err
	}

	// Применяем бонусы от навыков
	player.ApplySkillBonuses()

	log.Printf("💾 Игрок %s (chat_id: %d) загружен", player.Name, player.ChatID)
	return player, nil
}

// saveSkills сохраняет навыки игрока
func saveSkills(player *Player) error {
	query := `
		INSERT OR REPLACE INTO skills (chat_id, skill_id, level, unlocked)
		VALUES (?, ?, ?, ?)
	`

	for id, skill := range player.SkillTree.Skills {
		_, err := DB.Exec(query,
			player.ChatID,
			id,
			skill.Level,
			boolToInt(skill.Unlocked),
		)
		if err != nil {
			return fmt.Errorf("ошибка сохранения навыка %s: %w", id, err)
		}
	}

	return nil
}

// loadSkills загружает навыки игрока
func loadSkills(player *Player) error {
	query := `SELECT skill_id, level, unlocked FROM skills WHERE chat_id = ?`

	rows, err := DB.Query(query, player.ChatID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var skillID string
		var level int
		var unlocked int

		if err := rows.Scan(&skillID, &level, &unlocked); err != nil {
			return err
		}

		if skill, exists := player.SkillTree.Skills[skillID]; exists {
			skill.Level = level
			skill.Unlocked = intToBool(unlocked)
		}
	}

	return rows.Err()
}

// saveQuests сохраняет квесты игрока
func saveQuests(player *Player) error {
	query := `
		INSERT OR REPLACE INTO quests (chat_id, quest_id, progress, completed, deadline)
		VALUES (?, ?, ?, ?, ?)
	`

	for _, quest := range player.Quests.Quests {
		_, err := DB.Exec(query,
			player.ChatID,
			quest.ID,
			quest.Progress,
			boolToInt(quest.Completed),
			quest.Deadline,
		)
		if err != nil {
			return fmt.Errorf("ошибка сохранения квеста %s: %w", quest.ID, err)
		}
	}

	return nil
}

// loadQuests загружает квесты игрока
func loadQuests(player *Player) error {
	query := `SELECT quest_id, progress, completed, deadline FROM quests WHERE chat_id = ?`

	rows, err := DB.Query(query, player.ChatID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var questID string
		var progress int
		var completed int
		var deadline string

		if err := rows.Scan(&questID, &progress, &completed, &deadline); err != nil {
			return err
		}

		// Ищем квест в системе
		for _, quest := range player.Quests.Quests {
			if quest.ID == questID {
				quest.Progress = progress
				quest.Completed = intToBool(completed)
				quest.Deadline = deadline
				break
			}
		}
	}

	return rows.Err()
}

// saveAchievements сохраняет достижения
func saveAchievements(player *Player) error {
	// Сначала удаляем старые достижения
	_, err := DB.Exec("DELETE FROM achievements WHERE chat_id = ?", player.ChatID)
	if err != nil {
		return err
	}

	query := `INSERT INTO achievements (chat_id, achievement) VALUES (?, ?)`

	for _, achievement := range player.Achievements {
		_, err := DB.Exec(query, player.ChatID, achievement)
		if err != nil {
			return fmt.Errorf("ошибка сохранения достижения: %w", err)
		}
	}

	return nil
}

// loadAchievements загружает достижения
func loadAchievements(player *Player) error {
	query := `SELECT achievement FROM achievements WHERE chat_id = ? ORDER BY created_at`

	rows, err := DB.Query(query, player.ChatID)
	if err != nil {
		return err
	}
	defer rows.Close()

	player.Achievements = []string{}

	for rows.Next() {
		var achievement string
		if err := rows.Scan(&achievement); err != nil {
			return err
		}
		player.Achievements = append(player.Achievements, achievement)
	}

	return rows.Err()
}

// loadDayStreak загружает серию дней
func loadDayStreak(player *Player) error {
	query := `SELECT current_streak, best_streak, last_quest_date, total_quests_completed 
			  FROM day_streaks WHERE chat_id = ?`

	row := DB.QueryRow(query, player.ChatID)

	var currentStreak, bestStreak, totalCompleted int
	var lastQuestDate sql.NullString

	err := row.Scan(&currentStreak, &bestStreak, &lastQuestDate, &totalCompleted)
	if err == sql.ErrNoRows {
		// Создаём новую запись
		_, err = DB.Exec(
			"INSERT INTO day_streaks (chat_id, current_streak, best_streak, total_quests_completed) VALUES (?, 0, 0, 0)",
			player.ChatID,
		)
		return err
	}
	if err != nil {
		return err
	}

	player.Quests.DayStreak = currentStreak
	player.Quests.TotalCompleted = totalCompleted

	return nil
}

// saveDayStreak сохраняет серию дней
func saveDayStreak(player *Player) error {
	query := `
		INSERT OR REPLACE INTO day_streaks 
		(chat_id, current_streak, best_streak, total_quests_completed)
		VALUES (?, ?, 
			COALESCE((SELECT MAX(best_streak, ?) FROM day_streaks WHERE chat_id = ?), ?),
			?)
	`

	_, err := DB.Exec(query,
		player.ChatID,
		player.Quests.DayStreak,
		player.Quests.DayStreak,
		player.ChatID,
		player.Quests.DayStreak,
		player.Quests.TotalCompleted,
	)

	return err
}

// saveGameSession сохраняет игровую сессию (историю дня)
func saveGameSession(chatID int64, dayNumber, score, playTime int, bossDefeated bool, questsCompleted int) error {
	query := `
		INSERT INTO game_sessions 
		(chat_id, day_number, score, boss_defeated, quests_completed, play_time)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := DB.Exec(query,
		chatID,
		dayNumber,
		score,
		boolToInt(bossDefeated),
		questsCompleted,
		playTime,
	)

	return err
}

// GetLeaderboard возвращает таблицу лидеров
func GetLeaderboard(limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT name, level, experience, go_knowledge, 
		       (go_knowledge * 10 + focus * 5 + willpower * 3) as rating
		FROM players
		ORDER BY rating DESC
		LIMIT ?
	`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leaderboard := []map[string]interface{}{}

	for rows.Next() {
		var name string
		var level, experience, goKnowledge, rating int

		if err := rows.Scan(&name, &level, &experience, &goKnowledge, &rating); err != nil {
			return nil, err
		}

		leaderboard = append(leaderboard, map[string]interface{}{
			"name":     name,
			"level":    level,
			"rating":   rating,
			"knowledge": goKnowledge,
		})
	}

	return leaderboard, rows.Err()
}

// GetTotalPlayers возвращает общее количество игроков
func GetTotalPlayers() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM players").Scan(&count)
	return count, err
}

// GetPlayerStats возвращает статистику игрока
func GetPlayerStats(chatID int64) (map[string]interface{}, error) {
	// Количество преодолённых искушений
	var temptationsCount int
	err := DB.QueryRow(
		"SELECT COUNT(*) FROM temptations WHERE chat_id = ?",
		chatID,
	).Scan(&temptationsCount)
	if err != nil {
		return nil, err
	}

	// Количество дней в игре
	var sessionsCount int
	err = DB.QueryRow(
		"SELECT COUNT(*) FROM game_sessions WHERE chat_id = ?",
		chatID,
	).Scan(&sessionsCount)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"temptations_resisted": temptationsCount,
		"days_played":          sessionsCount,
	}, nil
}

// Вспомогательные функции

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func intToBool(i int) bool {
	return i != 0
}

// PlayerToJSON экспортирует игрока в JSON
func PlayerToJSON(player *Player) ([]byte, error) {
	data := map[string]interface{}{
		"chat_id":      player.ChatID,
		"name":         player.Name,
		"level":        player.Level,
		"experience":   player.Experience,
		"go_knowledge": player.GoKnowledge,
		"focus":        player.Focus,
		"willpower":    player.Willpower,
		"money":        player.Money,
		"dopamine":     player.Dopamine,
		"play_time":    player.PlayTime,
		"days_played":  player.DaysPlayed,
		"achievements": player.Achievements,
		"temptations":  player.Temptations,
	}

	// Сериализуем навыки
	skills := make(map[string]int)
	for id, skill := range player.SkillTree.Skills {
		skills[id] = skill.Level
	}
	data["skills"] = skills

	return json.Marshal(data)
}

// PlayerFromJSON импортирует игрока из JSON
func PlayerFromJSON(data []byte) (*Player, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}

	player := &Player{
		ChatID:      int64(jsonData["chat_id"].(float64)),
		Name:        jsonData["name"].(string),
		Level:       int(jsonData["level"].(float64)),
		Experience:  int(jsonData["experience"].(float64)),
		GoKnowledge: int(jsonData["go_knowledge"].(float64)),
		Focus:       int(jsonData["focus"].(float64)),
		Willpower:   int(jsonData["willpower"].(float64)),
		Money:       int(jsonData["money"].(float64)),
		Dopamine:    int(jsonData["dopamine"].(float64)),
		PlayTime:    int(jsonData["play_time"].(float64)),
		DaysPlayed:  int(jsonData["days_played"].(float64)),
	}

	// Восстанавливаем достижения
	if achievements, ok := jsonData["achievements"].([]interface{}); ok {
		for _, a := range achievements {
			player.Achievements = append(player.Achievements, a.(string))
		}
	}

	// Восстанавливаем искушения
	if temptations, ok := jsonData["temptations"].([]interface{}); ok {
		for _, t := range temptations {
			player.Temptations = append(player.Temptations, t.(string))
		}
	}

	// Инициализируем системы
	player.SkillTree = NewSkillTree()
	player.Quests = NewQuestSystem()

	// Восстанавливаем навыки
	if skills, ok := jsonData["skills"].(map[string]interface{}); ok {
		for id, level := range skills {
			if skill, exists := player.SkillTree.Skills[id]; exists {
				skill.Level = int(level.(float64))
				if level.(float64) > 0 {
					skill.Unlocked = true
				}
			}
		}
	}

	return player, nil
}
