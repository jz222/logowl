package logs

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/jz222/loggy/models"
)

func SaveError(errorLog models.Error) {
	hash := md5.Sum([]byte(errorLog.Message + errorLog.Stacktrace))
	errorLog.Fingerprint = hex.EncodeToString(hash[:])

	data, _ := json.MarshalIndent(errorLog, "", "  ")
	fmt.Println(string(data))
}
