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
	defaultPollFreq = 3 * time.Second
	bankDetails     = `
Stichting Oekraïners in Nederland
IBAN: NL 97 INGB 0006 5104 66
SWIFT/BIC: INGBNL2AXXX
Bank: ING Bank
Address: Bijlmerdreef 106, 1102 MG Amsterdam, Netherlands
`
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
				e.SetCaption(c, bankDetails)
				return menu.Forward
			}),
			flow.NewNode("humanitarian", forward).
				Add("clothes", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "ОДЯГ ТА ВЗУТТЯ ПОКИ НЕ ПРИЙМАЮТЬСЯ")
					return menu.Forward
				}).
				AddWith("essentials", forward,
					flow.NewNode("meds", func(e *menu.Node, c *tb.Callback) int {
						e.SetCaption(c, "https://docs.google.com/spreadsheets/d/1lcT6IhtkT-1G6rKi7wmxoTe_EhYHnHRtIAjLQATi1bQ/")
						return menu.Forward
					}),
					flow.NewNode("hygiene-food", func(e *menu.Node, c *tb.Callback) int {
						e.SetCaption(c, "https://docs.google.com/spreadsheets/d/1x320lSGqgqeCpiCBbeGxR2RtYcqJhh7d30XNYGomogc/")
						return menu.Forward
					}),
					flow.NewNode("ammunition", func(e *menu.Node, c *tb.Callback) int {
						e.SetCaption(c, "https://docs.google.com/spreadsheets/d/1nJWUQYcH3qkzC0fP7Jy8YYntW_fcPQDOXRIBKnVAlfc/")
						return menu.Forward
					}),
					flow.NewNode("reception", func(e *menu.Node, c *tb.Callback) int {
						e.SetCaption(c, "https://docs.google.com/spreadsheets/d/1W7fxCS8ZAbhsNrGeRzseCN0GD9cJnVMA/")
						return menu.Forward
					}),
					flow.NewNode("back", back),
				).
				Add("back", back),
			flow.NewNode("other", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "https://forms.gle/fDZHTH7pnPCNUfJdA")
				return menu.Forward
			}),
			flow.NewNode("back", back),
		).
		AddWith("volunteer", forward,
			flow.NewNode("humanitarian", forward).
				Add("support", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "https://docs.google.com/forms/d/e/1FAIpQLSfIYchO_qFxPFLti-0287FMr7B5ue5_laRq4TG3uSwsN_ZpwQ/viewform")
					return menu.Forward
				}).
				Add("reception", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "https://docs.google.com/forms/d/e/1FAIpQLSenM3-eAq7zj1VfmpYoFSnYFC7qdILmYQ6XY-hxWXbK36FI7w/viewform")
					return menu.Forward
				}).
				Add("logistics", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "https://docs.google.com/forms/d/e/1FAIpQLSe8IzPERiS33i6xJYO3cXWEo3bM4ig9rb8wINzBU49SMr9luQ/viewform")
					return menu.Forward
				}).
				Add("sorting", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "https://docs.google.com/forms/d/e/1FAIpQLSfU1JIBAE1Fe7XKg6ewAgTOF7LjnlipwXa8xFS2Xx9UkGZ93w/viewform")
					return menu.Forward
				}).
				Add("requests-ua", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "https://docs.google.com/forms/d/e/1FAIpQLSfdT5nlzQiq2SgxEYK5sZ60ZmirugGsn85crxcOc4Fv-cYyVw/viewform")
					return menu.Forward
				}).
				Add("back", back),
			flow.NewNode("coordination", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "https://forms.gle/KKCoYbeqp5PbUxA27")
				return menu.Forward
			}),
			flow.NewNode("events", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "https://forms.gle/Fh8dxRmjaMaKW5rz6")
				return menu.Forward
			}),
			flow.NewNode("pr", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "https://forms.gle/j5BuFL7apSpZQg2s9")
				return menu.Forward
			}),
			flow.NewNode("refugees", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "https://forms.gle/ihRiN5LwS7Pfzj6C6")
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
