# This rule removes all directories named '.tmp' under the project root.
clean:
	@find ./../ -name .tmp -type d -exec rm -rf {} +

# This rule removes all directories named '.db-data' under the project root.
purge:
	@find ./../ -name .db-data -type d -exec rm -rf {} +

# This rule builds the Docker containers defined in the 'docker-compose.yaml' file without using cache.
# Before building the containers, it performs the 'clean' action to remove any '.tmp' directories.
build: clean
	@docker-compose build --no-cache

# This rule starts the Docker containers. Before starting, it ensures that the containers are built using the 'build' rule.
run: build
	@docker-compose up

# This rule is a shortcut for cleaning up database data, building the Docker containers, and starting them.
rerun:
	@make purge
	@make build
	@docker-compose up

# This rule stops running Docker containers and removes them along with their networks and volumes.
stop:
	@docker-compose down

# This rule runs unit tests in the 'eth-graph-api' directory using Go's built-in testing tool.
test:
	@cd ../eth-graph-api && go test -count=1 ./...