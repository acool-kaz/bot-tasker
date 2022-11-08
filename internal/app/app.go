package app

import (
	"log"
	"time"

	"github.com/acool-kaz/bot-tasker/internal/config"
	"github.com/acool-kaz/bot-tasker/internal/delivery"
	"github.com/acool-kaz/bot-tasker/internal/models"
	"github.com/acool-kaz/bot-tasker/internal/repository"
	"github.com/acool-kaz/bot-tasker/internal/service"
	"gopkg.in/telebot.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	config *config.Config
	bot    *telebot.Bot
}

func NewApp(cfg *config.Config) *App {
	log.Println("init telebot")
	pref := telebot.Settings{
		Token:  cfg.Bot.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	return &App{
		config: cfg,
		bot:    b,
	}
}

func (a *App) Run() {
	log.Println("init gorm db")
	db, err := gorm.Open(sqlite.Open(a.config.Gorm.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(models.User{}, models.Task{}); err != nil {
		log.Fatal(err)
	}
	// if err := db.Migrator().CreateConstraint(&models.User{}, "Tasks"); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := db.Migrator().CreateConstraint(&models.User{}, "fk_users_tasks"); err != nil {
	// 	log.Fatal(err)
	// }

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := delivery.NewHandler(service)

	a.bot.Handle("/start", handler.StartHandler)
	a.bot.Handle("/addTask", handler.AddTaskHandler)
	a.bot.Handle("/tasks", handler.AllTasksHandler)
	a.bot.Handle("/deleteTask", handler.DeleteTaskHandler)
	a.bot.Handle(telebot.OnText, handler.NewMsgHandler)

	log.Println("start bot")
	a.bot.Start()
}
