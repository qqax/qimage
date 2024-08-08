package qimage

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

func Insert(ctx context.Context, dbPool *pgxpool.Pool, i *Imager, sql string) error {
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Insert Image error")
		return err
	}

	err = InsertInTx(ctx, tx, i, sql)
	if err != nil {
		log.Error().Err(err).Msg("Insert Image error")
		e := tx.Rollback(ctx)
		if e != nil {
			log.Error().Err(err).Msg("Insert Image error")
			return e
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Insert Image error")
		e := tx.Rollback(ctx)
		if e != nil {
			log.Error().Err(err).Msg("Insert Image error")
			return e
		}
		return err
	}

	return nil
}
func Delete(ctx context.Context, dbPool *pgxpool.Pool, i *Imager) error {
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Delete Image error: " + err.Error())
		return err
	}

	err = DeleteInTx(ctx, tx, i)
	if err != nil {
		log.Error().Err(err).Msg("Delete Image error: " + err.Error())
		e := tx.Rollback(ctx)
		if e != nil {
			log.Error().Err(err).Msg("Delete Image error: " + e.Error())
			return e
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Delete Image error: " + err.Error())
		e := tx.Rollback(ctx)
		if e != nil {
			log.Error().Err(err).Msg("Delete Image error: " + e.Error())
			return e
		}
		return err
	}

	return nil
}
