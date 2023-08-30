package redis

// LICENSE MIT
// copied almost verbatim from https://github.com/rbcervilla/redisstore
// in order to add the session encryption present in boj/redistore (which
// needs redigo)

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/gob"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	otrace "github.com/opentracing/opentracing-go"
)

const (
	sessionSize        = 1024 * 16
	clientTimeout      = 10 * time.Second
	namespaceSeparator = ":"
)

// GorillaStore stores gorilla sessions in Redis
//
//nolint:golint
type GorillaStore struct {
	// Client to connect to redis
	Client RedisClient

	Codecs []securecookie.Codec
	// default Options to use when a new session is created
	Options   sessions.Options
	MaxLength int
	// key prefix with which the session will be stored
	keyPrefix string
	// key generator
	keyGen KeyGenFunc
	// session serialiser
	serialiser SessionSerialiser
}

// KeyGenFunc defines a function used by store to generate a key
type KeyGenFunc func() (string, error)

// NewRedisStore returns a new RedisStore with default configuration
func NewRedisStore(cfg RedisConfig, keyPairs ...[]byte) (*GorillaStore, error) {

	client, err := NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}

	rs := &GorillaStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: sessions.Options{
			Path:   "/",
			MaxAge: 86400 * 30,
		},
		Client:     client,
		MaxLength:  sessionSize,
		keyPrefix:  cfg.Namespace() + namespaceSeparator,
		keyGen:     generateRandomKey,
		serialiser: GobSerialiser{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	return rs, rs.Client.Ping(ctx).Err()
}

// Get returns a session for the given name after adding it to the registry.
func (s *GorillaStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
func (s *GorillaStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := s.Options
	session.Options = &opts
	session.IsNew = true

	c, err := r.Cookie(name)
	if err != nil {
		return session, nil
	}
	session.ID = c.Value

	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	err = s.load(ctx, session)
	if err == nil {
		session.IsNew = false
	} else if errors.Is(err, redis.Nil) {
		err = nil // no data stored
	}
	return session, err
}

// Save adds a single session to the response.
//
// If the Options.MaxAge of the session is <= 0 then the session file will be
// deleted from the store. With this process it enforces the properly
// session cookie handling so no need to trust in the cookie management in the
// web browser.
func (s *GorillaStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {

	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	// Delete if max-age is <= 0
	if session.Options.MaxAge <= 0 {
		if err := s.delete(ctx, session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		id, err := s.keyGen()
		if err != nil {
			return errors.New("redisstore: failed to generate session id")
		}
		session.ID = id
	}

	ctx, cancel = context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	if err := s.save(ctx, session); err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), session.ID, session.Options))
	return nil
}

// KeyPrefix sets the key prefix to store session in Redis
func (s *GorillaStore) KeyPrefix(keyPrefix string) {
	s.keyPrefix = keyPrefix
}

// KeyGen sets the key generator function
func (s *GorillaStore) KeyGen(f KeyGenFunc) {
	s.keyGen = f
}

// Serialiser sets the session serialiser to store session
func (s *GorillaStore) Serialiser(ss SessionSerialiser) {
	s.serialiser = ss
}

// Close closes the Redis store
func (s *GorillaStore) Close() error {
	return s.Client.Close()
}

// save writes session in Redis
func (s *GorillaStore) save(ctx context.Context, session *sessions.Session) error {

	b, err := s.serialiser.Serialise(session)
	if err != nil {
		return err
	}
	if s.MaxLength != 0 && len(b) > s.MaxLength {
		return errors.New("sessionStore: the value to store is too big")
	}
	span, ctx := otrace.StartSpanFromContext(ctx, "redis.redistore.Set")
	defer span.Finish()
	return s.Client.Set(ctx, s.keyPrefix+session.ID, b, time.Duration(session.Options.MaxAge)*time.Second).Err()
}

// load reads session from Redis
func (s *GorillaStore) load(ctx context.Context, session *sessions.Session) error {

	span, ctx := otrace.StartSpanFromContext(ctx, "redis.redistore.Set")
	cmd := s.Client.Get(ctx, s.keyPrefix+session.ID)
	span.Finish()

	if cmd.Err() != nil {
		return cmd.Err()
	}

	b, err := cmd.Bytes()
	if err != nil {
		return err
	}

	return s.serialiser.Deserialise(b, session)
}

// delete deletes session in Redis
func (s *GorillaStore) delete(ctx context.Context, session *sessions.Session) error {
	return s.Client.Del(ctx, s.keyPrefix+session.ID).Err()
}

// SessionSerialiser provides an interface for serialise/deserialise a session
type SessionSerialiser interface {
	Serialise(s *sessions.Session) ([]byte, error)
	Deserialise(b []byte, s *sessions.Session) error
}

// Gob serialiser
type GobSerialiser struct{}

func (gs GobSerialiser) Serialise(s *sessions.Session) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s.Values)
	if err == nil {
		return buf.Bytes(), nil
	}
	return nil, err
}

func (gs GobSerialiser) Deserialise(d []byte, s *sessions.Session) error {
	dec := gob.NewDecoder(bytes.NewBuffer(d))
	return dec.Decode(&s.Values)
}

// generateRandomKey returns a new random key
func generateRandomKey() (string, error) {
	k := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return "", err
	}
	return strings.TrimRight(base32.StdEncoding.EncodeToString(k), "="), nil
}
