package usecase

import (
	"context"
	"errors"
	"study/internal"
)

type RepoMethod interface {
	CreateTable(ctx context.Context) error
	SavePlayer(ctx context.Context, bass internal.Player) error
	GetPlayer(ctx context.Context, id int) (internal.Player, error)
	ProcesPayment(ctx context.Context, id, amount int64) error
	UpdateLastSeen(ctx context.Context, id int64) error
}

type Repocontr struct {
	repo RepoMethod
}

func NewRepoContr(uc RepoMethod) *Repocontr {
	return &Repocontr{repo: uc}
}

func CheckName(name string) error {
	if name == "" {
		return errors.New("имя не можеть быть пустым")
	}
	return nil
}
func CheckBalance(balance int64) error {
	if balance < 0 {
		return errors.New("баланс не может быть меньше нуля!")
	}
	return nil
}
func CheckAmount(amount int64) error {
	if amount < 0 {
		return errors.New("перевод не может быть на сумму меньше нуля")
	}
	return nil
}
func (u *Repocontr) CreatePlayer(ctx context.Context, name string, balance int64) error {
	errs := make(chan error, 2)

	go func() { errs <- CheckName(name) }()
	go func() { errs <- CheckBalance(balance) }()

	for i := 0; i < 2; i++ {
		if err := <-errs; err != nil {
			return err
		}
	}
	return u.repo.SavePlayer(ctx, internal.Player{
		Name:    name,
		Balance: balance,
	})
}
