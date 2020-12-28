module github.com/tophergopher/easymongo

go 1.15

require (
	github.com/stretchr/testify v1.6.1
	github.com/tophergopher/mongotest v0.0.5
	go.mongodb.org/mongo-driver v1.4.3
)
replace github.com/tophergopher/mongotest => ../mongotest
