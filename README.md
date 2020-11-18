# 0.x project - not ready for the light of day
Coming soon to an IDE near you...

## Easy Mongo
This project aims to be a friendly (somewhat opinionated) abstraction on top of mongo-go-driver. The official driver provides much more fine-grained control than this project attempts to solve.

#### Examples

```go
easymongo.CreateCreds()
conn := easymongo.Connect("mongo://127.0.0.1")
type myObj struct {}
foo := myObj{}
// TODO: Update this to the new syntax
err = conn.GetDatabase("my_db").GetCollection("my_coll").Insert(&foo)
```

#### Why use `easymongo`?
You should use `easymongo` if:
- You are planning on connecting primarily to a single cluster or instance
- Are looking to make common query operations with fewer lines of code
- Are okay with some logic being extrapolated and defaults assumed
- Accept one developer's opinionated approach to mongo in golang ;-)

The mongo-go-driver was written with the mindset of allowing complete control of every part of the query process. This project aims to build on that power to make life simple.

This can be run alongside mongo-go-driver in the same project without fear of conflicts. In fact, you can call `easymongo.ConnectUsingMongoClient(client *mongo.Client)` if you already have a pre-initialized `mongo.Client`. You will be able to start consuming the `easymongo.Connection` methods immediately.

## How are things different from mongo-go-driver?
Some functionality (such as the ability to directly set the context on a query) has been abstracted, while other functionality has intentionally been ignored in order to keep it simple for the core use cases.

This project is built on top of mongo-go-driver, so if you desire additional power (and really, who doesn't?), you can call `MongoDriverClient()` on the `Connection` object to obtain direct access to the `*mongo.Client` object. From there, the world is your oyster.

The decision (which I expect to see at least 1 GitHub issue debating) was made to cache the most recent mongo connection in a global variable under the covers in `easymongo`. This enables the following ease-of-use functions that would otherwise not be available:


## Contributors
Anyone is welcome to submit PRs. Please ensure there is test coverage before submitting the request.

## TODO
- [ ] Support new tag - jbson - which is a common tag that can be used for both json and bson
- [X] Support find
- [X] Support update/upsert/replace
- [X] Support count
- [ ] Support delete
- [ ] Support create
- [ ] Support gridfs
- [ ] Support watching
- [ ] Support collection ease-of-use enhancements
- [ ] Support indices ease-of-use enhancements
- [ ] Prevent users from being able to reconnect over and over to cluster - cache connection
- [ ] Add helpers for letting users run tests that:
  - automatically spin up a docker container if a DB isn't available locally
  - allow mocking of commands
  - support exporting data from the DB
  - support easy preloading of exported data
- [ ] Make running general mongo commands easier
- [ ] Add a tool for converting mongo queries to golang mongo queries
- [ ] Support optionally returning ErrNotFound using a global option when no update/upsert/delete results are found.
- [ ] Handle nil slice/map auto-initialization / type
- [ ] Handle datetime objects being strings
- [ ] Add DistinctStrings(), DistinctInts(), Distinct()
- [ ] Document bson.M vs bson.D
- [ ] Add the option to create easymongo from an existing mongo-go-driver session
- [ ] Support read/write concern
- [ ] Add option for ErrNotFound
- [ ] Create docs on using bson tags
- [ ] Link to mongo-go-driver docs where appropriate
- [ ] Change Do() to Execute()

#### The Genesis
I started out with Golang and mongo by consuming the mgo library. I enjoyed it, but found the fact that it was radically out of date to be off-putting. When we were looking at migrating to DocumentDB with AWS, we realized that the mgo driver we were using was not compatible for a lot of our common operations, so it was obvious it was time to upgrade. I helped my company rewrite a massive library to consume mongo-go-driver. I found the experience...frustrasting to say the least. With the excellent support of the mongo-go-driver team though (through Slack, email and video calls), we were able to complete the transition and get things working properly. All the lines of code were more verbose and felt more error prone. (It was easy to miss checking for a context error even though one checked for the cursor Decode error for example.)

Frustration breeds creativity though - after a lot of careful thought and some back and forth with the mongo-go-driver core development team, I set forth to create a balance. I wanted something simpler - less lines of code to consume but still leveraging the power of the new driver. Thus `easymongo` was born.