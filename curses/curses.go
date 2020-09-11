package curses

// struct ldat{};
// #define _Bool int
// #define NCURSES_OPAQUE 1
// #include <curses.h>
// #cgo LDFLAGS: -lncurses
import "C"

import (
	"fmt"
	"unsafe"
)

type void unsafe.Pointer
type Window C.WINDOW

type CursesError struct {
	message string
}

func (ce CursesError) Error() string {
	return ce.message
}

// Cursor options.
const (
	CURS_HIDE = iota
	CURS_NORM
	CURS_HIGH
)

// Pointers to the values in curses, which may change values.
var Cols *int = nil
var Rows *int = nil

var Colors *int = nil
var ColorPairs *int = nil

var Tabsize *int = nil

// The window returned from C.initscr()
var Stdwin *Window = nil

// Initializes gocurses
func init() {
	Cols = (*int)(void(&C.COLS))
	Rows = (*int)(void(&C.LINES))

	Colors = (*int)(void(&C.COLORS))
	ColorPairs = (*int)(void(&C.COLOR_PAIRS))

	Tabsize = (*int)(void(&C.TABSIZE))
}

func Initscr() (*Window, error) {
	Stdwin = (*Window)(C.initscr())

	if Stdwin == nil {
		return nil, CursesError{"Initscr failed"}
	}

	return Stdwin, nil
}

func Newwin(rows int, cols int, starty int, startx int) (*Window, error) {
	nw := (*Window)(C.newwin(C.int(rows), C.int(cols), C.int(starty), C.int(startx)))

	if nw == nil {
		return nil, CursesError{"Failed to create window"}
	}

	return nw, nil
}

func (win *Window) Delwin() error {
	if int(C.delwin((*C.WINDOW)(win))) != OK {
		return CursesError{"delete failed"}
	}
	return nil
}

func (win *Window) Subwin(rows int, cols int, starty int, startx int) (*Window, error) {
	sw := (*Window)(C.subwin((*C.WINDOW)(win), C.int(rows), C.int(cols), C.int(starty), C.int(startx)))

	if sw == nil {
		return nil, CursesError{"Failed to create window"}
	}

	return sw, nil
}

func (win *Window) Derwin(rows int, cols int, starty int, startx int) (*Window, error) {
	dw := (*Window)(C.derwin((*C.WINDOW)(win), C.int(rows), C.int(cols), C.int(starty), C.int(startx)))

	if dw == nil {
		return nil, CursesError{"Failed to create window"}
	}

	return dw, nil
}

func Start_color() error {
	if int(C.has_colors()) == 0 {
		return CursesError{"terminal does not support color"}
	}
	if int(C.start_color()) != OK {
		return CursesError{"Start_color failed"}
	}

	return nil
}

func Init_pair(pair int, fg int, bg int) {
	C.init_pair(C.short(pair), C.short(fg), C.short(bg))
}

func Color_pair(pair int) int32 {
	return int32(C.COLOR_PAIR(C.int(pair)))
}

func Noecho() error {
	if int(C.noecho()) != OK {
		return CursesError{"Noecho failed"}
	}
	return nil
}

func DoUpdate() error {
	if int(C.doupdate()) != OK {
		return CursesError{"Doupdate failed"}
	}
	return nil
}

func Echo() error {
	if int(C.noecho()) != OK {
		return CursesError{"Echo failed"}
	}
	return nil
}

func Curs_set(c int) error {
	if C.curs_set(C.int(c)) == ERR {
		return CursesError{"Curs_set failed"}
	}
	return nil
}

func Nocbreak() error {
	if C.nocbreak() != OK {
		return CursesError{"Nocbreak failed"}
	}
	return nil
}

func Cbreak() error {
	if C.cbreak() != OK {
		return CursesError{"Cbreak failed"}
	}
	return nil
}

func Raw() error {
	if C.raw() != OK {
		return CursesError{"Raw failed"}
	}
	return nil
}

func Noraw() error {
	if C.noraw() != OK {
		return CursesError{"Noraw failed"}
	}
	return nil
}

func Endwin() error {
	if C.endwin() != OK {
		return CursesError{"Endwin failed"}
	}
	return nil
}

// Leaving the error handling on this one up to the user
func (win *Window) Getch() int {
	return int(C.wgetch((*C.WINDOW)(win)))
}

func (win *Window) Addch(x, y int, c int32, flags int32) error {
	if C.mvwaddch((*C.WINDOW)(win), C.int(y), C.int(x), C.chtype(c)|C.chtype(flags)) != OK {
		return CursesError{"Addch failed"}
	}
	return nil
}

// Since CGO currently can't handle varg C functions we'll mimic the
// ncurses addstr functions.
func (win *Window) Addstr(x, y int, str string, flags int32, v ...interface{}) {
	var newstr string
	if v != nil {
		newstr = fmt.Sprintf(str, v...)
	} else {
		newstr = str
	}

	win.Move(x, y)

	for i := 0; i < len(newstr); i++ {
		C.waddch((*C.WINDOW)(win), C.chtype(newstr[i])|C.chtype(flags))
	}
}

func (win *Window) Move(y, x int) error {
	if C.wmove((*C.WINDOW)(win), C.int(y), C.int(x)) != OK {
		return CursesError{"Move failed"}
	}
	return nil
}

func (win *Window) Resize(rows, cols int) error {
	if C.wresize((*C.WINDOW)(win), C.int(rows), C.int(cols)) != OK {
		return CursesError{"Resize failed"}
	}
	return nil
}

func (w *Window) Keypad(tf bool) error {
	outint := 0
	if tf {
		outint = 1
	}
	if C.keypad((*C.WINDOW)(w), C.int(outint)) != OK {
		return CursesError{"Keypad failed"}
	}
	return nil
}

func (win *Window) Refresh() error {
	if C.wrefresh((*C.WINDOW)(win)) != OK {
		return CursesError{"Refresh failed"}
	}
	return nil
}

func (win *Window) Redrawln(beg_line, num_lines int) error {
	if C.wredrawln((*C.WINDOW)(win), C.int(beg_line), C.int(num_lines)) != OK {
		return CursesError{"Redrawln failed"}
	}
	return nil
}

func (win *Window) Redraw() error {
	if C.redrawwin((*C.WINDOW)(win)) != OK {
		return CursesError{"Redraw failed"}
	}
	return nil
}

func (win *Window) Clear() error {
	if C.wclear((*C.WINDOW)(win)) != OK {
		return CursesError{"Clear failed"}
	}
	return nil
}

func (win *Window) Erase() error {
	if C.werase((*C.WINDOW)(win)) != OK {
		return CursesError{"Erase failed"}
	}
	return nil
}

func (win *Window) Clrtobot() error {
	if C.wclrtobot((*C.WINDOW)(win)) != OK {
		return CursesError{"Clrtobot failed"}
	}
	return nil
}

func (win *Window) Clrtoeol() error {
	if C.wclrtoeol((*C.WINDOW)(win)) != OK {
		return CursesError{"Clrtoeol failed"}
	}
	return nil
}

func (win *Window) Box(verch, horch int) {
	// Guaranteed to return OK, so we won't check this one
	C.box((*C.WINDOW)(win), C.chtype(verch), C.chtype(horch))
}

func (win *Window) Background(colour int32) {
	// Guaranteed to return OK, so we won't check this one
	C.wbkgd((*C.WINDOW)(win), C.chtype(colour))
}

func (win *Window) Attron(flags int32) error {
	if C.wattron((*C.WINDOW)(win), C.int(flags)) != OK {
		return CursesError{"Attron failed"}
	}
	return nil
}

func (win *Window) Attroff(flags int32) error {
	if C.wattroff((*C.WINDOW)(win), C.int(flags)) != OK {
		return CursesError{"Attron failed"}
	}
	return nil
}
