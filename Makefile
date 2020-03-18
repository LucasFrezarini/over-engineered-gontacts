integration-test:
	-docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down --volumes
integration-dependencies:
	-docker-compose -f docker-compose.test.yml up
	docker-compose -f docker-compose.test.yml down --volumes