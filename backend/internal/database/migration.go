package database

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// RunMigrations 执行数据库迁移
func RunMigrations(db *gorm.DB) error {
	// 获取迁移文件目录
	migrationDir := "migrations"
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		log.Println("No migrations directory found, skipping migrations")
		return nil
	}

	// 获取所有SQL文件
	var files []string
	err := filepath.WalkDir(migrationDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".sql") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// 按文件名排序
	sort.Strings(files)

	// 执行每个迁移文件
	for _, file := range files {
		// 对于包含 seed 的迁移文件，采用安全策略：
		// 仅在目标表为空时执行，以避免每次服务重启时重新插入被删除的数据。
		isSeed := strings.Contains(strings.ToLower(filepath.Base(file)), "seed")
		if isSeed {
			var cnt int64
			// 检查 prompt_templates 表是否已有数据
			if err := db.Table("prompt_templates").Count(&cnt).Error; err != nil {
				return fmt.Errorf("failed to check prompt_templates count: %w", err)
			}
			if cnt > 0 {
				log.Printf("Skipping seed migration %s because prompt_templates has %d rows", file, cnt)
				continue
			}
		}

		log.Printf("Executing migration: %s", file)

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// 执行SQL
		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		log.Printf("Migration completed: %s", file)
	}

	log.Println("All migrations completed successfully")
	return nil
}
