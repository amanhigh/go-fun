# Resources
# OpenAPI 3.0 Tutorial - https://golangexample.com/openapi-client-and-server-code-generator/
# Swagger Editor - https://editor.swagger.io/
# oapi-codegen - https://github.com/deepmap/oapi-codegen
# Expanded PetStore Eg - https://github.com/deepmap/oapi-codegen/blob/master/examples/petstore-expanded/petstore-client.gen.go
#
# Installation: go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest


echo "\033[1;32m Generate Server (Gin) Code & Types \033[0m \n";
oapi-codegen -generate types,gin -package openapi -o server.gen.go person.yaml
echo "\033[1;34m Now Implement Server Interface \033[0m \n";

echo "\033[1;32m Generate Server Client SDK & Types \033[0m \n";
oapi-codegen -generate types,client -package openapi -o client.gen.go person.yaml
echo "\033[1;34m Use NewClientWithResponse() and Methods in it \033[0m \n";
