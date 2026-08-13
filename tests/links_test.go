// Post-build validation: every href in the generated site resolves to a real
// file, carries the path prefix, and every song source produced a page.
//
// Run: mise run build && mise run test
package tests

import (
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Build output and content locations, relative to the repo root.
const (
	defaultSiteDir   = "public"
	defaultSongsDir  = "content/songs"
	defaultPrefix    = "/karaokay/"
	prefixEnvVar     = "KARAOKAY_PATH_PREFIX"
	siteDirEnvVar    = "KARAOKAY_SITE_DIR"
	rootMarkerFile   = "mise.toml"
	maxRootAscension = 5
)

var (
	hrefPattern      = regexp.MustCompile(`href="([^"]+)"`)
	externalPattern  = regexp.MustCompile(`^(https?:|mailto:|tel:|javascript:|#)`)
	titlePattern     = regexp.MustCompile(`(?m)^title:\s*["']?(.+?)["']?\s*$`)
	songTitlePattern = regexp.MustCompile(`(?s)<h1 class="song-title">(.*?)</h1>`)
	tagPattern       = regexp.MustCompile(`<[^>]*>`)
)

// repoRoot walks up from the test's working directory (which go test sets to
// the package directory, not the repo root) until it finds the root marker.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for i := 0; i < maxRootAscension; i++ {
		if _, err := os.Stat(filepath.Join(dir, rootMarkerFile)); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("could not locate repo root (no %s found above %s)", rootMarkerFile, dir)
	return ""
}

func siteDir(t *testing.T) string {
	t.Helper()
	dir := defaultSiteDir
	if override := os.Getenv(siteDirEnvVar); override != "" {
		dir = override
	}
	full := dir
	if !filepath.IsAbs(full) {
		full = filepath.Join(repoRoot(t), dir)
	}
	if _, err := os.Stat(full); err != nil {
		t.Fatalf("site directory %s not found — build first: %v", full, err)
	}
	return full
}

func pathPrefix() string {
	if override := os.Getenv(prefixEnvVar); override != "" {
		return override
	}
	return defaultPrefix
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(data)
}

// textOf strips tags and unescapes entities, treating </p> as a stanza break.
func textOf(fragment string) string {
	fragment = strings.ReplaceAll(fragment, "</p>", "\n\n")
	return strings.TrimSpace(html.UnescapeString(tagPattern.ReplaceAllString(fragment, "")))
}

// pagesUnder returns the index.html of each direct child directory of
// <root>/<section>, e.g. every /songs/<slug>/index.html.
func pagesUnder(t *testing.T, root, section string) []string {
	t.Helper()
	matches, err := filepath.Glob(filepath.Join(root, section, "*", "index.html"))
	if err != nil {
		t.Fatalf("globbing %s/%s: %v", root, section, err)
	}
	return matches
}

// collectHTMLFiles returns every .html file beneath root.
func collectHTMLFiles(t *testing.T, root string) []string {
	t.Helper()
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walking %s: %v", root, err)
	}
	if len(files) == 0 {
		t.Fatalf("no HTML files found under %s", root)
	}
	return files
}

// internalHrefs returns the hrefs in an HTML file that point within the site.
func internalHrefs(t *testing.T, file string) []string {
	t.Helper()
	var hrefs []string
	for _, match := range hrefPattern.FindAllStringSubmatch(readFile(t, file), -1) {
		if externalPattern.MatchString(match[1]) {
			continue
		}
		hrefs = append(hrefs, match[1])
	}
	return hrefs
}

func TestLinksResolve(t *testing.T) {
	site := siteDir(t)
	prefix := pathPrefix()
	checked := 0

	for _, file := range collectHTMLFiles(t, site) {
		rel, _ := filepath.Rel(site, file)
		t.Run(rel, func(t *testing.T) {
			for _, href := range internalHrefs(t, file) {
				checked++
				target := strings.TrimSuffix(href, "/")
				if strings.HasPrefix(href, prefix) {
					target = "/" + strings.TrimPrefix(href, prefix)
				}
				fsPath := filepath.Join(site, target)
				if strings.HasSuffix(href, "/") {
					fsPath = filepath.Join(fsPath, "index.html")
				}
				if _, err := os.Stat(fsPath); err != nil {
					t.Errorf("broken link %q → %s", href, fsPath)
				}
			}
		})
	}
	t.Logf("checked %d internal links", checked)
}

func TestLinksHavePrefix(t *testing.T) {
	site := siteDir(t)
	prefix := pathPrefix()
	checked := 0

	for _, file := range collectHTMLFiles(t, site) {
		rel, _ := filepath.Rel(site, file)
		t.Run(rel, func(t *testing.T) {
			for _, href := range internalHrefs(t, file) {
				checked++
				if !strings.HasPrefix(href, prefix) {
					t.Errorf("href %q is missing the path prefix %q", href, prefix)
				}
			}
		})
	}
	t.Logf("checked %d internal links", checked)
}

// TestSongPagesExist matches sources to pages by rendered title rather than by
// slug, so it never has to reimplement the generator's URL algorithm.
func TestSongPagesExist(t *testing.T) {
	root := repoRoot(t)
	site := siteDir(t)

	// Rendered title → number of pages carrying it.
	rendered := make(map[string]int)
	for _, page := range pagesUnder(t, site, "songs") {
		if match := songTitlePattern.FindStringSubmatch(readFile(t, page)); match != nil {
			rendered[textOf(match[1])]++
		}
	}

	entries, err := os.ReadDir(filepath.Join(root, defaultSongsDir))
	if err != nil {
		t.Fatalf("reading songs: %v", err)
	}

	sources := 0
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".md") || strings.HasPrefix(name, "_") {
			continue
		}
		sources++
		t.Run(name, func(t *testing.T) {
			match := titlePattern.FindStringSubmatch(readFile(t, filepath.Join(root, defaultSongsDir, name)))
			if match == nil {
				t.Fatalf("no title in frontmatter")
			}
			switch count := rendered[match[1]]; {
			case count == 0:
				t.Errorf("no page rendered for title %q", match[1])
			case count > 1:
				t.Errorf("title %q renders on %d pages — slug collision", match[1], count)
			}
		})
	}

	if pages := len(pagesUnder(t, site, "songs")); pages != sources {
		t.Errorf("song count mismatch: %d sources, %d pages", sources, pages)
	}
	t.Logf("checked %d songs", sources)
}
