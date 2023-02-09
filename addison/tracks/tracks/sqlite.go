package tracks

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	trackserrors "tracks/errors"

	_ "github.com/mattn/go-sqlite3"
)

const (
	errUniqueConstraintFailedString = "UNIQUE constraint failed: track.id"
	createDbQuery                   = `CREATE TABLE IF NOT EXISTS track (id TEXT PRIMARY KEY, audio TEXT)`
	readTrackQuery                  = `SELECT id, audio FROM track WHERE id = ?`
	listTrackQuery                  = `SELECT id, audio FROM track`
	createTrackQuery                = `INSERT INTO track (id, audio) VALUES (?, ?)`
	updateTrackQuery                = `UPDATE track SET audio = ? WHERE id= ?`
	deleteTrackQuery                = `DELETE FROM track WHERE id = ?`
)

type TrackStore struct {
	db *sql.DB
}

func NewTrackStore(dbName string) (*TrackStore, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createDbQuery)
	if err != nil {
		return nil, err
	}

	return &TrackStore{
		db: db,
	}, nil
}

func (t *TrackStore) Read(id string) (*Track, error) {
	track := &Track{}
	row := t.db.QueryRow(readTrackQuery, id)
	err := row.Scan(&track.Id, &track.Audio)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &trackserrors.ErrTrackNotFound{}
		}
		return nil, err
	}

	return track, nil
}

func (t *TrackStore) List() ([]*Track, error) {
	rows, err := t.db.Query(listTrackQuery)
	if err != nil {
		return nil, err
	}

	allTracks := make([]*Track, 0)
	for rows.Next() {
		track := &Track{}
		err := rows.Scan(&track.Id, &track.Audio)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, &trackserrors.ErrTrackNotFound{}
			}
			return nil, err
		}

		allTracks = append(allTracks, track)
	}

	return allTracks, nil
}

func (t *TrackStore) Create(id, audio string) (err error) {
	tx, err := t.db.BeginTx(context.Background(), nil /* sql.TxOptions */)
	if err != nil {
		return err
	}

	// Try to insert a new record, if there is already an entry with the same
	// Id, then we update the current entry's audio and return
	// TrackAlreadyExists error.
	_, err = tx.Exec(createTrackQuery, id, audio)
	if err != nil {
		if !strings.Contains(err.Error(), errUniqueConstraintFailedString) {
			return err
		}
		_, err = tx.Exec(updateTrackQuery, audio, id)
		if err != nil {
			return err
		}
		if err = tx.Commit(); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
		}
		return &trackserrors.ErrTrackAlreadyExists{}
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}

	return nil
}

func (t *TrackStore) Delete(id string) error {
	tx, err := t.db.BeginTx(context.Background(), nil /* sql.TxOptions */)
	if err != nil {
		return err
	}

	result, err := tx.Exec(deleteTrackQuery, id)
	if err != nil {
		return err
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numRows != 1 {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return &trackserrors.ErrTrackNotFound{}
	}

	if err = tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}

	return nil
}
