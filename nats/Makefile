.PHONY: server pub sub

# https://hub.docker.com/_/nats-streaming
server:
	docker run -d \
		-p 4223:4223 \
		-p 8223:8223 \
		$(params) \
		nats-streaming -p 4223 -m 8223

pub: 
	go run pub.go

sub:
	go run sub.go
