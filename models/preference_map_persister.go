package models

import (
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// MapSessionPersister assists with persisting session in a badger store
type MapSessionPersister struct {
	db *sync.Map
}

// NewMapSessionPersister creates a new MapSessionPersister instance
func NewMapSessionPersister() (*MapSessionPersister, error) {
	return &MapSessionPersister{
		db: &sync.Map{},
	}, nil
}

// ReadFromPersister reads the session data for the given userID
func (s *MapSessionPersister) ReadFromPersister(userID string) (*Session, error) {
	data := &Session{}

	if s.db == nil {
		return nil, errors.New("Connection to DB does not exist.")
	}

	if userID == "" {
		return nil, errors.New("User ID is empty.")
	}

	dataCopyB, ok := s.db.Load(userID)
	if ok {
		logrus.Debugf("retrieved session for user with id: %s", userID)
		newData, ok1 := dataCopyB.(*Session)
		if ok1 {
			logrus.Debugf("session for user with id: %s was read in tact.", userID)
			data = newData
		} else {
			logrus.Warnf("session for user with id: %s was NOT read in tact.", userID)
		}
	} else {
		logrus.Warnf("unable to find session for user with id: %s.", userID)
	}
	return data, nil
}

// WriteToPersister persists session for the user
func (s *MapSessionPersister) WriteToPersister(userID string, data *Session) error {
	if s.db == nil {
		return errors.New("connection to DB does not exist")
	}

	if userID == "" {
		return errors.New("User ID is empty.")
	}

	if data == nil {
		return errors.New("Given config data is nil.")
	}
	data.UpdatedAt = time.Now()
	newSess := &Session{}
	if err := copier.Copy(newSess, data); err != nil {
		logrus.Errorf("session copy error: %v", err)
		return err
	}

	s.db.Store(userID, newSess)

	return nil
}

// DeleteFromPersister removes the session for the user
func (s *MapSessionPersister) DeleteFromPersister(userID string) error {
	if s.db == nil {
		return errors.New("Connection to DB does not exist.")
	}

	if userID == "" {
		return errors.New("User ID is empty.")
	}
	s.db.Delete(userID)
	return nil
}

// ClosePersister closes the DB
func (s *MapSessionPersister) ClosePersister() {
	s.db = nil
}