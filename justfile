run:
    just run-server     
run-server:
    just server/init
    just server/run
run-web:
    just obcsapi-web/run
doc:
    cd docs && python3 -m http.server 8888
build:
    bash step1.sh
    bash step2.sh
    just server/build
    echo "Now push the Docker image to DockerHub"