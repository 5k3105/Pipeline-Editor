package graph //meta

import (
	"time"

	"github.com/emirpasic/gods/maps/treemap"
)

type Meta struct {
	Type *treemap.Map
	Data *treemap.Map
}

type MetaType struct {
	Type    string
	Version string
	Author  string
	Date    time.Time
}

func NewMeta() *Meta {
	var m = &Meta{
		Type: treemap.NewWithIntComparator(),
		Data: treemap.NewWithIntComparator(),
	}

	return m
}

func (m *Meta) AddMeta(id int, i interface{}, t string, v string, a string, yr, mon, day int) {

	var mt = &MetaType{
		Type:    t,
		Version: v,
		Author:  a,
		Date:    time.Time{},
	}

	mt.Date.AddDate(yr, mon, day)

	m.Type.Put(id, mt)

	m.Data.Put(id, i)

}
