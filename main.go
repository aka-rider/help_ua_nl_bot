package main

import (
	"log"
	"os"
	"time"

	"github.com/kisulken/go-telegram-flow/menu"

	"github.com/tucnak/tr"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	defaultLocale   = "ua"
	greetings       = "Привіт! Ми — українці в Нідерландах\nБажаєте допомогти Україні?"
	bye             = "Слава Україні "
	defaultPollFreq = 3 * time.Second
)

func Run(token string) {
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: defaultPollFreq},
	})
	if err != nil {
		panic(err)
	}

	if err := tr.Init("lang", defaultLocale); err != nil {
		panic(err)
	}

	flow, err := menu.NewMenuFlow("flow", bot, tr.DefaultEngine)
	if err != nil {
		panic(err)
	}

	flow.GetRoot().
		AddWith("once", forward,
			flow.NewNode("finance", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "https://bank.gov.ua/ua/news/all/natsionalniy-bank-vidkriv-spetsrahunok-dlya-zboru-koshtiv-na-potrebi-armiyi")
				return menu.Forward
			}),
			flow.NewNode("humanitarian", forward).
				Add("clothes", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "Червоний хрест в Нідерландах https://www.rodekruis.nl/")
					return menu.Forward
				}).
				Add("meds_food_hygiene", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "TODO: посилання на Excel")
					return menu.Forward
				}).
				Add("back", back),
			flow.NewNode("other", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "Звяжіться з Іриною +31684326175")
				return menu.Forward
			}),
			flow.NewNode("back", back),
		).
		AddWith("volunteer", forward,
			flow.NewNode("coordination", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "Звяжіться з Маєю +31648269715")
				return menu.Forward
			}),
			flow.NewNode("events", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "Звяжіться з Олександром +31615834846")
				return menu.Forward
			}),
			flow.NewNode("humanitarian", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "Звяжіться з Даною +31644535721 ")
				return menu.Forward
			}),
			flow.NewNode("pr", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "Звяжіться з Анастасіює +31643290389")
				return menu.Forward
			}),
			flow.NewNode("refugees", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "Звяжіться з Мариною +31646004341")
				return menu.Forward
			}),
			flow.NewNode("back", back),
		)

	flow.Build(defaultLocale)

	bot.Handle("/start", func(m *tb.Message) {
		err = flow.Start(m.Sender, greetings, defaultLocale)
		if err != nil {
			log.Println("failed to display the menu:", err)
		}
	})
	// bot.Handle("/stop", func(m *tb.Message) {
	// 	err = flow.Stop(m.Sender, bye, defaultLocale)
	// 	if err != nil {
	// 		log.Println("failed to stop the flow:", err)
	// 	}
	// })

	log.Println("running", bot.Me.Username, "...")
	bot.Start()
}

func userOrderSushi(e *menu.Node, c *tb.Callback) int {
	log.Println(c.Sender.Recipient(), "press", e.GetText())
	e.SetCaption(c, "Added "+e.GetText()+" to your order")
	return menu.Forward
}

func userPressLanguage(e *menu.Node, c *tb.Callback) int {
	log.Println(c.Sender.Recipient(), "press", e.GetText())
	if e.GetLanguage(c) == "en" {
		e.SetLanguage(c, "ru")
	} else {
		e.SetLanguage(c, "en")
	}
	return menu.Forward // continue
}

func forward(e *menu.Node, c *tb.Callback) int {
	log.Println(c.Sender.Recipient(), "->", e.GetText())
	e.SetCaption(c, "...")
	return menu.Forward
}

func back(e *menu.Node, c *tb.Callback) int {
	log.Println(c.Sender.Recipient(), "<-", e.GetText())
	e.SetCaption(c, "...")
	return menu.Back
}

func main() {
	token := os.Getenv("HELP_UA_NL_BOT_TOKEN")
	Run(token)
}
