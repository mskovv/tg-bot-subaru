package migrations

import (
	"github.com/mskovv/tg-bot-subaru96/internal/models"
	"gorm.io/gorm"
	"log"
)

func AddInitialUser(db *gorm.DB) {
	user := models.User{
		Name:         "Андрей",
		LastName:     "Плесовских",
		Nickname:     "Дюха",
		Role:         "Директор",
		Appointments: nil,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Не удалось добавить пользователя: %v", err)
	} else {
		log.Printf("Пользователь %s успешно добавлен с ролью %s.", user.Name, user.Role)
	}
}
