# 0.x project - not ready for the light of day
Coming soon to an IDE near you...
![Build Status](https://github.com/TopherGopher/easymongo/workflows/Go/badge.svg?branch=master) [![Coverage Status](https://coveralls.io/repos/github/TopherGopher/easymongo/badge.svg?branch=master)](https://coveralls.io/github/TopherGopher/easymongo?branch=master) [![Go Reference Docs](https://pkg.go.dev/badge/github.com/tophergopher/easymongo.svg)](https://pkg.go.dev/github.com/tophergopher/easymongo)

## Easy Mongo
This project aims to be a friendly (somewhat opinionated) abstraction on top of mongo-go-driver. The official driver provides much more fine-grained control than this project attempts to solve.

## TODO: C.R.U.D. Examples
Create/Read/Update/Destroy!
#### Get Data into the Database
Let's talk about first how to get data into the database.
#### Find it and query it back
#### Modify it
#### Delete it
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
