## Context

See proposal.md — Why. The relevant current state: the site deploys to GitHub Pages through
`.github/workflows/deploy.yml`, which builds on push to `master`, runs the Go tests, and publishes
`public/`. Pull requests run build and test with the deploy job gated to `master`. `hugo.toml` sets
`baseURL` to the Pages project URL, so every generated address carries a `/karaokay/` prefix, and
`tests/links_test.go` asserts that prefix on every internal href.

Two constraints shape the approach. First, GitHub Pages state lives in repository settings, not in the
repo — deleting the workflow does not unpublish anything. Second, `master` has no branch protection,
so the pull-request step was convention plus a workflow trigger, with nothing else enforcing it.

## Goals / Non-Goals

**Goals:**

- Leave no mechanism in the repo capable of publishing the site.
- Retire the public address rather than freezing a stale copy at it.
- Preserve local build, preview, and validation exactly as they work today.

**Non-Goals:**

- Replacing CI with a local equivalent (git hooks, pre-push validation). Validation stays manual and
  on-demand; adding an automatic local gate is a separate decision.
- Making the built output portable enough to open from the filesystem. See the third decision below.
- Removing the git remote. The repo stays on GitHub as off-machine backup; only publishing goes away.

## Decisions

**Disable the Pages site explicitly, rather than only deleting the workflow.**
Deleting the workflow stops future builds but leaves the last deployed build served at
`carlos.github.io/karaokay` forever. That is the worse of both outcomes: a public, permanently stale
copy of a copyrighted lyric corpus that no longer reflects the catalog and that nobody is maintaining.
Disabling Pages is a single API call (`DELETE /repos/carlos/karaokay/pages`) and is reversible from
repository settings. *Alternative considered:* leave it published as a read-only snapshot — rejected,
since it inherits every downside of publishing with none of the benefit.

**Set `baseURL` to `/` rather than keeping `/karaokay/` or removing the setting.**
The prefix exists solely because GitHub Pages serves project sites under a repository subpath. With
publishing gone it is vestigial, and it would make every local URL carry a meaningless segment.
Templates already resolve every internal link through `.RelPermalink` and `site.GetPage`, so they
adapt with no edits — this is precisely the property that made the earlier `relURL` bug worth fixing
properly. *Alternative considered:* remove `baseURL` entirely — rejected, because Hugo still needs a
value and an explicit `/` documents the intent.

**Access the site through `mise run dev`, keeping links root-relative.**
Root-relative addresses (`/songs/roar/`) require something to serve the site; they do not work when
`public/index.html` is opened directly from the filesystem. Since the site is browsed through the
development server, this is the correct trade. *Alternative considered:* `relativeURLs = true`, which
would make `public/` self-contained and copyable to a phone or USB stick — deferred, not rejected. It
remains a one-line change if portable output is ever wanted, and nothing in this change forecloses it.

**Keep the internal-href check, with `/` as the expected prefix.**
The check was added to catch hrefs that fail to carry the site's base path — the exact bug that would
have shipped a fully broken navigation. With `/` as the prefix it still catches relative and malformed
hrefs, which is most of its value. *Alternative considered:* delete it alongside the deployment it was
written for — rejected, because the failure it detects is a template bug, not a hosting concern.

**Delete the entire `.github/` directory, not just the workflow file.**
It contains nothing else. Leaving an empty `workflows/` directory invites something to be dropped back
into it.

## Risks / Trade-offs

**Nothing runs the tests automatically any more.** A broken build or dead link can now sit unnoticed
until the next manual run. → Accepted deliberately (see proposal.md — Impact). Mitigated by keeping
`mise run test` working and documenting it as the check to run after editing content. A git pre-push
hook would close this, and is left as a follow-up rather than smuggled into this change.

**Pages state is invisible from the repository.** A future reader cannot tell from the files that
publishing was disabled on purpose, and may re-add a workflow assuming it was simply missing. →
Mitigated by stating in `CLAUDE.md` that the site is deliberately unpublished and local-only.

**Every previously published URL breaks.** → Accepted; this is the point of the change. No redirects
are provided, since the destination is being retired rather than moved.

## Migration Plan

Order matters in one place: disable Pages **before or independently of** deleting the workflow. If the
workflow were deleted first and a push to `master` raced it, nothing would deploy anyway — but leaving
Pages enabled with no workflow is the ambiguous half-state worth avoiding.

1. Disable the Pages site.
2. Delete `.github/`.
3. Change `baseURL` and the test's expected prefix together — they must move in lockstep, or the link
   test fails against a correct build.
4. Rebuild, run the tests, and confirm generated hrefs carry no `/karaokay/` segment.
5. Update `CLAUDE.md`.

**Rollback:** re-enable Pages in repository settings, restore `.github/workflows/deploy.yml` from git
history, and revert `baseURL` and the test prefix. Every step is a revert of a tracked file except the
Pages toggle, which is a settings change.
