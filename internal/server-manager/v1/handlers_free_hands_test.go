package managerv1_test

import (
	"errors"
	"net/http"

	internalerrors "github.com/lapitskyss/chat-service/internal/errors"
	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
	"github.com/lapitskyss/chat-service/internal/types"
	canreceiveproblems "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/lapitskyss/chat-service/internal/usecases/manager/free-hands"
)

func (s *HandlersSuite) TestGetFreeHandsBtnAvailability_Usecase_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getFreeHandsBtnAvailability", "")
	s.canReceiveProblemsUseCase.EXPECT().Handle(eCtx.Request().Context(), canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(canreceiveproblems.Response{}, errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostGetFreeHandsBtnAvailability(eCtx, managerv1.PostGetFreeHandsBtnAvailabilityParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetFreeHandsBtnAvailability_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/getFreeHandsBtnAvailability", "")
	s.canReceiveProblemsUseCase.EXPECT().Handle(eCtx.Request().Context(), canreceiveproblems.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(canreceiveproblems.Response{Result: true}, nil)

	// Action.
	err := s.handlers.PostGetFreeHandsBtnAvailability(eCtx, managerv1.PostGetFreeHandsBtnAvailabilityParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
    "data":
    {
        "available": true
    }
}`, resp.Body.String())
}

func (s *HandlersSuite) TestFreeHands_Usecase_Error() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/freeHands", "")
	s.freeHandsUseCase.EXPECT().Handle(eCtx.Request().Context(), freehands.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(errors.New("something went wrong"))

	// Action.
	err := s.handlers.PostFreeHands(eCtx, managerv1.PostFreeHandsParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestGetHistory_Usecase_InvalidRequest() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/freeHands", "")
	s.freeHandsUseCase.EXPECT().Handle(eCtx.Request().Context(), freehands.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(freehands.ErrManagerOverloaded)

	// Action.
	err := s.handlers.PostFreeHands(eCtx, managerv1.PostFreeHandsParams{XRequestID: reqID})

	// Assert.
	s.Require().Error(err)
	s.EqualValues(managerv1.ErrorCodeManagerOverloaded, internalerrors.GetServerErrorCode(err))
	s.Empty(resp.Body)
}

func (s *HandlersSuite) TestFreeHands_Usecase_Success() {
	// Arrange.
	reqID := types.NewRequestID()
	resp, eCtx := s.newEchoCtx(reqID, "/v1/freeHands", "")
	s.freeHandsUseCase.EXPECT().Handle(eCtx.Request().Context(), freehands.Request{
		ID:        reqID,
		ManagerID: s.managerID,
	}).Return(nil)

	// Action.
	err := s.handlers.PostFreeHands(eCtx, managerv1.PostFreeHandsParams{XRequestID: reqID})

	// Assert.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.Code)
	s.JSONEq(`
{
    "data": null
}`, resp.Body.String())
}
