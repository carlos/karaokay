# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Karaokay is a static karaoke song lyrics catalog built with Eleventy (11ty) 3.x. Songs are Markdown files with YAML frontmatter (`title`, `artist`), and the site generates song pages, an all-songs index, artist listing, and per-artist pages.

## Commands

- `npm run dev` — dev server with hot reload (port 8081 via `.claude/launch.json`)
- `npm run build` — production build to `_site/`

There are no tests or linting configured.

## Architecture

- **Source:** `src/` → **Output:** `_site/`
- **Config:** `eleventy.config.js` defines two collections (`songs` sorted by title, `artists` grouped by artist) and an `inlineCSS` shortcode that embeds CSS directly into HTML
- **Templates:** Nunjucks (`.njk`). All pages extend `src/_includes/layouts/base.njk` (sidebar + nav + inline search script). Song pages use `src/_includes/layouts/song.njk`
- **Songs:** Add a `.md` file to `src/songs/` with `title` and `artist` frontmatter. Layout and permalink are set by `src/songs/songs.json`
- **Artist pages:** `src/artist-pages.njk` uses Eleventy pagination over `collections.artists` to generate `/artists/{slug}/`
- **Styling:** Single file `src/css/styles.css`, inlined at build time. Purple accent (#7c3aed), responsive at 768px, print-optimized for A4
- **Interactivity:** Vanilla JS client-side search filtering only, no frameworks or bundlers
- **Dependencies:** Only `@11ty/eleventy` (dev dependency). ES modules (`"type": "module"`)
