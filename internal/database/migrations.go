package database

import (
	"fmt"
	"log"
)

// Миграции базы данных
// Каждая миграция имеет версию и функцию применения

type Migration struct {
	Version int
	Name    string
	Up      func() error
}

// Миграции для применения
var migrations = []Migration{
	{
		Version: 1,
		Name:    "create_tables",
		Up:      createTables,
	},
	{
		Version: 2,
		Name:    "add_temptations_table",
		Up: func() error {
			query := `
				CREATE TABLE IF NOT EXISTS temptations_resisted (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					chat_id INTEGER NOT NULL,
					temptation_name TEXT NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE
				);
				CREATE INDEX IF NOT EXISTS idx_temptations_chat ON temptations_resisted(chat_id);
			`
			_, err := DB.Exec(query)
			return err
		},
	},
	{
		Version: 3,
		Name:    "add_leaderboard_index",
		Up: func() error {
			query := `
				CREATE INDEX IF NOT EXISTS idx_players_rating 
				ON players((go_knowledge * 10 + focus * 5 + willpower * 3));
			`
			_, err := DB.Exec(query)
			return err
		},
	},
	{
		Version: 4,
		Name:    "add_notification_settings",
		Up: func() error {
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
		},
	},
	{
		Version: 5,
		Name:    "add_skill_trees",
		Up: func() error {
			queries := []string{
				`CREATE TABLE IF NOT EXISTS skill_trees (
					chat_id INTEGER PRIMARY KEY,
					skill_points INTEGER DEFAULT 0,
					total_points INTEGER DEFAULT 0,
					FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE
				)`,
				`CREATE TABLE IF NOT EXISTS skills (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					chat_id INTEGER NOT NULL,
					skill_id TEXT NOT NULL,
					level INTEGER DEFAULT 0,
					unlocked INTEGER DEFAULT 0,
					FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE,
					UNIQUE(chat_id, skill_id)
				)`,
				`CREATE INDEX IF NOT EXISTS idx_skills_chat ON skills(chat_id)`,
			}
			for _, query := range queries {
				_, err := DB.Exec(query)
				if err != nil {
					return err
				}
			}
			return nil
		},
	},
	{
		Version: 6,
		Name:    "add_achievements",
		Up: func() error {
			query := `
				CREATE TABLE IF NOT EXISTS achievements (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					chat_id INTEGER NOT NULL,
					achievement_id TEXT NOT NULL,
					unlocked_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (chat_id) REFERENCES players(chat_id) ON DELETE CASCADE,
					UNIQUE(chat_id, achievement_id)
				);
				CREATE INDEX IF NOT EXISTS idx_achievements_chat ON achievements(chat_id);
			`
			_, err := DB.Exec(query)
			return err
		},
	},
}

// InitMigrations инициализирует таблицу миграций и применяет новые
func InitMigrations() error {
	// Создаём таблицу для хранения версий миграций
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы миграций: %w", err)
	}

	log.Println("✅ Таблица миграций создана")

	// Применяем все неприменённые миграции
	for _, migration := range migrations {
		if !isMigrationApplied(migration.Version) {
			log.Printf("🔄 Применение миграции v%d: %s", migration.Version, migration.Name)

			if err := migration.Up(); err != nil {
				return fmt.Errorf("ошибка применения миграции v%d %s: %w",
					migration.Version, migration.Name, err)
			}

			if err := markMigrationAsApplied(migration.Version); err != nil {
				return fmt.Errorf("ошибка записи миграции v%d: %w", migration.Version, err)
			}

			log.Printf("✅ Миграция v%d применена", migration.Version)
		}
	}

	// Получаем текущую версию
	currentVersion, err := getCurrentMigrationVersion()
	if err != nil {
		return err
	}

	log.Printf("✅ База данных актуальна (версия схемы: %d)", currentVersion)
	return nil
}

// isMigrationApplied проверяет, применена ли миграция
func isMigrationApplied(version int) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = ?)"
	err := DB.QueryRow(query, version).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// markMigrationAsApplied записывает информацию о применённой миграции
func markMigrationAsApplied(version int) error {
	query := "INSERT INTO schema_migrations (version) VALUES (?)"
	_, err := DB.Exec(query, version)
	return err
}

// getCurrentMigrationVersion возвращает текущую версию схемы
func getCurrentMigrationVersion() (int, error) {
	var version int
	query := "SELECT COALESCE(MAX(version), 0) FROM schema_migrations"
	err := DB.QueryRow(query).Scan(&version)
	return version, err
}

// GetMigrationVersion возвращает текущую версию миграции
func GetMigrationVersion() int {
	version, _ := getCurrentMigrationVersion()
	return version
}

// RollbackMigration откатывает последнюю миграцию (для отладки)
func RollbackMigration() error {
	currentVersion, err := getCurrentMigrationVersion()
	if err != nil {
		return err
	}

	if currentVersion == 0 {
		return fmt.Errorf("нет миграций для отката")
	}

	// Находим миграцию для отката
	var migrationToRollback *Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if migrations[i].Version == currentVersion {
			migrationToRollback = &migrations[i]
			break
		}
	}

	if migrationToRollback == nil {
		return fmt.Errorf("миграция v%d не найдена", currentVersion)
	}

	log.Printf("🔄 Откат миграции v%d: %s", currentVersion, migrationToRollback.Name)

	// Удаляем запись о миграции
	query := "DELETE FROM schema_migrations WHERE version = ?"
	_, err = DB.Exec(query, currentVersion)
	if err != nil {
		return fmt.Errorf("ошибка удаления записи о миграции: %w", err)
	}

	log.Printf("✅ Миграция v%d откатана", currentVersion)
	return nil
}

// ListMigrations выводит список всех миграций
func ListMigrations() {
	log.Println("📋 Доступные миграции:")
	for _, m := range migrations {
		status := "⏳"
		if isMigrationApplied(m.Version) {
			status = "✅"
		}
		log.Printf("  %s v%d: %s", status, m.Version, m.Name)
	}
}
