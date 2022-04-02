package etf2l

import "github.com/PuerkitoBio/goquery"

func (c Client) ResolvePlayerSteamID(playerURL string) (string, error) {
	playerPage, err := c.getHtml(playerURL)
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
