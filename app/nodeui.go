package app

import (
	"fmt"
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
	case *node.CombineRows:
		doCombineRowsUpdate(n, d)
	case *node.Order:
		doOrderUpdate(n, d)
	case *node.Join:
		doJoinUpdate(n, d)
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
	case *node.CombineRows:
		doCombineRowsUI(n, d)
	case *node.Order:
		doOrderUI(n, d)
	case *node.Join:
		doJoinUI(n, d)
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

const UIFieldHeight = 24
const UIFieldSpacing = 4

func doPickColumnsUpdate(n *node.Node, p *node.PickColumns) {
	uiHeight := 0
	for range p.Entries {
		uiHeight += UIFieldHeight
		uiHeight += UIFieldSpacing
	}
	uiHeight += UIFieldHeight // for buttons

	n.UISize = rl.Vector2{300, float32(uiHeight)}

	opts := columnNameDropdownOpts(n.Inputs[0])
	for _, entry := range p.Entries {
		entry.ColDropdown.SetOptions(opts...)
	}
}

func doPickColumnsUI(n *node.Node, p *node.PickColumns) {
	openDropdown, isOpen := raygui.GetOpenDropdown(p.ColDropdowns())
	if isOpen {
		raygui.Disable()
		defer raygui.Enable()
	}

	// Render bottom to top to avoid overlap issues with dropdowns

	fieldY := n.UIRect.Y + n.UIRect.Height - UIFieldHeight
	if raygui.Button(rl.Rectangle{
		n.UIRect.X,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "+") {
		p.Entries = append(p.Entries, &node.PickColumnsEntry{})
	}
	if raygui.Button(rl.Rectangle{
		n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "-") {
		if len(p.Entries) > 1 {
			p.Entries = p.Entries[:len(p.Entries)-1]
		}
	}

	for i := len(p.Entries) - 1; i >= 0; i-- {
		fieldY -= UIFieldSpacing + UIFieldHeight
		func() {
			entry := p.Entries[i]
			if openDropdown == &entry.ColDropdown {
				raygui.Enable()
				defer raygui.Disable()
			}

			col := entry.ColDropdown.Do(rl.Rectangle{
				n.UIRect.X,
				fieldY,
				n.UIRect.Width/2 - UIFieldSpacing/2,
				UIFieldHeight,
			})
			entry.Col, _ = col.(string)

			aliasRect := rl.Rectangle{
				n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
				fieldY,
				n.UIRect.Width/2 - UIFieldSpacing/2,
				UIFieldHeight,
			}
			rl.DrawRectangleRec(aliasRect, rl.White)
			entry.Alias = entry.AliasTextbox.Do(aliasRect, entry.Alias, 100)
		}()
	}
}

const orderDirectionWidth = 70

func doOrderUpdate(n *node.Node, o *node.Order) {
	uiHeight := 0
	for range o.Cols {
		uiHeight += UIFieldHeight
		uiHeight += UIFieldSpacing
	}
	uiHeight += UIFieldHeight // for buttons

	n.UISize = rl.Vector2{300, float32(uiHeight)}

	// This will obliterate existing selections on resize,
	// but this shouldn't happen anyway if we're resizing correctly.
	if len(o.Cols) != len(o.ColDropdowns) {
		o.ColDropdowns = make([]*raygui.DropdownEx, len(o.Cols))
	}

	opts := columnNameDropdownOpts(n.Inputs[0])
	for _, dropdown := range o.ColDropdowns {
		dropdown.SetOptions(opts...)
	}
}

func doOrderUI(n *node.Node, o *node.Order) {
	openDropdown, isOpen := raygui.GetOpenDropdown(o.ColDropdowns)
	if isOpen {
		raygui.Disable()
		defer raygui.Enable()
	}

	// Render bottom to top to avoid overlap issues with dropdowns

	fieldY := n.UIRect.Y + n.UIRect.Height - UIFieldHeight
	if raygui.Button(rl.Rectangle{
		n.UIRect.X,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "+") {
		o.Cols = append(o.Cols, node.OrderColumn{})
		o.ColDropdowns = append(o.ColDropdowns, &raygui.DropdownEx{})
	}
	if raygui.Button(rl.Rectangle{
		n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
		fieldY,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	}, "-") {
		if len(o.Cols) > 1 {
			o.Cols = o.Cols[:len(o.Cols)-1]
			o.ColDropdowns = o.ColDropdowns[:len(o.ColDropdowns)-1]
		}
	}

	for i := len(o.ColDropdowns) - 1; i >= 0; i-- {
		fieldY -= UIFieldSpacing + UIFieldHeight
		func() {
			col := &o.Cols[i]
			dropdown := o.ColDropdowns[i]
			if openDropdown == dropdown {
				raygui.Enable()
				defer raygui.Disable()
			}

			colName := dropdown.Do(rl.Rectangle{
				n.UIRect.X,
				fieldY,
				n.UIRect.Width - orderDirectionWidth - UIFieldSpacing,
				UIFieldHeight,
			})
			col.Col, _ = colName.(string)

			activeSort := 0
			if col.Descending {
				activeSort = 1
			}
			newSort := raygui.ComboBox(rl.Rectangle{
				n.UIRect.X + n.UIRect.Width - orderDirectionWidth,
				fieldY,
				orderDirectionWidth,
				UIFieldHeight,
			}, "A-Z;Z-A", activeSort)
			switch newSort {
			case 1:
				col.Descending = true
			default:
				col.Descending = false
			}
		}()
	}
}

func doCombineRowsUpdate(n *node.Node, c *node.CombineRows) {
	n.UISize = rl.Vector2{X: 200, Y: float32(48)}
	c.Dropdown.SetOptions(combineRowsOpts...)
}

var combineRowsOpts = []raygui.DropdownExOption{
	{"UNION", node.Union},
	{"UNION ALL", node.UnionAll},
	{"INTERSECT", node.Intersect},
	{"EXCEPT", node.Except},
}

func doCombineRowsUI(n *node.Node, d *node.CombineRows) {
	chosen := d.Dropdown.Do(n.UIRect)
	d.CombinationType = chosen.(node.CombineType)
}

func doJoinUpdate(n *node.Node, j *node.Join) {
	n.UISize = rl.Vector2{X: 200, Y: float32(48)}

	// TODO probably put this somewhere else?
	for i := range n.Inputs[1:] {
		fmt.Println(i)
		if len(j.Conditions) < i+1 {
			j.Conditions = append(j.Conditions, &node.JoinCondition{
				Type:      node.InnerJoin,
				Condition: "???",
				TextBox:   &raygui.TextBoxEx{},
			})
		}
	}

}

func doJoinUI(n *node.Node, j *node.Join) {
	for _, condition := range j.Conditions {
		condition.Condition = condition.TextBox.Do(n.UIRect, condition.Condition, 100)
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

var errorOpts = []raygui.DropdownExOption{{"ERROR", "ERROR"}}

// Gets dropdown options for the table produced by the given node.
// Returns default options if no schema can be found.
func columnNameDropdownOpts(inputNode *node.Node) []raygui.DropdownExOption {
	if inputNode == nil {
		return errorOpts
	}

	var opts []raygui.DropdownExOption
	schemaCols, err := getSchema(inputNode)
	if err == nil {
		for _, col := range schemaCols {
			opts = append(opts, raygui.DropdownExOption{
				Name:  col,
				Value: col,
			})
		}
	} else {
		log.Print(err)
		return errorOpts
	}

	return opts
}
