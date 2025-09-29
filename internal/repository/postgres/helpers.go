package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// checkGetQueryResult обрабатывает результат GORM-запроса и возвращает notFound, если запись не найдена
func checkGetQueryResult(result *gorm.DB, notFound error) error {
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return notFound
	}

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

// postgresDuplicate проверяет, вызвана ли ошибка дублирования уникального значения
func postgresDuplicate(result *gorm.DB, ErrIsExists error) error {
	if err := result.Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrIsExists
		}

		return err
	}

	return nil
}

// postgresForeignKeyViolation проверяет, была ли нарушена ссылка внешнего ключа
func postgresForeignKeyViolation(result *gorm.DB, ErrInUse error) error {
	if err := result.Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return ErrInUse
		}

		return err
	}

	return nil
}
