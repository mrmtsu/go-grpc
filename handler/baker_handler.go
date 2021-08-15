package handler

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"pancake.maker/pb"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type BakerHandler struct {
	report *report
}

type report struct {
	sync.Mutex
	data map[pb.Pancake_Menu]int
}

// BakerHandlerの初期化
func NewBakerHandler() *BakerHandler {
	return &BakerHandler{
		report: &report{
			data: make(map[pb.Pancake_Menu]int),
		},
	}
}

func (h *BakerHandler) Bake(ctx context.Context, req *pb.BakeRequest) (*pb.BakeResponse, error) {
	if req.Menu == pb.Pancake_UNKNOWN || req.Menu == pb.Pancake_SPICY_CURRY {
		return nil, status.Errorf(codes.InvalidArgument, "パンケーキを選んでください！")
	}

	now := time.Now()
	h.report.Lock()
	h.report.data[req.Menu] = h.report.data[req.Menu] + 1
	h.report.Unlock()

	return &pb.BakeResponse{
		Pancake: &pb.Pancake{
			Menu:           req.Menu,
			ChefName:       "gami",
			TechnicalScore: rand.Float32(),
			CreateTime: &timestamp.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		},
	}, nil
}

func (h *BakerHandler) Report(ctx context.Context, req *pb.ReportRequest) (*pb.ReportResponse, error) {
	counts := make([]*pb.Report_BakeCount, 0)

	h.report.Lock()
	for k, v := range h.report.data {
		counts = append(counts, &pb.Report_BakeCount{
			Menu:  k,
			Count: int32(v),
		})
	}
	h.report.Unlock()

	return &pb.ReportResponse{
		Report: &pb.Report{
			BakeCounts: counts,
		},
	}, nil
}
