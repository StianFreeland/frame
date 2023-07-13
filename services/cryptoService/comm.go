package cryptoService

import (
	"crypto/sha256"
	"encoding/hex"
	"frame/services/zlog"
	"go.uber.org/zap"
	"os"
)

const moduleName = "crypto service"
const rootPwd = "YLRLM1FAkRRAhNGqNrn88pAFmLCX0phO"
const secretKey = "CGMjNNSsEh1E07bSNrSjfmCI2kLjTQz2"

var TokenKey = []byte("Oy3gxdwZMa3mlzF73WS1wvV5Vc5WH6hc")
var PvtKey []byte

func Init() {
	zlog.Warn(moduleName, "init ...")
	initPvtKey()
}

func GetRootPwdSum() string {
	sum := sha256.Sum256([]byte(rootPwd + secretKey))
	return hex.EncodeToString(sum[:])
}

func GetPwdSum(pwd string) string {
	sum := sha256.Sum256([]byte(pwd + secretKey))
	return hex.EncodeToString(sum[:])
}

func initPvtKey() {
	content, err := os.ReadFile("certs/private.pem")
	if err != nil {
		zlog.Fatal(moduleName, "init private key", zap.Error(err))
	}
	PvtKey = content
	if len(PvtKey) == 0 {
		zlog.Fatal(moduleName, "init private key", zap.ByteString("pvt_key", PvtKey))
	}
}
