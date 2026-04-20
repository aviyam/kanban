# Kanban TUI

A standalone Kanban board application for the terminal built with Go. This tool is designed for keyboard-driven efficiency, featuring SQLite persistence, customizable shortcuts, and automatic link handling.

## Features

- Three-column Kanban layout: To Do, In Progress, and Done.
- Keyboard-first navigation and task management.
- Persistent storage using a local SQLite database.
- Contextual help bar showing available shortcuts.
- Customizable shortcuts and column colors via JSON configuration.
- Automatic link extraction and opening in the default web browser.
- Automatic cleanup of completed tasks after a configurable period.
- Task reordering within columns.
- Detailed task view and deletion confirmation.

## Prerequisites

- Go 1.18 or higher.
- A terminal emulator that supports ANSI colors.

## Installation

1. Clone the repository or download the source code.
2. Initialize and build the application using the Makefile:
   ```bash
   make init
   ```
   Or manually:
   ```bash
   go mod tidy
   go build -o kanban-tui .
   ```

## Usage

Run the application using the Makefile:
```bash
make run
```
Or directly:
```bash
./kanban-tui
```

### Default Keyboard Shortcuts

#### Board Navigation
- h / l: Switch between columns.
- j / k: Navigate tasks within the focused column.
- Enter: View full task details.
- a: Add a new task.
- e: Edit the selected task.
- d: Delete the selected task (requires confirmation).
- o: Open URLs found in the task description.
- q: Quit the application.

#### Task Organization
- [ / ]: Move the selected task to the left or right column.
- Shift + J: Move the selected task down within the column.
- Shift + K: Move the selected task up within the column.

#### Form, Detail, and Delete Views
- Enter: Save form or return to board.
- Esc: Cancel action or return to board.
- Tab: Switch between fields in the task form.
- y / n: Confirm or cancel task deletion.

## Configuration

On the first run, the application generates a `config.json` file in the root directory. You can modify this file to customize your experience.

### Configuration Options

- shortcuts: Remap any board action to your preferred key.
- colors: Change the ANSI color codes for each column and the focused border.
- cleanup_done_after_days: Set the number of days after which completed tasks are automatically removed (set to 0 to disable).

### Example Configuration

```json
{
  "shortcuts": {
    "move_left": "h",
    "move_right": "l",
    "move_task_left": "[",
    "move_task_right": "]",
    "move_task_up": "K",
    "move_task_down": "J",
    "add_task": "a",
    "edit_task": "e",
    "delete_task": "d",
    "open_link": "o",
    "quit": "q"
  },
  "colors": {
    "todo_column": "204",
    "in_progress_column": "214",
    "done_column": "42",
    "focused_border": "62"
  },
  "cleanup_done_after_days": 1
}
```

## Maintenance

The project includes a Makefile for common tasks:
- `make build`: Compile the binary.
- `make run`: Build and run the application.
- `make test`: Run unit tests.
- `make clean`: Remove the binary and clean the cache.
- `make tidy`: Update go.mod and go.sum.

## Data Storage

The application stores all tasks in a local SQLite database named `kanban.db`. This file is created automatically in the project directory upon the first launch.
