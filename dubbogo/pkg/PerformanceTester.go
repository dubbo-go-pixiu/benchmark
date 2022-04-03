package pkg

const (
	PayloadType_COMPRESSABLE PayloadType = 0
)

type PayloadType int32

type StressRequest struct {
	ResponseType PayloadType `json:"response_type,omitempty"`
	// Desired payload size in the response from the server.
	ResponseSize int32 `json:"response_size,omitempty"`
	// Optional input payload sent along with the request.
	Payload *Payload `json:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

type Payload struct {

	// 请求包体的大小
	Type PayloadType `json:"type,omitempty"`
	// 包体序列化成字节数组的大小
	Body []byte `json:"body,omitempty"`
}
