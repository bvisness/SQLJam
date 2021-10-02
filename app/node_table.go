package app

import (
	"log"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Table struct {
	NodeData
	SqlSource
	Alias string
	Table string

	// UI data
	TableDropdown raygui.DropdownEx
}

func NewTable() *Node {
	return &Node{
		Title:   "Table",
		CanSnap: false,
		Color:   rl.NewColor(242, 201, 76, 255),
		Data: &Table{
			TableDropdown: raygui.NewDropdownEx(),
		},
	}
}

func (t *Table) SourceToSql(indent int) string {
	return t.Table
}

func (t *Table) SourceAlias() string {
	return "a"
}

func (t *Table) IsPure() bool {
	return true
}

func (t *Table) SourceTableName() string {
	return t.Table
}

func (t *Table) Update(n *Node) {
	// init dropdown
	if len(t.TableDropdown.GetOptions()) == 0 {
		updateTableDropdown(&t.TableDropdown)
	}

	n.UISize = rl.Vector2{X: 200, Y: 24}
}

func (t *Table) DoUI(n *Node) {
	if ival := t.TableDropdown.Do(n.UIRect); ival != nil {
		t.Table = ival.(string)
	}
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
