package main

import (
	"github.com/joho/godotenv"
	"github.com/mskovv/tg-bot-subaru96/internal/database"
	"github.com/mskovv/tg-bot-subaru96/internal/handler"
	"github.com/mskovv/tg-bot-subaru96/internal/repository"
	"github.com/mskovv/tg-bot-subaru96/internal/service"
	"github.com/mskovv/tg-bot-subaru96/internal/storage"
	"github.com/mymmrac/telego"
	"log"
	"os"

	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TG_BOT_TOKEN")
	redisAddr := os.Getenv("DOCKER_REDIS_PORT")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
	}

	me, err := bot.GetMe()
	log.Printf("Authorized on account  %v\n", me)

	db := database.ConnectDb()

	appointmentRepo := repository.NewAppointmentRepository(db)
	appointmentService := service.NewAppointmentService(appointmentRepo)
	redisStorage, err := storage.NewRedisStorage(redisAddr)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService, redisStorage, bot)

	commands := []telego.BotCommand{
		{Command: "create_appointment", Description: "Создать запись"},
		{Command: "update_appointment", Description: "Обновить запись"},
	}
	err = bot.SetMyCommands(&telego.SetMyCommandsParams{Commands: commands})
	if err != nil {
		log.Fatal(err)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)

	bh, _ := th.NewBotHandler(bot, updates)
	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		appointmentHandler.SendStartMessage(update)
	}, th.CommandEqual("start"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		appointmentHandler.HandleMessage(update)
	}, th.Any(), th.AnyCommand())

	bh.HandleCallbackQuery(func(bot *telego.Bot, callbackQuery telego.CallbackQuery) {
		appointmentHandler.HandleCallback(callbackQuery)
	}, th.AnyCallbackQuery())

	bh.Start()

}
