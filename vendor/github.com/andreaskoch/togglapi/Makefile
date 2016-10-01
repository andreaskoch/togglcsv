test:
	go test
	go test ./date

coverage:
	go test ./ -coverprofile=coverage-api.out && go tool cover -html=coverage-api.out
	go test ./date -coverprofile=coverage-date.out && go tool cover -html=coverage-date.out
