# Interactive Configuration Documentation

## Overview

The jenkinsfile-validator now supports interactive configuration with two commands:
- `config`: Set up or update all configuration fields
- `refresh`: Update credentials with individual field selection

## Config Command

The `config` command provides a streamlined way to configure all Jenkins settings.

### Usage

```bash
jenkinsfile-validator config
```

### Behavior

1. **First-time setup**: Prompts for all fields sequentially
2. **Existing configuration**: Shows current values and asks if you want to update

### Features
- Displays current configuration values
- Validates Jenkins URL format (must start with http:// or https://)
- Prevents empty usernames
- Masks API tokens for security

### Example

```
Configure Jenkins Validator Settings
=====================================

Existing configuration found:
Jenkins URL: https://jenkins.example.com
Username: myuser
Token: ************************

Do you want to update the configuration? (y/N) [n]: y

Jenkins URL [https://jenkins.example.com]: https://newjenkins.example.com
Username [myuser]: newuser
API Token: ****

Configuration saved successfully!
```

## Refresh Command

The `refresh` command provides fine-grained control over credential updates.

### Usage

```bash
jenkinsfile-validator refresh
```

### Features

1. **Interactive menu**: Choose to update all fields or individual fields
2. **Field selection**: Update only the fields you need to change
3. **Current value display**: Shows existing values before updates
4. **Validation**: Ensures valid URLs and non-empty usernames

### Update Options

When running `refresh`, you'll see:

```
Refresh Jenkins Credentials
===========================

Current configuration:
Jenkins URL: https://jenkins.example.com
Username: myuser
Token: ************************

What would you like to update?
1. All fields
2. Individual fields
3. Cancel
Choice (1-3) [2]: 
```

### Individual Field Updates

If you select option 2 (Individual fields), you get a menu:

```
Select field to update:
1. Jenkins URL
2. Username
3. API Token
4. Done
Choice (1-4): 
```

Each field update:
- Shows the current value
- Prompts for new value
- Validates input
- Confirms the update

### Example Individual Update

```
Select field to update:
1. Jenkins URL
2. Username
3. API Token
4. Done
Choice (1-4): 1

Current Jenkins URL: https://jenkins.example.com
New Jenkins URL [https://jenkins.example.com]: https://newjenkins.example.com
Jenkins URL updated.

Select field to update:
1. Jenkins URL
2. Username
3. API Token
4. Done
Choice (1-4): 4

Configuration saved successfully!
```

## Validation Rules

### Jenkins URL
- Cannot be empty
- Must be a valid URL format
- Must start with `http://` or `https://`
- Must contain a valid host

### Username
- Cannot be empty
- No special character restrictions

### API Token
- Can be left empty to keep existing token
- Not validated for format (Jenkins tokens vary)

## Security Notes

- API tokens are never displayed in plain text
- Tokens are masked with asterisks based on length
- Pressing Enter without input keeps the existing token
- Configuration is stored in `~/.validator_config.json` with file permissions 0644

## Error Handling

The interactive commands handle errors gracefully:

- **Invalid URL**: Prompts again with error message
- **Empty username**: Prompts again with error message
- **File permission issues**: Displays error and exits

## Tips

1. Use `refresh` when you only need to update one field (e.g., rotating API token)
2. Use `config` for initial setup or complete reconfiguration
3. Press Enter to keep existing values when prompted
4. The configuration file location is `~/.validator_config.json`