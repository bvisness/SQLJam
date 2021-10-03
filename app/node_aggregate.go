package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Aggregate struct {
	Aggregates []*AggregateColumn
	GroupBys   []*AggregateGroupBy
}

type AggregateColumn struct {
	Type         AggregateType
	Col          string
	Alias        string
	TypeDropdown raygui.DropdownEx
	ColDropdown  raygui.DropdownEx
	AliasTextbox raygui.TextBoxEx
}

type AggregateGroupBy struct {
	Col         string
	ColDropdown raygui.DropdownEx
}

type AggregateType int

const (
	Avg AggregateType = iota + 1
	Max
	Min
	Sum
	Count
	CountDistinct
)

func NewAggregate() *Node {
	return &Node{
		Title:   "Aggregate",
		CanSnap: true,
		Color:   rl.NewColor(182,213,60, 255),
		Inputs:  make([]*Node, 1),
		Data: &Aggregate{
			Aggregates: []*AggregateColumn{{}},
		},
	}
}

func (d *Aggregate) AllDropdowns() []*raygui.DropdownEx {
	res := make([]*raygui.DropdownEx, 0, 2*len(d.Aggregates)+len(d.GroupBys))
	for _, agg := range d.Aggregates {
		res = append(res, &agg.TypeDropdown)
		res = append(res, &agg.ColDropdown)
	}
	for _, gb := range d.GroupBys {
		res = append(res, &gb.ColDropdown)
	}
	return res
}

var aggregateTypeOpts = []raygui.DropdownExOption{
	{"AVG", Avg},
	{"MAX", Max},
	{"MIN", Min},
	{"SUM", Sum},
	{"COUNT", Count},
	{"COUNT DISTINCT", CountDistinct},
}

func (d *Aggregate) Update(n *Node) {
	height := 0

	// Aggregated columns
	for range d.Aggregates {
		height += UIFieldHeight + UIFieldSpacing
	}
	height += UIFieldHeight + UIFieldSpacing // for +/- buttons

	// Group bys
	height += UIFieldHeight + UIFieldSpacing // for "Group by" label
	for range d.GroupBys {
		height += UIFieldHeight + UIFieldSpacing // for group by rows
	}
	height += UIFieldHeight // for +/- buttons

	n.UISize = rl.Vector2{480, float32(height)}

	colOpts := columnNameDropdownOpts(n.Inputs[0])
	for _, agg := range d.Aggregates {
		agg.TypeDropdown.SetOptions(aggregateTypeOpts...)
		agg.ColDropdown.SetOptions(colOpts...)
	}
	for _, gb := range d.GroupBys {
		gb.ColDropdown.SetOptions(colOpts...)
	}
}

func (d *Aggregate) DoUI(n *Node) {
	const typeWidth = 140

	openDropdown, isOpen := raygui.GetOpenDropdown(d.AllDropdowns())
	if isOpen {
		raygui.Disable()
		defer raygui.Enable()
	}

	// I love rendering bottom to top!!
	fieldY := n.UIRect.Y + n.UIRect.Height

	// Group by
	{
		fieldY -= UIFieldHeight
		if raygui.Button(rl.Rectangle{
			n.UIRect.X,
			fieldY,
			n.UIRect.Width/2 - UIFieldSpacing/2,
			UIFieldHeight,
		}, "+") {
			d.GroupBys = append(d.GroupBys, &AggregateGroupBy{})
		}
		if raygui.Button(rl.Rectangle{
			n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
			fieldY,
			n.UIRect.Width/2 - UIFieldSpacing/2,
			UIFieldHeight,
		}, "-") {
			if len(d.GroupBys) > 0 {
				d.GroupBys = d.GroupBys[:len(d.GroupBys)-1]
			}
		}

		for i := len(d.GroupBys) - 1; i >= 0; i-- {
			func() {
				gb := d.GroupBys[i]

				if openDropdown == &gb.ColDropdown {
					raygui.Enable()
					defer raygui.Disable()
				}

				fieldY -= UIFieldSpacing + UIFieldHeight
				icol := gb.ColDropdown.Do(rl.Rectangle{n.UIRect.X, fieldY, n.UIRect.Width, UIFieldHeight})
				gb.Col, _ = icol.(string)
			}()
		}

		fieldY -= UIFieldSpacing + UIFieldHeight
		const textSize = 20
		drawBasicText("Group by", n.UIRect.X, fieldY+(UIFieldHeight-textSize), textSize, rl.Black)
	}

	// Aggregate columns
	{
		fieldY -= UIFieldSpacing + UIFieldHeight
		if raygui.Button(rl.Rectangle{
			n.UIRect.X,
			fieldY,
			n.UIRect.Width/2 - UIFieldSpacing/2,
			UIFieldHeight,
		}, "+") {
			d.Aggregates = append(d.Aggregates, &AggregateColumn{})
		}
		if raygui.Button(rl.Rectangle{
			n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
			fieldY,
			n.UIRect.Width/2 - UIFieldSpacing/2,
			UIFieldHeight,
		}, "-") {
			if len(d.Aggregates) > 1 {
				d.Aggregates = d.Aggregates[:len(d.Aggregates)-1]
			}
		}

		for i := len(d.Aggregates) - 1; i >= 0; i-- {
			func() {
				agg := d.Aggregates[i]

				if openDropdown == &agg.TypeDropdown || openDropdown == &agg.ColDropdown {
					raygui.Enable()
					defer raygui.Disable()
				}

				fieldY -= UIFieldSpacing + UIFieldHeight
				fieldX := n.UIRect.X

				itype := agg.TypeDropdown.Do(rl.Rectangle{fieldX, fieldY, typeWidth, UIFieldHeight})
				agg.Type, _ = itype.(AggregateType)
				fieldX += typeWidth + UIFieldSpacing

				remainingWidth := n.UIRect.X + n.UIRect.Width - fieldX
				fieldWidth := remainingWidth/2 - UIFieldSpacing/2

				icol := agg.ColDropdown.Do(rl.Rectangle{fieldX, fieldY, fieldWidth, UIFieldHeight})
				agg.Col, _ = icol.(string)
				fieldX += fieldWidth + UIFieldSpacing

				aliasRect := rl.Rectangle{fieldX, fieldY, fieldWidth, UIFieldHeight}
				agg.Alias, _ = agg.AliasTextbox.Do(aliasRect, agg.Alias, 100)
			}()
		}
	}
}

func (d *Aggregate) Serialize() string {
	res := ""
	for _, agg := range d.Aggregates {
		res += fmt.Sprintf("%v", agg.Type)
		res += agg.Col
		res += agg.Alias
	}
	for _, gb := range d.GroupBys {
		res += gb.Col
	}
	return res
}
