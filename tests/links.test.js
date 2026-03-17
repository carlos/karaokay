/**
 * Post-build test: ensures every href in the generated _site/ HTML
 * either points to an existing file/directory or is an external URL.
 *
 * Run: npm test  (after npm run build)
 */

import { readFileSync, existsSync, readdirSync, statSync } from "fs";
import { join, resolve } from "path";

const SITE_DIR = resolve("_site");
const PATH_PREFIX = process.env.ELEVENTY_PATH_PREFIX || "/karaokay/";

let passed = 0;
let failed = 0;
const errors = [];

function collectHtmlFiles(dir) {
  const files = [];
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    if (statSync(full).isDirectory()) {
      files.push(...collectHtmlFiles(full));
    } else if (entry.endsWith(".html")) {
      files.push(full);
    }
  }
  return files;
}

function extractHrefs(html) {
  const re = /href="([^"]+)"/g;
  const hrefs = [];
  let m;
  while ((m = re.exec(html)) !== null) {
    hrefs.push(m[1]);
  }
  return hrefs;
}

// --- Test 1: All internal links resolve to existing files ---
function testLinksResolve() {
  const htmlFiles = collectHtmlFiles(SITE_DIR);

  for (const file of htmlFiles) {
    const html = readFileSync(file, "utf-8");
    const hrefs = extractHrefs(html);

    for (const href of hrefs) {
      // Skip external links, anchors, mailto, tel, javascript
      if (/^(https?:|mailto:|tel:|javascript:|#)/.test(href)) continue;

      // Strip the path prefix to get the site-relative path
      let sitePath = href;
      if (sitePath.startsWith(PATH_PREFIX)) {
        sitePath = "/" + sitePath.slice(PATH_PREFIX.length);
      }

      // Resolve to filesystem path
      let fsPath = join(SITE_DIR, sitePath);

      // If it ends with /, look for index.html
      if (fsPath.endsWith("/")) {
        fsPath = join(fsPath, "index.html");
      }

      if (existsSync(fsPath)) {
        passed++;
      } else {
        failed++;
        const relFile = file.replace(SITE_DIR, "");
        errors.push(`  BROKEN: ${href} (in ${relFile})`);
      }
    }
  }
}

// --- Test 2: All internal links use the path prefix ---
function testLinksHavePrefix() {
  const htmlFiles = collectHtmlFiles(SITE_DIR);

  for (const file of htmlFiles) {
    const html = readFileSync(file, "utf-8");
    const hrefs = extractHrefs(html);

    for (const href of hrefs) {
      if (/^(https?:|mailto:|tel:|javascript:|#)/.test(href)) continue;

      if (href.startsWith(PATH_PREFIX)) {
        passed++;
      } else {
        failed++;
        const relFile = file.replace(SITE_DIR, "");
        errors.push(`  MISSING PREFIX: ${href} (in ${relFile}) — expected to start with ${PATH_PREFIX}`);
      }
    }
  }
}

// --- Test 3: Every song .md produced an HTML page ---
function testSongPagesExist() {
  const songsDir = resolve("src/songs");
  for (const entry of readdirSync(songsDir)) {
    if (!entry.endsWith(".md")) continue;

    // Read the markdown to get the title for the slug
    const md = readFileSync(join(songsDir, entry), "utf-8");
    const titleMatch = md.match(/^title:\s*["']?(.+?)["']?\s*$/m);
    if (!titleMatch) {
      failed++;
      errors.push(`  NO TITLE: ${entry} has no title frontmatter`);
      continue;
    }

    // Approximate Eleventy's slugify: normalize accents, lowercase, strip non-alphanum
    const title = titleMatch[1];
    const slug = title
      .normalize("NFD")
      .replace(/[\u0300-\u036f]/g, "")
      .toLowerCase()
      .replace(/['']/g, "")
      .replace(/[^a-z0-9]+/g, "-")
      .replace(/^-|-$/g, "");

    const pagePath = join(SITE_DIR, "songs", slug, "index.html");
    if (existsSync(pagePath)) {
      passed++;
    } else {
      failed++;
      errors.push(`  MISSING PAGE: ${entry} → expected _site/songs/${slug}/index.html`);
    }
  }
}

// Run all tests
testLinksResolve();
testLinksHavePrefix();
testSongPagesExist();

console.log(`\n  Links test: ${passed} passed, ${failed} failed\n`);
if (errors.length) {
  console.log(errors.join("\n"));
  console.log();
  process.exit(1);
}
