# CaseForge Acceptance Test Suite

> **Usage:** Run `scripts/acceptance.sh` to execute all scenarios and regenerate results.
> Update this document whenever a new feature is added.

---

## How to Run

```bash
# Full acceptance run (builds binary + executes all scenarios)
./scripts/acceptance.sh

# Individual scenario (manual smoke test)
/tmp/caseforge <command> [flags]
```

---

## Acceptance Scenarios

### Core / CLI

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-001 | `--version` flag | `caseforge --version` | prints `caseforge version <tag>` | ✅ PASS |
| AT-002 | All commands registered | `caseforge --help` | lists ask, completion, config, diff, doctor, explore, fake, gen, init, lint, mcp, onboard, pairwise, run | ✅ PASS |
| AT-003 | `init` creates config | `caseforge init` in empty dir | `.caseforge.yaml` created | ✅ PASS |

---

### `gen` — Test Case Generation

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-004 | gen hurl format | `caseforge gen --spec petstore.yaml --format hurl --output ./cases` | 17 test cases, `.hurl` files written | ✅ PASS |
| AT-005 | gen json format | `caseforge gen --spec petstore.yaml --format json --output ./cases` | valid JSON with `$schema`, `version`, `test_cases[]` | ✅ PASS |
| AT-006 | gen postman format | `caseforge gen --spec petstore.yaml --format postman --output ./cases` | Postman collection file written | ✅ PASS |
| AT-007 | gen k6 format | `caseforge gen --spec petstore.yaml --format k6 --output ./cases` | k6 JS script written | ✅ PASS |
| AT-008 | gen csv format | `caseforge gen --spec petstore.yaml --format csv --output ./cases` | CSV file written | ✅ PASS |
| AT-009 | gen markdown format | `caseforge gen --spec petstore.yaml --format markdown --output ./cases` | Markdown file written | ✅ PASS |
| AT-010 | gen --no-ai flag | `caseforge gen --spec petstore.yaml --no-ai --format hurl` | generates without LLM, same count | ✅ PASS |
| AT-011 | gen invalid spec path | `caseforge gen --spec nonexistent.yaml` | error: file not found | ✅ PASS |

---

### `gen` — CLI Flags (P1-1 to P1-4)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-067 | gen --technique filters output | `caseforge gen --spec petstore.yaml --no-ai --technique equivalence_partitioning --output ./cases` | only `equivalence_partitioning` cases in index.json | ✅ PASS |
| AT-068 | gen --priority filters output | `caseforge gen --spec petstore.yaml --no-ai --priority P1 --output ./cases` | index.json contains only P0/P1 cases | ✅ PASS |
| AT-069 | gen --operations filters spec | `caseforge gen --spec petstore.yaml --no-ai --operations listPets --output ./cases` | only cases for listPets operationId | ✅ PASS |
| AT-070 | gen --concurrency flag accepted | `caseforge gen --help` | `--concurrency` listed | ✅ PASS |

---

### `gen` — index.json Metadata (P1-6 to P1-10)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-071 | index.json contains meta object | `caseforge gen --spec petstore.yaml --no-ai --output ./cases` | `meta` key present in index.json | ✅ PASS |
| AT-072 | meta.spec_hash is non-empty SHA256 | same as AT-071 | `meta.spec_hash` is 64-char hex string | ✅ PASS |
| AT-073 | meta.caseforge_version present | same as AT-071 | `meta.caseforge_version` is non-empty | ✅ PASS |
| AT-074 | meta.by_technique counts match cases | same as AT-071 | `meta.by_technique` sums to total case count | ✅ PASS |
| AT-075 | meta.by_kind counts match cases | same as AT-071 | `meta.by_kind` sums to total case count | ✅ PASS |

---

### `gen` — Assertion Operators (P1-11 to P1-13)

| ID | Scenario | Expected | Status |
|----|----------|----------|--------|
| AT-076 | `exists` operator used in response assertions | Generated cases for an endpoint with object response include `exists` assertions | ✅ PASS |
| AT-077 | `is_uuid` operator used for uuid-format fields | Cases for endpoint with `format: uuid` response field use `is_uuid` operator | ✅ PASS |
| AT-078 | `is_iso8601` operator used for date-time fields | Cases for endpoint with `format: date-time` response field use `is_iso8601` operator | ✅ PASS |

---

### `rbt --generate` — High-Risk Auto-Generation (2.2)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-079 | `--generate` flag registered | `caseforge rbt --help` | `--generate` listed in help | ✅ PASS |
| AT-080 | `--generate --dry-run` prints "ignored" info message | `caseforge rbt --spec petstore.yaml --dry-run --generate` | output contains "ignored with" | ✅ PASS |
| AT-081 | `--generate` writes index.json for real high-risk op | git repo with changed handler.go mapped via caseforge-map.yaml, run `rbt --generate --no-ai` | `index.json` created in cases dir | ✅ PASS |

### `rbt index --strategy embed` — Embed Phase (2.3)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-082 | `rbt index --strategy embed` writes map file (regex fallback without API key) | `caseforge rbt index --spec openapi.yaml --src /tmp/src --strategy embed` with no `OPENAI_API_KEY` | map file written with `mappings:` | ✅ PASS |

### `caseforge export` — Platform Adapters (3.2)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-083 | `export` command registered | `caseforge --help` | `export` listed | ✅ PASS |
| AT-084 | `--format allure` creates Allure result file | `caseforge export --cases ./cases --format allure --output /tmp/out` | `*-result.json` in `/tmp/out/allure/` | ✅ PASS |
| AT-085 | `--format xray` creates xray-import.json | `caseforge export --cases ./cases --format xray --output /tmp/out` | `xray-import.json` in `/tmp/out/xray/` | ✅ PASS |
| AT-086 | `--format testrail` creates testrail-import.csv | `caseforge export --cases ./cases --format testrail --output /tmp/out` | `testrail-import.csv` in `/tmp/out/testrail/` | ✅ PASS |
| AT-087 | `--technique example_extraction` generates cases from spec examples | `caseforge gen --spec example-spec.yaml --no-ai --technique example_extraction` | Output contains `example_extraction` | ✅ PASS |
| AT-088 | Example extraction produces valid (P1) and invalid (P2) cases | Run `--technique example_extraction` on spec with named examples | `.hurl` output contains example name `valid_widget` | ✅ PASS |

---

### `caseforge diff --gen-cases` — Auto-generate for Breaking Changes (3.3)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-089 | `--gen-cases` flag registered | `caseforge diff --help` | `--gen-cases` listed | ✅ PASS |
| AT-090 | breaking changes → `index.json` written | `caseforge diff --old v1.yaml --new v2.yaml --gen-cases /tmp/gen` | `index.json` with `test_cases` array in gen dir | ✅ PASS |

---

<!-- AT-091–AT-092 reserved for future features -->

### `caseforge suite` — TestSuite Orchestration (3.6)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-093 | `suite` command registered | `caseforge --help` | `suite` listed | ✅ PASS |
| AT-094 | `suite create` writes valid suite.json | `caseforge suite create --id S --title T --kind chain --cases TC-001,TC-002` | `suite.json` with `$schema` and `cases` array | ✅ PASS |
| AT-095 | `suite validate` confirms valid suite | `caseforge suite validate --suite suite.json` | `valid ✓` output | ✅ PASS |

---

<!-- AT-091–AT-092 reserved for future features -->

### Assertion Operator Rendering (1.3 completeness)

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-096 | `gen` produces index.json with assertions | `caseforge gen --no-ai` on numeric+uuid+datetime spec | `assertions` key present in index.json | ✅ PASS |
| AT-097 | Hurl output has no unrendered assertions | `caseforge gen --no-ai --format hurl` | No `# unrendered assertion` comment in any `.hurl` file | ✅ PASS |
| AT-098 | k6 output has no unrendered assertions | `caseforge gen --no-ai --format k6` | No `// unrendered:` comment in k6 output | ✅ PASS |

---

### Phase 2 CLI commands — watch / stats / ci

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-099 | `stats` command registered | `caseforge --help` | `stats` listed | ✅ PASS |
| AT-100 | `stats` reads index.json and prints summary | `caseforge stats --cases <dir>` with valid index.json | Output contains total count and `方法论` | ✅ PASS |
| AT-101 | `stats --format json` outputs valid JSON | `caseforge stats --cases <dir> --format json` | Valid JSON with `total` field | ✅ PASS |
| AT-102 | `watch` command registered | `caseforge --help` | `watch` listed | ✅ PASS |
| AT-103 | `ci` command and `ci init` subcommand registered | `caseforge ci --help` | `init` listed | ✅ PASS |
| AT-104 | `ci init --platform github-actions` generates workflow | `caseforge ci init --platform github-actions --output <file>` | File contains `caseforge lint` and `caseforge gen` | ✅ PASS |

### MCP tools & assertion enhancements

| ID | Scenario | Command / Setup | Expected | Status |
|----|----------|-----------------|----------|--------|
| AT-105 | MCP server has `lint_spec` tool | `caseforge mcp --help` | `lint_spec` in output | ✅ PASS |
| AT-106 | MCP server has `lint_spec` and `ask_test_cases` tools registered | `go test ./internal/mcp/... -run TestServerHas` | All tool registration tests pass | ✅ PASS |
| AT-107 | Email format field maps to `matches` assertion | `go test ./internal/assert/... -run TestSchemaAssertions_EmailFormatUsesMatches` | PASS | ✅ PASS |
| AT-108 | Schema `minimum`/`maximum` constraints generate `gte`/`lte` assertions | `go test ./internal/assert/... -run TestRangeAssertions` | PASS | ✅ PASS |
| AT-109 | `classification_tree` technique applies when enum/boolean params present | `go test ./internal/methodology/... -run TestClassificationTree` | All classification tree tests pass | ✅ PASS |
| AT-110 | `orthogonal_array` technique generates L4/L8/L27 arrays for 3–13 params | `go test ./internal/methodology/... -run TestOrthogonalArray\|TestSelectOA\|TestExtractOA\|TestLevelTo` | All orthogonal array tests pass | ✅ PASS |
| AT-111 | DEA seeds array constraints (minItems/maxItems) and format violations | `go test ./internal/dea/... -run TestSeedHypotheses_Array\|TestSeedHypotheses_Format\|TestSeedHypotheses_RequiredQuery` | All new seeder tests pass | ✅ PASS |
| AT-112 | DEA infers rules for array, required query param, and format violation hypotheses | `go test ./internal/dea/... -run TestInferRule_Array\|TestInferRule_Required\|TestInferRule_Format` | All new inferencer tests pass | ✅ PASS |
| AT-113 | TUI shows completed operations list (scrolls last 12 rows) | `go test ./internal/tui/... -run TestProgressModel_ViewShows\|TestProgressModel_ViewScrolls\|TestProgressModel_WindowSize\|TestProgressModel_OperationDone` | All TUI enhanced tests pass | ✅ PASS |
| AT-114 | Checkpoint Manager saves / loads / deletes state.json | `go test ./internal/checkpoint/... -v` | All 8 checkpoint tests pass | ✅ PASS |
| AT-115 | gen --resume flag and --operations/--technique/--format tab completion registered | `go test ./cmd/... -run TestGenResume\|TestGenCompletion` | All gen UX tests pass | ✅ PASS |
| AT-116 | `score` command registered | `caseforge --help` | `score` listed | ✅ PASS |
| AT-117 | `score` scores test cases across four dimensions | `go test ./cmd/... -run TestScoreCommand_TerminalOutput` | Output contains `Overall:`, `Coverage Breadth`, `Boundary Coverage`, `Security Coverage`, `Executability` | ✅ PASS |
| AT-118 | `score --format json` outputs valid JSON report | `go test ./cmd/... -run TestScoreCommand_JSONOutput` | Valid JSON with `overall`, `dimensions`, `total_cases` fields | ✅ PASS |
| AT-119 | `score` generates improvement suggestions for missing security/boundary cases | `go test ./cmd/... -run TestScoreCommand_OutputContainsSuggestions` | Output contains `Suggestions` and `owasp` | ✅ PASS |
| AT-120 | gen flag behavioral tests (--no-ai, --technique, --priority, --operations, --resume) | `go test ./cmd/... -run 'TestGen_NoAI\|TestGen_Technique\|TestGen_Priority\|TestGen_Operations\|TestGen_Resume\|TestGen_CombinedFlags\|TestGen_Format'` | All 19 gen e2e behavioral tests pass | ✅ PASS |
| AT-121 | webhook package unit tests (sender retry, HMAC signing, event filtering) | `go test ./internal/webhook/... -v` | All 14 webhook unit tests pass | ✅ PASS |
| AT-122 | gen fires on_generate and on_run_complete webhook events | `go test ./cmd/... -run 'TestGenWebhook'` | All 4 webhook integration tests pass | ✅ PASS |
| AT-123 | isolated_negative generates one-invalid-field cases | `go test ./cmd/... -run 'TestGen_IsolatedNegative'` | All isolated_negative cases generated, technique field set | ✅ PASS |
| AT-124 | schema_violation generates comprehensive constraint cases | `go test ./cmd/... -run 'TestGen_SchemaViolation'` | All schema_violation cases generated with 422 assertions | ✅ PASS |
| AT-125 | variable_irrelevance detects dependency groups and generates NA cases | `go test ./cmd/... -run 'TestGen_VariableIrrelevance'` | No error even when technique doesn't apply | ✅ PASS |
| AT-126 | pairwise --tuple-level 3 generates 3-way combinations | `go test ./cmd/... -run 'TestGen_TupleLevel3'` | --tuple-level=3 accepted without error | ✅ PASS |
| AT-127 | --seed produces deterministic output across runs | `go test ./cmd/... -run 'TestGen_Seed_Deterministic'` | Same seed produces same number of cases | ✅ PASS |
| AT-128 | pairwise filters infeasible cross-variable combinations | `go test ./internal/methodology/... -run 'TestPairwise_Filter'` | Infeasible sort=false+sort_field combinations removed | ✅ PASS |
| AT-129 | mutation technique generates per-field invalid-value cases | `go test ./internal/methodology/... -run 'TestMutationTechnique'` | All 4 mutation tests pass | ✅ PASS |
| AT-130 | auth_chain technique generates login→token→use chain cases | `go test ./internal/methodology/... -run 'TestAuthChainTechnique'` | All 6 auth_chain tests pass | ✅ PASS |
| AT-131 | engine maxCasesPerOp ceiling truncates by priority | `go test ./internal/methodology/... -run 'TestEngine_MaxCasesPerOp'` | Ceiling enforced, P0 prioritized | ✅ PASS |
| AT-132 | chain command registers and has required flags | `go test ./cmd/... -run 'TestChainCommand_IsRegistered\|TestChainCommand_HasRequiredFlags'` | chain command present with spec/depth/output flags | ✅ PASS |
| AT-133 | chain depth-2 generates multi-step chain cases | `go test ./cmd/... -run 'TestChainCommand_GeneratesChainCases'` | chain cases with ≥2 steps generated | ✅ PASS |
| AT-134 | chain depth-1 generates single-op cases | `go test ./cmd/... -run 'TestChainCommand_Depth1_SingleOpCases'` | Each case has exactly 1 step | ✅ PASS |
| AT-135 | chain invalid depth exits non-zero | `go test ./cmd/... -run 'TestChainCommand'` | Error returned for depth 0 | ✅ PASS |
| AT-136 | N-step chain includes update step when PUT present | `go test ./internal/methodology/... -run 'TestChainTechnique_NStepChain'` | 4-step chain: setup→update→test→teardown | ✅ PASS |
| AT-137 | gen registers mutation and auth_chain techniques without error | `go test ./cmd/... -run 'TestGen_Seed_DeterministicOutput'` | Deterministic output with new techniques | ✅ PASS |
| AT-138 | OpenAPI Links parsed into Operation.Links | `go test ./internal/spec/... -run 'TestParsedSpec_LinksPopulated'` | Links slice populated with name, operationId, parameters | ✅ PASS |
| AT-139 | OpenAPI Links create dep-graph edges | `go test ./internal/methodology/... -run 'TestBuildDepGraph_OpenAPILinks'` | Edge with correct creator/consumer/pathParam/captureFrom | ✅ PASS |
| AT-140 | BFS chain appends DELETE teardown for non-DELETE consumers | `go test ./cmd/... -run TestChainCommand_AddsTeardownForNonDeleteChains` | Chain case contains step with type=teardown | ✅ PASS |
| AT-141 | DataPool Add/ValueFor/Save/Load/Merge unit tests pass | `go test ./internal/datagen/... -run TestDataPool` | All 5 DataPool tests pass | ✅ PASS |
| AT-142 | explore --export-pool writes pool JSON in dry-run | `go test ./cmd/... -run TestExploreCommand_ExportPool_DryRun` | pool.json created | ✅ PASS |
| AT-143 | chain --data-pool loads pool without error | `go test ./cmd/... -run TestChainCommand_DataPool_Loaded` | index.json produced, no error | ✅ PASS |
| AT-144 | score includes Status Coverage dimension | `go test ./internal/score/... -run TestComputeStatusCoverage` | Status Coverage dimension present with correct score | ✅ PASS |
| AT-145 | Postman collection parsing extracts body fields into DataPool | `go test ./internal/datagen/... -run TestParsePostmanCollection` | All 3 postman tests pass | ✅ PASS |
| AT-146 | chain --seed-postman loads collection without error | `go test ./cmd/... -run TestChainCommand_SeedPostman` | index.json produced | ✅ PASS |
| AT-147 | explore --prioritize-uncovered dry-run reports probes | `go test ./internal/dea/... -run TestExplorer_PrioritizeUncovered_DryRun` | TotalProbes > 0 | ✅ PASS |
| AT-148 | explore --prioritize-uncovered flag accepted without error | `go test ./cmd/... -run TestExploreCommand_PrioritizeUncoveredFlag` | No error returned | ✅ PASS |
| AT-149 | FilterSet unit tests pass | `go test ./internal/spec/... -run TestFilterSet` | All FilterSet tests pass | ✅ PASS |
| AT-150 | gen --include-path filters operations | `go test ./cmd/... -run TestBuildFilterSet` | buildFilterSet returns correct FilterSet | ✅ PASS |
| AT-151 | gen --exclude-tag flag accepted without error | `go test ./cmd/... -run TestGenCommand_HasFilterFlags` | flags registered on genCmd | ✅ PASS |
| AT-152 | response_check unit tests pass | `go test ./internal/dea/... -run TestFindResponseSchema\|TestCheckResponseBody\|TestValidateProbeResponse` | All response check tests pass | ✅ PASS |
| AT-153 | explore discovers response schema mismatch rule | `go test ./internal/dea/... -run TestExplorer_ResponseSchemaMismatch_ProducesRule` | DiscoveredRule with category response_schema_mismatch | ✅ PASS |
| AT-154 | constraint_mutation generates null injection cases | `go test ./internal/methodology/... -run TestConstraintMutationTechnique_Generate_NullInjection` | null case present with status_code eq 422 | ✅ PASS |
| AT-155 | constraint_mutation generates wrong-content-type case | `go test ./internal/methodology/... -run TestConstraintMutationTechnique_Generate_WrongContentType` | case with Content-Type: text/plain, expects 415 | ✅ PASS |
| AT-156 | boundary_value cases carry named Scenario field | `go test ./internal/methodology/... -run TestBoundaryGeneratesMinMaxCases` | cases have non-empty Source.Scenario | ✅ PASS |
| AT-157 | score Boundary Coverage detail shows scenario info | `go test ./internal/score/... -run TestBoundaryDetail_IncludesScenarioInfo` | detail contains covered scenario names | ✅ PASS |
| AT-158 | explore --max-failures stops after N rules | `go test ./internal/dea/... -run TestExplorer_MaxFailures_StopsEarly` | len(Rules) ≤ MaxFailures | ✅ PASS |
| AT-159 | explore deduplicates rules with same operation+category+fieldPath | `go test ./internal/dea/... -run TestExplorer_RuleDeduplication` | no duplicate (op, category, fieldPath) triplets in report | ✅ PASS |
| AT-160 | type_coercion Applies() for op with typed fields | `go test ./internal/methodology/... -run TestTypeCoercionTechnique_Applies_True` | PASS | ✅ PASS |
| AT-161 | type_coercion generates WRONG_TYPE cases for all typed fields | `go test ./internal/methodology/... -run TestTypeCoercionTechnique_Generate_StringField` | PASS | ✅ PASS |
| AT-162 | unicode_fuzzing Applies for op with string field | `go test ./internal/methodology/... -run 'TestUnicodeFuzzingTechnique_Applies_True'` | PASS | ✅ PASS |
| AT-163 | unicode_fuzzing generates 5 cases per string field | `go test ./internal/methodology/... -run 'TestUnicodeFuzzingTechnique_Generate_ProducesExactly5PerStringField'` | PASS | ✅ PASS |
| AT-164 | mass_assignment Applies for op with request body | `go test ./internal/methodology/... -run 'TestMassAssignmentTechnique_Applies_True'` | PASS | ✅ PASS |
| AT-165 | mass_assignment generates 4 category cases | `go test ./internal/methodology/... -run 'TestMassAssignmentTechnique_Generate_ProducesExactly4Cases'` | PASS | ✅ PASS |
| AT-166 | idor Applies for op with integer path param | `go test ./internal/methodology/... -run 'TestIDORTechnique_Applies_IntPathParam'` | PASS | ✅ PASS |
| AT-167 | idor generates 2 cases for integer ID param | `go test ./internal/methodology/... -run 'TestIDORTechnique_Generate_IntegerID_Produces2Cases'` | PASS | ✅ PASS |
| AT-168 | explore seeds KindTypeCoercion hypotheses for typed fields | `go test ./internal/dea/... -run 'TestSeedHypotheses_IncludesTypeCoercion' -count=1` | PASS | ✅ PASS |
| AT-169 | explore seeds KindMassAssignment hypothesis for body ops | `go test ./internal/dea/... -run 'TestSeedHypotheses_IncludesMassAssignment' -count=1` | PASS | ✅ PASS |
| AT-170 | datagen generates pattern-matching strings for simple patterns | `go test ./internal/datagen/... -run 'TestGenerateByPattern_Digits' -count=1` | PASS | ✅ PASS |
| AT-171 | datagen falls back gracefully on invalid patterns | `go test ./internal/datagen/... -run 'TestGenerateByPattern_InvalidPattern' -count=1` | PASS | ✅ PASS |
| AT-172 | semantic_annotation generates cases for nullable fields | `caseforge gen --spec cmd/testdata/semantic.yaml --no-ai --technique semantic_annotation \| grep NULLABLE_ACCEPTANCE` | output contains NULLABLE_ACCEPTANCE | ✅ PASS |
| AT-173 | nullable acceptance case expects 2xx | `go test ./internal/methodology/... -run 'TestSemanticAnnotationTechnique_Generate_NullableCase_Expects2xx' -count=1` | PASS | ✅ PASS |
| AT-174 | readOnly write rejection case expects 4xx | `go test ./internal/methodology/... -run 'TestSemanticAnnotationTechnique_Generate_ReadOnlyCase_Expects4xx' -count=1` | PASS | ✅ PASS |
| AT-175 | writeOnly read suppression case has jsonpath assertion | `go test ./internal/methodology/... -run 'TestSemanticAnnotationTechnique_Generate_WriteOnlyCase_FieldAbsent' -count=1` | PASS | ✅ PASS |
| AT-176 | schema ReadOnly field parsed from spec | `go test ./internal/spec/... -run 'TestSemanticAnnotation_ReadOnly_Parsed' -count=1` | PASS | ✅ PASS |
| AT-177 | schema WriteOnly field parsed from spec | `go test ./internal/spec/... -run 'TestSemanticAnnotation_WriteOnly_Parsed' -count=1` | PASS | ✅ PASS |
| AT-178 | field_boundary Applies for op with constrained fields | `caseforge gen --spec cmd/testdata/field_boundary.yaml --no-ai --technique field_boundary \| grep field_boundary` | output contains field_boundary | ✅ PASS |
| AT-179 | field_boundary generates cases for constrained fields | `go test ./internal/methodology/... -run 'TestFieldBoundaryTechnique_Generate_4CasesPerConstrainedField' -count=1` | PASS | ✅ PASS |
| AT-180 | field_boundary generates nested field path cases | `caseforge gen --spec cmd/testdata/field_boundary.yaml --no-ai --technique field_boundary \| grep address.zip` | output contains address.zip | ✅ PASS |
| AT-181 | field_boundary valid cases expect 2xx assertions | `go test ./internal/methodology/... -run 'TestFieldBoundaryTechnique_Generate_ValidBoundaryExpects2xx' -count=1` | PASS | ✅ PASS |
| AT-182 | required_omission Applies for op with required fields | `caseforge gen --spec cmd/testdata/required_omission.yaml --no-ai --technique required_omission \| grep required_omission` | output contains required_omission | ✅ PASS |
| AT-183 | required_omission generates one case per required field | `go test ./internal/methodology/... -run 'TestRequiredOmissionTechnique_Generate_OneCasePerRequiredField' -count=1` | PASS | ✅ PASS |
| AT-184 | required_omission case has field absent (REQUIRED_OMISSION scenario) | `caseforge gen --spec cmd/testdata/required_omission.yaml --no-ai --technique required_omission \| grep REQUIRED_OMISSION` | output contains REQUIRED_OMISSION | ✅ PASS |
| AT-185 | required_omission cases expect 4xx | `go test ./internal/methodology/... -run 'TestRequiredOmissionTechnique_Generate_Expects4xx' -count=1` | PASS | ✅ PASS |
| AT-186 | positive_examples Applies for op with parameter examples | `caseforge gen --spec cmd/testdata/positive_examples.yaml --no-ai --technique positive_examples \| grep positive_examples` | output contains positive_examples | ✅ PASS |
| AT-187 | positive_examples generates one case per named example | `go test ./internal/methodology/... -run 'TestPositiveExamplesTechnique_Generate_OneCasePerNamedExample' -count=1` | PASS | ✅ PASS |
| AT-188 | positive_examples substitutes path param value in URL | `caseforge gen --spec cmd/testdata/positive_examples.yaml --no-ai --technique positive_examples \| grep /users/42` | output contains /users/42 | ✅ PASS |
| AT-189 | positive_examples expects 2xx assertions | `go test ./internal/methodology/... -run 'TestPositiveExamplesTechnique_Generate_Expects2xx' -count=1` | PASS | ✅ PASS |
| AT-190 | chain_crud detects POST+GET+DELETE chain | `caseforge gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_crud \| grep chain_crud` | output contains chain_crud | ✅ PASS |
| AT-191 | chain_crud generates kind=chain test case | `caseforge gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_crud \| grep '"kind":"chain"'` | output contains "kind":"chain" | ✅ PASS |
| AT-192 | chain_crud setup step captures created id | `caseforge gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_crud \| grep '"captures"'` | output contains captures | ✅ PASS |
| AT-193 | chain_crud test step uses captured id in path | `caseforge gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_crud \| grep '{{id}}'` | output contains {{id}} | ✅ PASS |
| AT-194 | chain_crud teardown step is DELETE | `caseforge gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_crud \| grep '"type":"teardown"'` | output contains "type":"teardown" | ✅ PASS |
| AT-195 | chain_crud source scenario is CRUD_FLOW | `go test ./internal/methodology/... -run 'TestChainTechnique_Source_ScenarioCRUDFlow' -count=1` | PASS | ✅ PASS |

---

### `import har` — HAR Traffic Replay

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-196 | import subcommand registered | `caseforge import --help` | help output contains "import" | ✅ PASS |
| AT-197 | import har subcommand registered | `caseforge import har --help` | help output contains "har" | ✅ PASS |
| AT-198 | import har parses entries from HAR file | `caseforge import har cmd/testdata/sample.har --output /tmp/har_test_198` | output directory contains files | ✅ PASS |
| AT-199 | import har strips noise headers | `caseforge import har cmd/testdata/sample.har --output /tmp/har_test_199` | user-agent absent from generated output | ✅ PASS |
| AT-200 | import har deduplicates identical METHOD+PATH entries | `caseforge import har cmd/testdata/sample.har --output /tmp/har_test_200` | exactly 2 output files (POST /users deduplicated) | ✅ PASS |
| AT-201 | import har writes test cases to output directory | `caseforge import har cmd/testdata/sample.har --output /tmp/har_test_201` | output directory is non-empty | ✅ PASS |
| AT-202 | score --format json includes conformance block | `caseforge score --cases cmd/testdata/score_cases --format json` | output contains `"conformance"` | ✅ PASS |
| AT-203 | score --format json conformance has trend field | `caseforge score --cases cmd/testdata/score_cases --format json` | output contains `"trend"` | ✅ PASS |
| AT-204 | score terminal output shows conformance trend | `caseforge score --cases cmd/testdata/score_cases` | output contains `"trend:"` or `"Conformance"` | ✅ PASS |
| AT-205 | score --min-score passes when score meets threshold | `caseforge score --cases cmd/testdata/score_cases --min-score 0` | exit code 0 | ✅ PASS |
| AT-206 | score --min-score fails when score below threshold | `caseforge score --cases cmd/testdata/score_cases --min-score 200` | exit code non-zero | ✅ PASS |
| AT-207 | score --save-history writes .caseforge-conformance.json | `cd /tmp && caseforge score --cases /Users/yuchou/Github/yuchou87/caseforge/cmd/testdata/score_cases --save-history && test -f .caseforge-conformance.json` | file exists | ✅ PASS |
| AT-208 | gen --auth-bootstrap exits 0 (spec without security) | `caseforge gen --spec cmd/testdata/crud.yaml --no-ai --technique equivalence_partitioning --auth-bootstrap --output /tmp/at208` | exit code 0 | ✅ PASS |
| AT-209 | gen --auth-bootstrap skips when no auth op in spec | `caseforge gen --spec cmd/testdata/field_boundary.yaml --no-ai --auth-bootstrap --output /tmp/at209` | exits 0, no crash | ✅ PASS |
| AT-210 | gen --auth-bootstrap output directory non-empty | `caseforge gen --spec cmd/testdata/crud.yaml --no-ai --technique equivalence_partitioning --auth-bootstrap --output /tmp/at210` | output directory non-empty | ✅ PASS |
| AT-211 | gen --auth-bootstrap preserves non-secured op cases | `caseforge gen --spec cmd/testdata/field_boundary.yaml --no-ai --auth-bootstrap --output /tmp/at211` | output directory non-empty (cases generated normally) | ✅ PASS |
| AT-212 | classifyFailure unit tests pass | `go test ./cmd/ -run TestClassifyFailure -count=1` | all 4 classification tests pass | ✅ PASS |
| AT-213 | classifyFailure server_error for happy-path technique | `go test ./cmd/ -run TestClassifyFailure_ServerError -count=1` | PASS | ✅ PASS |
| AT-214 | classifyFailure missing_validation for mutation technique | `go test ./cmd/ -run TestClassifyFailure_MissingValidation -count=1` | PASS | ✅ PASS |
| AT-215 | classifyFailure security_regression for owasp technique | `go test ./cmd/ -run TestClassifyFailure_SecurityRegression -count=1` | PASS | ✅ PASS |

---

### `gen` — Technique Coverage

| ID | Scenario | Expected Techniques | Status |
|----|----------|---------------------|--------|
| AT-012 | Equivalence partitioning cases | `equivalence_partitioning` cases generated for GET/POST/DELETE | ✅ PASS |
| AT-013 | OWASP API Top 10 security cases | `owasp_api_top10` and `owasp_api_top10_spec` cases generated | ✅ PASS |
| AT-014 | Idempotency chain cases | `idempotency` technique → `kind: "chain"` with 2 steps (setup + test), `step-test.depends_on = ["step-setup"]` | ✅ PASS |
| AT-015 | CRUD chain cases | `chain_crud` technique generated for POST+GET+DELETE operations | ✅ PASS |

---

### `lint` — Spec Linting

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-016 | lint valid spec | `caseforge lint --spec petstore.yaml` | score reported, warnings for missing descriptions | ✅ PASS |
| AT-017 | lint missing operationId | `caseforge lint --spec bad.yaml` (no operationId) | `[L001]` warning reported | ✅ PASS |

---

### `diff` — Spec Diff

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-018 | diff identical specs | `caseforge diff --old spec.yaml --new spec.yaml` | `No changes detected.` | ✅ PASS |
| AT-019 | diff with breaking changes | `caseforge diff --old v1.yaml --new v2.yaml` | BREAKING: removed endpoints listed, NON-BREAKING: new endpoints listed | ✅ PASS |

---

### `doctor` — Environment Check

| ID | Scenario | Expected | Status |
|----|----------|----------|--------|
| AT-020 | doctor checks tools and API keys | reports hurl/k6 found/missing, API key set/not set per provider | ✅ PASS |

---

### `fake` — Data Generation

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-021 | fake from JSON schema | `caseforge fake --schema '{"type":"object","properties":{"name":{"type":"string"},"age":{"type":"integer"}}}'` | valid JSON object with name (string) and age (integer) | ✅ PASS |

---

### `pairwise` — Pairwise Combinations

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-022 | pairwise 3 params | `caseforge pairwise --params "browser:chrome,firefox os:win,mac lang:en,zh"` | 4 combinations (< full factorial 8), all pairs covered | ✅ PASS |

---

### `completion` — Shell Completion

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-023 | bash completion | `caseforge completion bash` | bash completion script output | ✅ PASS |
| AT-024 | zsh completion | `caseforge completion zsh` | zsh completion script output | ✅ PASS |
| AT-025 | fish completion | `caseforge completion fish` | fish completion script output | ✅ PASS |

---

### `config show` — Configuration Display

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-026 | config show defaults (no file) | `caseforge config show` in dir without `.caseforge.yaml` | shows defaults: provider=noop, format=hurl | ✅ PASS |
| AT-027 | config show with file | `caseforge config show` in dir with `.caseforge.yaml` | shows config values, API key masked as `sk-ant...` (first 6 chars + `...`) | ✅ PASS |

---

### `ask` — NL-to-Test-Case

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-028 | ask with no LLM (noop provider) | `caseforge ask "POST /users"` with noop config | error: AI provider unavailable | ✅ PASS |
| AT-029 | ask with live LLM (gemini) | `caseforge ask "POST /users create user with name and email"` | generates test cases, writes to `./cases` | ✅ PASS |

---

### `explore` — Dynamic Exploration Agent

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-035 | explore command registered | `caseforge --help` | `explore` listed | ✅ PASS |
| AT-036 | explore --dry-run produces report | `caseforge explore --spec petstore.yaml --dry-run --output ./reports` | dea-report.json written, planned rules listed | ✅ PASS |
| AT-037 | explore missing --spec returns error | `caseforge explore --target http://x` | error: --spec is required | ✅ PASS |
| AT-038 | explore missing --target returns error | `caseforge explore --spec petstore.yaml` | error: --target is required (or use --dry-run) | ✅ PASS |

---

### `rbt` — Regression-Based Testing

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-039 | rbt command registered | `caseforge --help` | `rbt` listed | ✅ PASS |
| AT-040 | missing --spec returns error | `caseforge rbt` | error message | ✅ PASS |
| AT-041 | --format json + --dry-run produces valid JSON | see script | rbt-report.json with diff_base field | ✅ PASS |
| AT-042 | --fail-on high, dry-run → exit 0 | see script | exit 0 | ✅ PASS |
| AT-043 | --dry-run skips git/tree-sitter | see script | no git errors | ✅ PASS |
| AT-044 | doctor shows tree-sitter status | `caseforge doctor` | tree-sitter line present | ✅ PASS |
| AT-045 | rbt index command registered | `caseforge rbt --help` | `index` listed | ✅ PASS |
| AT-046 | rbt index --strategy llm writes map file | see script | map.yaml created with mappings: | ✅ PASS |
| AT-044b | rbt index --out existing without --overwrite fails | see script | error: already exists | ✅ PASS |

---

### `rbt` — Call Graph (V2)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-061 | --depth flag registered on rbt index | `caseforge rbt index --help` | `--depth` listed | ✅ PASS |
| AT-062 | rbt --dry-run exits 0 | `caseforge rbt --spec petstore.yaml --dry-run` | exit 0, report generated | ✅ PASS |
| AT-063 | --depth default is 0 on rbt index | `caseforge rbt index --help` | output contains `depth int` | ✅ PASS |

---

### `rbt` — Call Graph V3 (Go type-aware)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-064 | --algo flag registered on rbt index | `caseforge rbt index --help` | `--algo` listed | ✅ PASS |
| AT-065 | rbt index hybrid runs without error (no Go module) | `caseforge rbt index --spec petstore.yaml --strategy hybrid --src /tmp` | exit 0, map file written | ✅ PASS |
| AT-066 | --algo vta flag accepted | `caseforge rbt index --help` | `vta` mentioned in --algo description | ✅ PASS |

---

### `dedupe` — Duplicate Test Case Detection

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-047 | dedupe command registered | `caseforge --help` | `dedupe` listed | ✅ PASS |
| AT-048 | no cases dir returns error | `caseforge dedupe --cases /nonexistent/xyz/cases` | error: cases | ✅ PASS |
| AT-049 | no duplicates exits 0 | `caseforge dedupe --cases <unique-cases-dir>` | exit 0 | ✅ PASS |
| AT-050 | exact duplicate reports group | `caseforge dedupe --cases <dup-cases-dir>` | output contains `Group 1` | ✅ PASS |
| AT-051 | --dry-run exits 0 and files still exist | `caseforge dedupe --cases <dup-cases-dir> --dry-run` | exit 0, both files present | ✅ PASS |
| AT-052 | --merge exits 0 and deletes lower-scoring file | `caseforge dedupe --cases <dup-cases-dir> --merge` | exit 0, lower-scoring file removed | ✅ PASS |

---

### `onboard` — Setup Wizard

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-030 | onboard --yes (non-interactive) | `ANTHROPIC_API_KEY=x caseforge onboard --yes` | auto-selects anthropic, writes `.caseforge.yaml` with `provider: anthropic`, prints Next steps | ✅ PASS |
| AT-031 | onboard skip existing config | `echo n \| caseforge onboard` with existing `.caseforge.yaml` | prompts overwrite, answers n → file unchanged | ✅ PASS |

---

### `run` — Test Execution

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-032 | run hurl (no server) | `caseforge run --cases ./cases --format hurl` | hurl exits with `base_url not set` error (expected without server) | ✅ PASS |
| AT-033 | run k6 (no server) | `caseforge run --cases ./cases --format k6` | k6 exits with connection refused (expected without server) | ✅ PASS |
| AT-034 | run non-existent dir | `caseforge run --cases /nonexistent --format k6` | error: file not found | ✅ PASS |
| AT-053 | run --target injects BASE_URL | `caseforge run --cases ./cases --target http://localhost:9999` | BASE_URL injected (hurl error mentions base_url) | ✅ PASS |
| AT-054 | run --output writes run-report.json | `caseforge run --cases ./cases --target http://localhost:9999 --output ./reports` | `run-report.json` created | ✅ PASS |

---

### `lint` — Enhancement (AT-055–AT-060)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-055 | lint --format json outputs valid JSON | `caseforge lint --spec petstore.yaml --format json` | parseable JSON with `score` and `issues` keys | ✅ PASS |
| AT-056 | lint --output writes lint-report.json | `caseforge lint --spec petstore.yaml --output /tmp/lr` | `/tmp/lr/lint-report.json` created | ✅ PASS |
| AT-057 | lint --skip-rules suppresses rule | `caseforge lint --spec petstore.yaml --skip-rules L014 --format json` | L014 absent from issues | ✅ PASS |
| AT-058 | .caseforgelint.yaml skip_rules respected | `.caseforgelint.yaml` with `skip_rules: [L014]`, run lint | L014 absent from output | ✅ PASS |
| AT-059 | L016 duplicate operationId detected | spec with two operations sharing same operationId | error L016 reported | ✅ PASS |
| AT-060 | L020 sensitive query param detected | spec with `?token` query parameter | error L020 reported | ✅ PASS |

---

### Exit Codes (P1-15, P1-16)

| ID | Scenario | Command | Expected | Status |
|----|----------|---------|----------|--------|
| AT-071 | lint exits 3 when errors found | `caseforge lint --spec <spec-with-errors>` (spec with duplicate operationId) | exit code 3 | ✅ PASS |
| AT-072 | gen exits 4 when LLM unavailable without --no-ai | `caseforge gen --spec petstore.yaml` (no API key, no --no-ai) | exit code 4, error message about --no-ai | ✅ PASS |

---


## Summary (last run: 2026-04-12)

| Category | Total | Pass | Fail |
|----------|-------|------|------|
| Core / CLI | 3 | 3 | 0 |
| gen — formats | 7 | 7 | 0 |
| gen — techniques | 4 | 4 | 0 |
| gen — CLI flags | 4 | 4 | 0 |
| gen — metadata | 5 | 5 | 0 |
| gen — assertion operators | 3 | 3 | 0 |
| lint | 2 | 2 | 0 |
| lint enhancement | 6 | 6 | 0 |
| diff | 2 | 2 | 0 |
| doctor | 1 | 1 | 0 |
| fake | 1 | 1 | 0 |
| pairwise | 1 | 1 | 0 |
| completion | 3 | 3 | 0 |
| config show | 2 | 2 | 0 |
| ask | 2 | 2 | 0 |
| explore | 4 | 4 | 0 |
| rbt | 15 | 15 | 0 |
| dedupe | 6 | 6 | 0 |
| onboard | 2 | 2 | 0 |
| run | 9 | 9 | 0 |
| exit codes | 2 | 2 | 0 |
| example_extraction | 2 | 2 | 0 |
| score | 10 | 10 | 0 |
| idor | 2 | 2 | 0 |
| datagen pattern | 2 | 2 | 0 |
| positive_examples | 4 | 4 | 0 |
| import har | 6 | 6 | 0 |
| gen — auth-bootstrap | 4 | 4 | 0 |
| **Total** | **126** | **126** | **0** |

---

## Issues Found During This Run

| Issue | Root Cause | Fix Applied |
|-------|-----------|-------------|
| `ask`, `config`, `completion` commands missing from main | PR #13 code lost when main was force-rebuilt | Restored from `feature/ask-config-completion` worktree, committed `07faf10` |
| Idempotency cases were `kind: "single"` with 1 step instead of `kind: "chain"` with 2 steps | PR #12 chain implementation lost in main rebuild | Restored from `feature/idempotent-chain` worktree, committed `bf5c25c` |

---

## Maintenance Guide

**When adding a new feature:**
1. Add scenario row(s) to the relevant table above
2. Run `./scripts/acceptance.sh` (or manually run the scenario)
3. Update the Summary table counts
4. Update the "last run" date

**When a scenario fails:**
- Code bug → fix code, update status to ✅ PASS after fix
- Environment issue (no server, no API key) → document in Expected column as "expected failure"
- Infrastructure issue → note in Issues Found table
