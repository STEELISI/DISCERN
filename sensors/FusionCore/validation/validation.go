package validation;

import "sync"
import "context"

import bridge "FusionBridge/validation"

type SendAndRecvServer struct {
    bridge.UnimplementedSendAndRecvServer
    mu sync.Mutex
}

func NewServer() *SendAndRecvServer {
    s := &SendAndRecvServer{}
    return s
}

func (s *SendAndRecvServer) HelloWorldBasic(ctx context.Context, MSG *bridge.HelloWorld) (*bridge.HelloWorldACK, error) {
    return &bridge.HelloWorldACK{Tmp: MSG.Tmp - 10}, nil
}

