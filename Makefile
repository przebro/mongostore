dockerbuild:
	docker build -t databazaar/mongodrv -f docker/Dockerfile .
start:
	docker run -d -l mongobzr1 -p20017:27017 -v ${PWD}/docker/etc:/etc/mongo \
		-e MONGO_INITDB_ROOT_USERNAME=admin \
		-e MONGO_INITDB_ROOT_PASSWORD=notsecure \
		databazaar/mongodrv --config /etc/mongo/config.yml
stop:
	docker rm -f $$( docker ps -qaf "label=mongobzr1")
tests:
	go test ./... -covermode=count --coverprofile='coverage.out'
	go tool cover -html coverage.out 