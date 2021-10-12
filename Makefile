docker-build:
	bash ./docker-build.sh

docker-push:
	bash ./docker-build.sh push

kube-apply:
	kubectl apply -f k8s/data-sweeper.yaml

kube-delete:
	kubectl delete -f k8s/data-sweeper.yaml
