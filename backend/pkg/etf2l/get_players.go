package etf2l

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var ErrMatchIsNotComplete = errors.New("match is not completed")

type Match struct {
	Players  []string
	Maps     []string
	PlayedAt time.Time
}

func (client ETF2L) ParseMatchPage(matchId int) (*Match, error) {
	url := fmt.Sprintf("https://etf2l.org/matches/%d/", matchId)

	matchPage, err := client.getHtml(url)
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

	doc.Find("span.winr").Each(func(i int, selection *goquery.Selection) {
		playerURL, _ := selection.Find("a").Attr("href")
		playerURLs = append(playerURLs, playerURL)
	})

	doc.Find("span.looser").Each(func(i int, selection *goquery.Selection) {
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

	matchDate := doc.Find("span.date").Get(5).FirstChild.Data
	if matchDate == "" {
		return nil, ErrMatchIsNotComplete
	}

	playedAt, err := time.Parse("2 Jan 2006", matchDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %v", err)
	}

	return &Match{
		Players:  playerURLs,
		Maps:     matchMaps,
		PlayedAt: playedAt,
	}, nil
}

func (client ETF2L) ResolvePlayerSteamID(playerID string) (string, error) {
	playerPage, err := client.getHtml(playerID)
	if err != nil {
		return "", err
	}
	defer playerPage.Close()

	playerDoc, err := goquery.NewDocumentFromReader(playerPage)
	if err != nil {
		return "", err
	}

	steamID, err := playerDoc.Find("table.playerinfo").Find("tbody").Find("a").Html()
	if err != nil {
		return "", err
	}
	return steamID, nil
}

func (client ETF2L) getHtml(url string) (io.ReadCloser, error) {
	resp, err := client.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	return resp.Body, nil
}
