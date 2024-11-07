package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mskovv/tg-bot-subaru96/internal/database"
	"github.com/mskovv/tg-bot-subaru96/internal/handler"
	"github.com/mskovv/tg-bot-subaru96/internal/repository"
	"github.com/mskovv/tg-bot-subaru96/internal/service"
	"log"
	"os"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	db := database.ConnectDb()
	appointmentRepo := repository.NewAppointmentRepository(db)
	appointmentService := service.NewAppointmentService(appointmentRepo)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		switch update.Message.Command() {
		case "remove_appointment":
			appointmentHandler.RemoveAppointment(update)
		}
	}

}
