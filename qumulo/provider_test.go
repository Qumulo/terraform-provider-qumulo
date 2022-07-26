package qumulo

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"qumulo": testAccProvider,
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("QUMULO_HOST"); v == "" {
		t.Fatal("QUMULO_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("QUMULO_PORT"); v == "" {
		t.Fatal("QUMULO_PORT must be set for acceptance tests")
	}
	if v := os.Getenv("QUMULO_USERNAME"); v == "" {
		t.Fatal("QUMULO_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("QUMULO_PASSWORD"); v == "" {
		t.Fatal("QUMULO_PASSWORD must be set for acceptance tests")
	}
}
