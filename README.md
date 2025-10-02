# baseline

A Go program that creates a baseline of Git repositories for easy searching.

I did [a little write up](https://dev.to/jonasbn/baseline-a-parallel-git-universeconcept-for-easy-searching-k1m)
for this repository on dev.to

## Overview

baseline integrates with GitHub and Bitbucket to clone repositories into a specified
directory, setting permissions to disallow write access, making them suitable for
searching rather than active development.

The tool uses `git clone` which creates repositories without working
directories - perfect for searching through code without the overhead of
checked-out files. The repositories are also set to read-only permissions
to prevent accidental modifications.

## Features

- Clone repositories from GitHub and Bitbucket organizations
- Concurrent cloning with configurable thread count
- Repository cloning for optimal searching
- Read-only permissions to prevent accidental modifications
- Update existing repositories with latest changes
- Discover available repositories before cloning
- Support for both public and private repositories (with authentication)

## Installation

### From Source

```bash
go install github.com/jonasbn/baseline@latest
```

### Building Locally

```bash
git clone https://github.com/jonasbn/baseline.git
cd baseline
go build
```

## Usage

### Commands

- `init`: Initialize the baseline by creating the target directory
- `discover`: List repositories available in the specified source
- `clone`: Clone repositories from the specified source into the target directory
- `update`: Update repositories in the target directory from the specified source

### Global Options

- `-d, --directory`: Target directory for the baseline (default: `./baseline`)
- `-g, --github-token`: GitHub token for accessing private repositories
- `-b, --bitbucket-token`: Bitbucket API token for accessing private repositories
- `-o, --organization`: Organization to fetch repositories from (default: `jonasbn`)
- `-s, --source`: Source platform, either `github` or `bitbucket` (default: `github`)
- `-v, --verbose`: Enable verbose output for debugging
- `-t, --threads`: Number of concurrent threads for cloning/updating (default: `4`)

### Examples

#### Initialize baseline directory

```bash
baseline init -d ./my-baseline
```

#### Discover repositories

```bash
# List public repositories from GitHub
baseline discover -o myorg

# List repositories from Bitbucket with authentication
baseline discover -s bitbucket -b your_api_token -o myorg
```

#### Clone repositories

```bash
# Clone all public repositories from GitHub organization
baseline clone -o myorg -d ./baseline

# Clone with authentication and custom thread count
baseline clone -g mytoken -o myorg -d ./baseline -t 8 -v

# Clone from Bitbucket
baseline clone -s bitbucket -b your_api_token -o myorg
```

#### Update repositories

```bash
# Update all existing repositories
baseline update -o myorg -d ./baseline

# Update with verbose output
baseline update -o myorg -d ./baseline -v
```

## Authentication

### GitHub

Create a personal access token at [GitHub](https://github.com/settings/tokens) with
appropriate repository access permissions.

```bash
baseline clone -g your_github_token -o organization_name
```

### Bitbucket

Create an API token at [Bitbucket](https://bitbucket.org/account/settings/access-management/api-tokens) with
repository read permissions. You can create either:

- **Repository Access Token**: For specific repositories
- **Project Access Token**: For all repositories in a project  
- **Workspace Access Token**: For all repositories in a workspace

```bash
baseline clone -s bitbucket -b your_api_token -o organization_name
```

**Note:** App passwords are deprecated by Bitbucket in favor of API tokens for better security and granular permissions.

## Directory Structure

Repositories are organized in the following structure:

```text
baseline/
├── owner1/
│   ├── repo1.git/
│   ├── repo2.git/
│   └── repo3.git/
└── owner2/
    ├── repo4.git/
    └── repo5.git/
```

Each repository is cloned as a bare repository (`.git` extension) with read-only permissions.

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file
for details.
