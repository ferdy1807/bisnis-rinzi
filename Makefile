.PHONY: start stop

SERVICES = cash finance inventory pos rental

start:
	@echo "Starting all services..."
	@mkdir -p log
	@echo "Starting gateway..."
	@nohup go run gateway/main.go > log/gateway.log 2>&1 &
	@echo "Starting auth..."
	@nohup go run auth/cmd/main.go > log/auth.log 2>&1 &
	@for service in $(SERVICES); do \
		echo "Starting $$service..."; \
		nohup go run services/$$service/cmd/main.go > log/$$service.log 2>&1 & \
	done
	@echo "All services started. Logs are saved in the 'log' directory."

stop:
	@echo "Stopping all processes on ports 8080-8086..."
	@fuser -k -9 8080/tcp 8081/tcp 8082/tcp 8083/tcp 8084/tcp 8085/tcp 8086/tcp 2>/dev/null || true
	@echo "Cleaning up log directory..."
	@rm -rf log
	@echo "All services stopped."
