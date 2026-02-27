package repo_internal

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	polls_models "oprosdom.ru/core/internal/polls/models"
)

// обязательно юзаем pgxpool так как нам нужен именно pool для веб-сервера, иначе будет очередь и slow

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, conn string) (*Postgres, error) {

	// таймаут 30 сек если ctx не придет быстрее
	ctxTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, errors.New("failed to parse db string with config")
	}

	pool, err := pgxpool.NewWithConfig(ctxTimeout, config)
	if err != nil {
		return nil, errors.New("failed to create pool")
	}

	// нужно проверить реальную связь с пулом

	if err := pool.Ping(ctxTimeout); err != nil {
		return nil, errors.New("pool ping failed")
	}

	return &Postgres{pool: pool}, nil

}

func (p *Postgres) Close() {
	p.pool.Close()
}

// func (p *Postgres) PhoneSend(ctx context.Context, v *polls_models.ValidatedPhoneSendReq) error {
// 	// обращаем внимание на кавычки!!! они в sql специфические
// 	const query = `INSERT INTO phonesend (phone, phone_type, useragent, ip) VALUES ($1, $2, $3, $4) RETURNING time`

// 	var time int

// 	// брать соединение из пула и возвращать - не нужно, QueryRow делает это сама. Если несколько операций или транзакция, то вручную тогда берем и возвращаем
// 	// QueryRow используется только если ожидается 0 или 1 строка!
// 	err := p.pool.QueryRow(ctx, query, v.Phone, v.PhoneType, v.UserAgent, v.IP).Scan(&time)
// 	if err != nil {
// 		return errors.New("insert kvartira failed")
// 	}

// 	log.Printf("phonesend inserted phone %v on %v", v.Phone, time)

// 	return nil
// }

func (p *Postgres) GetPolls(ctx context.Context) ([]*polls_models.Poll, error) {

	rows, err := p.pool.Query(ctx, `SELECT id, title FROM polls`)
	if err != nil {
		return nil, errors.New("error while get polls")
	}
	defer rows.Close() // не путаем с pool.Close() - это тут другое. Здесь освобождаем ресурсы чтоб возврат подключения произошел в пул.

	polls, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*polls_models.Poll, error) {
		var k polls_models.Poll
		if err := row.Scan(&k.Id, &k.Title); err != nil {
			return nil, errors.New("error scanning poll row")
		}
		return &k, nil
	})
	if err != nil {
		return nil, errors.New("error collect poll rows")
	}

	// ниже мы делаем проверку потому, что данные из базы идут порциями и если не будет этой проверки, то может прийти 100 из 1000 строк и мы не узнаем об этом думая что все строки получены
	if err := rows.Err(); err != nil {
		return nil, errors.New("not all rows was received from polls")
	}

	return polls, nil
}

func (p *Postgres) Vote(ctx context.Context, m *polls_models.ValidVoteReq) error {
	// обращаем внимание на кавычки!!! они в sql специфические
	// если пара poll_id и jti будут в таблице, то сущ vote заменится на новый
	// const - гарантия неизменяемости
	const query = `
        INSERT INTO poll_votes (poll_id, jti, vote)
        VALUES ($1, $2, $3)
        ON CONFLICT (poll_id, jti) 
        DO UPDATE SET vote = EXCLUDED.vote
    `

	// брать соединение из пула и возвращать - не нужно, Exec тоже делает это сама. Если несколько операций или транзакция, то вручную тогда берем и возвращаем
	// QueryRow используется только если ожидается 0 или 1 строка!
	// здесь юзаем Exec так как нам Returning Id не нужен
	if _, err := p.pool.Exec(ctx, query, m.PollID, m.Jti, m.Vote); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) PollStats(ctx context.Context) ([]*polls_models.PollStats, error) {

	// const - гарантия неизменяемости
	const query = `
		SELECT 
			poll_id,
			COUNT(*) FILTER (WHERE vote = 'za') AS za_count,
			COUNT(*) FILTER (WHERE vote = 'protiv') AS protiv_count
		FROM poll_votes
		GROUP BY poll_id
		ORDER BY poll_id
	`

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, errors.New("error while get pollstats")
	}
	defer rows.Close() // не путаем с pool.Close() - это тут другое. Здесь освобождаем ресурсы чтоб возврат подключения произошел в пул.

	pollStats, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*polls_models.PollStats, error) {
		p := &polls_models.PollStats{}
		// обязательно пишем по указателям
		if err := row.Scan(&p.PollID, &p.ZaCount, &p.ProtivCount); err != nil {
			return nil, errors.New("error scanning pollstats row")
		}
		return p, nil
	})

	if err != nil {
		return nil, errors.New("error collectrows pollstats")
	}

	// ниже мы делаем проверку потому, что данные из базы идут порциями и если не будет этой проверки, то может прийти 100 из 1000 строк и мы не узнаем об этом думая что все строки получены
	if err := rows.Err(); err != nil {
		return nil, errors.New("not all rows was received from pollstats")
	}

	return pollStats, nil

}
