# Проектная работа "Торговая площадка на основе потоковой обработки данных"

Работа выполнена для курса ["Microservice architecture"](https://otus.ru/lessons/microservice-architecture/)

## TODO

* [ ] добавить портфель и историю в схемы
* [x] регистрация брокера
  * [x] user role
  * [x] user role in gateway
* [ ] начисление комиссии брокера
  * [x] начислять комиссию
  * [ ] идемпотентные методы
* [ ] сервис торговой площадки
  * [x] код
  * [x] helm
  * [x] api gateway mapping
  * [ ] идемпотентные методы
* [x] сервис статистики
* [x] сервис истории сделок
* [x] сервис уведомлений (сделки)
* [ ] тесты
  * [x] trading
  * [ ] negative cases (not enough money)
  * [x] notifications
  * [x] history
  * [x] stats
* [ ] мониторинг
* [ ] performance тесты

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

## Отладка Apache Kafka

Kafka can be accessed by consumers via port 9092 on the following DNS name from within your cluster:

    kafka.otus.svc.cluster.local

Each Kafka broker can be accessed by producers via port 9092 on the following DNS name(s) from within your cluster:

    kafka-0.kafka-headless.otus.svc.cluster.local:9092

To create a pod that you can use as a Kafka client run the following commands:

    kubectl run kafka-client --restart='Never' --image docker.io/bitnami/kafka:3.1.0-debian-10-r49 --namespace otus --command -- sleep infinity
    kubectl exec --tty -i kafka-client --namespace otus -- bash

    PRODUCER:
        kafka-console-producer.sh \
            
            --broker-list kafka-0.kafka-headless.otus.svc.cluster.local:9092 \
            --topic test

    CONSUMER:
        kafka-console-consumer.sh \
            
            --bootstrap-server kafka.otus.svc.cluster.local:9092 \
            --topic test \
            --from-beginning
