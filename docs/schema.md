# Data Schema

```mermaid
erDiagram
    Account {
        int id PK
        string number
        string name
    }

    Transaction {
        int id PK
        decimal amount
        string description
        string transaction_type
        date transaction_date
        string account_id FK
    }

    Account ||--o{ Transaction : "has"
