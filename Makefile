.PHONY: all orders payments products records users appointments notifications run_gateway

all:
	@echo "Starting all services..."
	@make -j orders payments products records users appointments

orders:
	@echo "Starting Orders Service..."
	cd orders && go run . &

payments:
	@echo "Starting Payments Service..."
	cd payments && go run . &

products:
	@echo "Starting Products Service..."
	cd products && go run . &

records:
	@echo "Starting Records Service..."
	cd records && go run . &

users:
	@echo "Starting Users Service..."
	cd users && go run . &

appointments:
	@echo "Starting Appointment Service..."
	cd appointments && go run . &
notifications:
	@echo "Starting Notification Service..."
	cd notification && go run . &
stop:
	@echo "Stopping all services..."
	@pkill -f "go run"

run_gateway:
	@go run gateway/main.go