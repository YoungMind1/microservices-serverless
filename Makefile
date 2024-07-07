clean:
	@go clean
	@rm -rf ./bin

build: clean
	cd functions/events/show && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/event-show
	cd functions/events/store && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/event-store
	cd functions/events/update && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/event-update
	cd functions/events/destroy && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/event-destroy
	cd functions/users/show && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/user-show
	cd functions/users/store && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/user-store
	cd functions/users/update && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/user-update
	cd functions/users/destroy && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/user-destroy
	cd functions/registrations/show && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/registration-show
	cd functions/registrations/store && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/registration-store
	cd functions/registrations/cancelRegistration && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../../../bin/registration-cancelRegistratoins
start:
	sudo sls offline --useDocker start --host 0.0.0.0
