package app

import (
	"fmt"
	"reflect"

	"github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var ChartColor = rl.NewColor(155, 171, 178, 255)

type Chart struct {
	ValueCol         string
	LabelCol         string
	ValueColDropdown raygui.DropdownEx
	LabelColDropdown raygui.DropdownEx

	Size      rl.Vector2 // this will be applied to UISize which will determine the node Size. Make sense???
	StartSize rl.Vector2

	QueryResult *queryResult
}

var _ NodeData = &Chart{}

func NewChart() *Node {
	return &Node{
		Title:   "Chart",
		CanSnap: true,
		Color:   ChartColor,
		Inputs:  make([]*Node, 1),
		Data: &Chart{
			Size: rl.Vector2{600, 400},
		},
	}
}

func (c *Chart) Update(n *Node) {
	if n.Schema == nil {
		c.QueryResult = doQuery(n.GenerateSql(true))
		n.Schema = getSchema(n)
	}

	opts := columnNameDropdownOpts(n.Inputs[0])
	c.ValueColDropdown.SetOptions(opts...)
	c.LabelColDropdown.SetOptions(opts...)

	n.UISize = c.Size
}

func (c *Chart) DoUI(n *Node) {
	{
		var labelIndex int
		var valueIndex int
		for i, col := range c.QueryResult.Columns {
			if col == c.LabelCol {
				labelIndex = i
			}
			if col == c.ValueCol {
				valueIndex = i
			}
		}

		var series []barChartSeries

		for _, row := range c.QueryResult.Rows {
			label := fmt.Sprintf("%v", row[labelIndex])

			var value float64
			rt := reflect.TypeOf(value)
			rv := reflect.ValueOf(row[valueIndex])
			if rv.CanConvert(rt) {
				value = rv.Convert(rt).Float()
			}

			series = append(series, barChartSeries{
				Label:  label,
				Values: []float64{value}, // TODO: moar values!!
			})
		}

		const topPadding = 20

		drawBarChart(rl.Rectangle{
			n.UIRect.X,
			n.UIRect.Y + UIFieldHeight + topPadding,
			n.UIRect.Width,
			n.UIRect.Height - UIFieldHeight - topPadding,
		}, series, n.Color)
	}

	valueCol := c.ValueColDropdown.Do(rl.Rectangle{
		n.UIRect.X,
		n.UIRect.Y,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	})
	c.ValueCol, _ = valueCol.(string)

	labelCol := c.LabelColDropdown.Do(rl.Rectangle{
		n.UIRect.X + n.UIRect.Width/2 + UIFieldSpacing/2,
		n.UIRect.Y,
		n.UIRect.Width/2 - UIFieldSpacing/2,
		UIFieldHeight,
	})
	c.LabelCol, _ = labelCol.(string)

	// dragging
	{
		bottomRight := rl.Vector2{n.Pos.X + n.Size.X, n.Pos.Y + n.Size.Y}
		resizeRect := rl.Rectangle{bottomRight.X - 20, bottomRight.Y - 20, 20, 20}

		drawResizeHandle(bottomRight, n.Color)

		resizeDragKey := fmt.Sprintf("resize: %p", c)
		if tryStartDrag(resizeDragKey, resizeRect, rl.Vector2{}) {
			c.StartSize = c.Size
		}

		if resizingThis, _, canceled := dragState(resizeDragKey); resizingThis {
			if canceled {
				c.Size = c.StartSize
			} else {
				newSize := rl.Vector2Add(c.StartSize, dragOffset())
				if newSize.X < previewMinWidth {
					newSize.X = previewMinWidth
				}
				if newSize.Y < previewMinHeight {
					newSize.Y = previewMinHeight
				}
				c.Size = newSize
			}
		} else {
			c.Size = rl.Vector2{n.UIRect.Width, n.UIRect.Height}
		}
	}
}

func (c *Chart) Serialize() (string, bool) {
	return "", false
}

func (c *Chart) Dropdowns() []*raygui.DropdownEx {
	var res []*raygui.DropdownEx
	res = append(res, &c.ValueColDropdown)
	res = append(res, &c.LabelColDropdown)
	return res
}

type barChartSeries struct {
	Label  string
	Values []float64
}

func drawBarChart(bounds rl.Rectangle, series []barChartSeries, nodeColor rl.Color) {
	const spacingBetweenSeries = 0.5 // times width of series

	color := Brightness(nodeColor, 0.4)

	var maxVal float64
	for _, s := range series {
		for _, v := range s.Values {
			if v > maxVal {
				maxVal = v
			}
		}
	}

	totalAbstractWidth := float32(len(series)) + spacingBetweenSeries*float32(len(series)-1)
	seriesWidth := bounds.Width / totalAbstractWidth
	spacingWidth := seriesWidth * spacingBetweenSeries

	x := bounds.X
	for _, s := range series {
		barWidth := seriesWidth / float32(len(s.Values))

		var textHeight float32
		textY := bounds.Y + bounds.Height
		if seriesWidth > 30 {
			textHeight = UIFieldHeight
			textY -= textHeight

			const textSize = 20
			textMeasured := measureBasicText(s.Label, textSize)

			var size float32 = textSize
			if textMeasured.X > seriesWidth {
				resizeRatio := (seriesWidth / textMeasured.X)
				size = size * resizeRatio
				textMeasured.X = seriesWidth
				textMeasured.Y = textMeasured.Y * resizeRatio
			}

			drawBasicText(
				s.Label,
				x+seriesWidth/2-textMeasured.X/2,
				textY+textHeight/2-textMeasured.Y/2,
				size, color,
			)
		}

		for _, v := range s.Values {
			barHeight := float32(v/maxVal) * (bounds.Height - textHeight)
			rl.DrawRectangleRec(rl.Rectangle{
				x,
				textY - barHeight,
				barWidth,
				barHeight,
			}, color)
			x += barWidth
		}
		x += spacingWidth
	}
}
