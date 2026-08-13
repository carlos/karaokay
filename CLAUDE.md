# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Karaokay is a static karaoke song lyrics catalog built with Hugo. Songs are Markdown files with YAML
frontmatter (`title`, `artists`, `album`, `track`), and the site generates song pages, an all-songs
index, artist listing, per-artist pages, album listing, and per-album pages.

## Commands

Tool versions are pinned in `mise.toml`; CI installs the same versions via `jdx/mise-action`.

- `mise run dev` — `hugo server` with hot reload on port 8081
- `mise run build` — production build to `public/`
- `mise run test` — post-build validation (`go test ./...`); requires a fresh build first, since it
  reads `public/`. Nothing runs this automatically — see "Local only" below.

There is no linter and no test framework beyond Go's own — `tests/links_test.go` is a stdlib-only Go
test with three checks: every internal href resolves to a real file, every internal href carries the
path prefix, and every `content/songs/*.md` produced exactly one page.

Hugo does **not** clean `public/` between builds. Stale output from a removed or renamed song will
linger and can mask a real failure — `rm -rf public` when a build's output looks wrong.

## Architecture

- **Source:** `content/` + `layouts/` + `assets/` → **Output:** `public/`. No dependencies, no
  bundler, no JavaScript build.
- **Config:** `hugo.toml`. `artist` is the only taxonomy; default `tags`/`categories` are off. RSS and
  sitemap are disabled to match the site's current surface.
- **Templates:** Go templates. All pages extend `layouts/baseof.html`; song pages use
  `layouts/songs/single.html`.
- **Artist pages:** Hugo taxonomy routing generates `/artists/` (`layouts/artists/terms.html`) and
  `/artists/<slug>/` (`layouts/artists/term.html`) — note the term-page template is `term.html`, not
  `taxonomy.html`.
- **Album pages:** `album` is a plain param, not a taxonomy, because a song has exactly one album and
  album pages are ordered by `track`. `content/albums/_content.gotmpl` is a content adapter that emits
  one page per distinct album, so importing a song with a new album needs no extra file. Content
  adapters run *before* the site is initialized, so it cannot use `site.RegularPages` — it reads the
  song front matter off disk instead. The adapter is silently ignored unless `content/albums/_index.md`
  exists.

## Conventions

- **Adding a song:** drop a `.md` file in `content/songs/`. **The URL slug comes from the filename**,
  not the title, so the filename is significant — name it after the song.
- **`artists`** is plural because Hugo taxonomies key off the plural form; a singular `artist:` key is
  silently ignored and the song will not appear under any artist. The value may be a string or a YAML
  list — Hugo normalizes both. Multi-artist values render joined with commas and a final `&`
  (`layouts/partials/artist-names.html`), in frontmatter order, and each artist links to its own page.
  Every distinct string becomes its own artist page, so spelling must be consistent across files.
- **`track`** is optional and only affects ordering on album pages; untracked songs sort last.
- **Lyrics** are the Markdown body. `.lyrics` uses `white-space: pre-wrap` and `text-transform:
  uppercase`, so write normal casing in source and let CSS uppercase it; blank lines become paragraphs
  (stanza breaks), single newlines are preserved as-is. Goldmark's typographer is **disabled** so
  apostrophes stay as written (`you're`, not `you’re`) — anything searching the lyrics has to match
  what a person actually types. Avoid lyric lines that start with `-`, `*`, `#`, `>`, or `N.`, and
  avoid `*`/`_` in lyrics: Markdown will render them as lists or emphasis.
- **Internal links:** the site is served from the root (`baseURL = "/"`), so pages are addressed as
  `/songs/roar/`. **Never hand-build an internal href.** `relURL` and `absURL` do not behave the way
  their names suggest — they silently drop any baseURL subpath, which once shipped a site whose every
  link was broken. Use `.RelPermalink`, `site.Home.RelPermalink`, or resolve the page first with
  `site.GetPage`. `mise run test` fails on any href that is relative, malformed, or points nowhere.
  Root-relative addressing means the built `public/` needs to be *served* — opening `index.html` from
  the filesystem will not work. Set `relativeURLs = true` if you ever want portable output.
- **Search** is a single inline script in `baseof.html`: it binds to `.search-input` and shows/hides
  elements marked `[data-searchable]` by text match. A new listing page gets filtering for free by
  using those two hooks — no per-page JS.
- **Page titles** come from `.Title`, rendered by `baseof.html` as `Karaokay — <title>`.
- **Styling:** one file, `assets/css/styles.css`, inlined into every page via `resources.Get`. Purple
  accent (`#7c3aed`), responsive breakpoint at 768px, print styles tuned for A4 song sheets.

## Local only — do not add a deployment

This site is **deliberately unpublished**. It is built and previewed on one machine and is not hosted
anywhere. GitHub Pages was disabled and the deployment workflow deleted on purpose — their absence is
a decision, not an omission, so do not restore a workflow or re-enable Pages assuming something went
missing.

The git remote is kept for off-machine backup only; pushing does not publish anything.

Because there is no CI, **nothing runs the tests automatically**. Run `mise run build && mise run
test` yourself after editing content or templates — that is now the only thing standing between a
broken link and not noticing.
