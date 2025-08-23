.PHONY: test
test:
	go test -v -count=1 ./internal/adapters/s3/ -run ^TestUploadLogo$

.PHONY: ti
ti:
	go test -v -count=1 ./test/integration/ -run ^"TestSetCompanyEmail/should return 400 for email exceeding 255 characters"$

.PHONY: db
db:
	go test -v -count=1 ./internal/adapters/postgres/... 

.PHONY: all
all:
	go test -v -count=1 ./test/integration/... 
