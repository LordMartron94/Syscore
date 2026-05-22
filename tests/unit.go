package tests

import (
	"shield"
	"syscore"
)

var standardRunCfg = shield.SHIELD_Testing_ScenarioRunConfig{
	MaxIterations: 1,
}

func init() {
	mainOperation := shield.SHIELD_Testing_OperationCreateStateless(
		"operation_dummy",
		"A dummy operation for testing/debugging some syscore stuff",
		func(_ struct{}, execCtx shield.SHIELD_Testing_ExecutionContext) []shield.SHIELD_Testing_ScenarioRunResult {
			return []shield.SHIELD_Testing_ScenarioRunResult{
				runWindowTestingScenario(execCtx),
			}
		},
		"SYSCORE", "Dummy",
	)

	shield.SHIELD_Registry_OperationRegister(mainOperation)
}

func runWindowTestingScenario(execCtx shield.SHIELD_Testing_ExecutionContext) shield.SHIELD_Testing_ScenarioRunResult {
	type scenarioInput struct{}
	type scenarioOutput struct{}

	scenario := shield.SHIELD_Testing_ScenarioCreate[scenarioInput, scenarioOutput](
		"scenario_dummy_windowing",
		"Debugs the available windowing APIs",
		[]shield.SHIELD_Testing_Guard[scenarioInput, scenarioOutput]{
			shield.SHIELD_Testing_GuardCreate(
				"guard_dummy",
				scenarioInput{},
				shield.SHIELD_Testing_GuardPolicyMustNotPanic[scenarioOutput](),
			),
		},
		func(input scenarioInput) (output scenarioOutput, error error) {
			syscore.SYSCORE_Window_WindowTest()
			return scenarioOutput{}, nil
		},
	)

	return shield.SHIELD_Testing_OperationRunScenario(
		scenario,
		execCtx,
		standardRunCfg,
	)
}
