package main

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/gdamore/tcell"
	"github.com/gliderlabs/ssh"
	"github.com/kr/pty"
	"github.com/rivo/tview"
)

// ssh interface
// support for more types: "numerics", nested structs, slices, channels?
// arguments support for methods
// xpath query system. ${Bazqux.Lastname}

var logView *tview.TextView

type Foobar struct {
	FirstCar string
	Computer string
	Banana   bool
}

func (f *Foobar) HelloWorld() {
	fmt.Fprintln(logView, "Hello world!")
}

func (f *Foobar) GoodbyeWorld() {
	fmt.Fprintln(logView, "Goodbye world!")
}

type Bazqux struct {
	LastName string
	Banana   bool
	Computer string
}

func (f *Bazqux) Greet() {
	fmt.Fprintf(logView, "Hello, Mr %s!\n", f.LastName)
}

type Object struct {
	Name    string
	Value   interface{}
	Fields  map[string]reflect.Value
	Methods map[string]reflect.Value
}

func NewObject(name string, value interface{}) *Object {
	obj := &Object{
		Name:    name,
		Value:   value,
		Fields:  make(map[string]reflect.Value),
		Methods: make(map[string]reflect.Value),
	}
	rval := reflect.ValueOf(value)
	vtyp := rval.Elem().Type()
	ptyp := rval.Type()
	for i := 0; i < vtyp.NumField(); i++ {
		f := vtyp.Field(i)
		obj.Fields[f.Name] = rval.Elem().FieldByName(f.Name)
	}
	for i := 0; i < rval.NumMethod(); i++ {
		m := ptyp.Method(i)
		fn := rval.Method(i)
		obj.Methods[m.Name] = fn
	}
	return obj
}

func drawForm(form *tview.Form, obj *Object) {
	form.Clear(true)
	//form.AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil)
	for name, field := range obj.Fields {
		switch field.Kind() {
		case reflect.String:
			f := field
			form.AddInputField(name, field.String(), 20, nil, func(text string) {
				f.SetString(text)
			})
		case reflect.Bool:
			f := field
			form.AddCheckbox(name, field.Bool(), func(checked bool) {
				f.SetBool(checked)
			})
		}
	}
	for name, fnval := range obj.Methods {
		fn := fnval
		form.AddButton(name, func() {
			fn.Call(nil)
		})
	}
}

func main() {
	// drawBox(screen, 1, 1, 42, 6, tcell.StyleDefault.
	// 	Foreground(tcell.ColorWhite).Background(tcell.ColorRed), ' ')
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNavy
	app := tview.NewApplication()
	//app.SetScreen(screen)
	logView = tview.NewTextView()

	objs := []*Object{
		NewObject("Foobar", &Foobar{}),
		NewObject("Bazqux", &Bazqux{}),
	}

	rows := tview.NewFlex()
	columns := tview.NewFlex()
	list := tview.NewList()
	form := tview.NewForm()

	list.SetBorder(true).SetTitle("Objects")

	for _, obj := range objs {
		list.AddItem(obj.Name, "", 0, nil)
	}
	list.SetChangedFunc(func(idx int, _ string, _ string, _ rune) {
		drawForm(form, objs[idx])
		form.SetTitle(objs[idx].Name)
	})

	columns.AddItem(list, 0, 1, true)

	form.SetBorder(true)

	columns.AddItem(form, 0, 2, false)

	rows.SetDirection(tview.FlexRow)
	rows.AddItem(columns, 0, 3, false)

	logView.SetBorder(true)
	rows.AddItem(logView, 0, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			app.SetFocus(list)
			app.Draw()
			return nil
		}
		if event.Key() == tcell.KeyRight {
			app.SetFocus(form)
			app.Draw()
			return nil
		}
		if event.Key() == tcell.KeyTab && app.GetFocus() == list {
			app.SetFocus(form)
			app.Draw()
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			app.Stop()
			os.Exit(0)
			//sess.Exit(0)
		}
		return event
	})

	log.Println("About to run")
	if err := app.SetRoot(rows, true).Run(); err != nil {
		panic(err)
	}
}

func rungui(screen tcell.Screen, sess ssh.Session) {
	// drawBox(screen, 1, 1, 42, 6, tcell.StyleDefault.
	// 	Foreground(tcell.ColorWhite).Background(tcell.ColorRed), ' ')
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNavy
	app := tview.NewApplication()
	app.SetScreen(screen)
	logView = tview.NewTextView()

	objs := []*Object{
		NewObject("Foobar", &Foobar{}),
		NewObject("Bazqux", &Bazqux{}),
	}

	rows := tview.NewFlex()
	columns := tview.NewFlex()
	list := tview.NewList()
	form := tview.NewForm()

	list.SetBorder(true).SetTitle("Objects")

	for _, obj := range objs {
		list.AddItem(obj.Name, "", 0, nil)
	}
	list.SetChangedFunc(func(idx int, _ string, _ string, _ rune) {
		drawForm(form, objs[idx])
		form.SetTitle(objs[idx].Name)
	})

	columns.AddItem(list, 0, 1, true)

	form.SetBorder(true)

	columns.AddItem(form, 0, 2, false)

	rows.SetDirection(tview.FlexRow)
	rows.AddItem(columns, 0, 3, false)

	logView.SetBorder(true)
	rows.AddItem(logView, 0, 1, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft {
			app.SetFocus(list)
			app.Draw()
			return nil
		}
		if event.Key() == tcell.KeyRight {
			app.SetFocus(form)
			app.Draw()
			return nil
		}
		if event.Key() == tcell.KeyTab && app.GetFocus() == list {
			app.SetFocus(form)
			app.Draw()
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			app.Stop()
			sess.Exit(0)
		}
		return event
	})

	log.Println("About to run")
	if err := app.SetRoot(rows, true).Run(); err != nil {
		panic(err)
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

type sshDriver struct {
	winWidth  int
	winHeight int
	winCh     chan os.Signal
	ptmx      *os.File
	ready     chan bool
}

func (d *sshDriver) Init(winCh chan os.Signal) (*os.File, *os.File, error) {
	ptmx, tty, err := pty.Open()
	if err != nil {
		return nil, nil, err
	}
	d.ptmx = ptmx
	d.winCh = winCh
	close(d.ready)
	return tty, tty, nil
}
func (d *sshDriver) Fini() {}
func (d *sshDriver) WinSize() (int, int, error) {
	return d.winWidth, d.winHeight, nil
}
