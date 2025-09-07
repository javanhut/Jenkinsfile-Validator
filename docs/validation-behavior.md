# Jenkinsfile Validation Behavior

## Overview

The Jenkinsfile validator uses Jenkins' built-in pipeline-model-converter API to validate Jenkinsfile syntax and structure.

## How It Works

1. The validator reads your Jenkinsfile and sends it to your Jenkins instance
2. Jenkins parses the file and checks for:
   - Syntax errors
   - Invalid pipeline structure
   - Missing required sections
   - Invalid directives or steps
   - Undefined variables or methods

## Response Handling

The validator now properly parses the Jenkins validation response:

### Successful Validation
- Status: "ok"
- Output: "✓ Jenkinsfile is valid"

### Failed Validation
- Status: "fail"
- Output: "✗ Jenkinsfile validation failed"
- Shows detailed error messages including:
  - Line numbers where errors occur
  - Specific syntax or structural issues
  - Missing or invalid pipeline components

## Example Outputs

### Valid Jenkinsfile
```
$ jenkinsfile-validator validate Jenkinsfile
✓ Jenkinsfile is valid
```

### Invalid Jenkinsfile
```
$ jenkinsfile-validator validate Jenkinsfile
✗ Jenkinsfile validation failed

Errors:
- WorkflowScript: 5: Missing required section "agent" @ line 5, column 1.
- WorkflowScript: 10: Unknown stage section "invalidSection". Valid sections: steps, post, agent, tools, environment, when, input @ line 10, column 9.
```

## Important Notes

- A 200 HTTP status code does NOT mean the Jenkinsfile is valid
- The actual validation result is contained in the JSON response body
- Always check the "status" field in the response for the true validation result