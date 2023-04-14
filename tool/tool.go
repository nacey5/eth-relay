package tool

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
)

func GetRealDecimalValue(value string, decimal int) string {
	if strings.Contains(value, ".") {
		//.little
		arr := strings.Split(value, ".")
		if len(arr) != 2 {
			return ""
		}
		num := len(arr[1])
		left := decimal - num
		return arr[0] + arr[1] + strings.Repeat("0", left)
	} else {
		return value + strings.Repeat("0", decimal)
	}
}

// BuildERC20TransferData build match the ERC20 standard
func BuildERC20TransferData(value, receiver string, decimal int) string {
	realValue := GetRealDecimalValue(value, decimal)
	valueBig, _ := new(big.Int).SetString(realValue, 10)
	methodId := "0xa9059cbb"
	param1 := common.HexToHash(receiver).String()[2:]
	param2 := common.Bytes2Hex(valueBig.Bytes())[2:]
	return methodId + param1 + param2
}
