package model

// ShortenRequest is the JSON payload for POST /api/shorten.
// URL must contain the original absolute URL to shorten.
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse is the JSON response for POST /api/shorten.
// Result contains the absolute short URL.
type ShortenResponse struct {
	Result string `json:"result"`
}

// BatchShortenRequestItem describes a single item in the batch shorten request.
// CorrelationID is echoed back in the response to keep client-side ordering.
type BatchShortenRequestItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// BatchShortenResponseItem describes a single result of the batch shorten operation.
// ShortURL is the absolute short URL.
type BatchShortenResponseItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UrlsByUserResponseItem represents an item in the response of GET /api/user/urls.
type UrlsByUserResponseItem struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
