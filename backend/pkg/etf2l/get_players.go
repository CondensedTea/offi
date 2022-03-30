package etf2l

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var ErrMatchIsNotCompleted = errors.New("match is not completed")

type Match struct {
	Players  []string
	Maps     []string
	PlayedAt time.Time
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

	// <time class="brief">
	//    <span> YYYY-MM-DD </span>
	//    <span> hh:mm </span>
	// </time>
	// Result: "YYYY-MM-DD hh:mm"
	matchDateNode := doc.Find("time.brief").Get(6)
	matchDate := matchDateNode.FirstChild.FirstChild.Data + " " + matchDateNode.LastChild.FirstChild.Data

	playedAt, err := parseMatchDate(matchDate)
	if err != nil {
		return nil, err
	}

	return &Match{
		Players:  playerURLs,
		Maps:     matchMaps,
		PlayedAt: playedAt,
	}, nil
}

func (c Client) ResolvePlayerSteamID(playerID string) (string, error) {
	playerPage, err := c.getHtml(playerID)
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

func parseMatchDate(matchDate string) (time.Time, error) {
	if matchDate == "" {
		return time.Time{}, ErrMatchIsNotCompleted
	}
	if matchDate == "Yesterday" {
		return time.Now().Add(-24 * time.Hour), nil
	}
	playedAt, err := time.Parse("2 Jan 2006 15:04", matchDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %v", err)
	}
	return playedAt, nil
}
