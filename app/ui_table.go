package app

import (
	"github.com/bvisness/SQLJam/node"
	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
)

func doTableUpdate(n *node.Node, t *node.Table) {
	// init dropdown
	if len(t.TableDropdown.GetOptions()) == 0 {
		updateTableDropdown(&t.TableDropdown)
	}

	n.UISize = rl.Vector2{X: 200, Y: 24}
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
