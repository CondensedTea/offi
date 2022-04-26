package etf2l

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nleeper/goment"
	"golang.org/x/net/html"
)

var timeNow = func() time.Time {
	return time.Now().UTC()
}

var reSubmittedDateTime = regexp.MustCompile(`Results submitted: (\d{1,2} \w+? \d{4}|Yesterday|Today), (\d{2}):(\d{2})`)

type Match struct {
	Players  []string
	Maps     []string
	PlayedAt time.Time

	ID          int
	Competition string
	Tier        string
	Stage       string
}

func (c Client) ParseMatchPage(matchId int) (*Match, error) {
	url := fmt.Sprintf("https://etf2l.org/matches/%d/", matchId)

	matchPage, err := c.getHtml(url)
	if err != nil {
		return nil, err
	}
	defer matchPage.Close()

	doc, err := goquery.NewDocumentFromReader(matchPage)
	if err != nil {
		return nil, err
	}

	var (
		playerURLs []string
		matchMaps  []string
	)

	if post := doc.Find("div.post").Find("p").Text(); post == "Invalid Match ID specified." {
		return nil, fmt.Errorf("invalid match ID")
	}

	doc.
		Find(".fix.match-players").
		Find("span.winr, span.looser").
		Each(func(i int, selection *goquery.Selection) {
			playerURL, _ := selection.Find("a").Attr("href")
			playerURLs = append(playerURLs, playerURL)
		})

	doc.Find("div.maps").Each(func(i int, selection *goquery.Selection) {
		selection.
			Find("div.map").
			Find("h2").Each(func(i int, selection *goquery.Selection) {
			matchMaps = append(matchMaps, selection.Text())
		})
	})

	competition := doc.Find("h1.c").Text()

	var stage string

	doc.Find("h3.c").Each(func(i int, selection *goquery.Selection) {
		switch i {
		case 0:
			stage = selection.Text()
		case 1:
			stage += " " + selection.Text()
		}
	})

	if len(doc.Nodes) < 3 {
		return nil, fmt.Errorf("too little nodes found in doc")
	}
	node := doc.Find("h4.c").Get(2)

	matchDate, err := parseMatchDate(node)
	if err != nil {
		return nil, fmt.Errorf("failed to parse match date: %v", err)
	}

	return &Match{
		Players:     playerURLs,
		Maps:        matchMaps,
		PlayedAt:    matchDate,
		ID:          matchId,
		Competition: competition,
		Stage:       stage,
	}, nil
}

func parseMatchDate(node *html.Node) (time.Time, error) {
	match := reSubmittedDateTime.FindStringSubmatch(goquery.NewDocumentFromNode(node).Text())
	if len(match) < 1 {
		return time.Time{}, fmt.Errorf("could not find correct date")
	}

	var (
		gm  *goment.Goment
		err error
	)

	switch match[1] {
	case "Today":
		year, month, day := timeNow().Date()
		gm, err = goment.New(goment.DateTime{Year: year, Month: int(month), Day: day})
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to apply today's date: %v", err)
		}
	case "Yesterday":
		year, month, day := timeNow().AddDate(0, 0, -1).Date()
		gm, err = goment.New(goment.DateTime{Year: year, Month: int(month), Day: day})
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to apply yesterday's date: %v", err)
		}
	default:
		gm, err = goment.New(match[1], "DD MMM YYYY")
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to parse date: %v", err)
		}
	}

	hour, err := strconv.Atoi(match[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse hours: %v", err)
	}
	gm.Set("hour", hour)

	minute, err := strconv.Atoi(match[3])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse minutes: %v", err)
	}
	gm.Set("minute", minute)

	return gm.ToTime(), nil
}

func (c Client) getHtml(url string) (io.ReadCloser, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	return resp.Body, nil
}
