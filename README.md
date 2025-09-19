# baseline

A Go program that creates a baseline of Git repositories for easy searching.

I did [a little write up](https://dev.to/jonasbn/baseline-a-parallel-git-universeconcept-for-easy-searching-k1m) for this repository on dev.to

## Overview

baseline integrates with GitHub and Bitbucket to clone repositories into a specified directory, setting permissions to disallow write access, making them suitable for searching rather than active development.

The tool uses `git clone --bare` which creates repositories without working directories - perfect for searching through code without the overhead of checked-out files. The bare repositories are also set to read-only permissions to prevent accidental modifications.

## Features

- Clone repositories from GitHub and Bitbucket organizations
- Concurrent cloning with configurable thread count
- Bare repository cloning (no working directory) for optimal searching
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
- `-b, --bitbucket-username`: Bitbucket username for accessing private repositories
- `-p, --bitbucket-password`: Bitbucket app password for accessing private repositories
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
baseline discover -s bitbucket -b myusername -p myapppassword -o myorg
```

#### Clone repositories
```bash
# Clone all public repositories from GitHub organization
baseline clone -o myorg -d ./baseline

# Clone with authentication and custom thread count
baseline clone -g mytoken -o myorg -d ./baseline -t 8 -v

# Clone from Bitbucket
baseline clone -s bitbucket -b myusername -p myapppassword -o myorg
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

Create a personal access token at https://github.com/settings/tokens with appropriate repository access permissions.

```bash
baseline clone -g your_github_token -o organization_name
```

### Bitbucket

Create an app password at https://bitbucket.org/account/settings/app-passwords/ with repository read permissions.

```bash
baseline clone -s bitbucket -b your_username -p your_app_password -o organization_name
```

## Directory Structure

Repositories are organized in the following structure:

```
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

## Pros and Cons of using `git clone --bare`

### Pros
- **No working directory**: Smaller disk footprint and no risk of accidental file edits
- **Optimal for searching**: Perfect for code search tools as there are no checked-out files to confuse search results
- **Faster operations**: No working tree means faster clones and updates
- **Safer permissions**: Read-only permissions prevent accidental modifications
- **Multiple worktrees**: Can create temporary worktrees elsewhere if needed for inspection

### Cons
- **No working directory**: Cannot directly edit files or run builds in the repository
- **Tool compatibility**: Some IDEs and Git tools expect a working tree
- **Limited inspection**: Need to use git commands like `git show` or create temporary worktrees to view files

For baseline's use case of creating searchable code repositories, the pros far outweigh the cons.

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

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
