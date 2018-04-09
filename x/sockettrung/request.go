package sockettrung

type Request struct {
	Data []byte
	Payload []byte
	URL string
	Client WsClient
}