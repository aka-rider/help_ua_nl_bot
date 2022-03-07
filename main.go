package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/kisulken/go-telegram-flow/menu"

	"github.com/tucnak/tr"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	defaultLocale = "ua"

	startMsg_ua = `Відправте /start щоб почати спочатку`
	startMsg_en = `Enter /start to begin or restart`

	greetings_ua = `
Привіт! Ми українці в Нідерландах. Бажаєте допомогти Україні?

❗️ Не поширюйте цей чат бот за межами Нідерландів
` + startMsg_ua

	greetings_en = `
Hello! We are Ukrainians in the Netherlands. Do you want to help Ukraine?

❗️ Please do not share this chat bot outside the Netherlands
` + startMsg_en

	bankDetails = `
Stichting Oekraïners in Nederland
IBAN: NL 97 INGB 0006 5104 66
SWIFT/BIC: INGBNL2AXXX
Bank: ING Bank
Address: Bijlmerdreef 106, 1102 MG Amsterdam, Netherlands
`

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
				e.SetCaption(c, bankDetails)
				return menu.Forward
			}),
			flow.NewNode("humanitarian", forward).
				Add("clothes", func(e *menu.Node, c *tb.Callback) int {
					if e.GetLanguage(c) == "ua" {
						e.SetCaption(c, "🚫 НАРАЗІ НЕМАЄ ПОТРЕБИ В ОДЯГУ, ВЗУТТІ ТА ЇЖІ. БУДЬ ЛАСКА, НЕ НЕСІТЬ.")
					} else {
						e.SetCaption(c, "🚫 THERE IS NO NEED FOR CLOTHES, SHOES, AND FOOD.")
					}
					return menu.Forward
				}).
				Add("collection", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "https://help-ukraine.nl/collection-points-in-the-netherlands")
					return menu.Forward
				}).
				AddWith("essentials", forward,
					flow.NewNode("medication_and_dosage", func(e *menu.Node, c *tb.Callback) int {
						e.SetCaption(c, "https://help-ukraine.nl/medication-dosage")
						return menu.Forward
					}),
					flow.NewNode("medical_consumables", func(e *menu.Node, c *tb.Callback) int {
						e.SetCaption(c, "https://help-ukraine.nl/medical-consumables")
						return menu.Forward
					}),
					flow.NewNode("medical_equipment", func(e *menu.Node, c *tb.Callback) int {
						e.SetCaption(c, "https://help-ukraine.nl/medical-equipment")
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
		AddWith("volunteer", volunteerMenu,
			flow.NewNode("humanitarian", forward).
				Add("drivers", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "https://docs.google.com/forms/d/e/1FAIpQLSe8IzPERiS33i6xJYO3cXWEo3bM4ig9rb8wINzBU49SMr9luQ/viewform")
					return menu.Forward
				}).
				Add("general-support", func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, "https://forms.gle/5Ep4k8KiBn48jyeF9")
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
			flow.NewNode("other", func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "https://docs.google.com/forms/d/e/1FAIpQLSfpmhjq1O92sD5GmvD_X4J1AkZ4j8Q-tPW_KBQiBdTKcvcjHw/viewform")
				return menu.Forward
			}),
			flow.NewNode("back", back),
		).
		Add("refugee", func(e *menu.Node, c *tb.Callback) int {
			e.SetCaption(c, "https://help-ukraine.nl/refugee")
			return menu.Forward
		}).
		Add("language", switchLanguage)

	flow.Build("ua").Build("en")

	bot.Handle("/start", func(m *tb.Message) {
		err = flow.Start(m.Sender, greetings_ua, defaultLocale)
		if err != nil {
			log.Println("failed to display the menu:", err)
		}
	})

	log.Println("running", bot.Me.Username, "...")
	bot.Start()
}

func switchLanguage(e *menu.Node, c *tb.Callback) int {
	log.Println(c.Sender.Recipient(), "press", e.GetText())
	if e.GetLanguage(c) == "en" {
		e.SetLanguage(c, "ua")
		e.SetCaption(c, greetings_ua)
	} else {
		e.SetLanguage(c, "en")
		e.SetCaption(c, greetings_en)
	}
	return menu.Forward
}

func forward(e *menu.Node, c *tb.Callback) int {
	log.Println(c.Sender.Recipient(), "->", e.GetText())
	if e.GetLanguage(c) == "ua" {
		e.SetCaption(c, startMsg_ua)
	} else {
		e.SetCaption(c, startMsg_en)
	}
	return menu.Forward
}

func back(e *menu.Node, c *tb.Callback) int {
	log.Println(c.Sender.Recipient(), "<-", e.GetText())
	if e.GetLanguage(c) == "ua" {
		e.SetCaption(c, startMsg_ua)
	} else {
		e.SetCaption(c, startMsg_en)
	}
	return menu.Back
}

func volunteerMenu(e *menu.Node, c *tb.Callback) int {
	if e.GetLanguage(c) != "ua" {
		e.SetCaption(c, "🚫 Sorry, at the moment, we need only Ukrainian-speaking volunteers")
		return menu.Back
	}

	return forward(e, c)
}

func main() {
	debug := flag.Bool("debug", false, "Use test bot token")
	flag.Parse()
	var token string
	if *debug {
		token = os.Getenv("TEST_HELP_UA_NL_BOT_TOKEN")
	} else {
		token = os.Getenv("HELP_UA_NL_BOT_TOKEN")
	}
	Run(token)
}
