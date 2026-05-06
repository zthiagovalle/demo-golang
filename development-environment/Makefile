.PHONY: up down logs psql pgadmin awslocal-list-queues awslocal-publish

up:
	docker compose up -d

pgadmin:
	open http://localhost:5050 || xdg-open http://localhost:5050

down:
	docker compose down -v

logs:
	docker compose logs -f

psql:
	docker exec -it demo-postgres psql -U demo -d demo

awslocal-list-queues:
	docker exec demo-localstack awslocal sqs list-queues --region us-east-1

awslocal-receive:
	docker exec demo-localstack awslocal sqs receive-message \
	  --queue-url http://localhost:4566/000000000000/product-created-queue \
	  --region us-east-1
