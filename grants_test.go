package iroha

import (
	"encoding/base64"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/ishinvin/irohasign"
	irohapb "github.com/ishinvin/irohasign/proto"
)

func TestSignGrantPermission(t *testing.T) {
	kp, err := irohasign.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	inst := &Instance{}

	result, err := inst.SignGrantPermission(
		kp.PublicHex(), kp.PrivateHex(), "alice@test", "bob@test",
		[]string{"can_add_my_signatory", "can_set_my_quorum"},
	)
	if err != nil {
		t.Fatalf("SignGrantPermission: %v", err)
	}
	if len(result.ReducedHash) != 64 {
		t.Fatalf("ReducedHash = %q, want 64 hex chars", result.ReducedHash)
	}

	txBytes, err := base64.StdEncoding.DecodeString(result.Transaction)
	if err != nil {
		t.Fatalf("Transaction is not valid base64: %v", err)
	}
	var tx irohapb.Transaction
	if err := proto.Unmarshal(txBytes, &tx); err != nil {
		t.Fatalf("Transaction is not a valid Transaction: %v", err)
	}

	commands := tx.GetPayload().GetReducedPayload().GetCommands()
	if len(commands) != 2 {
		t.Fatalf("commands = %d, want 2", len(commands))
	}
	want := []irohapb.GrantablePermission{
		irohapb.GrantablePermission_can_add_my_signatory,
		irohapb.GrantablePermission_can_set_my_quorum,
	}
	for i, cmd := range commands {
		grant := cmd.GetGrantPermission()
		if grant == nil {
			t.Fatalf("command %d is not a GrantPermission", i)
		}
		if grant.GetAccountId() != "bob@test" {
			t.Fatalf("command %d account = %q, want bob@test", i, grant.GetAccountId())
		}
		if grant.GetPermission() != want[i] {
			t.Fatalf("command %d permission = %v, want %v", i, grant.GetPermission(), want[i])
		}
	}
	if tx.GetPayload().GetOptionalBatchMeta() != nil {
		t.Fatalf("plain grant transaction should carry no batch metadata")
	}
}

func TestSignGrantPermission_UnknownPermission(t *testing.T) {
	kp, err := irohasign.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	inst := &Instance{}

	_, err = inst.SignGrantPermission(kp.PublicHex(), kp.PrivateHex(), "alice@test", "bob@test", []string{"not_a_real_permission"})
	if err == nil {
		t.Fatal("expected error for unknown permission, got nil")
	}
	if !strings.Contains(err.Error(), "unknown GrantablePermission") {
		t.Fatalf("error = %q, want it to mention unknown GrantablePermission", err)
	}
}

func TestSignGrantPermission_BadKey(t *testing.T) {
	inst := &Instance{}
	_, err := inst.SignGrantPermission("not-hex", "not-hex", "alice@test", "bob@test", []string{"can_add_my_signatory"})
	if err == nil {
		t.Fatal("expected error for invalid key hex, got nil")
	}
	if !strings.Contains(err.Error(), "parse key pair") {
		t.Fatalf("error = %q, want it to mention key parsing", err)
	}
}

func TestReSignWithBatchMeta(t *testing.T) {
	kp, err := irohasign.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	inst := &Instance{}

	plain, err := inst.SignGrantPermission(kp.PublicHex(), kp.PrivateHex(), "alice@test", "bob@test", []string{"can_add_my_signatory"})
	if err != nil {
		t.Fatalf("SignGrantPermission: %v", err)
	}

	reducedHashes := []string{plain.ReducedHash, "DEADBEEF"}
	resigned, err := inst.ReSignWithBatchMeta(kp.PublicHex(), kp.PrivateHex(), plain.Transaction, reducedHashes)
	if err != nil {
		t.Fatalf("ReSignWithBatchMeta: %v", err)
	}

	txBytes, err := base64.StdEncoding.DecodeString(resigned)
	if err != nil {
		t.Fatalf("result is not valid base64: %v", err)
	}
	var tx irohapb.Transaction
	if err := proto.Unmarshal(txBytes, &tx); err != nil {
		t.Fatalf("result is not a valid Transaction: %v", err)
	}

	batchMeta := tx.GetPayload().GetBatch()
	if batchMeta == nil {
		t.Fatalf("re-signed transaction carries no batch metadata")
	}
	if got := batchMeta.GetReducedHashes(); len(got) != 2 || got[0] != reducedHashes[0] || got[1] != reducedHashes[1] {
		t.Fatalf("batch reduced hashes = %v, want %v", got, reducedHashes)
	}
}

func TestReSignWithBatchMeta_InvalidBase64(t *testing.T) {
	kp, err := irohasign.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	inst := &Instance{}

	_, err = inst.ReSignWithBatchMeta(kp.PublicHex(), kp.PrivateHex(), "not-base64!!!", []string{"DEADBEEF"})
	if err == nil {
		t.Fatal("expected error for invalid base64, got nil")
	}
	if !strings.Contains(err.Error(), "decode transaction") {
		t.Fatalf("error = %q, want it to mention transaction decoding", err)
	}
}

func TestReSignWithBatchMeta_BadKey(t *testing.T) {
	inst := &Instance{}
	_, err := inst.ReSignWithBatchMeta("not-hex", "not-hex", base64.StdEncoding.EncodeToString([]byte("irrelevant")), []string{"DEADBEEF"})
	if err == nil {
		t.Fatal("expected error for invalid key hex, got nil")
	}
	if !strings.Contains(err.Error(), "parse key pair") {
		t.Fatalf("error = %q, want it to mention key parsing", err)
	}
}
