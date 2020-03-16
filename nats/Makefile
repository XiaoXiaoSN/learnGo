.PHONY: cluster pub sub


# https://hub.docker.com/_/nats-streaming
params ?= "-d"
clusterID := "test-cluster"
cluster:
	docker network create nats-work || true

	docker run $(params) --rm \
		-p 4223:4223 \
		-p 8223:8223 \
		--name nats1 \
		--network=nats-work \
		nats \
			-p 4223 -m 8223 \
			--cluster=nats://0.0.0.0:6223 \
			-c-id=$(test-cluster)
		 -DV 

	docker run $(params) --rm \
		-p 4224:4223 \
		-p 8224:8223 \
		--name nats2 \
		--network=nats-work \
		nats \
			-p 4223 -m 8223 \
			--cluster=nats://0.0.0.0:6223 \
			--routes=nats://nats1:6223,nats://nats3:6223 -DV 

	docker run $(params) --rm \
		-p 4225:4223 \
		-p 8225:8223 \
		--name nats3 \
		--network=nats-work \
		nats \
			-p 4223 -m 8223 \
			--cluster=nats://0.0.0.0:6223 \
			--routes=nats://nats2:6223,nats://nats3:6223 -DV 

pub: 
	go run pub.go

sub:
	go run sub.go
