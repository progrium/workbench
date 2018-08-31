package main

import (
	"fmt"
	"reflect"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// ssh interface
// support for more types: "numerics", nested structs, slices, channels?
// arguments support for methods
// xpath query system. ${Bazqux.Lastname}

var log *tview.TextView

type Foobar struct {
	FirstCar string
	Computer string
	Banana   bool
}

func (f *Foobar) HelloWorld() {
	fmt.Fprintln(log, "Hello world!")
}

func (f *Foobar) GoodbyeWorld() {
	fmt.Fprintln(log, "Goodbye world!")
}

type Bazqux struct {
	LastName string
	Banana   bool
	Computer string
}

func (f *Bazqux) Greet() {
	fmt.Fprintf(log, "Hello, Mr %s!\n", f.LastName)
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
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNavy
	app := tview.NewApplication()
	log = tview.NewTextView()

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

	log.SetBorder(true)
	rows.AddItem(log, 0, 1, false)

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

		}
		return event
	})
	if err := app.SetRoot(rows, true).Run(); err != nil {
		panic(err)
	}
}
