## Why

Karaokay is a personal catalog that runs on one machine. Publishing it to GitHub Pages adds a
deployment pipeline, a pull-request gate, and a public URL that nothing needs — and it puts a corpus
of copyrighted song lyrics on a public website. Removing deployment removes all three at once, and
lets the site's URLs shed the `/karaokay/` project-path prefix that only ever existed to satisfy
GitHub Pages.

## What Changes

- Delete `.github/workflows/deploy.yml`, and the `.github/` directory with it — it holds nothing else.
- Disable the GitHub Pages site via the API. **BREAKING**: `https://carlos.github.io/karaokay/` stops
  resolving. Deleting the workflow alone would leave the last build published indefinitely, so this is
  a required step, not a side effect.
- Change `baseURL` in `hugo.toml` from `https://carlos.github.io/karaokay/` to `/`. **BREAKING**:
  every page URL loses the `/karaokay/` prefix (`/karaokay/songs/roar/` → `/songs/roar/`). Templates
  already resolve links through `.RelPermalink`, so they adapt with no edits.
- Update the expected prefix in `tests/links_test.go` from `/karaokay/` to `/`. The check is kept, not
  dropped — with `/` it still catches relative or malformed hrefs, which is the failure it was
  originally added for.
- Rewrite the deployment and path-prefix sections of `CLAUDE.md`.
- The pull-request requirement disappears with the workflow. `master` has no branch protection, so
  nothing else enforces it and nothing else needs removing.

Not changing: the build itself, the content, the templates, the Go tests, or the mise toolchain
pinning. `mise run dev` / `build` / `test` keep working exactly as they do now.

## Capabilities

### New Capabilities

- `site-delivery`: how the catalog is built, served, and reached — local build and preview only, with
  no publishing pipeline and no public URL.

### Modified Capabilities

<!-- None. No specs exist yet; this change introduces the first one. -->

## Impact

**Code and config**

- `.github/workflows/deploy.yml` — deleted
- `hugo.toml` — `baseURL`
- `tests/links_test.go` — `defaultPrefix` constant
- `CLAUDE.md` — Deploy section removed, path-prefix convention rewritten

**Outside the repo**

- GitHub Pages disabled for `carlos/karaokay`; the live site stops resolving.
- No further CI runs. `mise run test` still exists but becomes manual — nothing will run it
  automatically on commit or push. This is the real cost of the change and is accepted deliberately.

**Not affected**

- Content (`content/songs/*.md`), layouts, styles, the Go test suite, `mise.toml`.
- The git remote stays; the repo continues to work as off-machine backup.
