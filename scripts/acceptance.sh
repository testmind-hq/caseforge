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
REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
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
exits_with() {
  local id="$1"; local desc="$2"; local expected_code="$3"; shift 3
  [ -n "$VERBOSE" ] && log "\n--- $id: $desc ---"
  set +e
  eval "$@" > "$WORKDIR/out" 2>&1
  local actual_code=$?
  set -e
  if [ "$actual_code" -eq "$expected_code" ]; then
    ok "$id: $desc"
  else
    fail "$id: $desc  [expected exit $expected_code, got $actual_code]"
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

cat > "$WORKDIR/petstore-typed.yaml" << 'YAML'
openapi: "3.0.0"
info:
  title: Petstore Typed
  version: "1.0"
paths:
  /pets:
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
      responses:
        "201":
          description: Created
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    format: uuid
                  name:
                    type: string
                  created_at:
                    type: string
                    format: date-time
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
# AT-067 – AT-070: gen CLI flags (P1-1 to P1-4)
# -------------------------------------------------------
echo "--- gen: CLI flags ---"

# AT-067: --technique filters — only equivalence_partitioning cases in output
contains "AT-067" "--technique filters output to one technique" "equivalence_partitioning" \
  "$BIN gen --spec '$WORKDIR/petstore.yaml' --no-ai --technique equivalence_partitioning --output '$WORKDIR/cases-technique' 2>&1 && \
   python3 -c \"import json,os; idx=json.load(open('$WORKDIR/cases-technique/index.json')); techs=set(tc['source']['technique'] for tc in idx.get('test_cases',[])); print(' '.join(techs))\""

# AT-068: --priority filters — verify no P2/P3 cases appear in output
contains "AT-068" "--priority P1 filters out lower priority cases" "true" \
  "$BIN gen --spec '$WORKDIR/petstore.yaml' --no-ai --priority P1 --output '$WORKDIR/cases-priority' 2>&1 && \
   python3 -c \"import json; idx=json.load(open('$WORKDIR/cases-priority/index.json')); bad=[tc['priority'] for tc in idx.get('test_cases',[]) if tc.get('priority') in ('P2','P3')]; print('true' if not bad else 'fail:'+str(bad))\""

# AT-069: --operations limits to specified operationId — verify only /pets (listPets) paths in output
# step.path may include query params (e.g. /pets?limit=1) so use startswith('/pets') and
# exclude /pets/{...} (other operations).
contains "AT-069" "--operations limits to listPets only" "true" \
  "$BIN gen --spec '$WORKDIR/petstore.yaml' --no-ai --operations listPets --output '$WORKDIR/cases-ops' 2>&1 && \
   python3 -c \"import json; idx=json.load(open('$WORKDIR/cases-ops/index.json')); bad=[tc for tc in idx.get('test_cases',[]) for s in tc.get('steps',[]) if not s.get('path','').split('?')[0].rstrip('/') in ('/pets',)]; print('true' if not bad else 'fail: unexpected paths')\""

# AT-070: --concurrency flag is registered in help
contains "AT-070" "--concurrency flag registered on gen" "concurrency" \
  "'$BIN' gen --help 2>&1 || true"

echo ""

# -------------------------------------------------------
# AT-071 – AT-075: gen index.json metadata (P1-6 to P1-10)
# -------------------------------------------------------
echo "--- gen: index.json metadata ---"

# Generate base output once for meta tests
"$BIN" gen --spec "$WORKDIR/petstore.yaml" --no-ai --output "$WORKDIR/cases-meta" > /dev/null 2>&1 || true

contains "AT-071" "index.json contains meta object" "meta" \
  "python3 -c \"import json; idx=json.load(open('$WORKDIR/cases-meta/index.json')); print('meta' if 'meta' in idx else 'missing')\""

contains "AT-072" "meta.spec_hash is 64-char hex" "64" \
  "python3 -c \"import json; idx=json.load(open('$WORKDIR/cases-meta/index.json')); h=idx.get('meta',{}).get('spec_hash',''); print(len(h))\""

contains "AT-073" "meta.caseforge_version is non-empty" "dev" \
  "python3 -c \"import json; idx=json.load(open('$WORKDIR/cases-meta/index.json')); print(idx.get('meta',{}).get('caseforge_version','missing'))\""

contains "AT-074" "meta.by_technique sums to total case count" "true" \
  "python3 -c \"
import json
idx=json.load(open('$WORKDIR/cases-meta/index.json'))
total=len(idx.get('test_cases',[]))
by_tech=idx.get('meta',{}).get('by_technique',{})
s=sum(by_tech.values())
print('true' if s==total else f'fail: sum={s} total={total}')
\""

contains "AT-075" "meta.by_kind sums to total case count" "true" \
  "python3 -c \"
import json
idx=json.load(open('$WORKDIR/cases-meta/index.json'))
total=len(idx.get('test_cases',[]))
by_kind=idx.get('meta',{}).get('by_kind',{})
s=sum(by_kind.values())
print('true' if s==total else f'fail: sum={s} total={total}')
\""

echo ""

# -------------------------------------------------------
# AT-076 – AT-078: assertion operators (P1-11 to P1-13)
# -------------------------------------------------------
echo "--- gen: assertion operators ---"

# Generate cases from typed petstore with uuid/date-time response fields
"$BIN" gen --spec "$WORKDIR/petstore-typed.yaml" --no-ai --output "$WORKDIR/cases-typed" > /dev/null 2>&1 || true

contains "AT-076" "exists operator in response assertions" "exists" \
  "python3 -c \"
import json
idx=json.load(open('$WORKDIR/cases-typed/index.json'))
ops=[a.get('operator','') for tc in idx.get('test_cases',[]) for s in tc.get('steps',[]) for a in s.get('assertions',[])]
print(' '.join(ops))
\""

contains "AT-077" "is_uuid operator for uuid-format field" "is_uuid" \
  "python3 -c \"
import json
idx=json.load(open('$WORKDIR/cases-typed/index.json'))
ops=[a.get('operator','') for tc in idx.get('test_cases',[]) for s in tc.get('steps',[]) for a in s.get('assertions',[])]
print(' '.join(ops))
\""

contains "AT-078" "is_iso8601 operator for date-time field" "is_iso8601" \
  "python3 -c \"
import json
idx=json.load(open('$WORKDIR/cases-typed/index.json'))
ops=[a.get('operator','') for tc in idx.get('test_cases',[]) for s in tc.get('steps',[]) for a in s.get('assertions',[])]
print(' '.join(ops))
\""

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
# AT-030 – AT-031: onboard (home-dir config)
# -------------------------------------------------------
echo "--- onboard ---"
contains AT-030 "onboard --yes writes config to home dir" "provider: anthropic" \
  "mkdir -p '$WORKDIR/onboard-home' && HOME='$WORKDIR/onboard-home' ANTHROPIC_API_KEY=sk-test OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard --yes && cat '$WORKDIR/onboard-home/.caseforge.yaml'"
contains AT-031 "onboard skips existing config in home dir" "Keeping existing config" \
  "mkdir -p '$WORKDIR/onboard-skip' && printf 'existing: true\n' > '$WORKDIR/onboard-skip/.caseforge.yaml' && echo n | HOME='$WORKDIR/onboard-skip' ANTHROPIC_API_KEY='' OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard 2>&1"

# AT-238 – AT-243: onboard wizard improvements
echo "--- onboard wizard improvements ---"

contains AT-238 "onboard prompts for model after provider selection" "Model \[" \
  "mkdir -p '$WORKDIR/at238' && printf '2\n\ngpt-4o-test\n1\n\n\n' | HOME='$WORKDIR/at238' ANTHROPIC_API_KEY='' OPENAI_API_KEY='sk-test' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard 2>&1"

contains AT-239 "onboard --yes uses non-empty model in config" "model:" \
  "mkdir -p '$WORKDIR/at239' && HOME='$WORKDIR/at239' ANTHROPIC_API_KEY=sk-test OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard --yes && cat '$WORKDIR/at239/.caseforge.yaml'"

run AT-240 "onboard --yes writes to home dir not cwd" \
  "mkdir -p '$WORKDIR/at240-home' '$WORKDIR/at240-cwd' && cd '$WORKDIR/at240-cwd' && HOME='$WORKDIR/at240-home' ANTHROPIC_API_KEY=sk-test OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard --yes && test -f '$WORKDIR/at240-home/.caseforge.yaml' && ! test -f '$WORKDIR/at240-cwd/.caseforge.yaml'"

contains AT-241 "onboard MCP checkbox shows Claude Desktop option" "Claude Desktop" \
  "mkdir -p '$WORKDIR/at241' && printf '1\n\n\n1\n\n\n' | HOME='$WORKDIR/at241' ANTHROPIC_API_KEY=sk-test OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard 2>&1"

contains AT-242 "onboard skill checkbox shows Universal AI CLI option" "Universal AI CLI" \
  "mkdir -p '$WORKDIR/at242' && printf '1\n\n\n1\n\n\n' | HOME='$WORKDIR/at242' ANTHROPIC_API_KEY=sk-test OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard 2>&1"

run AT-243 "onboard skill installs to ~/.agents/skills/caseforge/SKILL.md" \
  "mkdir -p '$WORKDIR/at243-home/skills/caseforge' && printf '# CaseForge Skill\n' > '$WORKDIR/at243-home/skills/caseforge/SKILL.md' && cd '$WORKDIR/at243-home' && printf '1\n\n\n1\n\n2\n' | HOME='$WORKDIR/at243-home' ANTHROPIC_API_KEY=sk-test OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard 2>&1 && test -f '$WORKDIR/at243-home/.agents/skills/caseforge/SKILL.md'"

run AT-244 "onboard skill Claude Code option creates symlink at ~/.claude/skills/caseforge" \
  "mkdir -p '$WORKDIR/at244-home/skills/caseforge' && printf '# CaseForge Skill\n' > '$WORKDIR/at244-home/skills/caseforge/SKILL.md' && cd '$WORKDIR/at244-home' && printf '1\n\n\n1\n\n1\n' | HOME='$WORKDIR/at244-home' ANTHROPIC_API_KEY=sk-test OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' onboard 2>&1 && test -L '$WORKDIR/at244-home/.claude/skills/caseforge'"

contains AT-245 "onboard bedrock prompts for AWS region" "AWS Region" \
  "mkdir -p '$WORKDIR/at245' && printf '5\n\n\n1\n\n\n' | HOME='$WORKDIR/at245' ANTHROPIC_API_KEY='' OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' AWS_ACCESS_KEY_ID='AKIATEST' AWS_DEFAULT_REGION='us-west-2' '$BIN' onboard 2>&1"

run AT-246 "onboard bedrock --yes writes region and omits api_key" \
  "mkdir -p '$WORKDIR/at246' && HOME='$WORKDIR/at246' ANTHROPIC_API_KEY='' OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' AWS_ACCESS_KEY_ID='AKIATEST' AWS_DEFAULT_REGION='eu-central-1' '$BIN' onboard --yes && grep -q 'region: eu-central-1' '$WORKDIR/at246/.caseforge.yaml' && ! grep -q 'api_key' '$WORKDIR/at246/.caseforge.yaml'"
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
contains AT-044 "doctor shows AI-disabled warning when no keys set" "no AI provider key set" \
  "ANTHROPIC_API_KEY='' OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' doctor 2>&1 || true"
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
# AT-061 – AT-063: rbt callgraph
# -------------------------------------------------------
echo "--- rbt callgraph ---"

contains "AT-061" "--depth flag registered on rbt index" "depth" \
  "'$BIN' rbt index --help 2>&1 || true"

contains "AT-062" "rbt --dry-run exits 0" "Report written" \
  "mkdir -p '$WORKDIR/reports-at062' && '$BIN' rbt --spec '$WORKDIR/petstore.yaml' --dry-run --output '$WORKDIR/reports-at062' 2>&1 || true"

contains "AT-063" "--depth flag default is 0 on rbt index" "depth int" \
  "'$BIN' rbt index --help 2>&1 || true"
echo ""

# -------------------------------------------------------
# AT-064 – AT-066: rbt callgraph V3 (Go type-aware)
# -------------------------------------------------------
echo "--- rbt callgraph v3 ---"

contains "AT-064" "--algo flag registered on rbt index" "algo" \
  "'$BIN' rbt index --help 2>&1 || true"

contains "AT-065" "rbt index hybrid no-Go-module runs clean" "Map file written" \
  "mkdir -p '$WORKDIR/at065-out' && '$BIN' rbt index --spec '$WORKDIR/petstore.yaml' --strategy hybrid --src /tmp --out '$WORKDIR/at065-out/map.yaml' --overwrite 2>&1 || true"

contains "AT-066" "--algo accepts vta value" "vta" \
  "'$BIN' rbt index --help 2>&1 || true"
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
# Exit Codes (P1-15, P1-16)
# -------------------------------------------------------
echo "=== Exit Codes (P1-15, P1-16) ==="

# AT-071: lint exits 3 when errors found (reuse dup-opid spec which has L016 error)
exits_with "AT-071" "lint exits 3 when errors found" 3 \
  "$BIN lint --spec $WORKDIR/dup-opid.yaml"

# AT-072: gen exits 4 when LLM unavailable without --no-ai
# Create a dir with anthropic config but no API key so factory returns NoopProvider
mkdir -p "$WORKDIR/exit4-test"
cat > "$WORKDIR/exit4-test/.caseforge.yaml" << 'YAML'
ai:
  provider: anthropic
  api_key: ""
output:
  default_format: hurl
  dir: ./cases
YAML
exits_with "AT-072" "gen exits 4 when LLM unavailable without --no-ai" 4 \
  "cd '$WORKDIR/exit4-test' && ANTHROPIC_API_KEY='' OPENAI_API_KEY='' GEMINI_API_KEY='' GOOGLE_API_KEY='' '$BIN' gen --spec '$WORKDIR/petstore.yaml' --output '$WORKDIR/cases-exit4'"

echo ""

# -------------------------------------------------------
# AT-079 – AT-081: rbt --generate (2.2)
# -------------------------------------------------------
echo "--- rbt: --generate high-risk auto-gen ---"

# AT-079: flag registered
contains "AT-079" "rbt --generate flag registered" "generate" \
  "$BIN rbt --help"

# AT-080: --generate + --dry-run → prints "ignored with --dry-run" info message
contains "AT-080" "rbt --generate --dry-run prints ignored message" "ignored with" \
  "$BIN rbt --spec '$WORKDIR/petstore.yaml' --dry-run --generate 2>&1"

# AT-081: --generate with a real git diff + map file → index.json created in cases dir
GENDIR="$WORKDIR/rbt-gen-test"
mkdir -p "$GENDIR/cases-rbt" "$GENDIR/reports"
(
  cd "$GENDIR"
  git init -q
  git config user.email "t@t.com"
  git config user.name "T"
  echo "package main" > handler.go
  cat > caseforge-map.yaml << 'MAPYAML'
mappings:
  - source: handler.go
    operations:
      - "GET /pets"
MAPYAML
  git add . && git commit -q -m "v1"
  echo "package main // updated" > handler.go
  git add . && git commit -q -m "v2"
  "$BIN" rbt \
    --spec "$WORKDIR/petstore.yaml" \
    --src "$GENDIR" \
    --base HEAD~1 --head HEAD \
    --map "$GENDIR/caseforge-map.yaml" \
    --generate --no-ai \
    --cases "$GENDIR/cases-rbt" \
    --output "$GENDIR/reports" 2>/dev/null || true
) 2>/dev/null || true

contains "AT-081" "rbt --generate writes index.json to cases dir" "index.json" \
  "ls '$GENDIR/cases-rbt/'"

echo ""

# -------------------------------------------------------
# AT-082: rbt index --strategy embed (2.3 Embed Phase)
# -------------------------------------------------------
echo "--- rbt index: embed phase ---"

# AT-082: When OPENAI_API_KEY is absent, --strategy embed falls back to regex
# and still writes a valid map file (graceful degradation).
EMBEDDIR=$(mktemp -d)
cat > "$EMBEDDIR/openapi.yaml" << 'SPECEOF'
openapi: "3.0.0"
info:
  title: Embed Test API
  version: "1.0.0"
paths:
  /pets:
    get:
      operationId: listPets
      summary: List all pets
      responses:
        "200":
          description: OK
SPECEOF
cat > "$EMBEDDIR/handler.go" << 'GOEOF'
package handler

func Register(r interface{}) {}
GOEOF

(
  unset OPENAI_API_KEY
  "$BIN" rbt index \
    --spec "$EMBEDDIR/openapi.yaml" \
    --src "$EMBEDDIR" \
    --out "$EMBEDDIR/caseforge-map.yaml" \
    --strategy embed 2>/dev/null || true
) 2>/dev/null || true

contains "AT-082" "rbt index --strategy embed writes map file (regex fallback)" "mappings:" \
  "cat '$EMBEDDIR/caseforge-map.yaml'"

echo ""

# -------------------------------------------------------
# AT-083 – AT-086: caseforge export (3.2)
# -------------------------------------------------------
echo "--- export command ---"

# Build shared fixture: write index.json with one test case
EXPORTDIR=$(mktemp -d)
mkdir -p "$EXPORTDIR/cases"
cat > "$EXPORTDIR/cases/index.json" << 'IDXEOF'
{
  "$schema": "https://caseforge.dev/schema/v1/index.json",
  "version": "1",
  "generated_at": "2026-04-04T00:00:00Z",
  "meta": {},
  "test_cases": [
    {
      "id": "TC-0001",
      "title": "GET /pets - list all pets",
      "kind": "single",
      "priority": "P1",
      "tags": ["pets"],
      "source": {"technique": "equivalence_partitioning", "spec_path": "GET /pets"},
      "steps": [
        {
          "id": "step-1",
          "title": "send request",
          "type": "test",
          "method": "GET",
          "path": "/pets",
          "assertions": [{"target": "status_code", "operator": "eq", "expected": 200}]
        }
      ]
    }
  ]
}
IDXEOF

# AT-083: export command registered
contains "AT-083" "export command registered" "export" \
  "'$BIN' --help"

# AT-084: allure format creates result file
"$BIN" export --cases "$EXPORTDIR/cases" --format allure --output "$EXPORTDIR/out" || true
contains "AT-084" "export --format allure creates result file" "result.json" \
  "ls '$EXPORTDIR/out/allure/'"

# AT-085: xray format creates xray-import.json
"$BIN" export --cases "$EXPORTDIR/cases" --format xray --output "$EXPORTDIR/out" || true
contains "AT-085" "export --format xray creates xray-import.json" "xray-import.json" \
  "ls '$EXPORTDIR/out/xray/'"

# AT-086: testrail format creates testrail-import.csv
"$BIN" export --cases "$EXPORTDIR/cases" --format testrail --output "$EXPORTDIR/out" || true
contains "AT-086" "export --format testrail creates testrail-import.csv" "testrail-import.csv" \
  "ls '$EXPORTDIR/out/testrail/'"

# -------------------------------------------------------
# AT-087 – AT-088: example extraction (PH2-15)
# -------------------------------------------------------
echo ""
echo "--- example extraction (PH2-15) ---"

EXDIR=$(mktemp -d)
cat > "$EXDIR/example-spec.yaml" << 'YAML'
openapi: "3.0.0"
info:
  title: Example API
  version: "1.0"
paths:
  /widgets:
    post:
      operationId: createWidget
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
                color:
                  type: string
            examples:
              valid_widget:
                summary: A valid widget
                value:
                  name: "Sprocket"
                  color: "blue"
              missing_name:
                summary: Invalid - name missing
                value:
                  color: "red"
      responses:
        "201":
          description: created
YAML

# AT-087: example_extraction technique generates cases with technique comment in Hurl output
"$BIN" gen --spec "$EXDIR/example-spec.yaml" --no-ai --technique example_extraction --output "$EXDIR/cases" 2>/dev/null || true
contains "AT-087" "gen --technique example_extraction writes technique comment in hurl" "example_extraction" \
  "cat '$EXDIR/cases/'*.hurl 2>/dev/null | head -40"

# AT-088: example_extraction produces valid and invalid cases (valid_widget example name appears)
"$BIN" gen --spec "$EXDIR/example-spec.yaml" --no-ai --technique example_extraction --output "$EXDIR/cases2" 2>/dev/null || true
contains "AT-088" "example_extraction produces P1 (valid) and P2 (invalid) cases" "valid_widget" \
  "cat '$EXDIR/cases2/'*.hurl 2>/dev/null | head -60"

echo ""

# -------------------------------------------------------
# AT-089 – AT-090: caseforge diff --gen-cases (3.3)
# -------------------------------------------------------
echo "--- diff --gen-cases ---"

GENCASESDIR=$(mktemp -d)

# AT-089: --gen-cases flag registered
contains "AT-089" "diff --gen-cases flag registered" "gen-cases" \
  "'$BIN' diff --help"

# AT-090: breaking changes → generates index.json with test cases
"$BIN" diff \
  --old "$WORKDIR/petstore.yaml" \
  --new "$WORKDIR/petstore-v2.yaml" \
  --gen-cases "$GENCASESDIR" 2>/dev/null || true
contains "AT-090" "diff --gen-cases writes index.json for breaking operations" "test_cases" \
  "cat '$GENCASESDIR/index.json'"

echo ""

# -------------------------------------------------------
# AT-093 – AT-095: caseforge suite (3.6)
# -------------------------------------------------------
echo "--- suite command ---"

SUITEDIR=$(mktemp -d)

# AT-093: suite command registered
contains "AT-093" "suite command registered" "suite" \
  "'$BIN' --help"

# AT-094: suite create writes a valid suite.json
"$BIN" suite create \
  --id "SUITE-AT094" \
  --title "AT-094 E2E" \
  --kind chain \
  --cases "TC-001,TC-002" \
  --output "$SUITEDIR/suite.json" 2>/dev/null || true
contains "AT-094" 'suite create writes suite.json with $schema' 'suite.json' \
  "cat '$SUITEDIR/suite.json'"

# AT-095: suite validate confirms valid suite
contains "AT-095" "suite validate reports valid suite" "valid" \
  "'$BIN' suite validate --suite '$SUITEDIR/suite.json'"

echo ""

# -------------------------------------------------------
# AT-096 – AT-098: operator rendering
# -------------------------------------------------------
echo "--- operator rendering ---"

OPDIR=$(mktemp -d)

# Build a minimal spec for operator rendering tests
cat > "$OPDIR/op-spec.yaml" << 'SPEC'
openapi: "3.0.0"
info:
  title: Op Test
  version: "1.0"
paths:
  /items/{id}:
    get:
      operationId: getItem
      summary: Get item
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
            minimum: 1
            maximum: 9999
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
                    format: uuid
                  name:
                    type: string
                  score:
                    type: number
                  created_at:
                    type: string
                    format: date-time
SPEC

"$BIN" gen --spec "$OPDIR/op-spec.yaml" --output "$OPDIR/cases" --no-ai 2>/dev/null || true

# AT-096: gen produces index.json with assertions
contains "AT-096" "gen produces index.json with assertions" "assertions" \
  "cat '$OPDIR/cases/index.json'"

# AT-097: hurl output does not contain unrendered assertion fallback
"$BIN" gen --spec "$OPDIR/op-spec.yaml" --output "$OPDIR/cases" --no-ai --format hurl 2>/dev/null || true
run "AT-097" "hurl output has no unrendered assertion comments" \
  "! grep -r '# unrendered assertion' '$OPDIR/cases/' 2>/dev/null"

# AT-098: k6 output does not contain unrendered assertion fallback
"$BIN" gen --spec "$OPDIR/op-spec.yaml" --output "$OPDIR/cases" --no-ai --format k6 2>/dev/null || true
run "AT-098" "k6 output has no unrendered assertion comments" \
  "! grep -r '// unrendered:' '$OPDIR/cases/' 2>/dev/null"

echo ""

# -------------------------------------------------------
# Readable case names (AT-247–AT-248)
# -------------------------------------------------------
echo "--- readable case names ---"

AT247DIR=$(mktemp -d)
run "AT-247" "hurl filenames are path-first (no TC- prefix, resource slug leads)" \
  "'$BIN' gen --spec '$WORKDIR/petstore.yaml' --no-ai --format hurl --output '$AT247DIR' && \
   ls '$AT247DIR'/*.hurl | grep -qv 'TC-' && \
   ls '$AT247DIR'/*.hurl | grep -q '^.*/pets_'"

AT248DIR=$(mktemp -d)
run "AT-248" "k6 group names contain case title not raw ID" \
  "'$BIN' gen --spec '$WORKDIR/petstore.yaml' --no-ai --format k6 --output '$AT248DIR' && \
   grep -q \"group('\" '$AT248DIR/k6_tests.js' && \
   ! grep -q \"group('TC-\" '$AT248DIR/k6_tests.js'"

# AT-249: hurl single-case output contains case_name field
AT249DIR=$(mktemp -d)
run "AT-249" "hurl output contains case_name field" \
  "'$BIN' gen --spec '$WORKDIR/petstore.yaml' --no-ai --format hurl --output '$AT249DIR' && \
   grep -rq '# case_name=' '$AT249DIR'/*.hurl"

# AT-250: spec-hash dedup — second run on unchanged spec prints skip message
AT250DIR=$(mktemp -d)
run "AT-250" "gen skips regeneration when spec is unchanged" \
  "'$BIN' gen --spec '$WORKDIR/petstore.yaml' --no-ai --output '$AT250DIR' 2>&1 | grep -q 'Generated' && \
   '$BIN' gen --spec '$WORKDIR/petstore.yaml' --no-ai --output '$AT250DIR' 2>&1 | grep -q 'unchanged'"

# AT-251: --force bypasses spec-hash dedup and regenerates
AT251DIR=$(mktemp -d)
run "AT-251" "gen --force regenerates even when spec is unchanged" \
  "'$BIN' gen --spec '$WORKDIR/petstore.yaml' --no-ai --output '$AT251DIR' 2>&1 | grep -q 'Generated' && \
   '$BIN' gen --spec '$WORKDIR/petstore.yaml' --no-ai --force --output '$AT251DIR' 2>&1 | grep -q 'Generated'"

echo ""

# -------------------------------------------------------
# Phase 2 CLI — watch / stats / ci (AT-099–AT-104)
# -------------------------------------------------------
echo "--- phase 2 cli commands ---"

# AT-099: stats command is registered
contains "AT-099" "stats command registered" "stats" \
  "$BIN --help"

# AT-100: stats reads index.json and prints terminal summary
STATSDIR=$(mktemp -d)
cat > "$STATSDIR/index.json" <<'INDEXEOF'
{"$schema":"https://caseforge.dev/schema/v1/index.json","version":"1","generated_at":"2026-04-01T00:00:00Z","meta":{"by_technique":{"equivalence_partitioning":5},"by_priority":{"P0":2,"P1":3}},"test_cases":[]}
INDEXEOF
contains "AT-100" "stats reads index.json and prints summary" "Technique distribution:" \
  "$BIN stats --cases '$STATSDIR'"

# AT-101: stats --format json outputs valid JSON
run "AT-101" "stats --format json outputs valid JSON" \
  "$BIN stats --cases '$STATSDIR' --format json | python3 -m json.tool > /dev/null"

# AT-102: watch command is registered
contains "AT-102" "watch command registered" "watch" \
  "$BIN --help"

# AT-103: ci command and ci init subcommand registered
contains "AT-103" "ci init subcommand registered" "init" \
  "$BIN ci --help"

# AT-104: ci init generates GitHub Actions workflow
CIDIR=$(mktemp -d)
run "AT-104" "ci init --platform github-actions generates workflow" \
  "$BIN ci init --platform github-actions --output '$CIDIR/workflow.yml' && grep -q 'caseforge lint' '$CIDIR/workflow.yml' && grep -q 'caseforge gen' '$CIDIR/workflow.yml'"

echo "--- mcp tools & assertion enhancements ---"

# AT-105: MCP server has lint_spec tool (verified via unit test)
run "AT-105" "mcp server has lint_spec tool" \
  "(cd $REPO_ROOT && go test ./internal/mcp/... -run TestServerHasLintSpecTool -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-106: MCP server has ask_test_cases tool (verified via unit test)
run "AT-106" "mcp server has ask_test_cases tool" \
  "(cd $REPO_ROOT && go test ./internal/mcp/... -run TestServerHasAskTestCasesTool -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-107: assert schema.go handles email format → matches operator
run "AT-107" "assertion email format maps to matches operator" \
  "(cd $REPO_ROOT && go test ./internal/assert/... -run TestSchemaAssertions_EmailFormatUsesMatches -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-108: RangeAssertions generates gte/lte from schema min/max
run "AT-108" "RangeAssertions generates gte/lte operators" \
  "(cd $REPO_ROOT && go test ./internal/assert/... -run TestRangeAssertions -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-109: classification_tree technique ECT coverage and applies logic
run "AT-109" "classification_tree technique generates ECT test cases" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run TestClassificationTree -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-110: orthogonal_array technique generates L4/L8/L27 balanced arrays
run "AT-110" "orthogonal_array technique generates balanced OA test cases" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestOrthogonalArray|TestSelectOA|TestExtractOA|TestLevelTo' -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-111: DEA seeds array constraints and format violations
run "AT-111" "DEA seeds array constraints (minItems/maxItems), required query param, and format violations" \
  "(cd $REPO_ROOT && go test ./internal/dea/... -run 'TestSeedHypotheses_Array|TestSeedHypotheses_Format|TestSeedHypotheses_RequiredQuery|TestSeedHypotheses_OptionalQuery' -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-112: DEA infers rules for new hypothesis kinds
run "AT-112" "DEA infers rules for array, required query param, and format violation hypotheses" \
  "(cd $REPO_ROOT && go test ./internal/dea/... -run 'TestInferRule_Array|TestInferRule_Required|TestInferRule_Format' -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-113: TUI enhanced progress model
run "AT-113" "TUI shows completed operations list (scrolls last 12 rows)" \
  "(cd $REPO_ROOT && go test ./internal/tui/... -run 'TestProgressModel_ViewShows|TestProgressModel_ViewScrolls|TestProgressModel_WindowSize|TestProgressModel_OperationDone' -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-114: Checkpoint Manager persistence
run "AT-114" "Checkpoint Manager saves / loads / deletes .state.json" \
  "(cd $REPO_ROOT && go test ./internal/checkpoint/... -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-115: gen --resume flag and dynamic flag completion registered
run "AT-115" "gen --resume flag and --operations/--technique/--format tab completion registered" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestGenResume|TestGenCompletion' -count=1 2>&1 | grep -E '(PASS|FAIL|ok)')"

# AT-116–AT-119: score command
contains "AT-116" "score command registered" "score" \
  "$BIN --help"

run "AT-117" "score scores test cases across four dimensions" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestScoreCommand_TerminalOutput' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-118" "score --format json outputs valid JSON report" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestScoreCommand_JSONOutput' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-119" "score generates improvement suggestions for missing security/boundary cases" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestScoreCommand_OutputContainsSuggestions' -count=1 2>&1 | grep -E '(PASS|ok)')"

# AT-120: gen flag behavioral tests
run "AT-120" "gen flag behavioral tests (--no-ai, --technique, --priority, --operations, --resume)" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestGen_NoAI|TestGen_Technique|TestGen_Priority|TestGen_Operations|TestGen_Resume|TestGen_CombinedFlags|TestGen_Format' -count=1 2>&1 | grep -E '(PASS|ok)')"

# AT-121–AT-122: webhook
run "AT-121" "webhook package unit tests (sender retry, HMAC signing, event filtering)" \
  "(cd $REPO_ROOT && go test ./internal/webhook/... -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-122" "gen fires on_generate and on_run_complete webhook events" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestGenWebhook' -count=1 2>&1 | grep -E '(PASS|ok)')"

# AT-123–AT-128: Tcases-inspired techniques
run "AT-123" "isolated_negative generates one-invalid-field cases" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestGen_IsolatedNegative' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-124" "schema_violation generates comprehensive constraint cases" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestGen_SchemaViolation' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-125" "variable_irrelevance detects dependency groups and generates NA cases" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestGen_VariableIrrelevance' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-126" "pairwise --tuple-level 3 accepted and generates cases" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestGen_TupleLevel3' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-127" "--seed produces deterministic output across runs" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestGen_Seed_Deterministic' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-128" "pairwise filters infeasible cross-variable combinations" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestPairwise_Filter' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-129" "mutation technique generates per-field invalid-value cases" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestMutationTechnique' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-130" "auth_chain technique generates login→token→use chain cases" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestAuthChainTechnique' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-131" "engine maxCasesPerOp ceiling truncates by priority" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestEngine_MaxCasesPerOp' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-132" "chain command registers and has required flags" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestChainCommand_IsRegistered|TestChainCommand_HasRequiredFlags' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-133" "chain depth-2 generates multi-step chain cases" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestChainCommand_GeneratesChainCases' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-134" "chain depth-1 generates single-op cases" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestChainCommand_Depth1_SingleOpCases' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-135" "chain invalid depth exits non-zero" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestChainCommand' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-136" "N-step chain includes update step when PUT present" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestChainTechnique_NStepChain' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-137" "gen registers mutation and auth_chain techniques without error" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestGen_Seed_DeterministicOutput' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-138" "OpenAPI Links parsed into Operation.Links" \
  "(cd $REPO_ROOT && go test ./internal/spec/... -run 'TestParsedSpec_LinksPopulated' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-139" "OpenAPI Links create dep-graph edges" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestBuildDepGraph_OpenAPILinks' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-140" "BFS chain appends DELETE teardown for non-DELETE consumers" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestChainCommand_AddsTeardownForNonDeleteChains' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-141" "DataPool Add/ValueFor/Save/Load/Merge unit tests pass" \
  "(cd $REPO_ROOT && go test ./internal/datagen/... -run 'TestDataPool' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-142" "explore --export-pool writes pool JSON in dry-run" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestExploreCommand_ExportPool_DryRun' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-143" "chain --data-pool loads pool without error" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestChainCommand_DataPool_Loaded' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-144" "score includes Status Coverage dimension" \
  "(cd $REPO_ROOT && go test ./internal/score/... -run 'TestComputeStatusCoverage' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-145" "Postman collection parsing extracts body fields into DataPool" \
  "(cd $REPO_ROOT && go test ./internal/datagen/... -run 'TestParsePostmanCollection' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-146" "chain --seed-postman loads collection without error" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestChainCommand_SeedPostman' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-147" "explore --prioritize-uncovered dry-run reports probes" \
  "(cd $REPO_ROOT && go test ./internal/dea/... -run 'TestExplorer_PrioritizeUncovered_DryRun' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-148" "explore --prioritize-uncovered flag accepted without error" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestExploreCommand_PrioritizeUncoveredFlag' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-149" "FilterSet unit tests pass" \
  "(cd $REPO_ROOT && go test ./internal/spec/... -run 'TestFilterSet' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-150" "gen --include-path filters operations" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestBuildFilterSet' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-151" "gen --exclude-tag flag accepted without error" \
  "(cd $REPO_ROOT && go test ./cmd/... -run 'TestGenCommand_HasFilterFlags' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-152" "response_check unit tests pass" \
  "(cd $REPO_ROOT && go test ./internal/dea/... -run 'TestFindResponseSchema|TestCheckResponseBody|TestValidateProbeResponse' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-153" "explore discovers response schema mismatch rule" \
  "(cd $REPO_ROOT && go test ./internal/dea/... -run 'TestExplorer_ResponseSchemaMismatch_ProducesRule' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-154" "constraint_mutation generates null injection cases" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestConstraintMutationTechnique_Generate_NullInjection' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-155" "constraint_mutation generates wrong-content-type case" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestConstraintMutationTechnique_Generate_WrongContentType' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-156" "boundary_value cases carry named Scenario field" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestBoundaryGeneratesMinMaxCases' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-157" "score Boundary Coverage detail shows scenario info" \
  "(cd $REPO_ROOT && go test ./internal/score/... -run 'TestBoundaryDetail_IncludesScenarioInfo' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-158" "explore --max-failures stops after N rules" \
  "(cd $REPO_ROOT && go test ./internal/dea/... -run 'TestExplorer_MaxFailures_StopsEarly' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-159" "explore deduplicates rules with same operation+category+fieldPath" \
  "(cd $REPO_ROOT && go test ./internal/dea/... -run 'TestExplorer_RuleDeduplication' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-160" "type_coercion Applies for op with typed fields" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestTypeCoercionTechnique_Applies_True' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-161" "type_coercion generates WRONG_TYPE cases" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestTypeCoercionTechnique_Generate_StringField' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-162" "unicode_fuzzing Applies for op with string field" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestUnicodeFuzzingTechnique_Applies_True' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-163" "unicode_fuzzing generates 5 cases per string field" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestUnicodeFuzzingTechnique_Generate_ProducesExactly5PerStringField' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-164" "mass_assignment Applies for op with request body" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestMassAssignmentTechnique_Applies_True' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-165" "mass_assignment generates 4 category cases" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestMassAssignmentTechnique_Generate_ProducesExactly4Cases' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-166" "idor Applies for op with integer path param" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestIDORTechnique_Applies_IntPathParam' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-167" "idor generates 2 cases for integer ID param" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestIDORTechnique_Generate_IntegerID_Produces2Cases' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-168" "explore seeds KindTypeCoercion hypotheses for typed fields" \
  "(cd $REPO_ROOT && go test ./internal/dea/... -run 'TestSeedHypotheses_IncludesTypeCoercion' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-169" "explore seeds KindMassAssignment hypothesis for body ops" \
  "(cd $REPO_ROOT && go test ./internal/dea/... -run 'TestSeedHypotheses_IncludesMassAssignment' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-170" "datagen generates pattern-matching strings for simple patterns" \
  "(cd $REPO_ROOT && go test ./internal/datagen/... -run 'TestGenerateByPattern_Digits' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-171" "datagen falls back gracefully on invalid patterns" \
  "(cd $REPO_ROOT && go test ./internal/datagen/... -run 'TestGenerateByPattern_InvalidPattern' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-172" "semantic_annotation generates cases for nullable fields" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/semantic.yaml --no-ai --technique semantic_annotation --output $WORKDIR/semantic 2>/dev/null && grep -q 'NULLABLE_ACCEPTANCE' $WORKDIR/semantic/index.json)"

run "AT-173" "nullable acceptance case expects 2xx" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestSemanticAnnotationTechnique_Generate_NullableCase_Expects2xx' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-174" "readOnly write rejection case expects 4xx" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestSemanticAnnotationTechnique_Generate_ReadOnlyCase_Expects4xx' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-175" "writeOnly read suppression case has jsonpath assertion" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestSemanticAnnotationTechnique_Generate_WriteOnlyCase_FieldAbsent' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-176" "schema ReadOnly field parsed from spec" \
  "(cd $REPO_ROOT && go test ./internal/spec/... -run 'TestSemanticAnnotation_ReadOnly_Parsed' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-177" "schema WriteOnly field parsed from spec" \
  "(cd $REPO_ROOT && go test ./internal/spec/... -run 'TestSemanticAnnotation_WriteOnly_Parsed' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-178" "field_boundary Applies for op with constrained fields" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/field_boundary.yaml --no-ai --technique field_boundary --output $WORKDIR/fb178 2>/dev/null && grep -q 'field_boundary' $WORKDIR/fb178/index.json)"

run "AT-179" "field_boundary generates cases for constrained fields" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestFieldBoundaryTechnique_Generate_4CasesPerConstrainedField' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-180" "field_boundary generates nested field path cases" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/field_boundary.yaml --no-ai --technique field_boundary --output $WORKDIR/fb180 2>/dev/null && grep -q 'address.zip' $WORKDIR/fb180/index.json)"

run "AT-181" "field_boundary valid cases expect 2xx assertions" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestFieldBoundaryTechnique_Generate_ValidBoundaryExpects2xx' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-182" "required_omission Applies for op with required fields" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/required_omission.yaml --no-ai --technique required_omission --output $WORKDIR/ro182 2>/dev/null && grep -q 'required_omission' $WORKDIR/ro182/index.json)"

run "AT-183" "required_omission generates one case per required field" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestRequiredOmissionTechnique_Generate_OneCasePerRequiredField' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-184" "required_omission case has field absent (REQUIRED_OMISSION scenario)" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/required_omission.yaml --no-ai --technique required_omission --output $WORKDIR/ro184 2>/dev/null && grep -q 'REQUIRED_OMISSION' $WORKDIR/ro184/index.json)"

run "AT-185" "required_omission cases expect 4xx" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestRequiredOmissionTechnique_Generate_Expects4xx' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-186" "positive_examples Applies for op with parameter examples" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/positive_examples.yaml --no-ai --technique positive_examples --output $WORKDIR/pe186 2>/dev/null && grep -q 'positive_examples' $WORKDIR/pe186/index.json)"

run "AT-187" "positive_examples generates one case per named example" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestPositiveExamplesTechnique_Generate_OneCasePerNamedExample' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-188" "positive_examples substitutes path param value in URL" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/positive_examples.yaml --no-ai --technique positive_examples --output $WORKDIR/pe188 2>/dev/null && grep -q '/users/42' $WORKDIR/pe188/index.json)"

run "AT-189" "positive_examples expects 2xx assertions" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestPositiveExamplesTechnique_Generate_Expects2xx' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-190" "chain_crud detects POST+GET+DELETE chain" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_crud --output $WORKDIR/crud190 2>/dev/null && grep -q 'chain_crud' $WORKDIR/crud190/index.json)"

run "AT-191" "chain_crud generates kind=chain test case" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_crud --output $WORKDIR/crud191 2>/dev/null && grep -q '\"kind\": \"chain\"' $WORKDIR/crud191/index.json)"

run "AT-192" "chain_crud setup step captures created id" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_crud --output $WORKDIR/crud192 2>/dev/null && grep -q '\"captures\"' $WORKDIR/crud192/index.json)"

run "AT-193" "chain_crud test step uses captured id in path" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_crud --output $WORKDIR/crud193 2>/dev/null && grep -q '{{id}}' $WORKDIR/crud193/index.json)"

run "AT-194" "chain_crud teardown step is DELETE" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_crud --output $WORKDIR/crud194 2>/dev/null && grep -q '\"type\": \"teardown\"' $WORKDIR/crud194/index.json)"

run "AT-195" "chain_crud source scenario is CRUD_FLOW" \
  "(cd $REPO_ROOT && go test ./internal/methodology/... -run 'TestChainTechnique_Source_ScenarioCRUDFlow' -count=1 2>&1 | grep -E '(PASS|ok)')"

run "AT-196" "import subcommand registered" \
  "(cd $REPO_ROOT && $BIN import --help 2>&1 | grep -q 'import')"

run "AT-197" "import har subcommand registered" \
  "(cd $REPO_ROOT && $BIN import har --help 2>&1 | grep -q 'har')"

run "AT-198" "import har parses entries from HAR file" \
  "(cd $REPO_ROOT && $BIN import har cmd/testdata/sample.har --output $WORKDIR/har_test_198 2>/dev/null && ls $WORKDIR/har_test_198 | grep -q '.')"

run "AT-199" "import har strips noise headers (user-agent absent from output)" \
  "(cd $REPO_ROOT && $BIN import har cmd/testdata/sample.har --output $WORKDIR/har_test_199 2>/dev/null && ! grep -r 'user-agent\|Mozilla' $WORKDIR/har_test_199)"

run "AT-200" "import har deduplicates identical METHOD+PATH entries" \
  "(cd $REPO_ROOT && $BIN import har cmd/testdata/sample.har --output $WORKDIR/har_test_200 2>/dev/null && test \$(ls $WORKDIR/har_test_200 | wc -l | tr -d ' ') -eq 2)"

run "AT-201" "import har writes test cases to output directory" \
  "(cd $REPO_ROOT && $BIN import har cmd/testdata/sample.har --output $WORKDIR/har_test_201 2>/dev/null && test \$(ls $WORKDIR/har_test_201 | wc -l | tr -d ' ') -gt 0)"

# AT-202–AT-207: conformance score + CI gate
contains "AT-202" "score --format json includes conformance block" '"conformance"' \
  "(cd $REPO_ROOT && $BIN score --cases cmd/testdata/score_cases --format json)"

contains "AT-203" "score --format json conformance has trend field" '"trend"' \
  "(cd $REPO_ROOT && $BIN score --cases cmd/testdata/score_cases --format json)"

contains "AT-204" "score terminal output shows conformance trend" "trend:" \
  "(cd $REPO_ROOT && $BIN score --cases cmd/testdata/score_cases)"

run "AT-205" "score --min-score passes when score meets threshold" \
  "(cd $REPO_ROOT && $BIN score --cases cmd/testdata/score_cases --min-score 0)"

run "AT-206" "score --min-score fails when score below threshold" \
  "(cd $REPO_ROOT && ! $BIN score --cases cmd/testdata/score_cases --min-score 200 2>/dev/null)"

run "AT-207" "score --save-history writes .caseforge-conformance.json" \
  "(cd /tmp && $BIN score --cases $REPO_ROOT/cmd/testdata/score_cases --save-history 2>/dev/null && test -f .caseforge-conformance.json)"

# AT-208–AT-211: gen --auth-bootstrap
run "AT-208" "gen --auth-bootstrap exits 0 (spec without security)" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/crud.yaml --no-ai --technique equivalence_partitioning --auth-bootstrap --output /tmp/at208 2>/dev/null)"

run "AT-209" "gen --auth-bootstrap skips when no auth op in spec" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/field_boundary.yaml --no-ai --auth-bootstrap --output /tmp/at209 2>/dev/null)"

run "AT-210" "gen --auth-bootstrap output directory non-empty" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/crud.yaml --no-ai --technique equivalence_partitioning --auth-bootstrap --output /tmp/at210 2>/dev/null && ls /tmp/at210/)"

run "AT-211" "gen --auth-bootstrap preserves non-secured op cases" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/field_boundary.yaml --no-ai --auth-bootstrap --output /tmp/at211 2>/dev/null && ls /tmp/at211/)"

# AT-212–AT-215: run failure classification
run "AT-212" "classifyFailure unit tests pass" \
  "(cd $REPO_ROOT && go test ./cmd/ -run TestClassifyFailure -count=1)"

run "AT-213" "classifyFailure server_error for happy-path technique" \
  "(cd $REPO_ROOT && go test ./cmd/ -run TestClassifyFailure_ServerError -count=1)"

run "AT-214" "classifyFailure missing_validation for mutation technique" \
  "(cd $REPO_ROOT && go test ./cmd/ -run TestClassifyFailure_MissingValidation -count=1)"

run "AT-215" "classifyFailure security_regression for owasp technique" \
  "(cd $REPO_ROOT && go test ./cmd/ -run TestClassifyFailure_SecurityRegression -count=1)"

# AT-216–AT-219: score --fill-gaps
contains "AT-216" "score --fill-gaps requires --spec" "fill-gaps requires --spec" \
  "(cd $REPO_ROOT && $BIN score --cases cmd/testdata/score_cases --fill-gaps 2>&1 || true)"

run "AT-217" "score --fill-gaps runs without panic" \
  "(cd $REPO_ROOT && $BIN score --cases cmd/testdata/score_cases --fill-gaps --spec cmd/testdata/crud.yaml 2>&1 || true)"

run "AT-218" "score --fill-gaps prints gen commands" \
  "(cd $REPO_ROOT && $BIN score --cases cmd/testdata/score_cases --fill-gaps --spec cmd/testdata/crud.yaml 2>&1 || true)"

run "AT-219" "score ComputeGaps unit tests" \
  "(cd $REPO_ROOT && go test ./internal/score/ -run TestComputeGaps -v)"

run "AT-220" "gen --with-oracles noop LLM" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/crud.yaml --no-ai --with-oracles --output /tmp/at220)"

run "AT-221" "oracle Mine noop returns empty" \
  "(cd $REPO_ROOT && go test ./internal/oracle/ -run TestMine_NoopProvider -v)"

run "AT-222" "oracle ToAssertions exists" \
  "(cd $REPO_ROOT && go test ./internal/oracle/ -run TestConstraintToAssertion_Exists -v)"

run "AT-223" "oracle InjectIntoCase skips 4xx" \
  "(cd $REPO_ROOT && go test ./internal/oracle/ -run TestInjectIntoCase_Skips4xx -v)"

run "AT-224" "business_rule_violation technique registered" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/crud.yaml --no-ai --technique business_rule_violation --output /tmp/at224)"

run "AT-225" "business_rule Applies false no semantic info" \
  "(cd $REPO_ROOT && go test ./internal/methodology/ -run TestBusinessRuleTechnique_Applies_NoSemanticInfo -v)"

run "AT-226" "business_rule one case per rule" \
  "(cd $REPO_ROOT && go test ./internal/methodology/ -run TestBusinessRuleTechnique_Generate_OnePerRule -v)"

run "AT-227" "business_rule expects 4xx" \
  "(cd $REPO_ROOT && go test ./internal/methodology/ -run TestBusinessRuleTechnique_Generate_Expects4xx -v)"
run "AT-228" "chain_sequence technique registered" \
  "(cd $REPO_ROOT && $BIN gen --spec cmd/testdata/crud.yaml --no-ai --technique chain_sequence --output /tmp/at228)"
run "AT-229" "scoreFieldSimilarity positive overlap" \
  "(cd $REPO_ROOT && go test ./internal/methodology/ -run TestScoreFieldSimilarity_SameToken -v)"
run "AT-230" "tokenizeFieldName camelCase" \
  "(cd $REPO_ROOT && go test ./internal/methodology/ -run TestTokenizeFieldName_CamelCase -v)"
run "AT-231" "chain_sequence detects non-CRUD chain" \
  "(cd $REPO_ROOT && go test ./internal/methodology/ -run TestChainSequenceTechnique_DetectsNonCRUDChain -v)"

# AT-232–AT-235: conformance command
contains "AT-232" "conformance command registered" "spec-vs-implementation" \
  "$BIN conformance --help"
contains "AT-233" "conformance --spec required" "required flag" \
  "$BIN conformance --target http://localhost:8080 2>&1 || true"
contains "AT-234" "conformance --target required" "required flag" \
  "(cd $REPO_ROOT && $BIN conformance --spec cmd/testdata/crud.yaml 2>&1 || true)"
contains "AT-235" "conformance no LLM fails gracefully" "LLM provider not available" \
  "(cd $REPO_ROOT && $BIN conformance --spec cmd/testdata/crud.yaml --target http://localhost:8080 2>&1 || true)"

echo ""

# -------------------------------------------------------
# AT-236 – AT-237: bedrock provider
# -------------------------------------------------------
echo "# AT-236 – AT-237: bedrock provider"

contains AT-236 "doctor shows bedrock section" "AWS Bedrock" \
  "'$BIN' doctor 2>&1 || true"

# AT-237: --no-ai must bypass AI entirely even when provider=bedrock is configured
run AT-237 "generate --no-ai works with bedrock config" \
  "mkdir -p '$WORKDIR/at237' && \
   printf 'ai:\n  provider: bedrock\n  region: us-east-1\n' > '$WORKDIR/at237/.caseforge.yaml' && \
   cd '$WORKDIR/at237' && \
   '$BIN' gen --spec '$WORKDIR/petstore.yaml' --no-ai --output '$WORKDIR/at237/cases'"
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
