# AWS SNS Cross-Account Subscription Tool

A command-line tool to analyze AWS SNS topics and identify cross-account Lambda subscriptions across different AWS accounts.

## Features

- List cross-account Lambda subscriptions for SNS topics
- Support for multiple AWS profiles
- Region selection support
- CSV export functionality
- Easy-to-use command-line interface with short and long-form flags

## Installation

### Using Go Install

If you have Go installed (version 1.21 or later), you can install directly using:

```bash
go install github.com/dhairya13703/sns-tool@latest
```

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/dhairya13703/sns-tool.git
cd sns-tool
```

2. Install dependencies and build:
```bash
go mod tidy
go build -o sns-tool
```

3. (Optional) Move the binary to your PATH:
```bash
sudo mv sns-tool /usr/local/bin/
```

## Usage

### Basic Usage

```bash
# Using default profile
sns-tool list -a YOUR_ACCOUNT_NUMBER

# Using a specific profile
sns-tool list -a YOUR_ACCOUNT_NUMBER -p YOUR_PROFILE

# Using a different region
sns-tool list -a YOUR_ACCOUNT_NUMBER -r us-west-2

# Export to CSV
sns-tool list -a YOUR_ACCOUNT_NUMBER -e csv -o subscriptions.csv
```

### Available Commands and Flags

- `list`: List cross-account SNS subscriptions

Flags:
| Flag      | Short | Description                    | Required | Default    |
|-----------|-------|--------------------------------|----------|------------|
| --account | -a    | AWS account number            | Yes      | -          |
| --profile | -p    | AWS profile to use            | No       | default    |
| --region  | -r    | AWS region                    | No       | us-east-1  |
| --export  | -e    | Export format (csv)           | No       | -          |
| --output  | -o    | Output file name              | No       | -          |

### AWS Credentials

The tool uses AWS credentials in the following order:
1. Specified profile using `-p` flag
2. Default AWS credentials
3. Environment variables
4. EC2 instance role

Make sure you have your AWS credentials configured in `~/.aws/credentials` or `~/.aws/config` for AWS SSO.

Example credentials file:
```ini
[default]
aws_access_key_id = YOUR_ACCESS_KEY
aws_secret_access_key = YOUR_SECRET_KEY

[dev]
aws_access_key_id = DEV_ACCESS_KEY
aws_secret_access_key = DEV_SECRET_KEY
```

Example AWS SSO config:
```ini
[profile dev]
sso_start_url = https://your-sso-portal.awsapps.com/start
sso_region = us-east-1
sso_account_id = 123456789012
sso_role_name = YourRole
```

## Requirements

- Go 1.21 or later
- AWS credentials configured
- Appropriate AWS permissions to list SNS topics and subscriptions

## AWS Permissions Required

The tool requires the following AWS permissions:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "sns:ListTopics",
                "sns:ListSubscriptionsByTopic",
                "sts:GetCallerIdentity"
            ],
            "Resource": "*"
        }
    ]
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.