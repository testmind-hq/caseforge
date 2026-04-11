// internal/score/scenarios.go
package score

// CoverageScenario is a named test scenario type that identifies the specific
// constraint or property being exercised. Populated in CaseSource.Scenario.
// Inspired by schemathesis CoverageScenario — allows the scorer to report
// exactly which scenario types are covered and which are missing.
type CoverageScenario = string

const (
	ScenarioStringMinLength  CoverageScenario = "STRING_MIN_LENGTH"  // string at declared minLength
	ScenarioStringMaxLength  CoverageScenario = "STRING_MAX_LENGTH"  // string at declared maxLength
	ScenarioStringBelowMin   CoverageScenario = "STRING_BELOW_MIN"   // string shorter than minLength (invalid)
	ScenarioStringAboveMax   CoverageScenario = "STRING_ABOVE_MAX"   // string longer than maxLength (invalid)
	ScenarioNumberMin        CoverageScenario = "NUMBER_MIN"         // number at declared minimum
	ScenarioNumberMax        CoverageScenario = "NUMBER_MAX"         // number at declared maximum
	ScenarioNumberBelowMin   CoverageScenario = "NUMBER_BELOW_MIN"   // number below minimum (invalid)
	ScenarioNumberAboveMax   CoverageScenario = "NUMBER_ABOVE_MAX"   // number above maximum (invalid)
	ScenarioMissingRequired  CoverageScenario = "MISSING_REQUIRED"   // required field absent
	ScenarioNullInjection    CoverageScenario = "NULL_INJECTION"     // null sent for non-nullable field
	ScenarioEnumInvalid      CoverageScenario = "ENUM_INVALID"       // value not in declared enum
	ScenarioWrongType        CoverageScenario = "WRONG_TYPE"         // value of wrong JSON type
	ScenarioArrayMinItems    CoverageScenario = "ARRAY_MIN_ITEMS"    // array below minItems (invalid)
	ScenarioArrayMaxItems    CoverageScenario = "ARRAY_MAX_ITEMS"    // array above maxItems (invalid)
	ScenarioWrongContentType CoverageScenario = "WRONG_CONTENT_TYPE" // request with unsupported content-type
)

// trackedScenarios is the canonical list of scenarios reported in score output.
var trackedScenarios = []CoverageScenario{
	ScenarioMissingRequired,
	ScenarioNullInjection,
	ScenarioStringBelowMin,
	ScenarioStringAboveMax,
	ScenarioNumberBelowMin,
	ScenarioNumberAboveMax,
	ScenarioEnumInvalid,
	ScenarioWrongType,
	ScenarioArrayMinItems,
	ScenarioArrayMaxItems,
	ScenarioWrongContentType,
}
