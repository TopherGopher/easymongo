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
err = conn.D("my_db").C("my_coll").Insert().One(&foo)
```
Don't want to go through the arduous process of setting up a local mongo environment?
You can spawn a container for the life of a test by using `github.com/tophergopher/mongotest`:
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

## C.R.U.D. Examples
Create/Read/Update/Destroy!
#### Get Data into the Database
Let's talk about first how to get data into the database. Create a struct and add some reflection tags. These tags map to which field gets written. For example:
```go
type Enemy struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          string             `bson:"name"`
	Notes         string             `bson:"notes,omitempty"`
	LastEncounter *time.Time         `bson:"lastEncounter"`
	Deceased      bool               `bson:"deceased"`
	TimesFought   int                `bson:"timesFought"`
	Evilness      float64            `bson:"evilness"`
}
var theJoker = Enemy{
  ID: primitive.NewObjectID(),
  Name: "The Joker",
  Notes: "Follow-up about his scars.",
  TimesFought: 3,
  Evilness: 92.6,
}
```
Cool - so now we have an object - let's insert it:
```go
  id, err := coll.Insert().One(theJoker)
```
Hooray! The document is in the database! Other typical endings for `Insert()` are `.Many()`, which will insert a slice of any type (but with a small O(N) overhead), and `.ManyFromInterfaceSlice()`, which doesn't incur the O(N) overhead but requires you to rewrite slices to a slice of interface.

#### Find it and query it back
Find the first matching document:
```go
  var eLookup Enemy
  err := coll.Find(bson.M{"name": "The Joker"}).One(&eLookup)
  fmt.Printf("%#v\n", eLookup)
```
If you examine `eLookup`, `eLookup` will be populated with all the initial metadata from `theJoker` object.

Here's an example which leverages more flags to look-up all entries in the database and loads them into the `enemies` slice.
```go
  var enemies []Enemy
  err := coll.Find(bson.M{}).Comment(
			"Isn't this a fun query?").BatchSize(5).Projection(
			bson.M{"name": 1}).Hint("name").Sort(
			"-name").Skip(0).Limit(0).Timeout(time.Hour).All(&enemies)
```

#### Modify it
Let's say Batman runs into the Joker, Alfred will need to update the last time they ran into eachother:
```go
	err := coll.Update(bson.M{"name": "The Joker"}, bson.M{
			"$set": bson.M{"lastEncounter": time.Now()}}).One()
```
You'll note that we can end with either `.One()` or `.Many()`.

#### Find and mutate
What about if you want to return a document AND mutate it in some groovy way? `Find.OneAnd()` is your friend!
```go
  var enemyBefore Enemy
  filter := bson.M{"name": "The Joker"}
  set := bson.M{"$set": bson.M{"notes": "What's his endgame?"}}
  err := coll.Find(filter).OneAnd(&enemyBefore).Update(set)
```
This would return the object prior to being updated. Other endings to `OneAnd` are `Replace()` and `Delete()`. If you want the resultant object _after_ the mutation takes place, then add `ReturnDocumentAfterModification()` to the chain.

#### Delete it
When The Joker breaks into The Bat Cave, he wants to delete all of Batman's data. 
This would delete all documents for enemies encountered in the last 6 months:
```go
  filter := bson.M{"lastEncounter": bson.M{
    "$gte": time.Now().AddDate(0, -6, 0),
    },
  },
  err := coll.Delete(filter).Many()
```
`.One()` is also available for individual document deletions using a filter.

#### Mutate by ID
A note that all of the various CRUD operators have ObjectID helpers
`FindByID`, `DeleteByID`, `UpsertByID`, `ReplaceByID`. If you have the ObjectID in memory, it is one of the fastest methods for updating an object.
```go
  replacementEntry := Enemy{ID: theJoker.ID, Name: "Who knows?"}
  err := coll.ReplaceByID(theJoker.ID, replacementEntry)
```

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
