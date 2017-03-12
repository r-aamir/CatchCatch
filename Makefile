INSTANCE_NAME := tile38-local
TEMPLATE := tile38/tile38
PORT := 9851

.PHONY=run
run:
	cd go && CompileDaemon -color -command "./go"

.PHONY=run-tile38
run-tile38:
	@if test "`docker ps -a | grep $(INSTANCE_NAME)`"; then \
		docker start $(INSTANCE_NAME); \
	else \
		docker run --rm --name $(INSTANCE_NAME) -v $$PWD:/data -p $(PORT):$(PORT) -P $(TEMPLATE);\
	fi
