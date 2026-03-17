# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Karaokay is a static karaoke song lyrics catalog built with Eleventy (11ty) 3.x. Songs are Markdown files with YAML frontmatter (`title`, `artist`, `album`), and the site generates song pages, an all-songs index, artist listing, per-artist pages, album listing, and per-album pages.

## Commands

- `npm run dev` — dev server with hot reload (port 8081 via `.claude/launch.json`)
- `npm run build` — production build to `_site/`
- `npm test` — post-build link validation (run after `npm run build`)

## Architecture

- **Source:** `src/` → **Output:** `_site/`
- **Config:** `eleventy.config.js` defines three collections (`songs` sorted by title, `artists` grouped by artist, `albums` grouped by album) and an `inlineCSS` shortcode that embeds CSS directly into HTML
- **Templates:** Nunjucks (`.njk`). All pages extend `src/_includes/layouts/base.njk` (sidebar + nav + inline search script). Song pages use `src/_includes/layouts/song.njk`
- **Songs:** Add a `.md` file to `src/songs/` with `title`, `artist` (string or array), and `album` frontmatter. Layout and permalink are set by `src/songs/songs.json`. Multiple artists are displayed joined with commas and `&`
- **Artist pages:** `src/artist-pages.njk` uses Eleventy pagination over `collections.artists` to generate `/artists/{slug}/`
- **Album pages:** `src/album-pages.njk` uses Eleventy pagination over `collections.albums` to generate `/albums/{slug}/`
- **Path prefix:** Production uses `/karaokay/` prefix (for GitHub Pages). All template hrefs must use the `| url` filter. Dev overrides to `/` via env var. CI runs `npm test` to catch broken/unprefixed links
- **Styling:** Single file `src/css/styles.css`, inlined at build time. Purple accent (#7c3aed), responsive at 768px, print-optimized for A4
- **Interactivity:** Vanilla JS client-side search filtering only, no frameworks or bundlers
- **Dependencies:** Only `@11ty/eleventy` (dev dependency). ES modules (`"type": "module"`)
