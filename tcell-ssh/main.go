package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"syscall"

	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
	"github.com/progrium/tcell"

	"github.com/mattn/go-runewidth"
)

var defStyle tcell.Style

type sshDriver struct {
	winWidth  int
	winHeight int
	winCh     chan os.Signal
	ptmx      *os.File
}

func (d *sshDriver) Init(winCh chan os.Signal) (*os.File, *os.File, error) {
	ptmx, tty, err := pty.Open()
	if err != nil {
		return nil, nil, err
	}
	d.ptmx = ptmx
	d.winCh = winCh
	return tty, tty, nil
}
func (d *sshDriver) Fini() {}
func (d *sshDriver) WinSize() (int, int, error) {
	return d.winWidth, d.winHeight, nil
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, r, nil, style)
		}
	}
}

type DriverSetter interface {
	SetDriver(tcell.TermDriver)
}

func main() {
	ssh.Handle(func(s ssh.Session) {
		sshPty, winCh, isPty := s.Pty()
		if !isPty {
			return
		}
		driver := &sshDriver{}
		driver.winWidth = sshPty.Window.Width
		driver.winHeight = sshPty.Window.Height
		screen, e := tcell.NewTerminfoScreen()
		if e != nil {
			fmt.Fprintf(os.Stderr, "%v\n", e)
			return
		}
		dset, ok := screen.(DriverSetter)
		if !ok {
			log.Fatal("Unable to set tcell driver")
			return
		}
		dset.SetDriver(driver)
		if e := screen.Init(); e != nil {
			fmt.Fprintf(os.Stderr, "%v\n", e)
			return
		}
		go termLoop(screen, s)
		go func() {
			for win := range winCh {
				driver.winWidth = win.Width
				driver.winHeight = win.Height
				driver.winCh <- syscall.SIGWINCH
			}
		}()
		go func() { _, _ = io.Copy(driver.ptmx, s) }()
		_, _ = io.Copy(s, driver.ptmx)
	})

	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

func termLoop(s tcell.Screen, sess ssh.Session) {
	defStyle = tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.Clear()

	white := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).Background(tcell.ColorRed)

	w, h := s.Size()

	for {
		drawBox(s, 1, 1, 42, 6, white, ' ')
		emitStr(s, 2, 2, white, "Press ESC to exit, C to clear.")

		s.Show()
		ev := s.PollEvent()
		st := tcell.StyleDefault.Background(tcell.ColorRed)
		w, h = s.Size()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			s.SetContent(w-1, h-1, 'R', nil, st)
		case *tcell.EventKey:
			s.SetContent(w-2, h-2, ev.Rune(), nil, st)
			s.SetContent(w-1, h-1, 'K', nil, st)
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				sess.Exit(0)
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else {
				if ev.Rune() == 'C' || ev.Rune() == 'c' {
					s.Clear()
				}
			}

		default:
			s.SetContent(w-1, h-1, 'X', nil, st)
		}

	}
}
