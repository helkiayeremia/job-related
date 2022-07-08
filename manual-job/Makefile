run:
	docker network create redis-network
	docker run --name redis --network redis-network -p 6380:6379 redis redis-server --save 3600 1 --loglevel warning
client:
	docker run --name redisclient --network redis-network -it --rm redis redis-cli -h redis