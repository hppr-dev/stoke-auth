package cluster

import (
	"context"
	"reflect"
	"testing"
)

func TestStaticDiscoverer_Peers(t *testing.T) {
	ctx := context.Background()

	t.Run("single URL", func(t *testing.T) {
		d := &StaticDiscoverer{URLs: []string{"http://a"}}
		got, err := d.Peers(ctx)
		if err != nil {
			t.Fatalf("Peers(): err = %v, want nil", err)
		}
		want := []string{"http://a"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Peers() = %v, want %v", got, want)
		}
	})

	t.Run("empty URLs", func(t *testing.T) {
		d := &StaticDiscoverer{URLs: []string{}}
		got, err := d.Peers(ctx)
		if err != nil {
			t.Fatalf("Peers(): err = %v, want nil", err)
		}
		if len(got) != 0 {
			t.Errorf("Peers() = %v, want empty slice", got)
		}
	})

	t.Run("multiple URLs", func(t *testing.T) {
		d := &StaticDiscoverer{URLs: []string{"http://a", "https://b:8080"}}
		got, err := d.Peers(ctx)
		if err != nil {
			t.Fatalf("Peers(): err = %v, want nil", err)
		}
		want := []string{"http://a", "https://b:8080"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Peers() = %v, want %v", got, want)
		}
	})

	t.Run("nil receiver", func(t *testing.T) {
		var d *StaticDiscoverer
		got, err := d.Peers(ctx)
		if err != nil {
			t.Fatalf("Peers() on nil receiver: err = %v, want nil", err)
		}
		if got != nil {
			t.Errorf("Peers() on nil receiver = %v, want nil", got)
		}
	})
}
