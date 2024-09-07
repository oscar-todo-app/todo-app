package steps

import (
	"e2e/helpers"

	"github.com/playwright-community/playwright-go"
)

type Entity struct {
	Pw      *playwright.Playwright
	Browser playwright.Browser
	Page    playwright.Page
	Task    string
	Cases   helpers.Case
}
