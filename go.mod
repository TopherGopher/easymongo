module github.com/TopherGopher

go 1.15

require (
	docker.io/go-docker v1.0.0
	github.com/Microsoft/go-winio v0.4.15 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker/internal/testutil v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	go.mongodb.org/mongo-driver v1.4.3
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/docker/docker/internal/testutil => gotest.tools/v3 v3.0.0
