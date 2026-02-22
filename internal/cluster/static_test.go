package cluster

import (
	"context"
	"reflect"
	"testing"
)

func TestStaticDiscoverer_Peers(t *testing.T) {
	ctx := context.Background()
	d := &StaticDiscoverer{URLs: []string{"http://a"}}
	got, err := d.Peers(ctx)
	if err != nil {
		t.Fatalf("Peers(): err = %v, want nil", err)
	}
	want := []string{"http://a"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Peers() = %v, want %v", got, want)
	}
}
