package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const lrFontSize = 20
const lrCellPaddingV = 5
const lrCellPaddingH = 8
const lrRowHeight = lrCellPaddingV + lrFontSize + lrCellPaddingV

var latestResult *queryResult
var latestResultsPanel raygui.ScrollPanelEx

var lrRows [][]string
var lrColWidths []float32

func setLatestResult(res *queryResult) {
	latestResult = res

	// measure text only once
	lrRows = nil
	for i := -1; i < len(latestResult.Rows); i++ {
		if i < 0 {
			// headers
			lrRows = append(lrRows, latestResult.Columns)
		} else {
			row := latestResult.Rows[i]
			valStrings := make([]string, len(row))
			for i, v := range row {
				valStrings[i] = fmt.Sprintf("%v", v)
			}
			lrRows = append(lrRows, valStrings)
		}
	}

	lrColWidths = make([]float32, len(lrRows[0]))
	for r := 0; r < len(lrRows); r++ {
		for c := 0; c < len(lrRows[0]); c++ {
			thisColWidth := lrCellPaddingH + measureBasicText(lrRows[r][c], lrFontSize).X + lrCellPaddingH
			if lrColWidths[c] < thisColWidth {
				lrColWidths[c] = thisColWidth
			}
		}
	}
}

func drawLatestResults() {
	DoPane(rl.Rectangle{screenWidth - sidebarWidth, 0, sidebarWidth, screenHeight - currentSQLHeight}, func(p Pane) {
		if latestResult == nil {
			return
		}

		var totalWidth float32
		for _, w := range lrColWidths {
			totalWidth += w
		}

		panelContents := rl.Rectangle{0, 0, totalWidth, float32(len(lrRows) * lrRowHeight)}
		latestResultsPanel.Do(p.Bounds, panelContents, func(scroll raygui.ScrollContext) {
			cellPos := scroll.Start
			for _, row := range lrRows {
				cellPos.X = scroll.Start.X
				for i, cell := range row {
					cellRec := rl.Rectangle{cellPos.X, cellPos.Y, cellPos.X + lrColWidths[i], cellPos.Y + lrRowHeight}
					if rl.CheckCollisionRecs(cellRec, scroll.View) {
						drawBasicText(cell, cellPos.X+lrCellPaddingH, cellPos.Y+lrCellPaddingV+1, lrFontSize, rl.Black)
					}
					cellPos.X += lrColWidths[i]
				}

				lineY := cellPos.Y + lrRowHeight
				if scroll.View.Y <= lineY && lineY <= scroll.View.Y+scroll.View.Height {
					rl.DrawLine(0, int32(lineY), int32(scroll.View.X+scroll.View.Width), int32(lineY), rl.LightGray)
				}

				cellPos.Y += lrRowHeight
			}

			gridX := scroll.Start.X
			for _, width := range lrColWidths {
				gridX += width
				rl.DrawLine(int32(gridX), int32(scroll.View.Y), int32(gridX), int32(scroll.View.Y+scroll.View.Height), rl.LightGray)
			}
		})
	})
}
