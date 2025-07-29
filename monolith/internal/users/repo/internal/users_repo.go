package users_repo_internal

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"oprosdom.ru/monolith/internal/users/models"
)

// обязательно юзаем pgxpool так как нам нужен именно pool для веб-сервера, иначе будет очередь и slow

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, conn string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, errors.New("failed to parse db string with config")
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.New("failed to create pool")
	}

	// нужно проверить реальную связь с рулом, но может быть timeout поэтому подстраховываемся контекстом
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctxTimeout); err != nil {
		return nil, errors.New("pool ping failed")
	}

	return &Postgres{pool: pool}, nil

}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) KvartiraAdd(ctx context.Context, k *models.Kvartira) error {
	// обращаем внимание на кавычки!!! они в sql специфические
	const query = `INSERT INTO kvartiras (number, komnat) VALUES ($1, $2) RETURNING uuid::text`

	// брать соединение из пула и возвращать - не нужно, QueryRow делает это сама. Если несколько операций или транзакция, то вручную тогда берем и возвращаем
	// QueryRow используется только если ожидается 0 или 1 строка!
	err := p.pool.QueryRow(ctx, query, k.Number, k.Komnat).Scan(&k.Id)
	if err != nil {
		return errors.New("insert kvartira failed")
	}

	return nil
}

func (p *Postgres) KvartiraGetById(ctx context.Context, id string) (*models.Kvartira, error) {
	// обращаем внимание на кавычки!!! они в sql специфические
	const query = `SELECT uuid::text, number, komnat FROM kvartiras WHERE uuid=$1`

	var k models.Kvartira

	// брать соединение из пула и возвращать - не нужно, QueryRow делает это сама. Если несколько операций или транзакция, то вручную тогда берем и возвращаем
	// QueryRow используется только если ожидается 0 или 1 строка!
	err := p.pool.QueryRow(ctx, query, id).Scan(&k.Id, &k.Number, &k.Komnat)
	if err != nil {
		return nil, errors.New("kvartira search by id failed")
	}

	return &k, nil
}

func (p *Postgres) KvartiraUpdate(ctx context.Context, k *models.Kvartira) error {

	// обращаем внимание на кавычки!!! они в sql специфические
	const query = `UPDATE kvartiras SET number=$1, komnat=$2 WHERE uuid=$3 RETURNING uuid::text` // если записи не будет то вернется ErrNoRows

	var updatedId string
	// брать соединение из пула и возвращать - не нужно, QueryRow делает это сама. Если несколько операций или транзакция, то вручную тогда берем и возвращаем
	// QueryRow используется только если ожидается 0 или 1 строка!
	err := p.pool.QueryRow(ctx, query, k.Number, k.Komnat, k.Id).Scan(&updatedId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // юзаем pgx вместо pgxpool, потому что в pgxpool не находит ErrNoRows. В документации pgx написано так надо делать при работе с pgxpool
			return errors.New("kvartira id is not exists")
		}
		return errors.New("kvartira update failed")
	}

	return nil

}

func (p *Postgres) KvartirasGet(ctx context.Context) ([]*models.Kvartira, error) {

	rows, err := p.pool.Query(ctx, `SELECT uuid::text, number, komnat FROM kvartiras`)
	if err != nil {
		return nil, errors.New("error while get kvartiras")
	}
	defer rows.Close() // не путаем с pool.Close() - это тут другое. Здесь освобождаем ресурсы чтоб возврат подключения произошел в пул.

	kvartiras, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.Kvartira, error) {
		var k models.Kvartira
		if err := row.Scan(&k.Id, &k.Number, &k.Komnat); err != nil {
			return nil, errors.New("error scanning kvartira row")
		}
		return &k, nil
	})
	if err != nil {
		return nil, errors.New("error collectrows kvartiras")
	}

	// ниже мы делаем проверку потому, что данные из базы идут порциями и если не будет этой проверки, то может прийти 100 из 1000 строк и мы не узнаем об этом думая что все строки получены
	if err := rows.Err(); err != nil {
		return nil, errors.New("not all rows was received from kvartiras")
	}

	return kvartiras, nil
}

func (p *Postgres) MemberAdd(ctx context.Context, m *models.Member) error {
	// обращаем внимание на кавычки!!! они в sql специфические
	const query = `INSERT INTO members (name, phone, community) VALUES ($1, $2, $3) RETURNING uuid::text`

	// брать соединение из пула и возвращать - не нужно, QueryRow делает это сама. Если несколько операций или транзакция, то вручную тогда берем и возвращаем
	// QueryRow используется только если ожидается 0 или 1 строка!
	err := p.pool.QueryRow(ctx, query, m.Name, m.Phone, m.Community).Scan(&m.Id)
	if err != nil {
		return errors.New("insert member failed")
	}

	return nil
}

func (p *Postgres) MemberGetById(ctx context.Context, id string) (*models.Member, error) {
	// обращаем внимание на кавычки!!! они в sql специфические
	const query = `SELECT uuid::text, name, phone, community FROM members WHERE uuid=$1`

	var m models.Member

	// брать соединение из пула и возвращать - не нужно, QueryRow делает это сама. Если несколько операций или транзакция, то вручную тогда берем и возвращаем
	// QueryRow используется только если ожидается 0 или 1 строка!
	err := p.pool.QueryRow(ctx, query, id).Scan(&m.Id, &m.Name, &m.Phone, &m.Community)
	if err != nil {
		return nil, errors.New("member search by id failed")
	}

	return &m, nil
}

func (p *Postgres) MemberUpdate(ctx context.Context, m *models.Member) error {

	// обращаем внимание на кавычки!!! они в sql специфические
	const query = `UPDATE members SET name=$1, phone=$2, community=$3 WHERE uuid=$4 RETURNING uuid::text` // если записи не будет то вернется ErrNoRows

	var updatedId string
	// брать соединение из пула и возвращать - не нужно, QueryRow делает это сама. Если несколько операций или транзакция, то вручную тогда берем и возвращаем
	// QueryRow используется только если ожидается 0 или 1 строка!
	err := p.pool.QueryRow(ctx, query, m.Name, m.Phone, m.Community, m.Id).Scan(&updatedId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // юзаем pgx вместо pgxpool, потому что в pgxpool не находит ErrNoRows. В документации pgx написано так надо делать при работе с pgxpool
			return errors.New("member id is not exists")
		}
		return errors.New("member update failed")
	}

	return nil

}

func (p *Postgres) MembersGet(ctx context.Context) ([]*models.Member, error) {

	rows, err := p.pool.Query(ctx, `SELECT uuid::text, name, phone, community FROM members`)
	if err != nil {
		return nil, errors.New("error while get members")
	}
	defer rows.Close() // не путаем с pool.Close() - это тут другое. Здесь освобождаем ресурсы чтоб возврат подключения произошел в пул.

	members, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*models.Member, error) {
		var m models.Member
		if err := row.Scan(&m.Id, &m.Name, &m.Phone, &m.Community); err != nil {
			return nil, errors.New("error scanning member row")
		}
		return &m, nil
	})
	if err != nil {
		return nil, errors.New("error collectrows members")
	}

	// ниже мы делаем проверку потому, что данные из базы идут порциями и если не будет этой проверки, то может прийти 100 из 1000 строк и мы не узнаем об этом думая что все строки получены
	if err := rows.Err(); err != nil {
		return nil, errors.New("not all rows was received from members")
	}

	return members, nil
}

func (p *Postgres) DeleteById(ctx context.Context, id string, mk string) error {

	var query string

	switch mk {
	case "kvartira":
		// обращаем внимание на кавычки!!! они в sql специфические
		query = `DELETE FROM kvartiras WHERE uuid=$1 RETURNING uuid::text` // если записи не будет то вернется ErrNoRows
	case "member":
		// обращаем внимание на кавычки!!! они в sql специфические
		query = `DELETE FROM members WHERE uuid=$1 RETURNING uuid::text` // если записи не будет то вернется ErrNoRows
	}

	var deletedId string
	// брать соединение из пула и возвращать - не нужно, QueryRow делает это сама. Если несколько операций или транзакция, то вручную тогда берем и возвращаем
	// QueryRow используется только если ожидается 0 или 1 строка!
	err := p.pool.QueryRow(ctx, query, id).Scan(&deletedId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // юзаем pgx вместо pgxpool, потому что в pgxpool не находит ErrNoRows. В документации pgx написано так надо делать при работе с pgxpool
			return errors.New("id is not exists")
		}
		return errors.New("deletion failed")
	}

	return nil

}

func (p *Postgres) PayDebt(ctx context.Context, r *models.PayDebtRequest) (*models.PayDebtResponse, error) {

	var response models.PayDebtResponse

	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return nil, errors.New("begin transaction failed")
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`WITH payment AS (
            UPDATE members 
            SET balance = balance - $1::DECIMAL(15,2)
            WHERE uuid = $2::UUID AND balance >= $1::DECIMAL(15,2)
            RETURNING balance
        ),
        debt_payment AS (
            UPDATE kvartiras
            SET debt = debt - $1::DECIMAL(15,2)
            WHERE uuid = $3::UUID 
            RETURNING debt
        )
        SELECT 
            (SELECT balance FROM payment),
            (SELECT debt FROM debt_payment),
            gen_random_uuid()`,
		r.Amount, r.MemberId, r.KvartiraId,
	).Scan(&response.NewBalance, &response.NewDebt, &response.PaymentId)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO payments 
         (member_uuid, kvartira_uuid, amount) 
         VALUES ($1, $2, $3)`,
		r.MemberId, r.KvartiraId, r.Amount,
	)
	if err != nil {
		return nil, errors.New("failed to save into payments")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.New("commit failed")
	}

	return &response, nil
}
