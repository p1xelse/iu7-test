.PHONY: test start-docker integration-test test-allure clean
test:
	go clean -testcache
	cd internal && go test $$(go list ./... | grep -v /mocks) -test_dsn="host=localhost user=test password=test database=postgres port=13081" -cover

test-allure: clean
	sh allure_test_report.sh
	allure serve allure

integration-test-db:
	docker compose down && docker compose up -d
	cd test_integration && go test .


start-docker:
	docker compose down && docker compose up -d

integration-test: start-docker
	newman run Timetracker.postman_collection.json

clean:
	rm -rf allure