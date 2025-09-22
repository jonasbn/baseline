# Change log for baseline

## v0.3.1 2025-09-22 - minor bug fix release, update recommended

- Removed `--bare` from git clone commands, this was not what I wanted

## v0.3.0 2025-09-22 - Go Implementation, update not required

- Update not required, but recommended since the Go client will be the future

- Complete rewrite in Go from original shell script
- Added support for concurrent cloning with configurable thread count
- Implemented Cobra CLI framework with proper commands and flags
- Added GitHub API integration for fetching repository lists
- Added Bitbucket API integration for fetching repository lists
- Implemented bare repository cloning for optimal searching
- Added read-only permission management for cloned repositories
- Added update functionality to fetch latest changes for existing repositories
- Added discover functionality to list available repositories before cloning
- Added comprehensive error handling and verbose logging
- Added unit tests for core functionality
- Added proper documentation and usage examples

### Commands Added

- `init`: Initialize baseline directory
- `discover`: List available repositories
- `clone`: Clone repositories with concurrent processing
- `update`: Update existing repositories

### Features Added

- Support for both GitHub and Bitbucket
- Authentication support for private repositories
- Configurable concurrency (default 4 threads)
- Verbose logging for debugging
- Proper error handling and user feedback
- Read-only permissions for cloned repositories
- Bare repository cloning (no working directory)

## v0.2.0 2023-09-26 - Initial release, update not required

- Update note required, but recommended, the permission handling is useful

- Added recursive setting of permissions for disallowing write access to the
  repositories, so they can be used for searching and not for actual work

## v0.1.0, 2023-09-26 - Initial release

- Shell script for creating a baseline of Git repositories for easy searching
