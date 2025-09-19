# Copilot Instructions

This file contains instructions for GitHub Copilot to help it generate better code suggestions for this repository.

baseline is a Go program that creates a baseline of Git repositories for easy searching.

It integrates with:

- GitHub
- Bitbucket

It checks out repositories into a specified directory, setting permissions to disallow write access, making them suitable for searching rather than active development.

The main features of baseline include:

- Cloning repositories from GitHub and Bitbucket
- Setting directory permissions to read-only
- Organizing repositories in a structured manner for easy access
- Providing a simple command-line interface for users to specify options

The commands should be:

- `init`: Initializes the baseline by creating the target directory if it doesn't exist.
- `discover`: lists repositories available in the specified source (GitHub or Bitbucket).
- `update`: updates repositories in the target directory from the specified source (GitHub or Bitbucket).
- `clone`: clones repositories from the specified source (GitHub or Bitbucket) into the target directory and updates existing ones.

The command line options are for all commands are:

- `-d <directory>`: Specifies the target directory for the baseline (default is `./baseline`).
- `-g <github_token>`: GitHub token for accessing private repositories.
- `-b <bitbucket_username>`: Bitbucket username for accessing private repositories.
- `-p <bitbucket_app_password>`: Bitbucket app password for accessing private repositories.
- `-o <organization>`: Specifies the organization to fetch repositories from (default is `jonasbn`).
- `-h`: Displays help information about the command line options.
- `-v`: Enables verbose output for debugging purposes.
- `-s <source>`: Specifies the source platform, either `github` or `bitbucket` (default is `github`).

This is specific to the `clone` and `update` commands:

- `-t <threads>`: Sets the number of concurrent threads for cloning repositories (default is `4`).

The cloning process involves:

1. Fetching the list of repositories from the specified source (GitHub or Bitbucket).
2. Cloning each repository into the target directory, skipping those that already exist
3. The cloning should be bare, meaning it does not include a working directory.
4. Setting the permissions of the cloned repositories to read-only recursively.

The update process involves:

1. Fetching the list of repositories from the specified source (GitHub or Bitbucket).
2. Setting the permissions of the updated repositories to write recursively.
3. Updating each repository in the target directory, pulling the latest changes from the source.
4. The update should be bare, meaning it does not include a working directory.
5. Setting the permissions of the updated repositories to read-only recursively.

About the implementation:

- The program is designed to be run from the command line, and it provides feedback on its progress, including any errors encountered during the cloning process.

- It should follow Go best practices for error handling, logging, and user feedback to ensure a smooth user experience.

- It should be written in Go and make use of relevant libraries for HTTP requests, file system operations, and command-line parsing.

- It should be modular, with functions dedicated to specific tasks such as fetching repository lists, cloning repositories, and setting permissions.

- It should include comments and documentation to explain the purpose of functions and important code sections.

- It should include unit tests for critical functions to ensure reliability and correctness.

- The command line options should be handled using the `cobra` package

- Create a `cmd` directory to store the command line commands

- Update the `README.md` file that provides an overview of the application, how to install it, and how to use it
- Update the `CHANGELOG.md` file to document changes and updates to the application over time