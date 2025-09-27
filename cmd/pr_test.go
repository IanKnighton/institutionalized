package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetPRTemplate(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to the temp directory
	os.Chdir(tempDir)

	// Test 1: No template file exists
	template, err := getPRTemplate()
	if err != nil {
		t.Errorf("Expected no error when no template exists, got: %v", err)
	}
	if template != "" {
		t.Errorf("Expected empty template when no file exists, got: %s", template)
	}

	// Test 2: Create .github directory and template file
	githubDir := filepath.Join(tempDir, ".github")
	os.MkdirAll(githubDir, 0755)

	templateContent := "# Test PR Template\n\n## Description\nTest content"
	templateFile := filepath.Join(githubDir, "pull_request_template.md")
	err = os.WriteFile(templateFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test template file: %v", err)
	}

	template, err = getPRTemplate()
	if err != nil {
		t.Errorf("Expected no error when template exists, got: %v", err)
	}
	if template != templateContent {
		t.Errorf("Expected template content %q, got %q", templateContent, template)
	}

	// Test 3: Test uppercase version
	os.Remove(templateFile)
	uppercaseFile := filepath.Join(githubDir, "PULL_REQUEST_TEMPLATE.md")
	err = os.WriteFile(uppercaseFile, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create uppercase template file: %v", err)
	}

	template, err = getPRTemplate()
	if err != nil {
		t.Errorf("Expected no error when uppercase template exists, got: %v", err)
	}
	if template != templateContent {
		t.Errorf("Expected template content %q, got %q", templateContent, template)
	}

	// Test 4: Priority - lowercase should take precedence
	err = os.WriteFile(templateFile, []byte("lowercase template"), 0644)
	if err != nil {
		t.Fatalf("Failed to create lowercase template file: %v", err)
	}

	template, err = getPRTemplate()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !strings.Contains(template, "lowercase template") {
		t.Errorf("Expected lowercase template to take precedence, got: %q", template)
	}
}
