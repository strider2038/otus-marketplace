# Проектная работа "Торговая площадка на основе потоковой обработки данных"

Работа выполнена для курса ["Microservice architecture"](https://otus.ru/lessons/microservice-architecture/)

## Описание проекта

1. [Требования к системе](docs/01_requirements.md)
2. [Проектирование микросервисной архитектуры на основе Event storming](docs/02_design.md)
3. [Проектирование микросервиса сделок](docs/03_trading_design.md)
4. [Спецификация OpenAPI](docs/public-api.yaml)
5. [Спецификация AsyncAPI](docs/async-api.yaml)
6. [Презентация по проекту](https://docs.google.com/presentation/d/1KrmSC7teapaxYjeN1FhxTqiPeEEtaix5PytTRu5WEwI/edit?usp=sharing)

## TODO

* [x] добавить портфель и историю в схемы
* [x] регистрация брокера
* [x] начисление комиссии брокера
* [x] сервис торговой площадки
* [x] сервис статистики
* [x] сервис истории сделок
* [x] сервис уведомлений (сделки)
* [x] тесты
* [x] мониторинг
* [ ] performance тесты
* [ ] github release

## Запуск приложения

```shell
# запуск minikube
# версия k8s v1.19, на более поздних есть проблемы с установкой ambassador
minikube start --cpus=8 --memory=24g --disk-size='80000mb' --vm-driver=virtualbox --cni=flannel --kubernetes-version="v1.19.0"

kubectl create namespace otus
kubectl config set-context --current --namespace=otus

# установка Ambassador
helm install aes datawire/ambassador -f deploy/ambassador-values.yaml

# установка Apache Kafka
helm install kafka bitnami/kafka -f deploy/kafka-values.yaml

# установка prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add stable https://charts.helm.sh/stable
helm install prom prometheus-community/kube-prometheus-stack -f deploy/prometheus.yaml --atomic

# установка ingress nginx
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install nginx ingress-nginx/ingress-nginx -f deploy/nginx-ingress.yaml --atomic

# запуск микросервисов
helm install --wait -f deploy/identity-values.yaml identity-service ./services/identity-service/deployments/identity-service --atomic
helm install --wait -f deploy/billing-values.yaml billing-service ./services/billing-service/deployments/billing-service --atomic
helm install --wait -f deploy/trading-values.yaml trading-service ./services/trading-service/deployments/trading-service --atomic
helm install --wait -f deploy/history-values.yaml history-service ./services/history-service/deployments/history-service --atomic
helm install --wait -f deploy/statistics-values.yaml statistics-service ./services/statistics-service/deployments/statistics-service --atomic
helm install --wait -f deploy/notification-values.yaml notification-service ./services/notification-service/deployments/notification-service --atomic

# применение настроек ambassador
kubectl apply -f services/ambassador/
```

## Мониторинг

```shell
# отладка prometheus
kubectl port-forward service/prom-kube-prometheus-stack-prometheus 9090

# отладка grafana
kubectl port-forward service/prom-grafana 9000:80
```

## Тестирование

Тесты Postman расположены в директории `test/postman`. Запуск тестов.

```bash
newman run ./test/postman/test.postman_collection.json
```

Или с использованием Docker.

```bash
docker run -v $PWD/test/postman/:/etc/newman --network host -t postman/newman:alpine run test.postman_collection.json
```
