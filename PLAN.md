# easymongo refactor plan

The mongotest side of this work is described in `mongotest/PLAN.md`.

## Where things stand

**Branch**: `claude/mongotest-easymongo-refactor-wu2vhm`, which currently
holds this plan and nothing else.
**Tracking issue**: #13. The mongotest side is TopherGopher/mongotest#48.

Nothing here has started, and it is blocked rather than merely unstarted.
Issue #7 needs `github.com/tophergopher/mongotest/v2`, which does not exist
yet: it is TopherGopher/mongotest#42, and that waits on the `mongod` package
(#32 there), which is the next thing being built. Until the v2 module is
tagged, #7 says to use a local `replace` to a sibling checkout.

The order in the tracking issue matters and is not the issue-number order.
The driver package cannot compile half-migrated, so #2 to #5 land together in
one pull request, and #7 and #8 must land first so the test suite stops
depending on the old `github.com/tophergopher/mongotest` v0.1.0 module, which
imports easymongo and breaks the moment easymongo's types change.

### What changed on the mongotest side that reaches here

The Docker client there is now an interface, `dockerclient.Client`, with
three implementations behind it: `dockerapi` (standard library, the default),
`mobyclient` (wrapping the official client, its own module), and `dockermock`
(test doubles). Issue #7's snippets were updated to name `dockerclient` for
`IsNotFound`; the old `dockerapi` names still compile as aliases, so nothing
here is broken by it.

Also relevant when the helpers get written: mongotest now finds a Podman
socket with no configuration, and TopherGopher/mongotest#52 covers reaching
mongo when the daemon is a sibling container or behind a socket proxy. If the
easymongo suite runs in CI with a proxied docker socket, that issue is the
one to watch.

## Goals

1. Migrate from `go.mongodb.org/mongo-driver` v1 to `go.mongodb.org/mongo-driver/v2`.
2. Take the mongotest helpers that used to live in mongotest
   (`EasyMongoWithContainer`, the easymongo-embedding `TestConnection`) and
   host them here, in a subpackage, so the dependency is one-way.
3. Latest Go (`go 1.27`), latest GitHub Actions, latest packages; drop the
   `replace` lines and the logrus dependency.
4. Test-driven; the existing test suite is the regression net for the
   driver migration and must stay green.

## Decisions taken (with the user)

| Topic | Decision |
|---|---|
| Helper location | New subpackage `github.com/tophergopher/easymongo/easymongotest` wrapping `mongotest/v2`. Production importers of `easymongo` never compile Docker code. |
| Test structure | Keep one container per test function (existing `setup(t)` pattern), now backed by `easymongotest`. |
| Logging | `DefaultLogger` reimplemented on `log/slog`; the exported `Logger` interface (`Debugf`, `Errorf`) is unchanged. |
| Versioning | easymongo is 0.x; the `primitive.ObjectID` to `bson.ObjectID` change in the public API ships as a minor bump, not an `easymongo/v2` module. |
| Go version | `go 1.27`. |

## Driver v2 migration inventory

From a full scan of the module. File references are current line numbers.

Connection (`connect.go`):
- `mongo.NewClient` + `client.Connect(ctx)` (281/290, 315/326) become
  `mongo.Connect(opts)`. Add `Connection.Disconnect(ctx)`; nothing calls
  `Disconnect` today.
- `EnableDebug` copies a client struct (`*conn.client = *client`). Replace
  with disconnect, reconnect, swap the pointer.
- `bson.NewRegistryBuilder` + `bsoncodec.NewSliceCodec(bsonoptions...)`
  (91-107) becomes `opts.SetBSONOptions(&options.BSONOptions{NilSliceAsEmpty: true})`.
- `writeconcern.New(writeconcern.W(n))`, `writeconcern.J`, `WMajority` and
  the `[]writeconcern.Option` pattern (116, 471-484, 520-535) become
  `&writeconcern.WriteConcern{W: n, Journal: &j}` / `writeconcern.Majority()`.
- `readconcern.*()` and `readpref.*()` constructors are unchanged.
- `event.CommandFailedEvent.Failure` is an `error` in v2 (152-154).
- New: `ConnectionBuilder.TLSConfig(*tls.Config)` so the TLS mongotest mode
  can be exercised.

Options passed as struct literals (`FindOneOptions`, `FindOptions`,
`CountOptions`, `DistinctOptions`, `AggregateOptions`, `UpdateOptions`,
`DeleteOptions`, `FindOneAndUpdateOptions`, `FindOneAndReplaceOptions`,
`FindOneAndDeleteOptions` in `find_query.go`, `aggregate.go`,
`update_query.go`, `delete_query.go`, `find_and_query.go`): v2 methods take
`options.Lister[T]`, so each becomes a builder chain that only sets non-nil
fields. `MaxTime` no longer exists; the per-query context timeout that
easymongo already builds in `query.go` covers it. `options.ArrayFilters`
struct becomes `[]any`. `DeleteOptions` splits into one/many variants.

BSON: `primitive.ObjectID`, `primitive.NewObjectID`,
`primitive.NewObjectIDFromTimestamp` (insert_query.go, delete_query.go,
collection.go, tests) become `bson.*`. `x/bsonx` in `index.go` (already
deprecated in v1) becomes `bson.D`.

Results: `Distinct` returns `*mongo.DistinctResult` with `Decode`
(find_distinct.go). `mongo.ErrNoDocuments` and `UpdateResult` fields are
unchanged. No sessions, transactions or change streams exist to migrate.

## easymongotest package

```go
package easymongotest // github.com/tophergopher/easymongo/easymongotest

type TestConnection struct {
    *easymongo.Connection
    Mongo *mongotest.Instance   // mongotest/v2
}
func Start(ctx context.Context, opts ...mongotest.Option) (*TestConnection, error)
func Run(tb testing.TB, opts ...mongotest.Option) *TestConnection
func NewTestConnection(spinupDockerContainer bool) (*TestConnection, error) // migrated wrapper
func (tc *TestConnection) KillMongoContainer() error                        // migrated wrapper
func WithContainer(f func(*easymongo.Connection) error) error               // was mongotest.EasyMongoWithContainer
```
`Start` builds the easymongo connection from `Instance.URI()` and, in TLS
mode, from `Container.TLSConfig()`. `common_test.go` switches to this package
with no change to the per-test container pattern.

## TDD sequence

1. Add `easymongotest` with tests (start, insert through the easymongo API,
   kill, `WithContainer`, TLS mode) against a mongotest/v2 prerelease or a
   local `replace` during development only.
2. Migrate `connect.go` first with `connect_test.go` extended for
   `Disconnect`, `TLSConfig`, `EnableDebug` reconnect.
3. Migrate query files one at a time, each with its existing `_test.go`
   green before moving on; add tests for the `Distinct` decode path and
   `ArrayFilters`.
4. Replace `DefaultLogger` internals with slog; test that `Debugf` and
   `Errorf` reach the handler.
5. Toolchain last: `go 1.27`, `go mod tidy`, remove `replace` lines, update
   `.github/workflows/go.yml` (checkout, setup-go, coveralls, gcov2lcov to
   current tags; `go mod download` instead of `go get -d`), add dependabot
   for github-actions, rewrite the `Makefile` release target (it no longer
   bumps mongotest's go.mod), refresh the README examples.

## Open items with proposed defaults

- Whether to align easymongo's `Logger` interface with mongotest's structured
  one so the logrus/zap/zerolog adapter modules serve both. Default: leave
  easymongo's `Debugf`/`Errorf` interface alone in this pass.
- Exact action tags confirmed at implementation time.
