test:
	TAGS=integration docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
test-dependencies:
	-docker-compose -f docker-compose.test.yml up contacts_mysql_test
	docker-compose -f docker-compose.test.yml down --volumes