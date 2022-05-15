package common

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	jsoniter "github.com/json-iterator/go"
	"github.com/monaco-io/request"
	"github.com/myOmikron/echotools/color"
	"github.com/pnp-zone/pkg-manager/task"
	"github.com/pnp-zone/pnp-manager/conf"
	"os"
	"time"
)

var json = jsoniter.Config{
	EscapeHTML:    true,
	CaseSensitive: true,
}.Froze()

func ReceiveIndex(config *conf.Config) ([]task.Package, error) {
	if keyring, err := CheckPKARecord(config); err != nil {
		return nil, err
	} else {
		c := request.Client{
			URL:     config.URL + "/packages/index.json.asc",
			Method:  "GET",
			Timeout: time.Second,
		}
		sigData := c.Send().String()
		signature, err := crypto.NewPGPSignatureFromArmored(sigData)
		if err != nil {
			return nil, ErrBadIndexSignature
		}

		c.URL = config.URL + "/packages/index.json"
		data := c.Send().Bytes()

		message := crypto.NewPlainMessage(data)

		if err := keyring.VerifyDetached(message, signature, time.Now().Unix()); err != nil {
			color.Println(color.RED, "Index has an invalid signature")
			os.Exit(1)
		}

		var f task.IndexResponse
		if err := json.Unmarshal(data, &f); err != nil {
			return nil, err
		}

		return f.Packages, nil
	}
}
