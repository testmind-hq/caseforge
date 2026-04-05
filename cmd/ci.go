// cmd/ci.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "CI/CD integration helpers",
}

var ciInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a CI workflow configuration for CaseForge",
	Long: `CI init generates a ready-to-use CI workflow file that covers the
full lint → gen → run pipeline for your chosen platform.

Supported platforms: github-actions, gitlab-ci, jenkins, shell

Examples:
  caseforge ci init --platform github-actions
  caseforge ci init --platform gitlab-ci --spec ./api/openapi.yaml
  caseforge ci init --platform shell --output ./scripts/api-test.sh`,
	RunE:         runCIInit,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(ciCmd)
	ciCmd.AddCommand(ciInitCmd)
	ciInitCmd.Flags().String("platform", "github-actions", "Target platform: github-actions|gitlab-ci|jenkins|shell")
	ciInitCmd.Flags().String("spec", "openapi.yaml", "Path to OpenAPI spec used in the generated workflow")
	ciInitCmd.Flags().String("output", "", "Output file path (default: platform-specific standard path)")
	ciInitCmd.Flags().Bool("force", false, "Overwrite existing file without prompting")
}

func runCIInit(cmd *cobra.Command, _ []string) error {
	platform, _ := cmd.Flags().GetString("platform")
	specPath, _ := cmd.Flags().GetString("spec")
	outputPath, _ := cmd.Flags().GetString("output")
	force, _ := cmd.Flags().GetBool("force")

	gen, ok := ciGenerators[platform]
	if !ok {
		return fmt.Errorf("unknown --platform %q: must be github-actions, gitlab-ci, jenkins, or shell", platform)
	}

	content := gen(specPath)

	if outputPath == "" {
		outputPath = ciDefaultPath[platform]
	}

	// Guard against silent overwrite of an existing file
	if _, statErr := os.Stat(outputPath); statErr == nil && !force {
		return fmt.Errorf("%s already exists — use --force to overwrite", outputPath)
	}

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", outputPath, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Generated CI config → %s\n", outputPath)
	return nil
}

var ciDefaultPath = map[string]string{
	"github-actions": ".github/workflows/api-test.yml",
	"gitlab-ci":      ".gitlab-ci.yml",
	"jenkins":        "Jenkinsfile",
	"shell":          "scripts/api-test.sh",
}

var ciGenerators = map[string]func(spec string) string{
	"github-actions": generateGitHubActions,
	"gitlab-ci":      generateGitLabCI,
	"jenkins":        generateJenkinsfile,
	"shell":          generateShellScript,
}

func generateGitHubActions(specPath string) string {
	return fmt.Sprintf(`# 由 caseforge ci init --platform github-actions 生成
name: API Test

on:
  push:
    paths: ['%s', 'openapi/**']
  pull_request:
    paths: ['%s', 'openapi/**']

jobs:
  api-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install CaseForge
        run: |
          curl -sSL https://github.com/testmind-hq/caseforge/releases/latest/download/install.sh | sh

      - name: Lint Spec
        run: caseforge lint --spec ./%s --fail-on error

      - name: Generate Cases
        run: |
          caseforge gen \
            --spec ./%s \
            --technique equivalence_partitioning,boundary_value,owasp_api_top10 \
            --format hurl \
            --output ./cases/
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}

      - name: Run Tests
        run: |
          caseforge run \
            --cases ./cases/ \
            --target ${{ vars.API_BASE_URL }}

      - name: Upload Cases
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: test-cases
          path: cases/
`, specPath, specPath, specPath, specPath)
}

func generateGitLabCI(specPath string) string {
	return fmt.Sprintf(`# 由 caseforge ci init --platform gitlab-ci 生成
stages:
  - lint
  - generate
  - test

variables:
  SPEC_PATH: ./%s
  CASES_DIR: ./cases

before_script:
  - curl -sSL https://github.com/testmind-hq/caseforge/releases/latest/download/install.sh | sh

lint-spec:
  stage: lint
  script:
    - caseforge lint --spec $SPEC_PATH --fail-on error

generate-cases:
  stage: generate
  script:
    - caseforge gen --spec $SPEC_PATH --technique equivalence_partitioning,boundary_value,owasp_api_top10 --format hurl --output $CASES_DIR
  artifacts:
    paths:
      - cases/
    expire_in: 1 week

run-tests:
  stage: test
  script:
    - caseforge run --cases $CASES_DIR --target $API_BASE_URL
  dependencies:
    - generate-cases
`, specPath)
}

func generateJenkinsfile(specPath string) string {
	return fmt.Sprintf(`// 由 caseforge ci init --platform jenkins 生成
pipeline {
    agent any

    environment {
        SPEC_PATH = './%s'
        CASES_DIR = './cases'
        ANTHROPIC_API_KEY = credentials('anthropic-api-key')
    }

    stages {
        stage('Install CaseForge') {
            steps {
                sh 'curl -sSL https://github.com/testmind-hq/caseforge/releases/latest/download/install.sh | sh'
            }
        }

        stage('Lint Spec') {
            steps {
                sh 'caseforge lint --spec $SPEC_PATH --fail-on error'
            }
        }

        stage('Generate Cases') {
            steps {
                sh 'caseforge gen --spec $SPEC_PATH --technique equivalence_partitioning,boundary_value,owasp_api_top10 --format hurl --output $CASES_DIR'
            }
        }

        stage('Run Tests') {
            steps {
                sh 'caseforge run --cases $CASES_DIR --target $API_BASE_URL'
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: 'cases/**', fingerprint: true
        }
    }
}
`, specPath)
}

func generateShellScript(specPath string) string {
	return fmt.Sprintf(`#!/usr/bin/env bash
# 由 caseforge ci init --platform shell 生成
set -euo pipefail

SPEC_PATH="./%s"
CASES_DIR="./cases"
BASE_URL="${API_BASE_URL:-http://localhost:8080}"

echo "==> Installing CaseForge..."
curl -sSL https://github.com/testmind-hq/caseforge/releases/latest/download/install.sh | sh

echo "==> Linting spec..."
caseforge lint --spec "$SPEC_PATH" --fail-on error

echo "==> Generating test cases..."
caseforge gen \
  --spec "$SPEC_PATH" \
  --technique equivalence_partitioning,boundary_value,owasp_api_top10 \
  --format hurl \
  --output "$CASES_DIR"

echo "==> Running tests..."
caseforge run \
  --cases "$CASES_DIR" \
  --var "base_url=$BASE_URL"

echo "==> Done."
`, strings.TrimPrefix(specPath, "./"))
}
