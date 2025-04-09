package genqr

import (
	"github.com/ducnpdev/vietqr"
	"strconv"
)

type QrService struct{}

var QrGenerator *QrService

func InitQrService() {
	QrGenerator = &QrService{}
}

func (QrService *QrService) GenerateQrCode(text string, amount int, description string) string {
	content := vietqr.GenerateViQR(vietqr.RequestGenerateViQR{
		MerchantAccountInformation: vietqr.MerchantAccountInformation{
			AcqID:     "686868",
			AccountNo: "058914618",
		},
		TransactionAmount: strconv.Itoa(amount),

		AdditionalDataFieldTemplate: vietqr.AdditionalDataFieldTemplate{
			Description: description,
		},
		Mcc:          "7832",
		ReceiverName: "Trương Hoàng Nguyên",
	})
	return content
}
