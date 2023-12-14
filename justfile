run:
    just run-server     
run-server:
    just server/init
    just server/run
run-web:
    just obcsapi-web/run
doc:
    cd obcsapi-docs/ && pnpm docs:dev 
doc-build:
    #!/bin/bash
    cd obcsapi-docs/ && pnpm docs:build
    rm -rf /home/kkbt/gitee-website/kkbt/obcsapi-go
    cp -r docs/.vuepress/dist /home/kkbt/gitee-website/kkbt/obcsapi-go
    cp -r docs/ /home/kkbt/gitee-website/kkbt/obcsapi-go/docsify/
doc-docsify:
    cd obcsapi-docs/docs && python3 -m http.server 8888
build:
    bash step1.sh
    bash step2.sh
    just server/build
    echo "Now push the Docker image to DockerHub"
update-web:
    cd obcsapi-web && npm-check -u 
update-server:
    cd server && go get -u && go mod tidy && go mod vendor