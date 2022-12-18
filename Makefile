.SILENT:

help:
	printf "Available targets\n\n"
	awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "%-30s %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)


.PHONY: py_check
# Internal helper target - check if python is installed
py_check:
	{ \
	if ( ! ( command -v python3 >/dev/null  )); then \
		echo "Seems like you don't have Python installed. Make sure you review docs/development/README.md before continuing"; \
		exit 1; \
	fi; \
	}


.PHONY: py_format
## Format the python code
py_format:
	isort .
	black .

# Variables for `p2p_test_generator`
numRainTreeNodes ?= 12 # This is the default value with a randomly selected value
rainTreeTestOutputFilename ?= "/tmp/answer.go" # This is the default file where the test will be written to

.PHONY: p2p_test_generator
## Generate a RainTree unit test configured for `numRainTreeNodes` and written to `rainTreeTestOutputFilename`
p2p_test_generator: py_check
	echo "See python/README.md for additional details"
	python3 python/main.py --num_nodes=${numRainTreeNodes} --output_file=${rainTreeTestOutputFilename}