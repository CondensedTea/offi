package etf2l

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s ETF2L) GetMatches(competitionId int, lastParsedMatch *int) ([]Result, error) {
	isLastPage := false
	matches := make([]Result, 0)

	url := etf2lApiURL + fmt.Sprintf("competition/%d/results.json", competitionId)

	for !isLastPage {
		cr, err := s.getCompetitionResults(url)
		if err != nil {
			return nil, fmt.Errorf("failed to get results for competition %d: %v", competitionId, err)
		}

		if cr.Page.NextPageUrl == "" {
			isLastPage = true
		} else {
			url = cr.Page.NextPageUrl
		}

		for _, result := range cr.Results {
			if lastParsedMatch != nil && result.Id == *lastParsedMatch {
				return matches, nil
			}
			matches = append(matches, result)
		}
	}
	return matches, nil
}

func (s ETF2L) getCompetitionResults(url string) (*CompetitionResults, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned bad staus: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var cr CompetitionResults

	if err = json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return nil, err
	}
	return &cr, nil
}
