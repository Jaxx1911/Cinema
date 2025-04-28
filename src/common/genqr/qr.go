package genqr

import (
	"TTCS/src/common/log"
	"context"
	"github.com/ducnpdev/vietqr"
	"strconv"
)

type QrService struct{}

var QrGenerator *QrService

func InitQrService() {
	QrGenerator = &QrService{}
}

func (QrService *QrService) GenerateQrCode(amount int, description string) string {
	content := vietqr.GenerateViQR(vietqr.RequestGenerateViQR{
		MerchantAccountInformation: vietqr.MerchantAccountInformation{
			AcqID:     "970423",
			AccountNo: "10001259387",
		},
		TransactionAmount: strconv.Itoa(amount),

		AdditionalDataFieldTemplate: vietqr.AdditionalDataFieldTemplate{
			Description: description,
		},
		Mcc:          "7832",
		ReceiverName: "Truong Hoang Nguyen",
	})
	log.Info(context.Background(), description)
	log.Info(context.Background(), content)
	return content
}
