gen:
	gqlgen generate

GQL_URL ?= http://localhost:8000/graphql
GQL_TIMEOUT ?= 10
KEEP_DATA ?= false

ifeq ($(OS),Windows_NT)
PYTHON ?= py -3
else
PYTHON ?= python3
endif

.PHONY: gen test

test:
	$(PYTHON) test_gql_endpoints.py --url $(GQL_URL) --timeout $(GQL_TIMEOUT) $(if $(filter true,$(KEEP_DATA)),--keep-data,)