package constant

// Field keys
const (
        MsgType     = "msgType"    
        TxType      = "txType"     
        PcposId     = "pcposId"    
        PcposTxnId  = "pcposTxnId" 
        Status      = "status"     
        QrCode      = "qrCode"     
        ErrorCode   = "errorCode"  
        TrmId       = "trmId"      
        TrmTrUniqNo = "trmTrUniqNo"
        AprvNo      = "aprvNo"     
        ErrorDetail = "errorDetail"
)

// Message Types
const (
        MsgTypeTxReq = "1"
        MsgTypeTxRes = "2"
)

// Transaction Status on Server/Terminal
const (
        TxStatusTxInitOnServer      = "1"
        TxStatusTxInitFailOnServer  = "2"
        TxStatusTxCreatedOnTerminal = "3"
        TxStatusTxSuccessOnTerminal = "4"
        TxStatusTxFailOnTerminal    = "5"
)

// Error Codes
const (
        ErrCodeNoErr              = "0"
        ErrCodeNotFoundPcposId    = "1"
        ErrCodeDuplicateTx        = "2"
        ErrCodeTcpServerError     = "3"
        ErrCodeNotFoundTx         = "4"
        ErrCodeFieldInvalid       = "5"
        ErrCodeTerminalNotConnect = "6"
        ErrCodeNotFoundOriginTx   = "7"
        ErrCodeTrmBusy            = "8"
        ErrCodeThereIsAnOngoingTx = "9"
        ErrCodeUnknown            = "10"
        ErrCodeTrmNotResponse     = "11"
        ErrCodeTxRefunded         = "12"
        ErrCodeTxNotSuccess       = "13"
        ErrCodeTxVoided           = "14"
        ErrCodeTxVoidNotSuccess   = "15"
        ErrCodeCannotMapping      = "16"
)

// Error Details
const (
        ErrDetailCode0  = "No error"
        ErrDetailCode1  = "Can not find PC-POS ID"
        ErrDetailCode2  = "Duplicate transaction"
        ErrDetailCode3  = "Can not send to Terminal TCP server error"
        ErrDetailCode4  = "Can not find the Transaction"
        ErrDetailCode5  = "Invalid field:"
        ErrDetailCode6  = "Terminal not connect"
        ErrDetailCode7  = "Origin transaction not found"
        ErrDetailCode8  = "Terminal is busy"
        ErrDetailCode9  = "There is an ongoing transaction on K9"
        ErrDetailCode10 = "Unknown Error"
        ErrDetailCode11 = "Terminal is not response"
        ErrDetailCode12 = "The transaction has been refunded"
        ErrDetailCode13 = "Original Tx is not success"
        ErrDetailCode14 = "The transaction has been voided"
        ErrDetailCode15 = "The VOID transaction was not successful"
        ErrDetailCode16 = "Cannot map fields between DTO and Model"
)

// Transaction Types
const (
        TxTypeSale        = "SALE"
        TxTypeVoid        = "VOID"
        TxTypeQR          = "QR"
        TxTypeInstallment = "INSTALLMENT"
        TxTypeQRRefund    = "QR_REFUND"
        TxTypeCheckStatus = "CHECK_STATUS"
        TxTypeCancel      = "CANCEL"
)

// Transaction Status general
const (
        TxStatusNone     = "NONE"
        TxStatusStarted  = "STARTED"
        TxStatusSuccess  = "SUCCESS"
        TxStatusGenQr    = "GEN_QR"
        TxStatusFailed   = "FAILED"
        TxStatusCanceled = "CANCELED"
        TxStatusRefunded = "REFUNDED"
        TxStatusVoided   = "VOIDED"
)