# Проектная работа "Торговая площадка на основе потоковой обработки данных"

Работа выполнена для курса ["Microservice architecture"](https://otus.ru/lessons/microservice-architecture/)

## TODO

* [x] регистрация брокера
  * [x] user role
  * [x] user role in gateway
* [ ] начисление комиссии брокера
  * [x] начислять комиссию
  * [ ] идемпотентные методы
* [ ] сервис торговой площадки
  * [x] код
  * [ ] helm
  * [ ] api gateway mapping
  * [ ] идемпотентные методы
* [ ] сервис статистики
* [ ] сервис уведомлений (сделки)
* [ ] тесты

* дополнительно
  * [ ] debezium

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

## запуск проекта
helm install --wait -f deploy/identity-values.yaml identity-service ./services/identity-service/deployments/identity-service --atomic
helm install --wait -f deploy/billing-values.yaml billing-service ./services/billing-service/deployments/billing-service --atomic
helm install --wait -f deploy/order-values.yaml order-service ./services/order-service/deployments/order-service --atomic
helm install --wait -f deploy/notification-values.yaml notification-service ./services/notification-service/deployments/notification-service --atomic

# применение настроек ambassador
kubectl apply -f services/ambassador/
```
