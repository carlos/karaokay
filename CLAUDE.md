# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Karaokay is a static karaoke song lyrics catalog built with Eleventy (11ty) 3.x. Songs are Markdown files with YAML frontmatter (`title`, `artist`, `album`, `track`), and the site generates song pages, an all-songs index, artist listing, per-artist pages, album listing, and per-album pages.

## Commands

- `npm run dev` — `eleventy --serve` with hot reload (default port 8080; `.claude/launch.json` runs it on 8081)
- `npm run build` — production build to `_site/`
- `npm test` — post-build validation (`tests/links.test.js`); requires a fresh `npm run build` first, since it reads `_site/`

There is no linter and no test framework — `tests/links.test.js` is a plain Node script with three checks: every internal href resolves to a real file, every internal href carries the path prefix, and every `src/songs/*.md` produced an HTML page.

## Architecture

- **Source:** `src/` → **Output:** `_site/`. Only dependency is `@11ty/eleventy` (dev). ES modules (`"type": "module"`), no bundler or framework.
- **Config:** `eleventy.config.js` defines three collections — `songs` (sorted by title), `artists` (grouped by artist, songs sorted by title), `albums` (grouped by album, songs sorted by `track` with untracked songs last) — plus `formatArtists`/`toArtistList` filters and an `inlineCSS` shortcode that embeds `src/css/styles.css` into every page.
- **Templates:** Nunjucks (`.njk`), Markdown rendered through Nunjucks. All pages extend `src/_includes/layouts/base.njk`; song pages use `src/_includes/layouts/song.njk`, which itself extends base.
- **Artist/album pages:** `src/artist-pages.njk` and `src/album-pages.njk` use Eleventy pagination (`size: 1`) over the grouped collections to emit `/artists/{slug}/` and `/albums/{slug}/`. Slugs come from the `slugify` filter applied to the artist/album name.

## Conventions

- **Adding a song:** drop a `.md` file in `src/songs/`. Layout, permalink, and `pageTitle` come from `src/songs/songs.json` — do not set them per-file. The URL slug is derived from `title`, not the filename, so the filename is cosmetic (keep it matching the slug anyway).
- **`artist`** may be a string or a YAML list. Multi-artist values render joined with commas and a final `&` (`formatArtists`), and each artist links to its own page. Every distinct string becomes its own artist page, so spelling must be consistent across files.
- **`track`** is optional and only affects ordering on album pages.
- **Lyrics** are the Markdown body. `.lyrics` uses `white-space: pre-wrap` and `text-transform: uppercase`, so write normal casing in source and let CSS uppercase it; blank lines become paragraphs (stanza breaks), single newlines are preserved as-is.
- **Path prefix:** builds use `/karaokay/` (GitHub Pages), overridable via `ELEVENTY_PATH_PREFIX`. Every internal href in a template **must** go through the `| url` filter — including hand-built paths, e.g. `{{ ('/artists/' + artist.slug + '/') | url }}`. `npm test` fails the CI build on any unprefixed or broken link.
- **Search** is a single inline script in `base.njk`: it binds to `.search-input` and shows/hides elements marked `[data-searchable]` by text match. A new listing page gets filtering for free by using those two hooks — no per-page JS.
- **Page titles** come from a `pageTitle` variable consumed by `base.njk`.
- **Styling:** one file, `src/css/styles.css`. Purple accent (`#7c3aed`), responsive breakpoint at 768px, print styles tuned for A4 song sheets.

## Deploy

`.github/workflows/deploy.yml` builds on push to `master`, runs `npm test`, and publishes `_site/` to GitHub Pages.
