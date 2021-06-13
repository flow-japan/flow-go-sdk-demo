package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/onflow/cadence"
	"google.golang.org/grpc"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
)

func main() {
	UserSignatureDemo()
}

// 署名を検証する Cadence コード
var script = []byte(`
import Crypto

pub fun main(rawPublicKeys: [String], weights: [UFix64], signatures: [String], signedData: [UInt8]): Bool {
  let keyList = Crypto.KeyList()
  var i = 0
  for rawPublicKey in rawPublicKeys {
    keyList.add(
      PublicKey(
        publicKey: rawPublicKey.decodeHex(),
        signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
      ),
      hashAlgorithm: HashAlgorithm.SHA3_256,
      weight: weights[i],
    )
    i = i + 1
  }

  let signatureSet: [Crypto.KeyListSignature] = []
  var j = 0
  for signature in signatures {
    signatureSet.append(
      Crypto.KeyListSignature(
        keyIndex: j,
        signature: signature.decodeHex()
      )
    )
    j = j + 1
  }

  return keyList.isValid(
    signatureSet: signatureSet,
    signedData: signedData,
  )
}
`)

func UserSignatureDemo() {
	// 署名対象（任意のメッセージ）
	message := "test test test"

	// 秘密鍵
	rawPrivateKey := "9a3259d7c18fd98ccf51356a48df5b63d7d544153db49079c46d152ea9739539"

	// 署名処理
	privateKey, _ := crypto.DecodePrivateKeyHex(crypto.ECDSA_P256, rawPrivateKey)
	signer := crypto.NewInMemorySigner(privateKey, crypto.SHA3_256)
	messageHex := []byte(message)
	signature, _ := flow.SignUserMessage(signer, messageHex)

	// 署名検証

	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	Handle(err)

	publicKeys := cadence.NewArray([]cadence.Value{
		cadence.String(hex.EncodeToString(privateKey.PublicKey().Encode())),
	})

	weight, _ := cadence.NewUFix64("1.0")
	weights := cadence.NewArray([]cadence.Value{
		weight,
	})

	signatures := cadence.NewArray([]cadence.Value{
		cadence.String(hex.EncodeToString(signature)),
	})

	signedData := toUInt8Array(messageHex)

	result, err := flowClient.ExecuteScriptAtLatestBlock(
		context.Background(),
		script,
		[]cadence.Value{
			publicKeys,
			weights,
			signatures,
			signedData,
		},
	)
	Handle(err)

	if result == cadence.NewBool(true) {
		fmt.Println("Signature verification succeeded")
	} else {
		fmt.Println("Signature verification failed")
	}
}

func toUInt8Array(bytes []byte) cadence.Value {
	cadenceUInt8Array := make([]cadence.Value, 0, len(bytes))
	for _, b := range bytes {
		cadenceUInt8Array = append(cadenceUInt8Array, cadence.NewUInt8(b))
	}
	return cadence.NewArray(cadenceUInt8Array)
}

func Handle(err error) {
	if err != nil {
		fmt.Println("err:", err.Error())
		panic(err)
	}
}
