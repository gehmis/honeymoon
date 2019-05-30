package bort

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gehmis/honeymoon/pkg/eve"
	"github.com/headzoo/surf"
)

const (
	moonListUrl = "https://docs.google.com/spreadsheets/d/e/2PACX-1vTKHcuQ6KtNEbsaSda0njYbaU324p1_OuA-7rSgUCcN2384tThjQNi5EG64fNZ7Y_xy_RDOl56F-nld/pubhtml?gid=229192924&single=true"
)

var (
	compositionMatcher = regexp.MustCompile(`([A-Za-z]{3}):(\d+)%`)
)

type Moon struct {
	Name        string
	Location    string
	Cycle       time.Time
	Composition map[int64]int
}

func MoonList() ([]Moon, error) {
	ms := []Moon{}

	os, err := eve.OreList()
	if err != nil {
		return nil, err
	}

	ores := map[string]int64{}

	for _, o := range os {
		ores[o.Name[0:3]] = o.ID
	}

	b := surf.NewBrowser()

	if err := b.Open(moonListUrl); err != nil {
		return nil, err
	}

	rows := b.Find("tr")

	rows.Each(func(i int, s *goquery.Selection) {
		if i < 2 {
			return
		}

		m := Moon{
			Name:        strings.TrimSpace(s.Find("td:nth-child(2)").Text()),
			Location:    strings.TrimSpace(s.Find("td:nth-child(3)").Text()),
			Composition: map[int64]int{},
		}

		cms := compositionMatcher.FindAllStringSubmatch(s.Find("td:nth-child(5)").Text(), -1)

		for _, cm := range cms {
			cmi, err := strconv.Atoi(cm[2])
			if err != nil {
				return
			}

			m.Composition[ores[cm[1]]] = cmi
		}

		cycle := fmt.Sprintf("%s-%d %s", strings.TrimSpace(s.Find("td:nth-child(6)").Text()), time.Now().UTC().Year(), strings.TrimSpace(s.Find("td:nth-child(7)").Text()))

		cd, err := time.Parse("02-Jan-2006 15:04", cycle)
		if err != nil {
			return
		}

		m.Cycle = cd

		if m.Name == "" {
			return
		}

		ms = append(ms, m)
	})

	sort.Slice(ms, func(i, j int) bool { return ms[i].Cycle.Before(ms[j].Cycle) })

	return ms, nil
}
