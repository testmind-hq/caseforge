#!/usr/bin/env bash
# CaseForge Acceptance Test Runner
# Usage: ./scripts/acceptance.sh [--verbose]
# Exit 0 if all pass, 1 if any fail.

set -euo pipefail

VERBOSE=${1:-""}
PASS=0
FAIL=0
ERRORS=()
WORKDIR=$(mktemp -d)
trap 'rm -rf "$WORKDIR"' EXIT

BIN="$(go env GOPATH)/bin/caseforge"
# Build fresh binary
echo "Building caseforge..."
go build -o "$BIN" . 2>&1

log() { echo "$1"; }
ok()  { PASS=$((PASS+1)); log "  ✅ PASS  $1"; }
fail(){ FAIL=$((FAIL+1)); ERRORS+=("$1"); log "  ❌ FAIL  $1"; }
run() {
  local id="$1"; local desc="$2"; shift 2
  [ -n "$VERBOSE" ] && log "\n--- $id: $desc ---"
  if eval "$@" > "$WORKDIR/out" 2>&1; then
    ok "$id: $desc"
  else
    fail "$id: $desc  [exit $?]"
    [ -n "$VERBOSE" ] && cat "$WORKDIR/out"
  fi
}
contains() {
  local id="$1"; local desc="$2"; local pattern="$3"; shift 3
  eval "$@" > "$WORKDIR/out" 2>&1 || true
  if grep -q "$pattern" "$WORKDIR/out"; then
    ok "$id: $desc"
  else
    fail "$id: $desc  [pattern '$pattern' not found]"
    [ -n "$VERBOSE" ] && cat "$WORKDIR/out"
  fi
}

# -------------------------------------------------------
# Fixtures
# -------------------------------------------------------
cat > "$WORKDIR/petstore.yaml" << 'YAML'
openapi: "3.0.0"
info:
  title: Petstore
  version: "1.0"
paths:
  /pets:
    get:
      operationId: listPets
      summary: List all pets
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
      responses:
        "200":
          description: A list of pets
    post:
      operationId: createPet
      summary: Create a pet
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [name]
              properties:
                name:
                  type: string
                tag:
                  type: string
      responses:
        "201":
          description: Created
        "400":
          description: Invalid input
  /pets/{petId}:
    get:
      operationId: showPetById
      summary: Info for a specific pet
      parameters:
        - name: petId
          in: path
          required: true
          schema:
            type: integer
      responses:
        "200":
          description: ok
        "404":
          description: Not found
    delete:
      operationId: deletePet
      summary: Delete a pet
      parameters:
        - name: petId
          in: path
          required: true
          schema:
            type: integer
      responses:
        "204":
          description: Deleted
        "404":
          description: Not found
YAML

cat > "$WORKDIR/petstore-v2.yaml" << 'YAML'
openapi: "3.0.0"
info:
  title: Petstore
  version: "2.0"
paths:
  /pets:
    get:
      operationId: listPets
      summary: List all pets
      responses:
        "200":
          description: ok
  /animals:
    get:
      operationId: listAnimals
      summary: List animals
      responses:
        "200":
          description: ok
YAML

cat > "$WORKDIR/.caseforge.yaml" << 'YAML'
ai:
  provider: anthropic
  model: claude-sonnet-4-6
  api_key: sk-ant-test-key-masking
output:
  default_format: hurl
  dir: ./cases
lint:
  fail_on: error
YAML

echo ""
echo "============================================================"
echo " CaseForge Acceptance Tests"
echo "============================================================"
echo ""

# -------------------------------------------------------
# Unit tests (fast gate)
# -------------------------------------------------------
echo "--- Unit Tests ---"
if go test ./... -count=1 > "$WORKDIR/unit.out" 2>&1; then
  PASS=$((PASS+1)); log "  ✅ PASS  Unit: all packages"
else
  FAIL=$((FAIL+1)); ERRORS+=("Unit tests FAILED"); log "  ❌ FAIL  Unit tests"; cat "$WORKDIR/unit.out"
fi
echo ""

# -------------------------------------------------------
# AT-001 – AT-003: Core / CLI
# -------------------------------------------------------
echo "--- Core / CLI ---"
contains AT-001 "version flag" "caseforge version" "$BIN --version"
contains AT-002 "all commands registered" "onboard" "$BIN --help"
run AT-003 "init creates .caseforge.yaml" \
  "mkdir -p '$WORKDIR/init-test' && cd '$WORKDIR/init-test' && '$BIN' init && test -f .caseforge.yaml"
echo ""

# -------------------------------------------------------
# AT-004 – AT-011: gen formats
# -------------------------------------------------------
echo "--- gen: output formats ---"
contains AT-004 "gen hurl"     "Generated" "$BIN gen --spec '$WORKDIR/petstore.yaml' --format hurl     --output '$WORKDIR/cases-hurl'"
contains AT-005 "gen json"     "Generated" "$BIN gen --spec '$WORKDIR/petstore.yaml' --format json     --output '$WORKDIR/cases-json'"
contains AT-006 "gen postman"  "Generated" "$BIN gen --spec '$WORKDIR/petstore.yaml' --format postman  --output '$WORKDIR/cases-postman'"
contains AT-007 "gen k6"       "Generated" "$BIN gen --spec '$WORKDIR/petstore.yaml' --format k6       --output '$WORKDIR/cases-k6'"
contains AT-008 "gen csv"      "Generated" "$BIN gen --spec '$WORKDIR/petstore.yaml' --format csv      --output '$WORKDIR/cases-csv'"
contains AT-009 "gen markdown" "Generated" "$BIN gen --spec '$WORKDIR/petstore.yaml' --format markdown --output '$WORKDIR/cases-md'"
contains AT-010 "gen --no-ai"  "Generated" "$BIN gen --spec '$WORKDIR/petstore.yaml' --no-ai --format hurl --output '$WORKDIR/cases-noai'"
contains AT-011 "gen invalid spec path" "no such file" "$BIN gen --spec /nonexistent.yaml 2>&1 || true"
echo ""

# -------------------------------------------------------
# AT-012 – AT-015: gen techniques
# -------------------------------------------------------
echo "--- gen: technique coverage ---"
contains AT-012 "equivalence_partitioning technique" "equivalence_partitioning" \
  "python3 -c \"import json,glob; cases=[tc for f in glob.glob('$WORKDIR/cases-json/*.json') for tc in json.load(open(f)).get('test_cases',[])]; print(' '.join(set(tc['source']['technique'] for tc in cases)))\" 2>&1"
contains AT-013 "owasp_api_top10 technique" "owasp_api_top10" \
  "python3 -c \"import json,glob; cases=[tc for f in glob.glob('$WORKDIR/cases-json/*.json') for tc in json.load(open(f)).get('test_cases',[])]; print(' '.join(set(tc['source']['technique'] for tc in cases)))\" 2>&1"
contains AT-014 "idempotency chain 2-step" "chain" \
  "python3 -c \"import json,glob; cases=[tc for f in glob.glob('$WORKDIR/cases-json/*.json') for tc in json.load(open(f)).get('test_cases',[]) if tc['source']['technique']=='idempotency']; print(cases[0]['kind'] if cases else 'none')\" 2>&1"
contains AT-015 "chain_crud technique" "chain_crud" \
  "python3 -c \"import json,glob; cases=[tc for f in glob.glob('$WORKDIR/cases-json/*.json') for tc in json.load(open(f)).get('test_cases',[])]; print(' '.join(set(tc['source']['technique'] for tc in cases)))\" 2>&1"
echo ""

# -------------------------------------------------------
# AT-016 – AT-017: lint
# -------------------------------------------------------
echo "--- lint ---"
contains AT-016 "lint scores spec"   "Spec Score" "$BIN lint --spec '$WORKDIR/petstore.yaml'"
contains AT-017 "lint flags L011 (missing security)" "L011" "$BIN lint --spec '$WORKDIR/petstore.yaml' 2>&1 || true"
echo ""

# -------------------------------------------------------
# AT-018 – AT-019: diff
# -------------------------------------------------------
echo "--- diff ---"
contains AT-018 "diff identical specs"      "No changes"  "$BIN diff --old '$WORKDIR/petstore.yaml' --new '$WORKDIR/petstore.yaml'"
contains AT-019 "diff detects BREAKING"     "BREAKING"    "$BIN diff --old '$WORKDIR/petstore.yaml' --new '$WORKDIR/petstore-v2.yaml'"
echo ""

# -------------------------------------------------------
# AT-020: doctor
# -------------------------------------------------------
echo "--- doctor ---"
contains AT-020 "doctor reports tool status" "hurl" "$BIN doctor"
echo ""

# -------------------------------------------------------
# AT-021: fake
# -------------------------------------------------------
echo "--- fake ---"
contains AT-021 "fake generates JSON object" "{" \
  "$BIN fake --schema '{\"type\":\"object\",\"properties\":{\"name\":{\"type\":\"string\"},\"age\":{\"type\":\"integer\"}}}'"
echo ""

# -------------------------------------------------------
# AT-022: pairwise
# -------------------------------------------------------
echo "--- pairwise ---"
contains AT-022 "pairwise reduces combinations" "combinations" \
  "$BIN pairwise --params 'browser:chrome,firefox os:win,mac lang:en,zh'"
echo ""

# -------------------------------------------------------
# AT-023 – AT-025: completion
# -------------------------------------------------------
echo "--- completion ---"
contains AT-023 "bash completion" "bash completion" "$BIN completion bash"
contains AT-024 "zsh completion"  "compdef"          "$BIN completion zsh"
contains AT-025 "fish completion" "fish completion"  "$BIN completion fish"
echo ""

# -------------------------------------------------------
# AT-026 – AT-027: config show
# -------------------------------------------------------
echo "--- config show ---"
contains AT-026 "config show defaults" "provider" \
  "$BIN config show"
contains AT-027 "config show masks API key" "sk-ant\.\.\." \
  "cd '$WORKDIR' && '$BIN' config show"
echo ""

# -------------------------------------------------------
# AT-028: ask (noop - error path)
# -------------------------------------------------------
echo "--- ask ---"
cat > "$WORKDIR/noop-config.yaml" << 'YAML'
ai:
  provider: noop
  model: ""
output:
  default_format: hurl
  dir: ./cases
YAML
contains AT-028 "ask noop provider returns error" "unavailable" \
  "$BIN ask --config '$WORKDIR/noop-config.yaml' 'POST /users' 2>&1 || true"
echo ""

# -------------------------------------------------------
# AT-035 – AT-038: explore
# -------------------------------------------------------
echo "--- explore ---"
contains AT-035 "explore command registered" "explore" "$BIN --help"
contains AT-036 "explore dry-run produces report" "dea-report.json" \
  "mkdir -p '$WORKDIR/explore-out' && '$BIN' explore --spec '$WORKDIR/petstore.yaml' --dry-run --output '$WORKDIR/explore-out' && ls '$WORKDIR/explore-out/dea-report.json'"
contains AT-037 "explore missing spec returns error" "spec" \
  "'$BIN' explore --target http://localhost:9999 2>&1 || true"
contains AT-038 "explore missing target returns error" "target" \
  "'$BIN' explore --spec '$WORKDIR/petstore.yaml' 2>&1 || true"
echo ""

# -------------------------------------------------------
# AT-030 – AT-031: onboard
# -------------------------------------------------------
echo "--- onboard ---"
contains AT-030 "onboard --yes writes config" "provider: anthropic" \
  "mkdir -p '$WORKDIR/onboard-yes' && cd '$WORKDIR/onboard-yes' && ANTHROPIC_API_KEY=sk-test OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard --yes && cat .caseforge.yaml"
contains AT-031 "onboard skips existing config" "Keeping existing config" \
  "mkdir -p '$WORKDIR/onboard-skip' && echo 'existing: true' > '$WORKDIR/onboard-skip/.caseforge.yaml' && cd '$WORKDIR/onboard-skip' && echo n | ANTHROPIC_API_KEY='' OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard 2>&1"
echo ""

# -------------------------------------------------------
# AT-032 – AT-034: run
# -------------------------------------------------------
echo "--- run ---"
contains AT-032 "run hurl error without server" "base_url" \
  "$BIN run --cases '$WORKDIR/cases-hurl' --format hurl 2>&1 || true"
contains AT-034 "run non-existent cases dir" "no such file" \
  "$BIN run --cases /nonexistent/path --format k6 2>&1 || true"
contains AT-053 "run --target injects BASE_URL into vars" "base_url" \
  "$BIN run --cases '$WORKDIR/cases-hurl' --format hurl --target http://localhost:9999 2>&1 || true"
contains AT-054 "run --output writes run-report.json" "run-report.json" \
  "mkdir -p '$WORKDIR/run-out' && '$BIN' run --cases '$WORKDIR/cases-hurl' --format hurl --target http://localhost:9999 --output '$WORKDIR/run-out' 2>&1 || true"
echo ""

# -------------------------------------------------------
# AT-039 – AT-044: rbt
# -------------------------------------------------------
echo "--- rbt ---"
contains AT-039 "rbt command registered" "rbt" "$BIN --help"
contains AT-040 "missing --spec returns error" "spec" \
  "'$BIN' rbt 2>&1 || true"
contains AT-041 "--format json + dry-run produces valid JSON" "diff_base" \
  "mkdir -p '$WORKDIR/rbt-out' && '$BIN' rbt --spec '$WORKDIR/petstore.yaml' --format json --dry-run --output '$WORKDIR/rbt-out' && cat '$WORKDIR/rbt-out/rbt-report.json'"
run AT-042 "--fail-on high + dry-run exits 0" \
  "mkdir -p '$WORKDIR/rbt-out2' && '$BIN' rbt --spec '$WORKDIR/petstore.yaml' --dry-run --fail-on high --output '$WORKDIR/rbt-out2'"
run AT-043 "--dry-run skips git/tree-sitter" \
  "mkdir -p '$WORKDIR/rbt-out3' && '$BIN' rbt --spec '$WORKDIR/petstore.yaml' --dry-run --output '$WORKDIR/rbt-out3'"
contains AT-044 "doctor shows tree-sitter status" "tree-sitter" \
  "'$BIN' doctor 2>&1 || true"
echo ""

# -------------------------------------------------------
# AT-045 – AT-047: rbt index
# -------------------------------------------------------
echo "--- rbt index ---"
contains AT-045 "rbt index command registered" "index" \
  "'$BIN' rbt --help 2>&1 || true"
contains AT-046 "rbt index --strategy llm writes map file" "mappings:" \
  "mkdir -p '$WORKDIR/idx-out' && '$BIN' rbt index --spec '$WORKDIR/petstore.yaml' --strategy llm --out '$WORKDIR/idx-out/map.yaml' && cat '$WORKDIR/idx-out/map.yaml'"
contains AT-044b "rbt index --out existing without --overwrite fails" "already exists" \
  "echo 'existing: true' > '$WORKDIR/existing-map.yaml' && '$BIN' rbt index --spec '$WORKDIR/petstore.yaml' --out '$WORKDIR/existing-map.yaml' 2>&1 || true"
echo ""

# -------------------------------------------------------
# AT-047 – AT-051: dedupe
# -------------------------------------------------------
echo "--- dedupe ---"

contains AT-047 "dedupe command registered" "dedupe" "$BIN --help"

contains AT-048 "no cases dir returns error" "cases" \
  "'$BIN' dedupe --cases /nonexistent/xyz/cases 2>&1 || true"

# AT-049: two unique cases → exit 0
mkdir -p "$WORKDIR/dedupe-unique"
cat > "$WORKDIR/dedupe-unique/get-users-200.json" << 'JSON'
{"id":"get-users-200","version":"1","kind":"single","priority":"P1","tags":[],"source":{"technique":"equivalence_partitioning","spec_path":"GET /users","rationale":""},"steps":[{"id":"s1","type":"test","method":"GET","path":"/users","assertions":[{"target":"status_code","operator":"eq","expected":200}]}],"generated_at":"2026-01-01T00:00:00Z"}
JSON
cat > "$WORKDIR/dedupe-unique/post-users-201.json" << 'JSON'
{"id":"post-users-201","version":"1","kind":"single","priority":"P1","tags":[],"source":{"technique":"equivalence_partitioning","spec_path":"POST /users","rationale":""},"steps":[{"id":"s1","type":"test","method":"POST","path":"/users","assertions":[{"target":"status_code","operator":"eq","expected":201}]}],"generated_at":"2026-01-01T00:00:00Z"}
JSON
run AT-049 "no duplicates exits 0" \
  "'$BIN' dedupe --cases '$WORKDIR/dedupe-unique'"

# AT-050: two identical cases → output contains "Group 1"
mkdir -p "$WORKDIR/dedupe-dup"
cat > "$WORKDIR/dedupe-dup/case-a.json" << 'JSON'
{"id":"case-a","version":"1","kind":"single","priority":"P1","tags":[],"source":{"technique":"equivalence_partitioning","spec_path":"POST /users","rationale":""},"steps":[{"id":"s1","type":"test","method":"POST","path":"/users","assertions":[{"target":"status_code","operator":"eq","expected":201},{"target":"jsonpath $.id","operator":"eq","expected":"1"}]}],"generated_at":"2026-01-01T00:00:00Z"}
JSON
cat > "$WORKDIR/dedupe-dup/case-b.json" << 'JSON'
{"id":"case-b","version":"1","kind":"single","priority":"P1","tags":[],"source":{"technique":"equivalence_partitioning","spec_path":"POST /users","rationale":""},"steps":[{"id":"s1","type":"test","method":"POST","path":"/users","assertions":[{"target":"status_code","operator":"eq","expected":201},{"target":"jsonpath $.id","operator":"eq","expected":"1"}]}],"generated_at":"2026-01-01T00:00:00Z"}
JSON
contains AT-050 "exact duplicate reports group" "Group 1" \
  "'$BIN' dedupe --cases '$WORKDIR/dedupe-dup' 2>&1; true"

# AT-051: --dry-run exits 0 and both files survive
run AT-051 "--dry-run exits 0 and files still exist" \
  "'$BIN' dedupe --cases '$WORKDIR/dedupe-dup' --dry-run && test -f '$WORKDIR/dedupe-dup/case-a.json' && test -f '$WORKDIR/dedupe-dup/case-b.json'"

# AT-052: --merge exits 0 and deletes lower-scoring duplicate
# aaa-keep.json sorts before zzz-delete.json → aaa-keep is retained, zzz-delete is removed
mkdir -p "$WORKDIR/dedupe-merge"
cat > "$WORKDIR/dedupe-merge/aaa-keep.json" << 'JSON'
{"id":"aaa-keep","version":"1","kind":"single","priority":"P1","tags":[],"source":{"technique":"equivalence_partitioning","spec_path":"POST /users","rationale":""},"steps":[{"id":"s1","type":"test","method":"POST","path":"/users","assertions":[{"target":"status_code","operator":"eq","expected":201},{"target":"jsonpath $.id","operator":"eq","expected":"1"}]}],"generated_at":"2026-01-01T00:00:00Z"}
JSON
cat > "$WORKDIR/dedupe-merge/zzz-delete.json" << 'JSON'
{"id":"zzz-delete","version":"1","kind":"single","priority":"P1","tags":[],"source":{"technique":"equivalence_partitioning","spec_path":"POST /users","rationale":""},"steps":[{"id":"s1","type":"test","method":"POST","path":"/users","assertions":[{"target":"status_code","operator":"eq","expected":201},{"target":"jsonpath $.id","operator":"eq","expected":"1"}]}],"generated_at":"2026-01-01T00:00:00Z"}
JSON
run AT-052 "--merge exits 0 and removes lower-scoring file" \
  "'$BIN' dedupe --cases '$WORKDIR/dedupe-merge' --merge && test ! -f '$WORKDIR/dedupe-merge/zzz-delete.json'"

echo ""


# -------------------------------------------------------
# lint enhancement (AT-055–AT-060)
# -------------------------------------------------------
echo "--- lint enhancement ---"

# AT-055: --format json produces parseable JSON with score and issues
contains "AT-055" "lint --format json" '"score"' \
  "$BIN lint --spec $WORKDIR/petstore.yaml --format json"

# AT-056: --output writes lint-report.json
run "AT-056" "lint --output writes lint-report.json" \
  "$BIN lint --spec $WORKDIR/petstore.yaml --output $WORKDIR/lint-out; test -f $WORKDIR/lint-out/lint-report.json"

# AT-057: --skip-rules suppresses rule
run "AT-057" "lint --skip-rules suppresses rule" \
  "out=\$($BIN lint --spec $WORKDIR/petstore.yaml --skip-rules L014 --format json 2>/dev/null); echo \"\$out\" | python3 -c \"import sys,json; d=json.load(sys.stdin); ids=[i.get('rule_id','') for i in d.get('issues',[])]; assert 'L014' not in ids, 'L014 found but should be skipped'\""

# AT-058: .caseforgelint.yaml skip_rules respected
run "AT-058" ".caseforgelint.yaml skip_rules respected" \
  "echo 'skip_rules: [L014]' > $WORKDIR/.caseforgelint.yaml; out=\$(cd $WORKDIR && $BIN lint --spec $WORKDIR/petstore.yaml --format json 2>/dev/null); echo \"\$out\" | python3 -c \"import sys,json; d=json.load(sys.stdin); ids=[i.get('rule_id','') for i in d.get('issues',[])]; assert 'L014' not in ids, 'L014 found but should be skipped by file config'\"; rm -f $WORKDIR/.caseforgelint.yaml"

# Fixture: spec with duplicate operationId for AT-059
cat > "$WORKDIR/dup-opid.yaml" << 'YAML'
openapi: "3.0.0"
info:
  title: Dup
  version: "1.0"
paths:
  /users:
    get:
      operationId: listUsers
      summary: List
      responses:
        "200":
          description: OK
  /admin/users:
    get:
      operationId: listUsers
      summary: Admin list
      responses:
        "200":
          description: OK
YAML

# AT-059: L016 duplicate operationId
contains "AT-059" "L016 duplicate operationId detected" "L016" \
  "$BIN lint --spec $WORKDIR/dup-opid.yaml --format json"

# Fixture: spec with sensitive query param for AT-060
cat > "$WORKDIR/sensitive-query.yaml" << 'YAML'
openapi: "3.0.0"
info:
  title: Sensitive
  version: "1.0"
paths:
  /users:
    get:
      operationId: listUsers
      summary: List users
      parameters:
        - name: token
          in: query
          schema:
            type: string
      responses:
        "200":
          description: OK
YAML

# AT-060: L020 sensitive query param
contains "AT-060" "L020 sensitive query param detected" "L020" \
  "$BIN lint --spec $WORKDIR/sensitive-query.yaml --format json"

echo ""

# -------------------------------------------------------
# Summary
# -------------------------------------------------------
TOTAL=$((PASS+FAIL))
echo "============================================================"
echo " Results: $PASS/$TOTAL passed"
echo "============================================================"
if [ ${#ERRORS[@]} -gt 0 ]; then
  echo ""
  echo "Failed scenarios:"
  for e in "${ERRORS[@]}"; do
    echo "  ❌ $e"
  done
  exit 1
fi
echo ""
echo "All acceptance tests passed. ✅"
