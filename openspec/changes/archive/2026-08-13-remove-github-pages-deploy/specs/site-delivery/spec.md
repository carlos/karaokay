## Purpose

Defines how the Karaokay catalog is built, served, and reached. The catalog is a personal tool that
runs on the owner's machine: it is generated locally and previewed locally, and is not published to
any public address.

## ADDED Requirements

### Requirement: Local-only delivery

The catalog SHALL be built and served entirely on the owner's machine. The project SHALL NOT contain
any automated pipeline that publishes the generated site to a public host, and no build output SHALL
be reachable from a public URL.

#### Scenario: Building the catalog

- **WHEN** the owner runs the project's build command
- **THEN** the generated site is written to a local output directory
- **AND** no artifact is uploaded, pushed, or published to any external host

#### Scenario: Previewing the catalog

- **WHEN** the owner runs the project's development command
- **THEN** the catalog is served from a local address on the owner's machine
- **AND** it is reachable only from that machine

#### Scenario: No publishing pipeline is present

- **WHEN** the repository is inspected for continuous-integration or deployment configuration
- **THEN** no workflow, action, or scheduled job exists that builds or publishes the site

#### Scenario: The former public address is retired

- **WHEN** a request is made to the address the catalog was previously published at
- **THEN** the request does not return the catalog

### Requirement: Root-relative site addressing

Because the catalog is not hosted under a project subpath, generated pages SHALL be addressed from
the site root. Internal links SHALL resolve correctly when the site is served from the root of a
local address.

#### Scenario: Page addresses carry no project prefix

- **WHEN** a song, artist, or album page is generated
- **THEN** its address begins at the site root and contains no project-name path segment

#### Scenario: Internal links resolve when served locally

- **WHEN** the built site is served from the root of a local address
- **AND** any internal link on any page is followed
- **THEN** it resolves to an existing page

#### Scenario: Malformed internal links are rejected

- **WHEN** the catalog's validation is run against a build containing an internal link that is not
  root-relative
- **THEN** validation fails and identifies the offending page and link

### Requirement: Contribution requires no review gate

Changes to the catalog SHALL be made directly by the owner without a mandatory review or approval
step. The project SHALL NOT require a pull request, status check, or branch protection rule in order
for a change to take effect.

#### Scenario: Committing a change directly

- **WHEN** the owner commits a change to the default branch
- **THEN** the change takes effect with no pull request, review, or automated approval required

### Requirement: Validation remains available on demand

Removing automated publishing SHALL NOT remove the ability to validate a build. The project SHALL
retain a command that checks the generated site, runnable by the owner at any time.

#### Scenario: Running validation manually

- **WHEN** the owner runs the project's test command after a build
- **THEN** the generated site is checked for broken and malformed internal links and for missing song
  pages
- **AND** the result is reported as a pass or a failure naming what broke
