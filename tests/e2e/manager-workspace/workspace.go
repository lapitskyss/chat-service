package managerworkspace

import (
	"container/list"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/onsi/ginkgo/v2"
	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"

	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/lapitskyss/chat-service/pkg/pointer"
	apimanagerevents "github.com/lapitskyss/chat-service/tests/e2e/api/manager/events"
	apimanagerv1 "github.com/lapitskyss/chat-service/tests/e2e/api/manager/v1"
)

const pageSize = 10

var (
	errUnknownChat = errors.New("unknown chat")

	errNoResponseBody   = errors.New("no response body")
	errNoDataInResponse = errors.New("no data field in response")
)

//go:generate options-gen -out-filename=workspace_options.gen.go -from-struct=Options
type Options struct {
	id    types.UserID                      `option:"mandatory" validate:"required"`
	token string                            `option:"mandatory" validate:"required"`
	api   *apimanagerv1.ClientWithResponses `option:"mandatory" validate:"required"`
}

type Workspace struct {
	Options

	chatsMu   *sync.RWMutex
	chatsByID map[types.ChatID]*Chat
	chats     *list.List

	canTakeMoreProblems atomic.Bool
}

func New(opts Options) (*Workspace, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	return &Workspace{
		Options:   opts,
		chatsMu:   new(sync.RWMutex),
		chatsByID: make(map[types.ChatID]*Chat),
		chats:     list.New(),
	}, nil
}

func (ws *Workspace) ManagerID() types.UserID {
	return ws.id
}

func (ws *Workspace) AccessToken() string {
	return ws.token
}

func (ws *Workspace) FirstChat() (Chat, bool) {
	ws.chatsMu.RLock()
	defer ws.chatsMu.RUnlock()

	if ws.chats.Len() == 0 {
		return Chat{}, false
	}
	return *ws.chats.Front().Value.(*Chat), true
}

func (ws *Workspace) LastChat() (Chat, bool) {
	ws.chatsMu.RLock()
	defer ws.chatsMu.RUnlock()

	if ws.chats.Len() == 0 {
		return Chat{}, false
	}
	return *ws.chats.Back().Value.(*Chat), true
}

func (ws *Workspace) Chats() []Chat {
	ws.chatsMu.RLock()
	defer ws.chatsMu.RUnlock()

	result := make([]Chat, 0, ws.chats.Len())
	for item := ws.chats.Front(); item != nil; item = item.Next() {
		result = append(result, *(item.Value.(*Chat)))
	}
	return result
}

func (ws *Workspace) ChatsCount() int {
	ws.chatsMu.RLock()
	defer ws.chatsMu.RUnlock()

	return ws.chats.Len()
}

// Refresh emulates the browser page reloading.
func (ws *Workspace) Refresh(ctx context.Context) error {
	ws.chatsMu.Lock()
	{
		ws.chats.Init()
		ws.chatsByID = make(map[types.ChatID]*Chat)
	}
	ws.chatsMu.Unlock()

	wg, gCtx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return ws.ReceiveNewProblemsAvailability(gCtx)
	})
	wg.Go(func() error {
		return ws.GetChats(gCtx)
	})
	if err := wg.Wait(); err != nil {
		return err
	}

	if first, ok := ws.FirstChat(); ok {
		return ws.GetChatHistory(ctx, first.ID)
	}
	return nil
}

func (ws *Workspace) CanTakeMoreProblems() bool {
	return ws.canTakeMoreProblems.Load()
}

func (ws *Workspace) GetChats(ctx context.Context) error {
	resp, err := ws.api.PostGetChatsWithResponse(ctx,
		&apimanagerv1.PostGetChatsParams{XRequestID: types.NewRequestID()},
	)
	if err != nil {
		return fmt.Errorf("post request: %v", err)
	}
	if resp.JSON200 == nil {
		return errNoResponseBody
	}
	if err := resp.JSON200.Error; err != nil {
		return fmt.Errorf("%v: %v", err.Code, err.Message)
	}

	data := resp.JSON200.Data
	if data == nil {
		return errNoDataInResponse
	}

	for _, chatItem := range data.Chats {
		ws.appendChat(chatItem.ChatId, chatItem.ClientId)
	}
	return nil
}

func (ws *Workspace) GetChatHistory(ctx context.Context, chatID types.ChatID) error {
	chatItem, ok := ws.getChat(chatID)
	if !ok {
		return fmt.Errorf("%v: %v", errUnknownChat, chatID)
	}

	resp, err := ws.api.PostGetChatHistoryWithResponse(ctx,
		&apimanagerv1.PostGetChatHistoryParams{XRequestID: types.NewRequestID()},
		apimanagerv1.PostGetChatHistoryJSONRequestBody{
			ChatId:   chatID,
			Cursor:   pointer.Ptr(chatItem.cursor),
			PageSize: pointer.Ptr(pageSize),
		},
	)
	if err != nil {
		return fmt.Errorf("post request: %v", err)
	}
	if resp.JSON200 == nil {
		return errNoResponseBody
	}
	if err := resp.JSON200.Error; err != nil {
		return fmt.Errorf("%v: %v", err.Code, err.Message)
	}

	data := resp.JSON200.Data
	if data == nil {
		return errNoDataInResponse
	}

	for _, m := range data.Messages {
		chatItem.pushToFront(NewMessage(
			m.Id,
			chatID,
			m.AuthorId,
			m.Body,
			m.CreatedAt,
		))
	}
	chatItem.cursor = data.Next

	return nil
}

func (ws *Workspace) ReceiveNewProblemsAvailability(ctx context.Context) error {
	resp, err := ws.api.PostGetFreeHandsBtnAvailabilityWithResponse(ctx,
		&apimanagerv1.PostGetFreeHandsBtnAvailabilityParams{XRequestID: types.NewRequestID()},
	)
	if err != nil {
		return fmt.Errorf("post request: %v", err)
	}
	if resp.JSON200 == nil {
		return errNoResponseBody
	}
	if err := resp.JSON200.Error; err != nil {
		return fmt.Errorf("%v: %v", err.Code, err.Message)
	}

	data := resp.JSON200.Data
	if data == nil {
		return errNoDataInResponse
	}
	if err := resp.JSON200.Error; err != nil {
		return fmt.Errorf("%v: %v", err.Code, err.Message)
	}

	ws.setCanTakeMoreProblemsFlag(data.Available)
	return nil
}

func (ws *Workspace) ReadyToNewProblems(ctx context.Context) error {
	if !ws.CanTakeMoreProblems() {
		return fmt.Errorf("manager %s cannot receive new problems", ws.id)
	}

	resp, err := ws.api.PostFreeHandsWithResponse(ctx,
		&apimanagerv1.PostFreeHandsParams{XRequestID: types.NewRequestID()},
	)
	if err != nil {
		return fmt.Errorf("post request: %v", err)
	}
	if resp.JSON200 != nil {
		if err := resp.JSON200.Error; err != nil {
			return fmt.Errorf("%v: %v", err.Code, err.Message)
		}
	}

	return nil
}

func (ws *Workspace) HandleEvent(_ context.Context, data []byte) error {
	ginkgo.GinkgoWriter.Printf("manager %s workspace: new event: %s\n", ws.id, string(data))

	if len(data) == 0 {
		return nil
	}

	var event apimanagerevents.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("unmarshal event: %v", err)
	}

	v, err := event.ValueByDiscriminator()
	if err != nil {
		return fmt.Errorf("unmarshal sub-event: %v", err)
	}

	switch vv := v.(type) {
	case apimanagerevents.NewChatEvent:
		ws.setCanTakeMoreProblemsFlag(vv.CanTakeMoreProblems)
		ws.appendChat(vv.ChatId, vv.ClientId)

	case apimanagerevents.NewMessageEvent:
		ws.pushMessageToBack(NewMessage(
			vv.MessageId,
			vv.ChatId,
			vv.AuthorId,
			vv.Body,
			vv.CreatedAt,
		))
	}

	return nil
}

func (ws *Workspace) setCanTakeMoreProblemsFlag(v bool) {
	ws.canTakeMoreProblems.Store(v)
}

func (ws *Workspace) appendChat(chatID types.ChatID, clientID types.UserID) {
	ws.chatsMu.Lock()
	defer ws.chatsMu.Unlock()

	if _, ok := ws.chatsByID[chatID]; !ok {
		chatItem := NewChat(chatID, clientID)
		chatItem.listItemRef = ws.chats.PushBack(chatItem)
		ws.chatsByID[chatItem.ID] = chatItem
	}
}

// func (ws *Workspace) removeChat(id types.ChatID) error {
//	if _, ok := ws.getChat(id); !ok {
//		return fmt.Errorf("%v: %v", errUnknownChat, id)
//	}
//
//	ws.chatsMu.Lock()
//	defer ws.chatsMu.Unlock()
//
//	item := ws.chatsByID[id]
//	delete(ws.chatsByID, id)
//	ws.chats.Remove(item.listItemRef)
//
//	return nil
//}

func (ws *Workspace) getChat(id types.ChatID) (*Chat, bool) {
	ws.chatsMu.RLock()
	defer ws.chatsMu.RUnlock()

	item, ok := ws.chatsByID[id]
	return item, ok
}

func (ws *Workspace) pushMessageToBack(msg *Message) {
	ws.chatsMu.RLock()
	defer ws.chatsMu.RUnlock()

	chatID := msg.ChatID

	if _, ok := ws.chatsByID[chatID]; !ok {
		ginkgo.GinkgoWriter.Printf(
			"manager %s workspace: skip message %s for chat %s because of no chat yet\n",
			ws.id, msg.ID, chatID,
		)
	} else {
		ws.chatsByID[chatID].pushToBack(msg)
	}
}
