package model
 
import "time"

const (
	HeaderContentType     = "Content-Type"
	HeaderContentEncoding = "Content-Encoding"
	HeaderAcceptEncoding  = "Accept-Encoding"
	HeaderLocation        = "Location"
	HeaderAuth            = "Authorization"
	ContentTypeText       = "text/plain"
	ContentTypeJSON       = "application/json"
	ContentTypeHTML       = "text/html"

	CookieAuth = "Auth"
	CookieUser = "User"

	ContextValueUser = "User"

	CompressTypeGzip    = "gzip"
	CompressTypeDeflate = "deflate"
	CompressTypeEmpty   = ""

	Expire = time.Hour * 12
)