package steps

import (
	"e2e/helpers"

	"github.com/cucumber/godog"
	"github.com/playwright-community/playwright-go"
)

func (e *Entity) OpenSite(params string) error {
	pw, err := playwright.Run()
	helpers.LogPanicln(err)
	browser, err := pw.Firefox.Launch()
	helpers.LogPanicln(err)

	e.Browser = browser

	page, err := e.Browser.NewPage()
	helpers.LogPanicln(err)

	e.Page = page

	_, err = page.Goto(params)
	helpers.LogPanicln(err)
	return nil

}

func (e *Entity) IfillForm(table *godog.Table) error {
	for _, row := range table.Rows {
		field := row.Cells[0]
		value := row.Cells[1]

		switch field.Value {
		case "todo":
			todo := e.Page.Locator("#todo")
			err := todo.Fill(value.Value)
			helpers.LogPanicln(err)

			task, err := todo.TextContent()
			helpers.LogPanicln(err)
			e.Task = task
		}
	}
	return nil
}

func (e *Entity) IClickSubmit() error {
	submitButton := e.Page.Locator("button[type=submit]")
	err := submitButton.Click()
	helpers.LogPanicln(err)
	return nil
}

func (e *Entity) VerifyResult(expected string) error {
	taskView := e.Page.Locator("text=New Todo Item")

	task, err := taskView.TextContent()
	helpers.LogPanicln(err)

	e.Cases.AssertEqual(expected, task)
	return nil
}
