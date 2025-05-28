package dto

import "github.com/google/uuid"

type InternalNotificationRequest struct {
	TemplateName string `validate:"required"`
	Recipient    string `validate:"required"`
	Placeholder  map[string]string
}

type GetNotificationByTypeResponse struct {
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"`
	Name string    `json:"name"`
}

type CreateNotificationRequest struct {
	Title   string `json:"title" validate:"required"`
	Detail  string `json:"detail" validate:"required"`
	Url     string `json:"url" validate:"required"`
	UserID  string `json:"user_id" validate:"required"`
	NTypeID string `json:"n_type_id" validate:"required"`
}

type GetListNotificationResponse struct {
	Notification []NotificationDetail `json:"notification"`
	TotalPages   int                  `json:"total_page"`
	CurrentPage  int                  `json:"current_page"`
	PageSize     int                  `json:"page_size"`
	TotalData    int                  `json:"total_data"`
}

type NotificationDetail struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Detail    *string `json:"detail"`
	Url       string  `json:"url"`
	NTypeID   string  `json:"n_type_id"`
	UserID    string  `json:"user_id"`
	CreatedAt string  `json:"created_at"`
}

type SendBatchFcmNotificationRequest struct {
	FcmToken []string `json:"fcm_token"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
}

type SendFcmNotificationRequest struct {
	FcmToken        string `json:"fcm_token"`
	Title           string `json:"title"`
	Body            string `json:"body"`
	UserID          string `json:"user_id"`
	IsStatusChanged bool   `json:"is_status_changed"`
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
}

type SendTransactionEmailRequest struct {
	TOName                 string      `json:"to_name"`
	TOEmail                string      `json:"to_email"`
	Items                  []OrderItem `json:"items"`
	TotalProductAmount     int         `json:"total_product_amount"`
	TotalTransactionAmount int         `json:"total_transaction_amount"`
	TotalDiscount          int         `json:"total_discount"`
	TotalQuantity          int         `json:"total_quantity"`
	TaxAmount              int         `json:"tax_amount"`
	TaxValue               int         `json:"tax_value"`
	IsStatusChanged        bool        `json:"is_status_changed"`
}

type OrderItem struct {
	ProductName  string `json:"product_name"`
	Quantity     int    `json:"quantity"`
	ProductPrice int    `json:"product_price"`
	TotalPrice   int    `json:"total_price"`
}
