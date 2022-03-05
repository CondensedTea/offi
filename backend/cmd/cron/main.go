package main

import (
	"context"
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
	log "github.com/sirupsen/logrus"
	database "offi/pkg/db"
	"offi/pkg/etf2l"
	"strings"
)

var k = koanf.New(".")

func main() {
	ctx := context.Background()

	err := k.Load(env.Provider("OFFI_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "OFFI_")), "_", ".", -1)
	}), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(k.String("db.dsn"))

	db, err := database.New(ctx, k.String("db.dsn"))
	if err != nil {
		panic(err)
	}

	scrapper, err := etf2l.New()
	if err != nil {
		panic(err)
	}

	comps, err := db.GetCompetitions(ctx)
	if err != nil {
		log.Errorf("failed to get competitions: %v", err)
	}

	for _, c := range comps {
		matches, err := scrapper.GetMatches(c.ID, c.LastMatchID)
		if err != nil {
			log.Errorf("failed to get matches: %v", err)

		}
		fmt.Println(matches)
	}
}
