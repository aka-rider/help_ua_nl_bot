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
	ua = "ua"
	en = "en"

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

	once                = "once"
	helpFinance         = "finance"
	bankTransfer        = "bank_transfer"
	onlinePayment       = "online_payment"
	helpHumanitarian    = "humanitarian"
	clothes             = "clothes"
	collection          = "collection"
	essentials          = "essentials"
	medicationAndDosage = "medication_and_dosage"
	medicalConsumables  = "medical_consumables"
	medicalEquipment    = "medical_equipment"
	volunteer           = "volunteer"
	refugee             = "refugee"
	hotline             = "hotline"
	language            = "language"
	helpOther           = "other"
	backLabel           = "back"
)

func Run(token string) {
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: defaultPollFreq},
	})
	if err != nil {
		panic(err)
	}

	if err := tr.Init("lang", ua); err != nil {
		panic(err)
	}

	flow, err := menu.NewMenuFlow("flow", bot, tr.DefaultEngine)
	if err != nil {
		panic(err)
	}

	flow.GetRoot().
		AddWith(once, forward,
			flow.NewNode(helpFinance, forward).
				Add(bankTransfer, func(e *menu.Node, c *tb.Callback) int {
					e.SetCaption(c, bankDetails)
					return menu.Forward
				}).
				Add(onlinePayment, forward).
				Add(backLabel, back),
			flow.NewNode(helpHumanitarian, forward),
			flow.NewNode(helpOther, func(e *menu.Node, c *tb.Callback) int {
				e.SetCaption(c, "https://forms.gle/fDZHTH7pnPCNUfJdA")
				return menu.Forward
			}),
			flow.NewNode(backLabel, back),
		).
		Add(volunteer, forward).
		Add(refugee, func(e *menu.Node, c *tb.Callback) int {
			e.SetCaption(c, "https://help-ukraine.nl/refugee")
			return menu.Forward
		}).
		Add(hotline, forward).
		Add(language, switchLanguage)

	nodeUrls := NodeKeyUrls{
		{hotline, "https://t.me/ukrainians_nl_support_bot"},
		{volunteer, "https://help-ukraine.nl/volunteers"},
		{onlinePayment, "https://useplink.com/payment/yIXtkzDtTlBGUBxiRQVaA/"},
		{helpHumanitarian, "https://help-ukraine.nl/humanitarian-aid"},
	}.toNodeUrls(flow)

	for _, locale := range locales() {
		flow.Build(locale)
	}

	nodeUrls.addUrls()

	bot.Handle("/start", func(m *tb.Message) {
		err = flow.Start(m.Sender, greetings_ua, ua)
		if err != nil {
			log.Println("failed to display the menu:", err)
		}
	})

	log.Println("running", bot.Me.Username, "...")
	bot.Start()
}

func locales() []string {
	return []string{ua, en}
}

func switchLanguage(e *menu.Node, c *tb.Callback) int {
	log.Println(c.Sender.Recipient(), "press", e.GetText())
	if e.GetLanguage(c) == en {
		e.SetLanguage(c, ua)
		e.SetCaption(c, greetings_ua)
	} else {
		e.SetLanguage(c, en)
		e.SetCaption(c, greetings_en)
	}
	return menu.Forward
}

func forward(e *menu.Node, c *tb.Callback) int {
	log.Println(c.Sender.Recipient(), "->", e.GetText())
	if e.GetLanguage(c) == ua {
		e.SetCaption(c, startMsg_ua)
	} else {
		e.SetCaption(c, startMsg_en)
	}
	return menu.Forward
}

func back(e *menu.Node, c *tb.Callback) int {
	log.Println(c.Sender.Recipient(), "<-", e.GetText())
	if e.GetLanguage(c) == ua {
		e.SetCaption(c, startMsg_ua)
	} else {
		e.SetCaption(c, startMsg_en)
	}
	return menu.Back
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
