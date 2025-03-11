module github.com/josephburgess/joeburgess.dev

go 1.24.1

require (
	github.com/gorilla/mux v1.8.1
	github.com/jarcoal/httpmock v1.3.1
	github.com/joho/godotenv v1.5.1
	github.com/josephburgess/glog v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.1
	go.uber.org/zap v1.27.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/yuin/goldmark v1.7.8 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/josephburgess/glog => ../glog
