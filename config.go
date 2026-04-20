package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Shortcuts struct {
		MoveLeft      string `json:"move_left"`
		MoveRight     string `json:"move_right"`
		MoveTaskLeft  string `json:"move_task_left"`
		MoveTaskRight string `json:"move_task_right"`
		MoveTaskUp    string `json:"move_task_up"`
		MoveTaskDown  string `json:"move_task_down"`
		AddTask       string `json:"add_task"`
		EditTask      string `json:"edit_task"`
		DeleteTask    string `json:"delete_task"`
		OpenLink      string `json:"open_link"`
		Quit          string `json:"quit"`
	} `json:"shortcuts"`
	Colors struct {
		TodoColumn       string `json:"todo_column"`
		InProgressColumn string `json:"in_progress_column"`
		DoneColumn       string `json:"done_column"`
		FocusedBorder    string `json:"focused_border"`
	} `json:"colors"`
	CleanupDoneAfterDays int `json:"cleanup_done_after_days"`
}

func DefaultConfig() Config {
	c := Config{}
	c.Shortcuts.MoveLeft = "h"
	c.Shortcuts.MoveRight = "l"
	c.Shortcuts.MoveTaskLeft = "["
	c.Shortcuts.MoveTaskRight = "]"
	c.Shortcuts.MoveTaskUp = "K"
	c.Shortcuts.MoveTaskDown = "J"
	c.Shortcuts.AddTask = "a"
	c.Shortcuts.EditTask = "e"
	c.Shortcuts.DeleteTask = "d"
	c.Shortcuts.OpenLink = "o"
	c.Shortcuts.Quit = "q"

	c.Colors.TodoColumn = "204"       // Pinkish
	c.Colors.InProgressColumn = "214" // Orange
	c.Colors.DoneColumn = "42"        // Green
	c.Colors.FocusedBorder = "62"     // Purple
	c.CleanupDoneAfterDays = 1
	return c
}

func LoadConfig() Config {
	file, err := os.ReadFile("config.json")
	if err != nil {
		c := DefaultConfig()
		data, _ := json.MarshalIndent(c, "", "  ")
		_ = os.WriteFile("config.json", data, 0644)
		return c
	}
	var c Config
	_ = json.Unmarshal(file, &c)
	return c
}
