package news

import (
	"fmt"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/jspc/jshell/apps"
	"github.com/mmcdole/gofeed"
)

var (
	source = "https://feeds.bbci.co.uk/news/rss.xml?edition=uk"

	header = ` o    o                    8 8  o
 8    8                    8 8
o8oooo8 .oPYo. .oPYo. .oPYo8 8 o8 odYo. .oPYo. .oPYo.
 8    8 8oooo8 .oooo8 8    8 8  8 8' '8 8oooo8 Yb..
 8    8 8.     8    8 8    8 8  8 8   8 8.       'Yb.
 8    8 'Yooo' 'YooP8 'YooP' 8  8 8   8 'Yooo' 'YooP'
:..:::..:.....::.....::.....:..:....::..:.....::.....:
::::::::::::::::::::::::::::::::::::::::::::::::::::::
::::::::::::::::::::::::::::::::::::::::::::::::::::::


`
)

type headline struct {
	Title     string
	Published time.Time
	Url       string
}

type News struct{}

func (News) Name() string        { return "Headlines" }
func (News) Description() string { return "🗞️  Latest headlines to your screen!" }
func (News) Cleanup() error      { return nil }

func (News) Run() (err error) {
	headlines := make([]headline, 0)

	p := gofeed.NewParser()
	feed, err := p.ParseURL(source)
	if err != nil {
		return
	}

	for idx, item := range feed.Items {
		if idx >= 10 {
			break
		}

		headlines = append(headlines, headline{
			Title:     item.Title,
			Published: *item.PublishedParsed,
			Url:       item.Link,
		})
	}

	headlines = dedupeHeadlines(headlines)

	sort.Slice(headlines, func(i, j int) bool {
		return headlines[j].Published.Before(headlines[i].Published)
	})

	fmt.Print("\033[H\033[2J")
	color.Cyan(header)

	for _, h := range headlines {
		fmt.Printf("%s - %s (%s)\n",
			color.HiMagentaString(apps.FormatDateAndTime(h.Published)),
			color.HiGreenString(h.Title),
			color.HiBlueString(h.Url),
		)
	}

	fmt.Println()
	fmt.Println()

	time.Sleep(time.Second)

	fmt.Println("press enter to return to the main menu")

	//#nosec
	fmt.Scanln()

	return
}

func dedupeHeadlines(headlines []headline) (deduped []headline) {
	deduped = make([]headline, 0)

	tmp := make(map[string]headline)
	for _, hl := range headlines {
		tmp[hl.Title] = hl
	}

	for _, hl := range tmp {
		deduped = append(deduped, hl)
	}

	return
}
