## Run Test Suite inside docker
tests:
	@docker-compose -f docker-compose.yml up --build --abort-on-container-exit
	@docker-compose -f docker-compose.yml down --volumes

## Run Integration Test
## Note: This command is intended to be executed within docker env
integration-tests:
	@sh -c "while ! curl -s http://rabbitmq:15672 > /dev/null; do echo waiting for 3s; sleep 3; done"
	go test -v ./...
