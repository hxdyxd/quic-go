package crypto

import (
	"crypto"
	"encoding/binary"

	"github.com/bifurcation/mint"
	"github.com/phuslu/quic-go/internal/protocol"
)

var quicVersion1Salt = []byte{0xaf, 0xc8, 0x24, 0xec, 0x5f, 0xc7, 0x7e, 0xca, 0x1e, 0x9d, 0x36, 0xf3, 0x7f, 0xb2, 0xd4, 0x65, 0x18, 0xc3, 0x66, 0x39}

func newNullAEADAESGCM(connectionID protocol.ConnectionID, pers protocol.Perspective) (AEAD, error) {
	clientSecret, serverSecret := computeSecrets(connectionID)

	var mySecret, otherSecret []byte
	if pers == protocol.PerspectiveClient {
		mySecret = clientSecret
		otherSecret = serverSecret
	} else {
		mySecret = serverSecret
		otherSecret = clientSecret
	}

	myKey, myIV := computeNullAEADKeyAndIV(mySecret)
	otherKey, otherIV := computeNullAEADKeyAndIV(otherSecret)

	return NewAEADAESGCM(otherKey, myKey, otherIV, myIV)
}

func computeSecrets(connectionID protocol.ConnectionID) (clientSecret, serverSecret []byte) {
	connID := make([]byte, 8)
	binary.BigEndian.PutUint64(connID, uint64(connectionID))
	handshakeSecret := mint.HkdfExtract(crypto.SHA256, quicVersion1Salt, connID)
	clientSecret = qhkdfExpand(handshakeSecret, "client hs", crypto.SHA256.Size())
	serverSecret = qhkdfExpand(handshakeSecret, "server hs", crypto.SHA256.Size())
	return
}

func computeNullAEADKeyAndIV(secret []byte) (key, iv []byte) {
	key = qhkdfExpand(secret, "key", 16)
	iv = qhkdfExpand(secret, "iv", 12)
	return
}
