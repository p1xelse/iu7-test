.PHONY: test start-docker integration-test test-allure clean
unit-test:
	cd internal && go test ./...

test-allure: clean
	sh allure_test_report.sh
	allure serve allure-results

test-allure-shuffle: clean
	sh allure_test_report_shuffle.sh
	allure serve allure-results

docker-unit-tests:
	docker build -t docker-tests-img -f tests.Dockerfile . && \
	docker run -it docker-tests-img go test -v ./internal/...

docker-integration-tests:
	docker build -t docker-tests-img -f tests.Dockerfile . && \
	docker run -it docker-tests-img go test -v ./integration_test/...

integration-test:
	sh allure_integration_test.sh
	allure serve allure-results

clean:
	rm -rf allure-results