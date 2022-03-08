# Проектирование системы сделок

## Регистрация пользователя

```mermaid
sequenceDiagram
    actor Client
    
    Client->>API Gateway: POST /api/v1/identity/register
    API Gateway->>Identity Service: POST /api/v1/register
    Identity Service->>Identity Service: create user
    Identity Service->>Message broker: publish UserCreated
    
    par Complete registration
        Identity Service-->>API Gateway: 201 Created (user)
        API Gateway-->>Client: 201 Created (user)
    and Create billing account
        Message broker-->>Billing Service: consume UserCreated
        Billing Service->>Billing Service: create billing account
    and Create notification profile
        Message broker-->>Notification Service: consume UserCreated
        Notification Service->>Notification Service: create notification profile
    end
```

## Диаграмма состояний заявки на покупку

```mermaid
stateDiagram-v2
    [*] --> Pending
    Pending --> Canceled: cancel by user
    Pending --> PaymentPending: sell order found
    PaymentPending --> PaymentSucceded
    PaymentPending --> PaymentFailed: not enough money
    PaymentSucceded --> Approved: after money transit
    Approved --> [*]
    Canceled --> [*]
    PaymentFailed --> [*]
```

## Диаграмма состояний заявки на продажу

```mermaid
stateDiagram-v2
    [*] --> Pending
    Pending --> Canceled: cancel by user
    Pending --> DealPending: purshase order found
    DealPending --> Pending: purshase order canceled
    DealPending --> AccrualPending: wait for money transit
    AccrualPending --> Approved
    AccrualPending --> Pending
    Approved --> [*]
    Canceled --> [*]
```

## Осуществление сделки

```mermaid
sequenceDiagram
    actor Seller
    actor Purchaser
    
    Seller->>API Gateway: POST /api/v1/trading/sell-orders
    API Gateway->>Trading Service: POST /api/v1/sell-orders
    Trading Service->>Trading Service: create sell order with status "pending"
    Trading Service-->>API Gateway: 202 Accepted
    API Gateway-->>Seller: 202 Accepted
    
    Purchaser->>API Gateway: POST /api/v1/trading/purchase-orders
    API Gateway->>Trading Service: POST /api/v1/purchase-orders
    Trading Service->>Trading Service: find appropriate sell order
    Trading Service->>Trading Service: update sell order with status "deal pending"
    Trading Service->>Trading Service: create purchase order with status "payment pending"
    Trading Service-->>API Gateway: 202 Accepted
    API Gateway-->>Purchaser: 202 Accepted
    Trading Service->>Message broker: publish CreatePayment
    Message broker-->>Billing Service: consume CreatePayment
    
    alt Deal succeded
        Billing Service->>Message broker: publish PaymentSucceded
        Message broker-->>Trading Service: consume PaymentSucceded
        Trading Service->>Trading Service: update purchase order with status "payment succeded"
        Trading Service->>Trading Service: update sell order with status "accrual pending"
        Trading Service->>Message broker: publish CreateAccrual
        
        Message broker->>Billing Service: consume CreateAccrual
        Billing Service-->>Message broker: publish AccrualApproved
        
        Trading Service->>Trading Service: update purchase order with status "approved"
        Trading Service->>Trading Service: update sell order with status "approved"
        
        Trading Service->>Message broker: publish DealSucceeded
        par Deal notification
            Message broker-->>Notification Service: consume DealSucceeded
            Notification Service->>Purchaser: Purchase order succeded email
            Notification Service->>Purchaser: Sell order succeded email
        and Deal push event
            Message broker-->>Push Service: consume DealSucceeded
            Push Service-->>Purchaser: HTTP/2 push event (purchase order succeded)
            Push Service-->>Seller: HTTP/2 push event (sell order succeded)
        and Statistics updated
            Message broker-->>Statistics Service: consume DealSucceeded
            Statistics Service->>Statistics Service: update stats
        end
    else Deal failed
        Billing Service->>Message broker: publish PaymentDeclined
        Message broker-->>Trading Service: consume PaymentDeclined
        Trading Service->>Trading Service: update purchase order with status "payment failed"
        Trading Service->>Trading Service: update sell order with status "pending"
        Trading Service->>Message broker: publish PurchaseFailed
        par Deal notification
            Message broker-->>Notification Service: consume PurchaseFailed
            Notification Service->>Purchaser: Purchase failed email
        and Deal push event
            Message broker-->>Push Service: consume PurchaseFailed
            Push Service-->>Purchaser: HTTP/2 push event (purchase failed)
        end
    end
```

## Структура сущностей

### Сущности торговой площадки

* виртуальный предмет (item)
  * id
  * name
  * initial price (стоимость при размещении)
  * commission (комиссия за сделку)
* заявка на продажу (sell order)
  * id
  * user id (nullable)
    * null в случае первичного размещения
  * item id
  * price - желаемая цена продажи
  * status
* заявка на покупку (purchase order)
  * id
  * user id
  * item id
  * price - желаемая цена покупки
  * status

### Сущности статистики торговой площадки

* сделки за день (daily stat)
  * date
  * item id
  * deals count - количество сделок за день
  * deals amount - общая сумма сделок за день
* топ сделок (top deals)
  * item id
  * deals count - количество сделок за день
  * deals amount - общая сумма сделок за день

### Сущности биллинга

* аккаунт (account)
  * id (= user id)
  * amount
* история баланса (balance history)
  * id
  * account id
  * type
    * deposit
    * withdraw
    * payment
    * accrual
  * amount 
