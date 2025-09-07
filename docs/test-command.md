# Test Command Documentation

## Overview

The `test` command verifies that your Jenkins connection is properly configured and working.

## Usage

```bash
jenkinsfile-validator test
```

## What it does

1. Loads the configuration from `~/.validator_config.json`
2. Verifies all required credentials are present (Jenkins URL, username, and API token)
3. Makes a GET request to the Jenkins API endpoint
4. Reports the connection status and basic Jenkins instance information

## Output

### Successful connection:
```
Testing connection to Jenkins...
Jenkins URL: https://jenkins.example.com
Username: myusername

Connection successful!
Connected to Jenkins node: Built-In Node
Jenkins mode: NORMAL
```

### Failed connection examples:

**Missing configuration:**
```
Error: Please configure Jenkins settings first using 'jenkinsfile-validator config'
```

**Authentication failure:**
```
Error: Authentication failed. Please check your username and API token
```

**Network/Connection issues:**
```
Error: Error connecting to Jenkins: <error details>
```

## Prerequisites

Before using the `test` command, ensure you have:

1. Configured your Jenkins settings using `jenkinsfile-validator config`
2. A valid Jenkins API token (not your password)
3. Network access to your Jenkins instance

## Troubleshooting

- **401 Unauthorized**: Check your username and API token
- **403 Forbidden**: Your user may not have sufficient permissions
- **Connection timeout**: Verify the Jenkins URL and network connectivity
- **SSL errors**: May need to configure SSL certificates or use HTTP for testing