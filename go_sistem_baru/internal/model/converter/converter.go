package converter

import (
	"capstone-qris/internal/entity"
	"capstone-qris/internal/model"
	"capstone-qris/internal/util"
)

func ToInquiryResponse(merchant *entity.Merchant, inquiryID string, parsed *util.QRISData) *model.InquiryResponse {
	fixedAmount := 0.0
	terminalID := ""
	if parsed != nil {
		fixedAmount = parsed.FixedAmount
		terminalID = parsed.TerminalID
	}
	return &model.InquiryResponse{
		MerchantID:   merchant.MerchantID,
		MerchantName: merchant.MerchantName,
		TerminalID:   terminalID,
		City:         merchant.City,
		FixedAmount:  fixedAmount,
		InquiryID:    inquiryID,
	}
}

func ToTransactionStatusResponse(tx *entity.Transaction) *model.TransactionStatusResponse {
	return &model.TransactionStatusResponse{
		TransactionID: tx.TransactionID,
		Status:        tx.Status,
		FinalBalance:  tx.Amount,
		Timestamp:     tx.UpdatedAt,
	}
}
