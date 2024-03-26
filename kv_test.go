package runn

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestKV(t *testing.T) {
	tests := []struct {
		in any
	}{
		{nil},
		{"str"},
		{3},
		{4.5},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			kv := newKV()
			kv.set("key", tt.in)
			got := kv.get("key")
			if diff := cmp.Diff(got, tt.in); diff != "" {
				t.Error(diff)
			}

			{
				kv.clear()
				got := kv.get("key")
				if got != nil {
					t.Errorf("got %v, want %v", got, nil)
				}
			}
		})
	}
}

func TestRunNWithKV(t *testing.T) {
	ctx := context.Background()
	book := "testdata/book/kv.yml"
	want := newRunNResult(t, 1, []*RunResult{
		{
			ID:          "6fdfa57431f3700a161b5ef02f945a117fd70216",
			Path:        "testdata/book/kv.yml",
			Err:         nil,
			StepResults: []*StepResult{{ID: "6fdfa57431f3700a161b5ef02f945a117fd70216?step=0", Key: "0", Err: nil}},
		},
	})
	ops, err := Load(book)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(50 * time.Millisecond)
		ops.SetKV("email", "test@example.com")
	}()
	if err := ops.RunN(ctx); err != nil {
		t.Error(err)
	}
	got := ops.Result()
	opts := []cmp.Option{
		cmpopts.IgnoreFields(runResultSimplified{}, "Elapsed"),
		cmpopts.IgnoreFields(stepResultSimplified{}, "Elapsed"),
	}
	if diff := cmp.Diff(got.simplify(), want.simplify(), opts...); diff != "" {
		t.Error(diff)
	}
}
