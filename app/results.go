package app

import (
	"fmt"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const resultFontSize = 20
const resultCellPaddingV = 5
const resultCellPaddingH = 8
const resultRowHeight = resultCellPaddingV + resultFontSize + resultCellPaddingV

type QueryResultPanel struct {
	QueryResult *queryResult
	ScrollPanel raygui.ScrollPanelEx

	Rows      [][]string
	ColWidths []float32
}

func (p *QueryResultPanel) Update(q *queryResult) {
	p.QueryResult = q

	// measure text only once
	p.Rows = nil
	for i := -1; i < len(p.QueryResult.Rows); i++ {
		if i < 0 {
			// headers
			p.Rows = append(p.Rows, p.QueryResult.Columns)
		} else {
			row := p.QueryResult.Rows[i]
			valStrings := make([]string, len(row))
			for i, v := range row {
				valStrings[i] = fmt.Sprintf("%v", v)
			}
			p.Rows = append(p.Rows, valStrings)
		}
	}

	p.ColWidths = make([]float32, len(p.Rows[0]))
	for r := 0; r < len(p.Rows); r++ {
		for c := 0; c < len(p.Rows[0]); c++ {
			thisColWidth := resultCellPaddingH + measureBasicText(p.Rows[r][c], resultFontSize).X + resultCellPaddingH
			if p.ColWidths[c] < thisColWidth {
				p.ColWidths[c] = thisColWidth
			}
		}
	}
}

func (p *QueryResultPanel) Draw(bounds rl.Rectangle) {
	if p.QueryResult == nil {
		return
	}

	var totalWidth float32
	for _, w := range p.ColWidths {
		totalWidth += w
	}

	panelContents := rl.Rectangle{0, 0, totalWidth, float32(len(p.Rows) * resultRowHeight)}
	p.ScrollPanel.Do(bounds, panelContents, func(scroll raygui.ScrollContext) {
		cellPos := scroll.Start
		for _, row := range p.Rows {
			cellPos.X = scroll.Start.X
			for i, cell := range row {
				cellRec := rl.Rectangle{cellPos.X, cellPos.Y, cellPos.X + p.ColWidths[i], cellPos.Y + resultRowHeight}
				if rl.CheckCollisionRecs(cellRec, scroll.View) {
					drawBasicText(cell, cellPos.X+resultCellPaddingH, cellPos.Y+resultCellPaddingV+1, resultFontSize, PaneFontColor)
				}
				cellPos.X += p.ColWidths[i]
			}

			lineY := cellPos.Y + resultRowHeight
			if scroll.View.Y <= lineY && lineY <= scroll.View.Y+scroll.View.Height {
				rl.DrawLine(0, int32(lineY), int32(scroll.View.X+scroll.View.Width), int32(lineY), PaneLineColor)
			}

			cellPos.Y += resultRowHeight
		}

		gridX := scroll.Start.X
		for _, width := range p.ColWidths {
			gridX += width
			rl.DrawLine(int32(gridX), int32(scroll.View.Y), int32(gridX), int32(scroll.View.Y+scroll.View.Height), PaneLineColor)
		}
	})
}

const resultsOpenDuration = 0.3

var resultsOpen bool
var resultsOpenFrac float32 // 0 for closed, 1 for open
var resultsCurrentHeight float32 = 0

var latestResults = &QueryResultPanel{}

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
		rl.NewColor(98, 85, 101, 255),
	)

	SetStyleColor(raygui.Default, raygui.BackgroundColorProp, MainColor())

	DoPane(rl.Rectangle{0, screenHeight - resultsCurrentHeight, screenWidth - currentSQLWidth, resultsMaxHeight}, func(p Pane) {
		latestResults.Draw(p.Bounds)
	})

	if rl.CheckCollisionPointRec(rl.GetMousePosition(), tabRect) && rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		resultsOpen = !resultsOpen
	}
}

func setResultsOpen(open bool) {
	resultsOpen = !resultsOpen
}
