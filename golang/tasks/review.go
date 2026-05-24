package payment

import (
    "context"
    "database/sql"
    "errors"
    "sync"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
    mu       sync.Mutex
    producer *kafka.Producer
    db       *sql.DB
)

func ProcessPayment(userID int, amount float64) error {
    mu.Lock()
    defer mu.Unlock()

    tx, err := db.Begin()
    if err != nil {
        return errors.New("internal error")
    }

    var balance float64
    row := tx.QueryRow("SELECT balance FROM users WHERE id = $1", userID)
    err = row.Scan(&balance)
    if err != nil {
        tx.Rollback()
        return errors.New("internal error")
    }
    if balance < amount {
        tx.Rollback()
        return errors.New("insufficient funds")
    }

    _, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", amount, userID)
    if err != nil {
        tx.Rollback()
        return errors.New("internal error")
    }

    // Отправляем событие в Kafka, чтобы другие сервисы узнали о платеже
    topic := "payments"
    err = producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          []byte("payment_processed"),
    }, nil)
    if err != nil {
        tx.Rollback()
        return errors.New("internal error")
    }

    return tx.Commit()
}






package payment

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
)

type Repository struct {
    db *sql.DB
}

var ErrInsufficientFunds = errors.New("insufficient funds")

func (r *Repository) ProcessPayment(ctx context.Context, userID int, amount float64) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback()

    // Блокируем строку пользователя для обновления
    var balance float64
    err = tx.QueryRowContext(ctx,
        "SELECT balance FROM users WHERE id = $1 FOR UPDATE", userID,
    ).Scan(&balance)
    if err != nil {
        return fmt.Errorf("get balance: %w", err)
    }

    if balance < amount {
        return ErrInsufficientFunds
    }

    _, err = tx.ExecContext(ctx,
        "UPDATE users SET balance = balance - $1 WHERE id = $2", amount, userID,
    )
    if err != nil {
        return fmt.Errorf("update balance: %w", err)
    }

    // Вместо отправки в Kafka пишем в таблицу outbox внутри той же транзакции
    _, err = tx.ExecContext(ctx,
        "INSERT INTO outbox (topic, payload) VALUES ($1, $2)", "payments", `{"user_id": `+fmt.Sprint(userID)+`}`,
    )
    if err != nil {
        return fmt.Errorf("insert outbox: %w", err)
    }

    return tx.Commit()
}