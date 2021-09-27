package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	gui2 "github.com/bvisness/SQLJam/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	_ "github.com/mattn/go-sqlite3"
)

const screenWidth = 800
const screenHeight = 450

var ballPosition = rl.Vector2{screenWidth / 2, screenHeight / 2}

var movies []string

func main() {
	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	fmt.Println(gui2.GetStyle(0, 0))

	db, err := sql.Open("sqlite3", "./sakila.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select title from film limit 5")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var title string
		err = rows.Scan(&title)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, title)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	for !rl.WindowShouldClose() {
		doFrame()
	}
}

//var colors = []string{"Maroon", "Blue", "Green"}
var colors = "Maroon;Blue;Green"
var rlColors = []rl.Color{rl.Maroon, rl.Blue, rl.Green}
var selectedColor = 0

var ballLabel = "Ball"

var textBoxActive = false
var dropdownOpen = false

var selectedMovie int
var movieDropdownOpen = false

func doFrame() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

	// ballPosition.X = gui.SliderBar(rl.Rectangle{600, 40, 120, 20}, ballPosition.X, 0, screenWidth)
	// ballPosition.Y = gui.SliderBar(rl.Rectangle{600, 70, 120, 20}, ballPosition.Y, 0, screenHeight)

	ballPosition := rl.GetMousePosition()

	var toggleTextBox bool
	if ballLabel, toggleTextBox = gui2.TextBox(rl.Rectangle{40, 40, 120, 20}, ballLabel, 100, textBoxActive); toggleTextBox {
		textBoxActive = !textBoxActive
	}

	//selectedColor = gui.ComboBox(rl.Rectangle{40, 40, 120, 20}, colors, selectedColor)
	selectedColor = gui2.ComboBox(rl.Rectangle{40, 70, 120, 20}, colors, selectedColor)
	if gui2.DropdownBox(rl.Rectangle{40, 100, 120, 20}, colors, &selectedColor, dropdownOpen) {
		dropdownOpen = !dropdownOpen
	}
	//ballLabel = gui.TextBox(rl.Rectangle{40, 70, 120, 20}, ballLabel)

	rl.ClearBackground(rl.RayWhite)
	rl.DrawCircleV(ballPosition, 50, rlColors[selectedColor])
	rl.DrawText(ballLabel, int32(ballPosition.X), int32(ballPosition.Y), 14, rl.DarkGray)

	if gui2.DropdownBox(rl.Rectangle{40, 140, 120, 20}, makeDropdownOptions(movies...), &selectedMovie, movieDropdownOpen) {
		movieDropdownOpen = !movieDropdownOpen
	}
}

func makeDropdownOptions(opts ...string) string {
	return strings.Join(opts, ";")
}
