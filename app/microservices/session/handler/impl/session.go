package impl

import (
	"backendServer/app/microservices/session/handler"
	"backendServer/app/microservices/session/usecase"
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/metrics"

	"golang.org/x/net/context"
)

type SessionCheckerServerImpl struct {
	SessionUseCase usecase.SessionUseCase
	handler.UnimplementedSessionCheckerServer
}

func CreateSessionCheckerServer(sessionUseCase usecase.SessionUseCase) handler.SessionCheckerServer {
	return &SessionCheckerServerImpl{SessionUseCase: sessionUseCase}
}

func (sessionChecker *SessionCheckerServerImpl) Create(ctx context.Context, in *handler.SessionInfo) (*handler.SessionID, error) {
	if in == nil {
		metrics.SessionHits.WithLabelValues("500", "create session; empty in").Inc()
		return &handler.SessionID{}, customErrors.ErrInternal
	}

	sessionID, err := sessionChecker.SessionUseCase.Create(in.UID)
	if err != nil {
		metrics.SessionHits.WithLabelValues("500", "create session; can't create").Inc()
		return &handler.SessionID{}, err
	}

	metrics.SessionHits.WithLabelValues("200", "create session").Inc()
	return &handler.SessionID{ID: sessionID}, nil
}

func (sessionChecker *SessionCheckerServerImpl) Get(ctx context.Context, in *handler.SessionID) (*handler.SessionInfo, error) {
	if in == nil {
		metrics.SessionHits.WithLabelValues("500", "get session; empty in").Inc()
		return &handler.SessionInfo{}, customErrors.ErrInternal
	}

	uid, err := sessionChecker.SessionUseCase.Get(in.ID)
	if err != nil {
		metrics.SessionHits.WithLabelValues("500", "get session; can't get").Inc()
		return &handler.SessionInfo{}, err
	}

	metrics.SessionHits.WithLabelValues("200", "get session").Inc()
	return &handler.SessionInfo{UID: uid}, nil
}

func (sessionChecker *SessionCheckerServerImpl) Delete(ctx context.Context, in *handler.SessionID) (*handler.Nothing, error) {
	if in == nil {
		metrics.SessionHits.WithLabelValues("500", "delete session; empty in").Inc()
		return &handler.Nothing{}, customErrors.ErrInternal
	}

	err := sessionChecker.SessionUseCase.Delete(in.ID)
	if err != nil {
		metrics.SessionHits.WithLabelValues("500", "delete session; can't delete").Inc()
		return &handler.Nothing{}, err
	}

	metrics.SessionHits.WithLabelValues("200", "delete session").Inc()
	return &handler.Nothing{}, nil
}
