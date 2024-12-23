package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mskovv/tg-bot-subaru96/internal/database"
	fsmstate "github.com/mskovv/tg-bot-subaru96/internal/fsm"
	"github.com/mskovv/tg-bot-subaru96/internal/handler/appointment"
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
	envMode := os.Getenv("ENV_MODE")

	bot, err := telego.NewBot(botToken)
	if err != nil {
		log.Fatal(err)
	}

	me, err := bot.GetMe()
	log.Printf("Authorized on account  %v\n", me)

	db := database.ConnectDb()

	appointmentRepo := repository.NewAppointmentRepository(db)
	appointmentService := service.NewAppointmentService(appointmentRepo)

	carDictionaryRepo := repository.NewCarDictionaryRepository(db)
	carDictionaryService := service.NewCarDictionaryService(carDictionaryRepo)

	redisStorage, err := storage.NewRedisStorage()
	fsmState := fsmstate.NewAppointmentFSM()
	appointmentHandler := appointment.NewAppointmentHandler(appointmentService, carDictionaryService, redisStorage, bot, fsmState)

	commands := []telego.BotCommand{
		{Command: "create_appointment", Description: "Создать запись"},
		{Command: "view_appointments", Description: "Посмотреть записи"},
		{Command: "update_appointment", Description: "Обновить запись"},
		{Command: "start", Description: "Запустить бота/Сбросить"},
	}
	err = bot.SetMyCommands(&telego.SetMyCommandsParams{Commands: commands})
	if err != nil {
		log.Fatal(err)
	}
	if envMode == "dev" {

		updates, _ := bot.UpdatesViaLongPolling(nil)

		bh, _ := th.NewBotHandler(bot, updates)
		defer bh.Stop()
		defer bot.StopLongPolling()

		setupHandlers(bh, appointmentHandler)
		bh.Start()
	} else if envMode == "prod" {
		webhookURL := os.Getenv("WEBHOOK_URL") + bot.Token() // URL для вебхука
		err = bot.SetWebhook(&telego.SetWebhookParams{
			URL: webhookURL,
		})
		if err != nil {
			fmt.Println(bot.GetWebhookInfo())
			log.Fatal(err)
		}

		log.Printf("Webhook set to: %s", webhookURL)

		go func() {
			err = bot.StartWebhook("localhost:443")
			if err != nil {
				webhook, _ := bot.GetWebhookInfo()
				fmt.Println("webhook", webhook)
				log.Fatal(err)
			}
		}()

		defer func() {
			_ = bot.StopWebhook()
		}()

		updates, err := bot.UpdatesViaWebhook(webhookURL) // Путь для вебхука
		if err != nil {
			log.Fatal(err)
		}
		bh, _ := th.NewBotHandler(bot, updates)
		defer bh.Stop()

		setupHandlers(bh, appointmentHandler)
		bh.Start()
	}

}

func setupHandlers(bh *th.BotHandler, appointmentHandler *appointment.Handler) {
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		appointmentHandler.SendStartMessage(update)
	}, th.CommandEqual("start"))

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		appointmentHandler.HandleCommand(update)
	}, th.AnyCommand())

	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		appointmentHandler.HandleMessage(message)
	}, th.AnyMessage())

	bh.HandleCallbackQuery(func(bot *telego.Bot, callbackQuery telego.CallbackQuery) {
		appointmentHandler.HandleCallback(callbackQuery)
	}, th.AnyCallbackQuery())
}
