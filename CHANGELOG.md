# Change log for baseline

## v0.5.0 2025-10-02 - Directory Structure Cleanup

### Changed

- **Directory Structure**: Removed `.git` postfix from repository directory names for cleaner organization
  - Repositories are now stored as `baseline/owner/repo-name/` instead of `baseline/owner/repo-name.git/`
  - Updated all cloning, updating, and existence checking functions to use the new naming convention
  - Provides a cleaner and more intuitive directory structure for users

### Fixed

- **Test Updates**: Updated unit tests to reflect the new directory naming convention

## v0.4.0 2025-10-02 - SSH Support and Bitbucket Authentication Updates

### Added

- **SSH Support**: Added `--ssh` flag to `clone` and `update` commands for using SSH URLs instead of HTTPS
- **Enhanced Bitbucket Authentication**: Updated to use username + API token combination (required by Bitbucket API)
- **Debug Logging**: Added comprehensive HTTP request/response debugging for Bitbucket API calls with verbose flag
- **Auto-detection**: GitHub client now automatically detects and handles both organization and personal account repositories

### Changed

- **Bitbucket Authentication**: Now requires both username (`-u`) and API token (`-b`) instead of just token
- **Command Flags**: Added `--bitbucket-username` flag alongside existing `--bitbucket-token` flag
- **Documentation**: Updated README with SSH usage examples and Bitbucket API token guidance

### Fixed

- **Bitbucket API Integration**: Corrected authentication method to match Bitbucket's actual API requirements
- **GitHub Repository Detection**: Fixed handling of personal accounts vs organizations

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
