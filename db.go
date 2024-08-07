package qimage

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func ReadInTx(ctx context.Context, tx pgx.Tx, i *Imager, sql string) error {
	var (
		raster      uint32
		size        int
		largeObject *pgx.LargeObject
	)

	los := tx.LargeObjects()

	err := tx.QueryRow(ctx, sql).
		Scan(&raster, &size)
	if err != nil {
		e := tx.Rollback(ctx)
		if e != nil {
			return e
		}
		return err
	}

	largeObject, err = los.Open(ctx, raster, 0x40000)
	if err != nil {
		e := tx.Rollback(ctx)
		if e != nil {
			return e
		}
		return err
	}

	raw := make([]byte, size)
	_, err = largeObject.Read(raw)
	if err != nil {
		e := tx.Rollback(ctx)
		if e != nil {
			return e
		}
		return err
	}

	err = largeObject.Close()
	if err != nil {
		e := tx.Rollback(ctx)
		if e != nil {
			return e
		}
		return err
	}

	(*i).SetRaw(raw)

	return nil
}
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
func InsertInTx(ctx context.Context, tx pgx.Tx, i *Imager, sql string) error {
	var (
		id          uint32
		oid         uint32
		name        string
		largeObject *pgx.LargeObject
	)

	los := tx.LargeObjects()

	oid, err := los.Create(ctx, 0)
	if err != nil {
		return err
	}

	largeObject, err = los.Open(ctx, oid, 0x20000)
	if err != nil {
		return err
	}

	size, err := largeObject.Write((*i).GetRaw())
	if err != nil {
		return err
	}

	err = largeObject.Close()
	if err != nil {
		return err
	}

	err = tx.QueryRow(ctx, sql, oid, size).
		Scan(&id, &name)
	if err != nil {
		return err
	}

	(*i).SetID(id)
	(*i).SetName(name)
	(*i).SetOID(oid)
	(*i).SetSize(size)

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
func DeleteInTx(ctx context.Context, tx pgx.Tx, i *Imager) error {
	los := tx.LargeObjects()

	err := los.Unlink(ctx, (*i).GetOID())
	if err != nil {
		return err
	}

	return nil
}
