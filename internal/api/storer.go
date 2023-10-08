package api

import (
	"errors"

	"github.com/hashicorp/go-memdb"
)

var (
	ErrEAStoreNotFound         = errors.New("EAStore  not found")
	ErrEAStoreAlreadyExists    = errors.New("EAStore already exists")
	ErrEHSClusterNotFound      = errors.New("EHSCluster not found")
	ErrEHSClusterAlreadyExists = errors.New("EHSCluster already exists")
	ErrAWNotFound              = errors.New("AW not found")
	ErrAWAlreadyExists         = errors.New("AW already exists")
	ErrAVNotFound              = errors.New("AV not found")
	ErrAVAlreadyExists         = errors.New("AV already exists")
)

type Storer struct {
	db *memdb.MemDB
}

func NewStorer() (*Storer, error) {
	db, err := memdb.NewMemDB(&memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
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
			"ehscluster": {
				Name: "ehscluster",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
			"aw": {
				Name: "aw",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
			"av": {
				Name: "av",
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

func (s *Storer) GetEHSCluster(id string) (EHSCluster, error) {
	txn := s.db.Txn(false)
	ap, err := txn.First("ehscluster", "id", id)
	if err != nil {
		return EHSCluster{}, err
	}
	if ap == nil {
		return EHSCluster{}, ErrEHSClusterNotFound
	}
	return *ap.(*EHSCluster), nil
}

func (s *Storer) CreateEHSCluster(ap EHSCluster) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("ehscluster", "id", ap.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrEHSClusterAlreadyExists
	}
	err = txn.Insert("ehscluster", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateEHSCluster(ap EHSCluster) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("ehscluster", "id", ap.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrEHSClusterNotFound
	}
	err = txn.Insert("ehscluster", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteEHSCluster(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("ehscluster", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrEHSClusterNotFound
	}
	err = txn.Delete("ehscluster", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) GetAW(id string) (AW, error) {
	txn := s.db.Txn(false)
	ap, err := txn.First("aw", "id", id)
	if err != nil {
		return AW{}, err
	}
	if ap == nil {
		return AW{}, ErrAWNotFound
	}
	return *ap.(*AW), nil
}

func (s *Storer) CreateAW(ap AW) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("aw", "id", ap.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrAWAlreadyExists
	}
	err = txn.Insert("aw", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateAW(ap AW) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("aw", "id", ap.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAWNotFound
	}
	err = txn.Insert("aw", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteAW(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("aw", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAWNotFound
	}
	err = txn.Delete("aw", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) GetAV(id string) (AV, error) {
	txn := s.db.Txn(false)
	ap, err := txn.First("av", "id", id)
	if err != nil {
		return AV{}, err
	}
	if ap == nil {
		return AV{}, ErrAWNotFound
	}
	return *ap.(*AV), nil
}

func (s *Storer) CreateAV(ap AV) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("av", "id", ap.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrAWAlreadyExists
	}
	err = txn.Insert("av", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateAV(ap AV) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("av", "id", ap.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAWNotFound
	}
	err = txn.Insert("av", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteAV(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("av", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAWNotFound
	}
	err = txn.Delete("av", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}
