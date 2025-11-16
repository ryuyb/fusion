package fusion

//	@title			Fusion API
//	@version		1.0
//	@description	This is a sever for Fusion

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securitydefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

//go:generate go run -mod=mod github.com/swaggo/swag/cmd/swag init -g ./docs.go -o ./docs/api
