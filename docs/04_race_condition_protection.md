# Защита от гонок

## Идемпотентность операций биллинга

```mermaid
sequenceDiagram
    actor Client
    
    Client->>Trading service: GET /api/v1/billing/account
    Trading service->>PostgreSQL: get account
    PostgreSQL-->>Trading service: account with updated at
    Trading service-->>Client: 200 OK
    Note over Trading service, Client: ETag header with hash = f(account.updatedAt)
    Client->>Trading service: POST /api/v1/billing/withdraw
    Note over Trading service, Client: With If-Match header with hash
    Trading service->>Redis: obtain lock
    alt Lock failed
        Redis-->>Trading service: already locked
        Trading service-->>Client: 409 Conflict
    else Lock succeeded
        Redis-->>Trading service: lock obtained
        Trading service->>PostgreSQL: get account
        PostgreSQL-->>Trading service: account
        Trading service->>Trading service: verify hash
        alt Hash mismatch
            Trading service-->>Client: 412 Precondition failed
        else Hash match
            Trading service->>PostgreSQL: withdraw money
            PostgreSQL-->>Trading service: ok
            Trading service->>Redis: free lock
            Redis-->>Trading service: ok
            Trading service-->>Client: 202 Accepted
        end
    end
```

## Идемпотентность операций торговой площадки

```mermaid
sequenceDiagram
    actor Client
    
    Client->>Trading service: GET /api/v1/trading/purchase-orders
    Trading service->>Redis: get timestamp (from last purchase)
    Redis-->>Trading service: timestamp (UUID)
    Trading service-->>Client: 200 OK
    Note over Trading service, Client: ETag header with timestamp
    Client->>Trading service: POST /api/v1/trading/purchase-orders
    Note over Trading service, Client: With If-Match header with timestamp
    Trading service->>Redis: obtain lock
    alt Lock failed
        Redis-->>Trading service: already locked
        Trading service-->>Client: 409 Conflict
    else Lock succeeded
        Redis-->>Trading service: lock obtained
        Trading service->>Redis: get timestamp (from last purchase)
        Redis-->>Trading service: timestamp
        Trading service->>Trading service: verify timestamp
        alt Timestamp mismatch
            Trading service-->>Client: 412 Precondition failed
        else Timestamp match
            Trading service->>Trading service: create order
            Trading service->>Redis: update timestamp
            Redis-->>Trading service: ok
            Trading service->>Redis: free lock
            Redis-->>Trading service: ok
            Trading service-->>Client: 202 Accepted
        end
    end
```

## Осуществление сделок

Для подбора наиболее подходящей сделки из "стакана" (очереди заявок на продажу/покупку) используется блокировка строки
`FOR UPDATE` с пропуском заблокированных строк `SKIP LOCKED`. Таким образом исключается возможность
подбора одной сделки разными потоками. При этом потоки не блокируют друг друга и конкуретный
поток подберет следующую незаблокированную запись.

```postgresql
SELECT id, user_id, item_id, payment_id, deal_id, price, commission, status, created_at, updated_at
FROM "purchase_order"
WHERE item_id = $1 AND price <= $2 AND user_id != $3
ORDER BY price DESC
LIMIT 1 
FOR UPDATE SKIP LOCKED;
```
