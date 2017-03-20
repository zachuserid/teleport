.PHONY:start
start:
	docker build -t eg_sshd .
	docker run -d -P --name test_sshd --net telenet eg_sshd

.PHONY:stop
stop:
	docker rm -f test_sshd

.PHONY:enter
enter:
	docker exec -ti test_sshd /bin/bash
