run-tests:
		docker volume prune -f && \
		docker-compose -f ./docker-compose.test.yaml build && \
		docker-compose -f ./docker-compose.test.yaml \
		run test-customer-service gotestsum --format testname ./... && \
		docker-compose -f ./docker-compose.test.yaml down

golangci:
	golangci-lint run