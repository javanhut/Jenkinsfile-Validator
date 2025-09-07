# Jenkinsfile-Validator

A command-line tool to validate Jenkinsfiles by connecting to your Jenkins instance.

## Installation

```bash
go install github.com/javanhut/Jenkinsfile-Validator
```

## Usage

### Configure Jenkins Connection

First, configure your Jenkins connection details:

```bash
jenkinsfile-validator config
```

This will prompt you for:
- Jenkins URL (e.g., https://jenkins.example.com)
- Username
- API Token

### Update Credentials

To update your credentials interactively:

```bash
jenkinsfile-validator refresh
```

This command allows you to:
- Update all fields at once
- Update individual fields selectively
- View current configuration values

### Test Connection

Verify that your credentials and Jenkins URL are correct:

```bash
jenkinsfile-validator test
```

This command will attempt to connect to your Jenkins instance and display connection status.

### Validate a Jenkinsfile

To validate a Jenkinsfile:

```bash
jenkinsfile-validator validate path/to/Jenkinsfile
```

## Commands

- `config` - Set up or update Jenkins connection details
- `refresh` - Update credentials with interactive field selection
- `test` - Test connection to Jenkins using configured credentials
- `validate` - Validate a Jenkinsfile against your Jenkins instance

## Configuration

Configuration is stored in `~/.validator_config.json` with the following structure:

```json
{
  "jenkins_url": "https://jenkins.example.com",
  "username": "your-username",
  "token": "your-api-token"
}
```
