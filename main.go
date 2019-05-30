package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/gehmis/honeymoon/pkg/bort"
	"github.com/gehmis/honeymoon/pkg/cerlestes"
	"github.com/gehmis/honeymoon/pkg/eve"
)

type material struct {
	name    string
	price   float64
	percent int
}

var (
	compositionMatcher = regexp.MustCompile(`([A-Za-z]{3}):(\d+)%`)
	moonMatcher        = regexp.MustCompile(`(.*?) \[([^]]*)\]`)
	totalOre           = float64(14400000)
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
}

func run() error {
	ms, err := bort.MoonList()
	if err != nil {
		return err
	}

	ps, err := cerlestes.PriceList()
	if err != nil {
		return err
	}

	os, err := eve.OreList()
	if err != nil {
		return err
	}

	ores := map[int64]eve.Ore{}

	for _, o := range os {
		ores[o.ID] = o
	}

	for _, m := range ms {
		mats := []material{}

		for id, pct := range m.Composition {

			o := ores[id]

			var price float64

			for _, m := range o.Materials {
				price += float64(m.Quantity) * ps[m.ID]
			}

			price /= o.Volume

			mats = append(mats, material{
				name:    o.Name[0:3],
				price:   price,
				percent: pct,
			})
		}

		sort.Slice(mats, func(i, j int) bool { return mats[i].price > mats[j].price })

		var total float64
		cvs := []string{}

		for _, m := range mats {
			total += float64(m.percent) * m.price
			cvs = append(cvs, fmt.Sprintf("%s:%0.1f(%d%%)", m.name, m.price, m.percent))
		}

		fmt.Printf("%-45s  %-20s  %-20s  %-6.1f  ", m.Name, m.Location, m.Cycle.Format("Jan 02 15:04"), total)

		for _, cv := range cvs {
			fmt.Printf("%-15s ", cv)
		}

		fmt.Println()
	}

	return nil
}
