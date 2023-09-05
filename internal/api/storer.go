package api

import (
	"errors"

	"github.com/hashicorp/go-memdb"
)

var (
	ErrTalkNotFound          = errors.New("talk not found")
	ErrTalkAlreadyExists     = errors.New("talk already exists")
	ErrSpeakerNotFound       = errors.New("speaker  not found")
	ErrSpeakerAlreadyExists  = errors.New("speaker already exists")
	ErrWorkshopNotFound      = errors.New("workshop  not found")
	ErrWorkshopAlreadyExists = errors.New("workshop already exists")

	ErrEAStoreNotFound      = errors.New("EAStore  not found")
	ErrEAStoreAlreadyExists = errors.New("EAStore already exists")
)

type Storer struct {
	db *memdb.MemDB
}

func NewStorer() (*Storer, error) {
	db, err := memdb.NewMemDB(&memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"talk": {
				Name: "talk",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
			"speaker": {
				Name: "speaker",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
			"workshop": {
				Name: "workshop",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
			"eastore": {
				Name: "eastore",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &Storer{
		db: db,
	}, nil
}

func (s *Storer) GetTalk(id string) (Talk, error) {
	txn := s.db.Txn(false)
	ap, err := txn.First("talk", "id", id)
	if err != nil {
		return Talk{}, err
	}
	if ap == nil {
		return Talk{}, ErrTalkNotFound
	}
	return *ap.(*Talk), nil
}

func (s *Storer) CreateTalk(ap Talk) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("talk", "id", ap.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrTalkAlreadyExists
	}
	err = txn.Insert("talk", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateTalk(ap Talk) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("talk", "id", ap.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTalkNotFound
	}
	err = txn.Insert("talk", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteTalk(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("talk", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTalkNotFound
	}
	err = txn.Delete("talk", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) GetSpeaker(id string) (Speaker, error) {
	txn := s.db.Txn(false)
	ap, err := txn.First("speaker", "id", id)
	if err != nil {
		return Speaker{}, err
	}
	if ap == nil {
		return Speaker{}, ErrSpeakerNotFound
	}
	return *ap.(*Speaker), nil
}

func (s *Storer) CreateSpeaker(ap Speaker) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("speaker", "id", ap.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrSpeakerAlreadyExists
	}
	err = txn.Insert("speaker", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateSpeaker(ap Speaker) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("speaker", "id", ap.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSpeakerNotFound
	}
	err = txn.Insert("speaker", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteSpeaker(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("speaker", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSpeakerNotFound
	}
	err = txn.Delete("speaker", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) GetWorkshop(id string) (Workshop, error) {
	txn := s.db.Txn(false)
	ap, err := txn.First("workshop", "id", id)
	if err != nil {
		return Workshop{}, err
	}
	if ap == nil {
		return Workshop{}, ErrWorkshopNotFound
	}
	return *ap.(*Workshop), nil
}

func (s *Storer) CreateWorkshop(ap Workshop) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("workshop", "id", ap.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrWorkshopAlreadyExists
	}
	err = txn.Insert("workshop", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateWorkshop(ap Workshop) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("workshop", "id", ap.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrWorkshopNotFound
	}
	err = txn.Insert("workshop", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteWorkshop(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("workshop", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrWorkshopNotFound
	}
	err = txn.Delete("workshop", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) GetEAStore(id string) (EAStore, error) {
	txn := s.db.Txn(false)
	ap, err := txn.First("eastore", "id", id)
	if err != nil {
		return EAStore{}, err
	}
	if ap == nil {
		return EAStore{}, ErrEAStoreNotFound
	}
	return *ap.(*EAStore), nil
}

func (s *Storer) CreateEAStore(ap EAStore) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("eastore", "id", ap.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrEAStoreAlreadyExists
	}
	err = txn.Insert("eastore", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateEAStore(ap EAStore) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("eastore", "id", ap.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrEAStoreNotFound
	}
	err = txn.Insert("eastore", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteEAStore(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("eastore", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrEAStoreNotFound
	}
	err = txn.Delete("eastore", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}
