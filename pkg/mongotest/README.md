# mongotest
`mongotest` was written with the goal of testing code that relies on mongo simpler. The only requirement to use mongotest is to have docker running (should you wish to use the docker container method).

Example:
```go
func TestFoo(t *testing.T) {
  useDockerContainer := true
  // conn is a mongotest.TestConnection object which embeds an easymongo.Connection object
  conn, err := mongotest.NewTestConnection(useDockerContainer)
  is.NoError(err)
  t.Cleanup(func() {
    conn.KillMongoContainer()
  })
  // Insert a document
  type enemy struct {
    ID   primitive.ObjectID  `bson:"_id"`
    Name string              `bson:"name"`
  }
  id, err := conn.Insert().One(&enemy{
    ID:   primitive.NewObjectID(),
    Name: "The Joker",
  })
}
```
The above code will spin-up a docker mongo container on a randomly assigned port, insert a document into the collection and when the test exits, the mongo container will be destroyed. In order to ensure that the docker container gracefully exits, it is recommended to run the `.KillMongoContainer()` command in a `t.Cleanup()` function.

If you choose to use a `defer` (rather than `t.Cleanup()`), note that it is (presently) not possible to automatically cleanup the created container should the test panic.

_* I was wondering if a compiler flag might be the way to go to always ensure clean-up, but I truly welcome input on how this might be accomplished cleanly._

# Cleaning up rogue containers
Containers are created with a label of `mongotest=regression`. If you run `docker ps` and note a lot of unreaped mongo containers, try running:

```shell
    docker rm --force $(docker ps -a -q --filter=label=mongotest=regression)
```

This will thwack any containers that were created via mongotest. 