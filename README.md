# 0.x project - not ready for the light of day
Coming soon to an IDE near you...
![Build Status](https://github.com/TopherGopher/easymongo/workflows/Go/badge.svg?branch=master) [![Coverage Status](https://coveralls.io/repos/github/TopherGopher/easymongo/badge.svg?branch=master)](https://coveralls.io/github/TopherGopher/easymongo?branch=master) [![Go Reference Docs](https://pkg.go.dev/badge/github.com/tophergopher/easymongo.svg)](https://pkg.go.dev/github.com/tophergopher/easymongo)

## Easy Mongo
This project aims to be a friendly (somewhat opinionated) abstraction on top of mongo-go-driver. The official driver provides much more fine-grained control than this project attempts to solve.

#### Get Connected
`easymongo` supports standard URI connection strings. All you have to do is call `Connect()`
(_A note that srv DNS records are not yet supported - this is currently
a go limitation_)
```go
// Connect to a locally running container
conn := easymongo.Connect("mongo://127.0.0.1")
type myObj struct {}
foo := myObj{}
err = conn.D("my_db").C("my_coll").Insert(&foo)
```
Don't want to go through the arduous process of setting up a local mongo environment?
You can spawn a container for the life of a test by using `tophergopher/mongotest`:
```go
useDockerContainer := true
// conn is a mongotest.TestConnection object which embeds an easymongo.Connection object
// When this command runs, a mongo docker container is spawned, bound on any available port
// and an easymongo connection is made. This connection object can be used just like a
// standard easymongo.Connection object
conn, err := mongotest.NewTestConnection(useDockerContainer)
// Find the first document with name sorted in descending order
err = conn.Find().Sort("-name").Skip(1).Limit(2).One(bson.M{"name": bson.M{"$ne": nil}})
// Then kill the test container
conn.KillMongoContainer()
```

## TODO: C.R.U.D. Examples
Create/Read/Update/Destroy!
#### Get Data into the Database
Let's talk about first how to get data into the database.
#### Find it and query it back
#### Modify it
#### Delete it
#### Examples



#### Why use `easymongo`?
You should use `easymongo` if:
- You are planning on connecting primarily to a single cluster or instance
- Are looking to make common query operations with fewer lines of code
- Are okay with some logic being extrapolated and defaults assumed
- You can accept one developer's opinionated approach to mongo in golang ;-)

The mongo-go-driver was written with the mindset of allowing complete control of every part of the query process. This project aims to build on that power, aiming to make integrating your project with mongo easier.

This can be run alongside mongo-go-driver in the same project without fear of conflicts. In fact, you can call `easymongo.ConnectUsingMongoClient(client *mongo.Client)` if you already have a pre-initialized `mongo.Client` and are considering trying `easymongo`. You will be able to start consuming the `easymongo.Connection` methods immediately.

## How are things different from mongo-go-driver?
Some functionality (such as the ability to directly set the context on a query) has been abstracted, while other functionality has intentionally been ignored in order to keep it simple for the core use cases. If you still want to set a context, you can explicitly set it with `SetContext(ctx)` on any Query object.

This project is built on top of mongo-go-driver, so if you desire additional power (and really, who doesn't?), you can call `MongoDriverClient()` on the `Connection` object to obtain direct access to the `*mongo.Client` object. From there, the world is your oyster.

The decision (which I expect to see at least 1 GitHub issue debating) was made to cache the most recent mongo connection in a global variable under the covers in `easymongo`. This means that once a connection is initialized, `conn := easymongo.GetCurrentConnection()` will return the most recent initialized DB connection object. This allows for easy conversion to mongo without the need to change dozens of function headers to expose the connection.

## Contributors
Anyone is welcome to submit PRs. Please ensure there is test coverage before submitting the request.
