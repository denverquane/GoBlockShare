package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"crypto/sha256"
	"encoding/hex"
)

type SignableTransaction struct {
	Origin      OriginInfo // needed to say who I am (WITHIN the transaction)
	Transaction TorrentTransaction
	R, S        *big.Int // signature of the transaction, should be separate from the actual "message" components
	TxID        string
}

func (st SignableTransaction) SetRS(r *big.Int, s *big.Int) SignableTransaction {
	st.R = r
	st.S = s
	return st
}

func (st SignableTransaction) GetRS() (*big.Int, *big.Int) {
	return st.R, st.S
}

func (st SignableTransaction) GetHash(haveRSbeenSet bool) []byte {
	h := sha256.New()
	h.Write(st.Origin.GetRawBytes())
	h.Write(st.Transaction.GetRawBytes())

	//Filters the cases where we just want the hash for non-signing purposes
	//(if the transaction hasn't been signed, we shouldn't hash R and S as they don't matter)
	if haveRSbeenSet {
		h.Write(st.R.Bytes())
		h.Write(st.S.Bytes())
	}
	return h.Sum(nil)
}

func (st SignableTransaction) GetOrigin() OriginInfo {
	return st.Origin
}

func (st SignableTransaction) ToString() string {
	return st.Origin.ToString() + "\"txref\":[],\n" +
		string(st.Transaction.GetRawBytes()) + "\",\n\"r\":" + st.R.String() + ",\n\"s\":" +
		st.S.String() + "\n}\n"
}


func (st SignableTransaction) SignAndSetTxID(priv *ecdsa.PrivateKey) SignableTransaction {
	hashed := st.GetHash(false)
	r, s, err := ecdsa.Sign(rand.Reader, priv, hashed)

	if err != nil {
		log.Println("Error when signing transaction!")
		return st
	}
	st = st.SetRS(r, s)
	st.TxID = hex.EncodeToString(st.GetHash(true))
	return st
}

func (st SignableTransaction) Verify() bool {
	origin := st.GetOrigin()
	key := ecdsa.PublicKey{AUTHENTICATION_CURVE, origin.PubKeyX, origin.PubKeyY}

	if st.VerifyWithKey(key) { //signed transaction isn't verified with the public key
		fmt.Println("Signed doesnt verify")
		return false
	} else if HashPublicToB64Address(key) != Base64Address(origin.Address) { //public key does not match up with the address
		fmt.Println("public doesnt match address")
		return false
	}
	return true
}

func (st SignableTransaction) VerifyWithKey(key ecdsa.PublicKey) bool {
	r, s := st.GetRS()
	return ecdsa.Verify(&key, st.GetHash(true), r, s)
}

type OriginInfo struct {
	PubKeyX *big.Int
	PubKeyY *big.Int
	Address Base64Address
}

func (oi OriginInfo) GetRawBytes() []byte {
	return []byte(string(oi.PubKeyX.Bytes()) + string(oi.PubKeyY.Bytes()) + string(oi.Address))
}

func (oi OriginInfo) ToString() string {
	return "\n{\n\"origin\":\n{\n\"address\":\"" + string(oi.Address) + "\",\n\"pubkeyx\":" + oi.PubKeyX.String() +
		",\n\"pubkeyy\":" + oi.PubKeyY.String() + "\n},\n"
}

func AddressToOriginInfo(address PersonalAddress) OriginInfo {
	return OriginInfo{address.PublicKey.X, address.PublicKey.Y, address.Address}
}
//
//type RESTWrappedFullTransaction struct {
//	Origin   OriginInfo
//	Txref    []string
//	Quantity float64
//	Payload  string
//	R        big.Int
//	S        big.Int
//	DestAddr string
//}
//
//func (rest RESTWrappedFullTransaction) ConvertToFull() (FullTransaction, error) {
//	var signed = SignableTransaction{rest.Origin, &rest.R, &rest.S}
//	var full = MakeFull(signed, rest.Txref)
//	return full, nil
//}