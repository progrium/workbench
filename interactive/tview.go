package main

import (
	"reflect"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Field struct {
	Name  string
	Value reflect.Value

	Val interface{}
}

type Object struct {
	Name   string
	Fields []*Field
}

func NewField(name string, value interface{}) *Field {
	f := &Field{
		Name: name,
		Val:  value,
	}
	f.Value = reflect.ValueOf(f.Val)
	return f
}

func drawForm(form *tview.Form, fields []*Field) {
	form.Clear(true)
	//form.AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil)
	for _, field := range fields {
		switch field.Value.Kind() {
		case reflect.String:
			f := field
			form.AddInputField(field.Name, field.Value.String(), 20, nil, func(text string) {
				f.Value.SetString(text)
			})
		case reflect.Bool:
			f := field
			form.AddCheckbox(field.Name, field.Value.Bool(), func(checked bool) {
				f.Value.SetBool(checked)
			})
		}
	}

	form.AddButton("Go", nil)
}

func main() {
	app := tview.NewApplication()

	objs := []*Object{
		&Object{
			Name: "Foobar",
			Fields: []*Field{
				NewField("First car", ""),
				NewField("Computer", ""),
				NewField("Banana", false),
			},
		},
		&Object{
			Name: "Bazqux",
			Fields: []*Field{
				NewField("Computer", ""),
				NewField("Banana", false),
				NewField("Last name", ""),
			},
		},
	}

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNavy
	flex := tview.NewFlex()
	list := tview.NewList()
	form := tview.NewForm()
	list.SetBorder(true).SetTitle("Objects")

	for _, obj := range objs {
		list.AddItem(obj.Name, "", 0, nil)
	}
	list.SetChangedFunc(func(idx int, _ string, _ string, _ rune) {
		drawForm(form, objs[idx].Fields)
		form.SetTitle(objs[idx].Name)
	})

	flex.AddItem(list, 0, 1, true)

	form.SetBorder(true)

	flex.AddItem(form, 0, 2, false)

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
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
