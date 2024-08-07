package qimage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"mime/multipart"
)

type Images []*Imager

func (is *Images) ReadFromMultipart(fileHeader *multipart.FileHeader, allowedFileTypes []string) error {
	for _, image := range *is {
		err := ReadFromMultipart(fileHeader, image, allowedFileTypes)
		if err != nil {
			return err
		}
	}
	return nil
}
func (is *Images) IsEmpty() bool {
	return is == nil || len(*is) == 0
}
func (is *Images) Insert(ctx context.Context, dbPool *pgxpool.Pool, sql string) error {
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Images.Insert error: " + err.Error())
		return err
	}

	for _, image := range *(is) {
		err = InsertInTx(ctx, tx, image, sql)
		if err != nil {
			log.Error().Err(err).Msg("Images.Insert error: " + err.Error())
			e := tx.Rollback(ctx)
			if e != nil {
				log.Error().Err(err).Msg("Images.Insert error: " + e.Error())
				return e
			}
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Images.Insert error: " + err.Error())
		e := tx.Rollback(ctx)
		if e != nil {
			log.Error().Err(err).Msg("Images.Insert error: " + e.Error())
			return e
		}
		return err
	}

	return nil
}
func (is *Images) Read(ctx context.Context, dbPool *pgxpool.Pool, sql string) error {
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("images Read error")
		return err
	}

	for _, image := range *is {
		err = ReadInTx(ctx, tx, image, sql)
		if err != nil {
			log.Error().Err(err).Msg("images Read error")
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("images Read error")
		e := tx.Rollback(ctx)
		if e != nil {
			log.Error().Err(err).Msg("images Read error")
			return e
		}
		return err
	}

	return nil
}
func (is *Images) DeleteLo(ctx context.Context, dbPool *pgxpool.Pool) error {
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("images DeleteLo error")
		return err
	}

	for _, image := range *is {
		err = DeleteInTx(ctx, tx, image)
		if err != nil {
			log.Error().Err(err).Msg("images DeleteLo error")
			e := tx.Rollback(ctx)
			if e != nil {
				log.Error().Err(err).Msg("images DeleteLo error")
				return e
			}
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("images DeleteLo error")
		e := tx.Rollback(ctx)
		if e != nil {
			log.Error().Err(err).Msg("images DeleteLo error")
			return e
		}
		return err
	}

	return nil
}
