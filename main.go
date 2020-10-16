package main

import (
    "fmt"
    "os"
    "time"

    "github.com/TheTitanrain/w32"
    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    choices  []string         // items on the to-do list
    cursor   int              // which to-do list item our cursor is pointing at
    selected map[int]struct{} // which to-do items are selected
}

func initialize() (tea.Model, tea.Cmd) {
    m := model{

        // Our to-do list is just a grocery list
        choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

        // A map which indicates which choices are selected. We're using
        // the  map like a mathematical set. The keys refer to the indexes
        // of the `choices` slice, above.
        selected: make(map[int]struct{}),
    }

    // Return the model and `nil`, which means "no I/O right now, please."
    return m, tick()
}

type tickMsg time.Time

var ii int

//var lines []string

var x, y int

func update(msg tea.Msg, mdl tea.Model) (tea.Model, tea.Cmd) {
    m, _ := mdl.(model)
    oldx := x
    oldy := y
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

        // These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit

        // The "up" and "k" keys move the cursor up
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        // The "down" and "j" keys move the cursor down
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }

        // The "enter" key and the spacebar (a literal space) toggle
        // the selected state for the item that the cursor is pointing at.
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }

            // case "w":
            //     y -= 1
            //     if y <= 0 {
            //         y = 0
            //     }
            // case "s":
            //     y += 1
            //     if y > 20 {
            //         y = 20
            //     }
            // case "a":
            //     x -= 1
            //     if x <= 0 {
            //         x = 0
            //     }
            // case "d":
            //     x += 1
            //     if x > 20 {
            //         x = 20
            //     }
        }

    case tickMsg:
        ii += 1
        // if m <= 0 {
        //     return m, tea.Quit
        // }
        return m, tick()

    }

    if oldx != x || oldy != y {
        // lines[oldx][oldy] = '.'
        // lines[x][y] = '#'
    }

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}

func tick() tea.Cmd {
    return tea.Tick(time.Millisecond*20, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func KeyIsPressing(key int) bool {
    keyState := w32.GetAsyncKeyState(key)

    return keyState&(1<<15) != 0
}
func view(mdl tea.Model) string {
    //m, _ := mdl.(model)
    // s := "What should we buy at the market?\n\n"
    // for i, choice := range m.choices {
    //     cursor := " " // no cursor
    //     if m.cursor == i {
    //         cursor = ">" // cursor!
    //     }
    //     checked := " " // not selected
    //     if _, ok := m.selected[i]; ok {
    //         checked = "x" // selected!
    //     }
    //     s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    // }
    // s += "\nPress q to quit.\n"

    if KeyIsPressing(w32.VK_UP) {

        y -= 1
        if y <= 0 {
            y = 0
        }
    }
    if KeyIsPressing(w32.VK_DOWN) {
        y += 1
        if y > 20 {
            y = 20
        }
    }
    if KeyIsPressing(w32.VK_LEFT) {
        x -= 1
        if x <= 0 {
            x = 0
        }
    }
    if KeyIsPressing(w32.VK_RIGHT) {
        x += 1
        if x > 20 {
            x = 20
        }
    }

    s := ""

    for i := 0; i < 20; i++ {
        for j := 0; j < 80; j++ {
            if j == x && i == y {
                s = s + "#"
            } else {
                s = s + "."
            }
        }
        s = s + "\n"
    }

    s += fmt.Sprintf("%d %d,%d", ii, x, y)
    keyState := w32.GetAsyncKeyState(w32.VK_SHIFT)
    s += fmt.Sprintf("%x", ii, x, y, keyState)

    // Send the UI for rendering
    return s
}

func main() {

    p := tea.NewProgram(initialize, update, view)
    if err := p.Start(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }

}
