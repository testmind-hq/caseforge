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
	ScenarioUnicodeInjection CoverageScenario = "UNICODE_INJECTION"  // unicode/control-char injection

	ScenarioMassAssignmentPrivilege CoverageScenario = "MASS_ASSIGNMENT_PRIVILEGE"  // privilege escalation via extra fields
	ScenarioMassAssignmentStatus    CoverageScenario = "MASS_ASSIGNMENT_STATUS"     // status manipulation via extra fields
	ScenarioMassAssignmentFinancial CoverageScenario = "MASS_ASSIGNMENT_FINANCIAL"  // financial manipulation via extra fields
	ScenarioMassAssignmentIdentity  CoverageScenario = "MASS_ASSIGNMENT_IDENTITY"   // ownership takeover via extra fields

	ScenarioIDORParam CoverageScenario = "IDOR_PARAM" // IDOR: substitute ID param with alternative value

	ScenarioNullableAcceptance CoverageScenario = "NULLABLE_ACCEPTANCE" // nullable field accepts null value
	ScenarioReadOnlyWrite      CoverageScenario = "READ_ONLY_WRITE"     // readOnly field rejected on write
	ScenarioWriteOnlyRead      CoverageScenario = "WRITE_ONLY_READ"     // writeOnly field absent from read response

	ScenarioFieldBoundaryValid   CoverageScenario = "FIELD_BOUNDARY_VALID"   // nested field at declared boundary (valid)
	ScenarioFieldBoundaryInvalid CoverageScenario = "FIELD_BOUNDARY_INVALID" // nested field outside declared boundary (invalid)
	ScenarioRequiredOmission     CoverageScenario = "REQUIRED_OMISSION"      // required field entirely absent from payload

	ScenarioPositiveExample CoverageScenario = "POSITIVE_EXAMPLE" // happy-path from spec examples

	ScenarioCRUDFlow CoverageScenario = "CRUD_FLOW" // multi-step CRUD lifecycle chain
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
	ScenarioUnicodeInjection,
	ScenarioMassAssignmentPrivilege,
	ScenarioMassAssignmentStatus,
	ScenarioMassAssignmentFinancial,
	ScenarioMassAssignmentIdentity,
	ScenarioIDORParam,
	ScenarioNullableAcceptance,
	ScenarioReadOnlyWrite,
	ScenarioWriteOnlyRead,
	ScenarioFieldBoundaryValid,
	ScenarioFieldBoundaryInvalid,
	ScenarioRequiredOmission,
	ScenarioPositiveExample,
	ScenarioCRUDFlow,
}
