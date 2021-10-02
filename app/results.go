package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var latestResult *queryResult
var latestResultsPanel raygui.ScrollPanelEx

func drawLatestResults() {
	DoPane(rl.Rectangle{screenWidth - sidebarWidth, 0, sidebarWidth, screenHeight - currentSQLHeight}, func(p Pane) {
		if latestResult == nil {
			return
		}

		const fontSize = 20
		const cellPaddingV = 5
		const cellPaddingH = 8
		const rowHeight = cellPaddingV + fontSize + cellPaddingV

		var rows [][]string
		for i := -1; i < len(latestResult.Rows); i++ {
			if i < 0 {
				// headers
				rows = append(rows, latestResult.Columns)
			} else {
				row := latestResult.Rows[i]
				valStrings := make([]string, len(row))
				for i, v := range row {
					valStrings[i] = fmt.Sprintf("%v", v)
				}
				rows = append(rows, valStrings)
			}
		}

		colWidths := make([]float32, len(rows[0]))
		for r := 0; r < len(rows); r++ {
			for c := 0; c < len(rows[0]); c++ {
				thisColWidth := cellPaddingH + measureBasicText(rows[r][c], fontSize).X + cellPaddingH
				if colWidths[c] < thisColWidth {
					colWidths[c] = thisColWidth
				}
			}
		}

		var totalWidth float32
		for _, w := range colWidths {
			totalWidth += w
		}

		panelContents := rl.Rectangle{0, 0, totalWidth, float32(len(rows) * rowHeight)}
		latestResultsPanel.Do(p.Bounds, panelContents, func(scroll raygui.ScrollContext) {
			cellPos := scroll.Start
			for _, row := range rows {
				cellPos.X = scroll.Start.X
				for i, cell := range row {
					cellRec := rl.Rectangle{cellPos.X, cellPos.Y, cellPos.X + colWidths[i], cellPos.Y + rowHeight}
					if rl.CheckCollisionRecs(cellRec, scroll.View) {
						drawBasicText(cell, cellPos.X+cellPaddingH, cellPos.Y+cellPaddingV+1, fontSize, rl.Black)
					}
					cellPos.X += colWidths[i]
				}

				lineY := cellPos.Y + rowHeight
				if scroll.View.Y <= lineY && lineY <= scroll.View.Y+scroll.View.Height {
					rl.DrawLine(0, int32(lineY), int32(scroll.View.X+scroll.View.Width), int32(lineY), rl.LightGray)
				}

				cellPos.Y += rowHeight
			}

			gridX := scroll.Start.X
			for _, width := range colWidths {
				gridX += width
				rl.DrawLine(int32(gridX), int32(scroll.View.Y), int32(gridX), int32(scroll.View.Y+scroll.View.Height), rl.LightGray)
			}
		})
	})
}
