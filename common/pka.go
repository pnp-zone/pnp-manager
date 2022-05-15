package common

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/miekg/dns"
	"github.com/monaco-io/request"
	"github.com/pnp-zone/pnp-manager/conf"
	"strings"
	"time"
)

func CheckPKARecord(config *conf.Config) (*crypto.KeyRing, error) {
	m := new(dns.Msg)
	m.SetQuestion(config.PKA, dns.TypeTXT)
	dnsClient := new(dns.Client)
	in, _, err := dnsClient.Exchange(m, "127.0.0.53:53")
	if err != nil {
		return nil, ErrNoDNS
	}

	var record string
	for _, rr := range in.Answer {
		record = strings.Replace(strings.Trim(rr.String(), rr.Header().String()), "\"", "", -1)
	}

	parts := strings.Split(record, ";")
	if len(parts) != 3 {
		return nil, ErrInvalidPKARecord
	}

	if parts[0] != "v=pka1" {
		return nil, ErrInvalidPKARecord
	}

	if !strings.HasPrefix(parts[1], "fpr=") {
		return nil, ErrInvalidPKARecord
	}

	fpr := strings.ToLower(strings.SplitN(parts[1], "=", 2)[1])

	if !strings.HasPrefix(parts[2], "uri=") {
		return nil, ErrInvalidPKARecord
	}

	uri := strings.SplitN(parts[2], "=", 2)[1]

	c := request.Client{
		URL:     uri,
		Method:  "GET",
		Timeout: time.Second,
	}

	armoredKey := strings.TrimSpace(c.Send().String())
	key, err := crypto.NewKeyFromArmored(armoredKey)
	if err != nil {
		return nil, ErrInvalidMasterKey
	}

	if key.IsRevoked() {
		return nil, ErrInvalidMasterKey
	}

	if key.IsExpired() {
		return nil, ErrInvalidMasterKey
	}

	if strings.ToLower(key.GetFingerprint()) != fpr {
		return nil, ErrInvalidMasterFingerprint
	}

	keyring, err := crypto.NewKeyRing(key)
	if err != nil {
		return nil, err
	}

	return keyring, nil
}
