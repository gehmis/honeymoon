package eve

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Material struct {
	Type
	Quantity int64
}

type Ore struct {
	Type
	Materials []Material
}

type Type struct {
	ID     int64
	Name   string
	Volume float64
}

func MaterialList(id int64) ([]Material, error) {
	ms := []Material{}

	materials, err := query("select * from invTypeMaterials where typeID = %d", id)
	if err != nil {
		return nil, err
	}

	for _, material := range materials {
		id := material["materialTypeID"].(int64)

		t, err := TypeGet(id)
		if err != nil {
			return nil, err
		}

		ms = append(ms, Material{
			Type:     *t,
			Quantity: material["quantity"].(int64),
		})
	}

	return ms, nil
}

func OreList() ([]Ore, error) {
	ts := []Type{}

	ots, err := typesByParentMarketGroup(54)
	if err != nil {
		return nil, err
	}

	ts = append(ts, ots...)

	gts, err := typesByParentMarketGroup(2395)
	if err != nil {
		return nil, err
	}

	ts = append(ts, gts...)

	os := []Ore{}

	for _, t := range ts {
		ms, err := MaterialList(t.ID)
		if err != nil {
			return nil, err
		}

		os = append(os, Ore{
			Type:      t,
			Materials: ms,
		})
	}

	return os, nil
}

func TypeGet(id int64) (*Type, error) {
	types, err := query("select * from invTypes where typeID = %d", id)
	if err != nil {
		return nil, err
	}

	if len(types) != 1 {
		return nil, fmt.Errorf("could not find type for id: %d", id)
	}

	t := typeLoad(types[0])

	return &t, nil
}

func typesByParentMarketGroup(id int64) ([]Type, error) {
	ts := []Type{}

	groups, err := query("select * from invMarketGroups where parentGroupId = %d", id)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		types, err := query("select * from invTypes where marketGroupID = %d", group["marketGroupID"])
		if err != nil {
			return nil, err
		}

		for _, typ := range types {
			ts = append(ts, typeLoad(typ))
		}
	}

	return ts, nil
}

func typeLoad(r row) Type {
	return Type{
		ID:     r["typeID"].(int64),
		Name:   r["typeName"].(string),
		Volume: r["volume"].(float64),
	}
}
