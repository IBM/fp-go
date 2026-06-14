// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mcp

import (
	"testing"
)

func TestParseSkillMetadata(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		expectedName string
		expectedDesc string
	}{
		{
			name: "valid frontmatter",
			content: `---
name: fp-go-logging
description: Use this skill when working with logging in fp-go
---
# Content here`,
			expectedName: "fp-go-logging",
			expectedDesc: "Use this skill when working with logging in fp-go",
		},
		{
			name: "no frontmatter",
			content: `# fp-go HTTP Requests

## Overview`,
			expectedName: "",
			expectedDesc: "",
		},
		{
			name: "frontmatter with only name",
			content: `---
name: test-skill
---
# Content`,
			expectedName: "test-skill",
			expectedDesc: "",
		},
		{
			name: "frontmatter with only description",
			content: `---
description: Test description
---
# Content`,
			expectedName: "",
			expectedDesc: "Test description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, desc := parseSkillMetadata([]byte(tt.content))
			if name != tt.expectedName {
				t.Errorf("parseSkillMetadata() name = %v, want %v", name, tt.expectedName)
			}
			if desc != tt.expectedDesc {
				t.Errorf("parseSkillMetadata() desc = %v, want %v", desc, tt.expectedDesc)
			}
		})
	}
}

func TestHandleUseSkill(t *testing.T) {
	tests := []struct {
		name        string
		skillName   string
		shouldError bool
	}{
		{
			name:        "valid skill - fp-go",
			skillName:   "fp-go",
			shouldError: false,
		},
		{
			name:        "valid skill - fp-go-http",
			skillName:   "fp-go-http",
			shouldError: false,
		},
		{
			name:        "valid skill - fp-go-lens",
			skillName:   "fp-go-lens",
			shouldError: false,
		},
		{
			name:        "valid skill - fp-go-logging",
			skillName:   "fp-go-logging",
			shouldError: false,
		},
		{
			name:        "valid skill - fp-go-pipe-flow",
			skillName:   "fp-go-pipe-flow",
			shouldError: false,
		},
		{
			name:        "invalid skill",
			skillName:   "nonexistent-skill",
			shouldError: true,
		},
		{
			name:        "empty skill name",
			skillName:   "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := UseSkillArgs{Name: tt.skillName}
			result, output, err := handleUseSkill(nil, nil, args)

			if tt.shouldError {
				if err == nil {
					t.Errorf("handleUseSkill() expected error but got none")
				}
				if !result.IsError {
					t.Errorf("handleUseSkill() expected IsError=true but got false")
				}
			} else {
				if err != nil {
					t.Errorf("handleUseSkill() unexpected error: %v", err)
				}
				if result.IsError {
					t.Errorf("handleUseSkill() unexpected IsError=true")
				}
				if output.Name != tt.skillName {
					t.Errorf("handleUseSkill() output.Name = %v, want %v", output.Name, tt.skillName)
				}
				if output.Content == "" {
					t.Errorf("handleUseSkill() output.Content is empty")
				}
			}
		})
	}
}

// Made with Bob
