module github.com/tophergopher/easymongo

go 1.16

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tophergopher/mongotest v0.0.27
	go.mongodb.org/mongo-driver v1.7.1
)

// replace github.com/tophergopher/mongotest v0.0.27 => ../mongotest
