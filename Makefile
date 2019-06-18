.PHONY: consumer

consumer:
	docker build -t lee0c/sbq-consumer ./consumer
	docker push lee0c/sbq-consumer

sessions:
	docker build -t lee0c/sbq-consumer-sessions ./session-consumer
	docker push lee0c/sbq-consumer-sessions