package app

import (
	"log"

	"github.com/bvisness/SQLJam/node"
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Before drawing. Sort out your data, set your layout info, etc.
func doNodeUpdate(n *node.Node) {
	switch d := n.Data.(type) {
	case *node.Table:
		doTableUpdate(n, d)
	case *node.Filter:
		doFilterUpdate(n, d)
	case *node.PickColumns:
		doPickColumnsUpdate(n, d)
	}
}

// Drawing and user input.
func doNodeUI(n *node.Node) {
	switch d := n.Data.(type) {
	case *node.Table:
		doTableUI(n, d)
	case *node.Filter:
		doFilterUI(n, d)
	case *node.PickColumns:
		doPickColumnsUI(n, d)
	}
}

func doTableUpdate(n *node.Node, t *node.Table) {
	// init dropdown
	if len(t.TableDropdown.GetOptions()) == 0 {
		updateTableDropdown(&t.TableDropdown)
	}

	n.UISize = rl.Vector2{200, 24}
}

func updateTableDropdown(dropdown *raygui.DropdownEx) {
	rows, err := db.Query(`
		SELECT name
		FROM sqlite_master
		WHERE
			type = 'table'
			AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`)
	if err != nil {
		log.Print(err)
		dropdown.SetOptions(raygui.DropdownExOption{"ERROR", nil})
		return
	}
	defer rows.Close()

	var opts []raygui.DropdownExOption
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Print(err)
			dropdown.SetOptions(raygui.DropdownExOption{"ERROR", nil})
			return
		}
		opts = append(opts, raygui.DropdownExOption{
			Name:  name,
			Value: name,
		})
	}

	err = rows.Err()
	if err != nil {
		log.Print(err)
		dropdown.SetOptions(raygui.DropdownExOption{"ERROR", nil})
		return
	}

	dropdown.SetOptions(opts...)
}

func doTableUI(n *node.Node, t *node.Table) {
	if ival := t.TableDropdown.Do(n.UIRect); ival != nil {
		t.Table = ival.(string)
	}
}

func doFilterUpdate(n *node.Node, f *node.Filter) {
	n.UISize = rl.Vector2{360, 24}
}

func doFilterUI(n *node.Node, f *node.Filter) {
	rl.DrawRectangleRec(n.UIRect, rl.White)
	f.Conditions = f.TextBox.Do(n.UIRect, f.Conditions, 100)
}

const pickColumnsFieldHeight = 24
const pickColumnsFieldSpacing = 4

func doPickColumnsUpdate(n *node.Node, p *node.PickColumns) {
	uiHeight := 0
	for range p.Cols {
		uiHeight += pickColumnsFieldHeight
		uiHeight += pickColumnsFieldSpacing
	}
	uiHeight += pickColumnsFieldHeight // for buttons

	n.UISize = rl.Vector2{300, float32(uiHeight)}

	// Get schema (if necessary) and update UI data
	// TODO: don't do this every frame plz

	// This will obliterate existing selections on resize,
	// but this shouldn't happen anyway if we're resizing correctly.
	if len(p.Cols) != len(p.ColDropdowns) {
		p.ColDropdowns = make([]raygui.DropdownEx, len(p.Cols))
	}

	if n.Inputs[0] != nil {
		var opts []raygui.DropdownExOption

		schemaCols, err := getSchema(n.Inputs[0])
		if err == nil {
			for _, col := range schemaCols {
				opts = append(opts, raygui.DropdownExOption{
					Name:  col,
					Value: col,
				})
			}
		} else {
			log.Print(err)
			opts = append(opts, raygui.DropdownExOption{"ERROR", "ERROR"})
		}

		for i := range p.ColDropdowns {
			dropdown := &p.ColDropdowns[i]
			dropdown.SetOptions(opts...)
		}
	}
}

func doPickColumnsUI(n *node.Node, p *node.PickColumns) {
	// TODO: Could we, like, do a closure or something
	var deferredDropdown *raygui.DropdownEx
	var deferredDropdownI int
	var deferredDropdownY float32

	fieldY := n.UIRect.Y
	for i := range p.ColDropdowns {
		dropdown := &p.ColDropdowns[i]
		if dropdown.Open && deferredDropdown == nil {
			deferredDropdown = dropdown
			deferredDropdownY = fieldY
		} else {
			col := dropdown.Do(rl.Rectangle{n.UIRect.X, fieldY, n.UIRect.Width, pickColumnsFieldHeight})
			p.Cols[i], _ = col.(string)
		}
		fieldY += pickColumnsFieldHeight + pickColumnsFieldSpacing
	}

	if raygui.Button(rl.Rectangle{
		n.UIRect.X,
		fieldY,
		n.UIRect.Width/2 - pickColumnsFieldSpacing/2,
		pickColumnsFieldHeight,
	}, "+") {
		p.Cols = append(p.Cols, "")
		p.ColDropdowns = append(p.ColDropdowns, raygui.DropdownEx{})
	}
	if raygui.Button(rl.Rectangle{
		n.UIRect.X + n.UIRect.Width/2 + pickColumnsFieldSpacing/2,
		fieldY,
		n.UIRect.Width/2 - pickColumnsFieldSpacing/2,
		pickColumnsFieldHeight,
	}, "-") {
		if len(p.Cols) > 1 {
			p.Cols = p.Cols[:len(p.Cols)-1]
			p.ColDropdowns = p.ColDropdowns[:len(p.ColDropdowns)-1]
		}
	}

	if deferredDropdown != nil {
		col := deferredDropdown.Do(rl.Rectangle{n.UIRect.X, deferredDropdownY, n.UIRect.Width, pickColumnsFieldHeight})
		p.Cols[deferredDropdownI], _ = col.(string)
	}
}

func getSchema(n *node.Node) ([]string, error) {
	rows, err := db.Query(n.GenerateSql() + " LIMIT 0") // TODO: The limit should be part of SQL generation, yeah?
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rows.Columns()
}
