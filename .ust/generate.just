[doc('Generate mocks, swagger, and templ files')]
[group('core')]
generate: _generate-mocks _generate-swagger _generate-templ

[group('generate')]
_generate-mocks:
    just _title "Generate" "Mocks"
    just _run "{{mockery}}"

[group('generate')]
_generate-swagger:
    just _title "Generate" "Swagger"
    just _run "cd components/fun-app && {{swag}} i --parseDependency true"
    just _info "Swagger" "http://localhost:8080/swagger/index.html"

[group('generate')]
_generate-templ:
    just _title "Generate" "Template Files"
    just _templ common/ui
    just _templ components/learn
    just _templ components/kohan
    just format
