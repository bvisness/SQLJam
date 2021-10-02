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

const resultsOpenDuration = 0.3

var resultsOpen bool
var resultsOpenFrac float32 // 0 for closed, 1 for open
var resultsCurrentHeight float32 = 0

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
	if resultsOpen {
		resultsOpenFrac += rl.GetFrameTime() / resultsOpenDuration
	} else {
		resultsOpenFrac -= rl.GetFrameTime() / resultsOpenDuration
	}
	resultsOpenFrac = Clamp(resultsOpenFrac, 0, 1)

	resultsCurrentHeight = EaseInOutCubic(resultsOpenFrac) * resultsMaxHeight

	var lineY float32 = screenHeight - resultsCurrentHeight - dividerThickness/2
	rl.DrawLineEx(
		rl.Vector2{0, lineY},
		rl.Vector2{screenWidth, lineY},
		dividerThickness, rl.Black,
	)

	const tabWidth = 60
	const tabHeight = 40

	var tabX float32 = screenWidth - tabWidth
	var tabY float32 = lineY - tabHeight
	tabRect := rl.Rectangle{tabX, tabY, tabWidth, tabHeight}
	rl.DrawRectangleRounded(tabRect, RoundnessPx(tabRect, 4), 5, rl.Black)

	const triangleSize = 12
	t1 := rl.Vector2{-triangleSize, triangleSize / 2}
	t2 := rl.Vector2{triangleSize, triangleSize / 2}
	t3 := rl.Vector2{0, -triangleSize / 2}

	trianglePos := rl.Vector2{tabX + tabWidth/2, tabY + tabHeight/2}
	triangleRot := resultsCurrentHeight / resultsMaxHeight * rl.Pi
	rl.DrawTriangle(
		rl.Vector2Add(Vector2Rotate(t1, triangleRot), trianglePos),
		rl.Vector2Add(Vector2Rotate(t2, triangleRot), trianglePos),
		rl.Vector2Add(Vector2Rotate(t3, triangleRot), trianglePos),
		rl.White,
	)

	DoPane(rl.Rectangle{0, screenHeight - resultsCurrentHeight, screenWidth - currentSQLWidth, resultsMaxHeight}, func(p Pane) {
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

	if rl.CheckCollisionPointRec(rl.GetMousePosition(), tabRect) && rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		resultsOpen = !resultsOpen
	}
}

func setResultsOpen(open bool) {
	resultsOpen = !resultsOpen
}
