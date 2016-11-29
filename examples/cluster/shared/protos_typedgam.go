// Code generated by protoc-gen-gogo.
// source: protos.proto
// DO NOT EDIT!

/*
Package shared is a generated protocol buffer package.

It is generated from these files:
	protos.proto

It has these top-level messages:
	HelloRequest
	HelloResponse
	AddRequest
	AddResponse
*/
package shared

import errors "errors"
import log "log"
import actor "github.com/AsynkronIT/gam/actor"
import cluster "github.com/AsynkronIT/gam/cluster"
import grain "github.com/AsynkronIT/gam/cluster/grain"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var xHelloFactory func() Hello

func HelloFactory(factory func() Hello) {
	xHelloFactory = factory
}

func GetHelloGrain(id string) *HelloGrain {
	return &HelloGrain{ID: id}
}

type Hello interface {
	Init(id string)

	SayHello(*HelloRequest) (*HelloResponse, error)

	Add(*AddRequest) (*AddResponse, error)
}
type HelloGrain struct {
	ID string
}

func (g *HelloGrain) SayHello(r *HelloRequest, options ...grain.GrainCallOption) (*HelloResponse, error) {
	conf := grain.ApplyGrainCallOptions(options)
	fun := func() (*HelloResponse, error) {
		pid, err := cluster.Get(g.ID, "Hello")
		if err != nil {
			return nil, err
		}
		bytes, err := proto.Marshal(r)
		if err != nil {
			return nil, err
		}
		request := &cluster.GrainRequest{Method: "SayHello", MessageData: bytes}
		response, err := pid.RequestFuture(request, conf.Timeout).Result()
		if err != nil {
			return nil, err
		}
		switch msg := response.(type) {
		case *cluster.GrainResponse:
			result := &HelloResponse{}
			err = proto.Unmarshal(msg.MessageData, result)
			if err != nil {
				return nil, err
			}
			return result, nil
		case *cluster.GrainErrorResponse:
			return nil, errors.New(msg.Err)
		default:
			return nil, errors.New("Unknown response")
		}
	}

	var res *HelloResponse
	var err error
	for i := 0; i < conf.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (g *HelloGrain) SayHelloChan(r *HelloRequest, options ...grain.GrainCallOption) (<-chan *HelloResponse, <-chan error) {
	c := make(chan *HelloResponse)
	e := make(chan error)
	go func() {
		res, err := g.SayHello(r, options...)
		if err != nil {
			e <- err
		} else {
			c <- res
		}
		close(c)
		close(e)
	}()
	return c, e
}

func (g *HelloGrain) Add(r *AddRequest, options ...grain.GrainCallOption) (*AddResponse, error) {
	conf := grain.ApplyGrainCallOptions(options)
	fun := func() (*AddResponse, error) {
		pid, err := cluster.Get(g.ID, "Hello")
		if err != nil {
			return nil, err
		}
		bytes, err := proto.Marshal(r)
		if err != nil {
			return nil, err
		}
		request := &cluster.GrainRequest{Method: "Add", MessageData: bytes}
		response, err := pid.RequestFuture(request, conf.Timeout).Result()
		if err != nil {
			return nil, err
		}
		switch msg := response.(type) {
		case *cluster.GrainResponse:
			result := &AddResponse{}
			err = proto.Unmarshal(msg.MessageData, result)
			if err != nil {
				return nil, err
			}
			return result, nil
		case *cluster.GrainErrorResponse:
			return nil, errors.New(msg.Err)
		default:
			return nil, errors.New("Unknown response")
		}
	}

	var res *AddResponse
	var err error
	for i := 0; i < conf.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		}
	}
	return nil, err
}

func (g *HelloGrain) AddChan(r *AddRequest, options ...grain.GrainCallOption) (<-chan *AddResponse, <-chan error) {
	c := make(chan *AddResponse)
	e := make(chan error)
	go func() {
		res, err := g.Add(r, options...)
		if err != nil {
			e <- err
		} else {
			c <- res
		}
		close(c)
		close(e)
	}()
	return c, e
}

type HelloActor struct {
	inner Hello
}

func (a *HelloActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.inner = xHelloFactory()
		a.inner.Init("abc")
	case *cluster.GrainRequest:
		switch msg.Method {

		case "SayHello":
			req := &HelloRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.SayHello(req)
			if err == nil {
				bytes, err := proto.Marshal(r0)
				if err != nil {
					log.Fatalf("[GRAIN] proto.Marshal failed %v", err)
				}
				resp := &cluster.GrainResponse{MessageData: bytes}
				ctx.Respond(resp)
			} else {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
			}

		case "Add":
			req := &AddRequest{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.Add(req)
			if err == nil {
				bytes, err := proto.Marshal(r0)
				if err != nil {
					log.Fatalf("[GRAIN] proto.Marshal failed %v", err)
				}
				resp := &cluster.GrainResponse{MessageData: bytes}
				ctx.Respond(resp)
			} else {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
			}

		}
	default:
		log.Printf("Unknown message %v", msg)
	}
}

func init() {

	cluster.Register("Hello", actor.FromProducer(func() actor.Actor {
		return &HelloActor{}
	}))

}

// type hello struct {
//	id	string
// }
// func (state *hello) Init(id string) {
// 	state.id = id
// }

// func (*hello) SayHello(r *HelloRequest) (*HelloResponse, error) {
// 	return &HelloResponse{}, nil
// }

// func (*hello) Add(r *AddRequest) (*AddResponse, error) {
// 	return &AddResponse{}, nil
// }

// func init() {
// 	//apply DI and setup logic

// 	HelloFactory(func() Hello { return &hello{} })

// }
