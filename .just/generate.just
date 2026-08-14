[doc('Generate mocks, swagger, and templ files')]
[group('core')]
generate: _generate-mocks _generate-fun-swagger _generate-templ

[group('generate')]
_generate-mocks:
    just _title "Generate" "Mocks"
    just _run "{{mockery}}"

[group('generate')]
_generate-fun-swagger:
    just components/fun-app/generate-swagger

[group('generate')]
_generate-templ:
    just _title "Generate" "Template Files"
    just _templ common/ui
    just _templ components/learn
    just _templ components/kohan
    just format
