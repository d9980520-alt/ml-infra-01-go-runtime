package postgres

import (
	"context"
	"study/internal"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(p *pgxpool.Pool) *Repo {
	return &Repo{pool: p}
}

func (p *Repo) CreateTable(ctx context.Context) error {
	querycreate := `CREATE TABLE IF NOT EXISTS "Player"(
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	balance BIGINT NOT NULL DEFAULT 0,
	lastseen TIMESTAMPZ
);`
	_, err := p.pool.Exec(ctx, querycreate)
	return err
}

func (p *Repo) SavePlayer(ctx context.Context, bass internal.Player) error {
	querysave := `INSERT INTO "Player"(name,balance)VALUES($1,$2);`

	_, err := p.pool.Exec(ctx, querysave, bass.Name, bass.Balance)
	return err
}

func (p *Repo) GetPlayer(ctx context.Context, id int) (internal.Player, error) {
	var bass internal.Player

	queryget := `SELECT id,name,balance,lastseen FROM "Player" WHERE id=$1;`

	err := p.pool.QueryRow(ctx, queryget, id).Scan(&bass.Id, &bass.Name, &bass.Balance, &bass.LastSeen)
	return bass, err
}

func (p *Repo) ProcesPayment(ctx context.Context, id, amount int64) error {
	queryupadte := `UPDATE "Player" SET balance=balance -$1 WHERE id=$2;`
	_, err := p.pool.Exec(ctx, queryupadte, amount, id)
	return err
}
func (p *Repo) UpdateLastSeen(ctx context.Context, id int64) error {
	query := `UPADTE "Player" SET lastseen = NOW() WHERE id=$1;`
	_, err := p.pool.Exec(ctx, query, id)
	return err
}
