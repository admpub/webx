package machinery

import (
	"context"
	"fmt"

	"github.com/RichardKnop/machinery/v1/backends/result"
	"github.com/RichardKnop/machinery/v1/tasks"
)

// Send 发送到队列
func (s *Server) Send(ctx context.Context, taskSignatures ...*tasks.Signature) ([]*result.AsyncResult, error) {
	server := s.Server
	result := []*result.AsyncResult{}
	for _, taskSignature := range taskSignatures {
		asyncResult, err := server.SendTaskWithContext(ctx, taskSignature)
		if err != nil {
			return result, fmt.Errorf("could not send task: %w", err)
		}
		result = append(result, asyncResult)
	}
	return result, nil
}

func (s *Server) SendGroup(ctx context.Context, sendConcurrency int, taskSignatures ...*tasks.Signature) ([]*result.AsyncResult, error) {
	server := s.Server

	group, err := tasks.NewGroup(taskSignatures...)
	if err != nil {
		return nil, fmt.Errorf("error creating group: %w", err)
	}

	if sendConcurrency <= 0 {
		sendConcurrency = 10
	}

	asyncResults, err := server.SendGroupWithContext(ctx, group, sendConcurrency)
	if err != nil {
		return asyncResults, fmt.Errorf("could not send group: %w", err)
	}
	/*

		for _, asyncResult := range asyncResults {
			results, err = asyncResult.Get(time.Duration(time.Millisecond * 5))
			if err != nil {
				return fmt.Errorf("Getting task result failed with error: %s", err.Error())
			}
			log.INFO.Printf(
				"%v + %v = %v\n",
				asyncResult.Signature.Args[0].Value,
				asyncResult.Signature.Args[1].Value,
				tasks.HumanReadableResults(results),
			)
		}

	*/
	return asyncResults, err
}

func (s *Server) SendGroupWithChord(ctx context.Context, sendConcurrency int, chordTask *tasks.Signature, taskSignatures ...*tasks.Signature) (*result.ChordAsyncResult, error) {
	server := s.Server

	group, err := tasks.NewGroup(taskSignatures...)
	if err != nil {
		return nil, fmt.Errorf("error creating group: %w", err)
	}

	chord, err := tasks.NewChord(group, chordTask)
	if err != nil {
		return nil, fmt.Errorf("error creating chord: %w", err)
	}

	if sendConcurrency <= 0 {
		sendConcurrency = 10
	}

	chordAsyncResult, err := server.SendChordWithContext(ctx, chord, 10)
	if err != nil {
		return chordAsyncResult, fmt.Errorf("could not send chord: %w", err)
	}
	return chordAsyncResult, err
}

func (s *Server) SendChain(ctx context.Context, taskSignatures ...*tasks.Signature) (*result.ChainAsyncResult, error) {
	server := s.Server

	chain, err := tasks.NewChain(taskSignatures...)
	if err != nil {
		return nil, fmt.Errorf("error creating chain: %w", err)
	}

	chainAsyncResult, err := server.SendChainWithContext(ctx, chain)
	if err != nil {
		return chainAsyncResult, fmt.Errorf("could not send chain: %w", err)
	}
	return chainAsyncResult, err
}
