package etf2l

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
)

func (s ETF2L) GetPlayers(matchId int) ([]string, error) {
	url := fmt.Sprintf("https://etf2l.org/matches/%d/", matchId)

	matchPage, err := s.getHtml(url)
	if err != nil {
		return nil, err
	}
	defer matchPage.Close()

	doc, err := goquery.NewDocumentFromReader(matchPage)
	if err != nil {
		return nil, err
	}

	playerUrls := make([]string, 0)

	var steamIDs []string

	doc.Find("span.winr").Each(func(i int, selection *goquery.Selection) {
		playerUrl, _ := selection.Find("a").Attr("href")
		playerUrls = append(playerUrls, playerUrl)
	})

	for _, player := range playerUrls {
		playerPage, err := s.getHtml(player)
		defer playerPage.Close()

		if err != nil {
			return nil, err
		}
		playerDoc, err := goquery.NewDocumentFromReader(playerPage)
		if err != nil {
			return nil, err
		}

		playerId, _ := playerDoc.Find("table.playerinfo").Find("tbody").Find("a").Html()
		if playerId != "" {
			steamIDs = append(steamIDs, playerId)
		}
	}
	return steamIDs, nil
}

func (s ETF2L) getHtml(url string) (io.ReadCloser, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	return resp.Body, nil
}
