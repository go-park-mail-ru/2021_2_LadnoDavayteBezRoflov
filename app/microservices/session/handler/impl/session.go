package impl

import (
	"backendServer/app/microservices/session/handler"
	"backendServer/app/microservices/session/usecase"
	customErrors "backendServer/pkg/errors"

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
		return &handler.SessionID{}, customErrors.ErrInternal
	}

	sessionID, err := sessionChecker.SessionUseCase.Create(in.UID)
	if err != nil {
		return &handler.SessionID{}, err
	}
	return &handler.SessionID{ID: sessionID}, nil
}

func (sessionChecker *SessionCheckerServerImpl) Get(ctx context.Context, in *handler.SessionID) (*handler.SessionInfo, error) {
	if in == nil {
		return &handler.SessionInfo{}, customErrors.ErrInternal
	}

	uid, err := sessionChecker.SessionUseCase.Get(in.ID)
	if err != nil {
		return &handler.SessionInfo{}, err
	}
	return &handler.SessionInfo{UID: uid}, nil
}

func (sessionChecker *SessionCheckerServerImpl) Delete(ctx context.Context, in *handler.SessionID) (*handler.Nothing, error) {
	if in == nil {
		return &handler.Nothing{}, customErrors.ErrInternal
	}

	err := sessionChecker.SessionUseCase.Delete(in.ID)
	if err != nil {
		return &handler.Nothing{}, err
	}
	return &handler.Nothing{}, nil
}
