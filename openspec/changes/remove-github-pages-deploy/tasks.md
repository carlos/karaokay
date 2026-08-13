## 1. Retire the public site

- [x] 1.1 Disable GitHub Pages for `carlos/karaokay` (`gh api -X DELETE repos/carlos/karaokay/pages`)
- [x] 1.2 Confirm `https://carlos.github.io/karaokay/` no longer returns the catalog
- [x] 1.3 Delete the `.github/` directory, including `workflows/deploy.yml`

## 2. Move the site to root addressing

- [x] 2.1 Change `baseURL` in `hugo.toml` from `https://carlos.github.io/karaokay/` to `/`
- [x] 2.2 Change `defaultPrefix` in `tests/links_test.go` from `/karaokay/` to `/`, and rename the
      `KARAOKAY_PATH_PREFIX` override if its name no longer reads accurately
      (kept — it still describes what it overrides, only the default changed)
- [x] 2.3 Rebuild from clean (`rm -rf public && mise run build`) and confirm generated hrefs contain
      no `/karaokay/` segment

## 3. Verify

- [x] 3.1 Run `mise run test` — all three checks pass against the root-addressed build
- [x] 3.2 Confirm the test still fails when given a malformed href, so the prefix check is doing real
      work at `/` and did not become vacuous
- [x] 3.3 Start `mise run dev` and confirm the site serves at `localhost:8081/` with working
      navigation, search, and a song page

## 4. Update documentation

- [x] 4.1 Remove the Deploy section from `CLAUDE.md`
- [x] 4.2 Rewrite the path-prefix convention in `CLAUDE.md` for root addressing, keeping the "never
      hand-build an internal href" rule — it guards a template bug, not a hosting concern
- [x] 4.3 State in `CLAUDE.md` that the site is deliberately local-only and unpublished, so a future
      reader does not restore a deployment workflow assuming one went missing
- [x] 4.4 Note that `mise run test` is now manual and should be run after content or template edits

## 5. Close out

- [x] 5.1 Commit directly to `master` — no pull request, matching the contribution requirement in the
      spec
- [x] 5.2 Confirm no workflow runs were triggered by the push
