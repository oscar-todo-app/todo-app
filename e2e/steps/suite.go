package steps

import (
	"context"

	"e2e/helpers"

	"github.com/cucumber/godog"
)

func InitializeScenario(sc *godog.ScenarioContext) {

	steps := &Entity{}

	sc.Step(`I open website "([^"]*)"`, steps.OpenSite)
	sc.Step(`I fill form in following information:`, steps.IfillForm)
	sc.Step(`I click submit button`, steps.IClickSubmit)
	sc.Step(`Verify result information "([^"]*)"`, steps.VerifyResult)

	sc.After(func(ctx context.Context, _ *godog.Scenario, err error) (context.Context, error) {
		if err := steps.Page.Close(); err != nil {
			helpers.LogPanicln(err)
		}
		if err := steps.Browser.Close(); err != nil {
			helpers.LogPanicln(err)
		}
		if err := steps.Pw.Stop(); err != nil {
			helpers.LogPanicln(err)
		}
		return ctx, nil
	})
}
