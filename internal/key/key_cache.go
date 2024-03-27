package key

import (
	"context"
	"encoding/json"
	"log"
	"stoke/internal/ent"
	"stoke/internal/ent/privatekey"
	"time"

	"entgo.io/ent/dialect/sql"
)

type KeyCache[P PrivateKey] struct {
	keys []KeyPair[P]
	KeyDuration time.Duration
	TokenDuration time.Duration
}

type publicJson struct {
	Text    string    `json:"text"`
	Expires int64     `json:"expires"`
	Renews  int64     `json:"renews"`
}

func (c *KeyCache[P]) GoManage() {
}

func (c *KeyCache[P]) CurrentKey() KeyPair[P] {
	return c.keys[len(c.keys) - 1]
}

func (c *KeyCache[P]) JSON() ([]byte, error) {
	out := make([]publicJson, len(c.keys))
	for i, k := range c.keys {
		out[i] = publicJson{
			Text:    k.Encode(),
			Expires: k.ExpiresAt().Unix(),
			Renews:  k.RenewsAt().Unix(),
		}
	}
	return json.Marshal(out)
}

func (c *KeyCache[P]) Generate() error {
	newKey := new(KeyPair[P])
	err := (*newKey).Generate()
	if err != nil {
		return err
	}

	(*newKey).SetExpires(time.Now().Add(c.KeyDuration))
	(*newKey).SetRenews(time.Now().Add(c.KeyDuration).Add(-c.TokenDuration))

	c.keys = append(c.keys, *newKey)

	return nil
}

func (c *KeyCache[P]) Bootstrap(db *ent.Client, pair KeyPair[P]) error {
	now := time.Now()
	pk, err := db.PrivateKey.Query().
		Order(privatekey.ByExpires(sql.OrderDesc())).
		First(context.Background())

	if err != nil || pk.Expires.Before(now) {
		log.Printf("Could not retrieve private key: %v", err)
		pair.Generate()

		pk = db.PrivateKey.Create().
			SetText(pair.Encode()).
			SetExpires(now.Add(c.KeyDuration)).
			SetRenews(now.Add(c.KeyDuration).Add(-c.TokenDuration)).
			SaveX(context.Background())
	} else {
		err := pair.Decode(pk.Text)
		if err != nil {
			return err
		}
	}

	pair.SetExpires(pk.Expires)
	pair.SetRenews(pk.Renews)

	c.keys = append(c.keys, pair)
	return nil
}

func (c *KeyCache[P]) Clean() {
	now := time.Now()
	var valid []KeyPair[P]
	for _, e := range c.keys {
		if e.ExpiresAt().Before(now) {
			valid = append(valid, e)
		}
	}
	c.keys = valid
}
