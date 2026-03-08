package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// CreateBackup создаёт резервную копию базы данных
func CreateBackup() (string, error) {
	// Создаём директорию для бэкапов
	backupDir := "backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", fmt.Errorf("ошибка создания директории бэкапов: %w", err)
	}

	// Имя файла с датой
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("focusgo_backup_%s.db", timestamp))

	// Открываем исходную БД
	src, err := os.Open("focusgo.db")
	if err != nil {
		return "", fmt.Errorf("ошибка открытия БД: %w", err)
	}
	defer src.Close()

	// Создаём файл бэкапа
	dst, err := os.Create(backupFile)
	if err != nil {
		return "", fmt.Errorf("ошибка создания файла бэкапа: %w", err)
	}
	defer dst.Close()

	// Копируем данные
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", fmt.Errorf("ошибка копирования данных: %w", err)
	}

	log.Printf("💾 Бэкап создан: %s", backupFile)
	return backupFile, nil
}

// AutoBackup запускает автоматические бэкапы
func AutoBackup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			backupPath, err := CreateBackup()
			if err != nil {
				log.Printf("❌ Ошибка автоматического бэкапа: %v", err)
			} else {
				log.Printf("✅ Автоматический бэкап: %s", backupPath)
			}

			// Удаляем старые бэкапы (храним последние 7)
			cleanupOldBackups(7)
		}
	}()
	log.Println("✅ Автоматические бэкапы запущены")
}

// cleanupOldBackups удаляет старые бэкапы, оставляя только последние n
func cleanupOldBackups(keep int) {
	backupDir := "backups"
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return
	}

	var backups []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && len(file.Name()) > 5 && file.Name()[len(file.Name())-5:] == ".db" {
			backups = append(backups, file)
		}
	}

	// Сортируем по имени (дата в имени)
	// Удаляем старые, оставляем keep последних
	if len(backups) > keep {
		for i := 0; i < len(backups)-keep; i++ {
			filePath := filepath.Join(backupDir, backups[i].Name())
			os.Remove(filePath)
			log.Printf("🗑️  Удалён старый бэкап: %s", filePath)
		}
	}
}

// ListBackups возвращает список доступных бэкапов
func ListBackups() ([]string, error) {
	backupDir := "backups"
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, err
	}

	var backups []string
	for _, file := range files {
		if !file.IsDir() && len(file.Name()) > 5 && file.Name()[len(file.Name())-5:] == ".db" {
			backups = append(backups, file.Name())
		}
	}

	return backups, nil
}

// RestoreBackup восстанавливает базу данных из бэкапа
func RestoreBackup(backupFile string) error {
	backupPath := filepath.Join("backups", backupFile)

	// Открываем бэкап
	src, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия бэкапа: %w", err)
	}
	defer src.Close()

	// Создаём временный файл
	tempFile := "focusgo_temp.db"
	dst, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("ошибка создания временного файла: %w", err)
	}
	defer dst.Close()

	// Копируем данные
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("ошибка копирования данных: %w", err)
	}

	// Закрываем текущую БД
	if DB != nil {
		DB.Close()
	}

	// Переименовываем файлы
	os.Rename("focusgo.db", "focusgo.old.db")
	os.Rename(tempFile, "focusgo.db")

	// Переоткрываем БД
	return InitDB("focusgo.db")
}
