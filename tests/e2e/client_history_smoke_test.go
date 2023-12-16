//go:build e2e

package e2e_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lapitskyss/chat-service/internal/types"
	clientchat "github.com/lapitskyss/chat-service/tests/e2e/client-chat"
	managerworkspace "github.com/lapitskyss/chat-service/tests/e2e/manager-workspace"
	wsstream "github.com/lapitskyss/chat-service/tests/e2e/ws-stream"
)

var _ = Describe("Client History Smoke", Ordered, func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc

		clientChat *clientchat.Chat

		managerWs          *managerworkspace.Workspace
		managerStream      *wsstream.Stream
		managerStreamErrCh = make(chan error, 1)
	)

	BeforeAll(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Setup client.
		clientChat = newClientChat(ctx, clientsPool.Get())

		// Setup manager.
		managerWs = newManagerWs(ctx, managersPool.Get())

		var err error
		managerStream, err = wsstream.New(wsstream.NewOptions(
			wsManagerEndpoint,
			wsManagerOrigin,
			wsManagerSecProtocol,
			managerWs.AccessToken(),
			managerWs.HandleEvent,
		))
		Expect(err).ShouldNot(HaveOccurred())
		go func() { managerStreamErrCh <- managerStream.Run(ctx) }()
	})

	AfterAll(func() {
		cancel()
		Expect(<-managerStreamErrCh).ShouldNot(HaveOccurred())
	})

	It("no chat messages at the start of communication", func() {
		err := clientChat.Refresh(ctx)
		Expect(err).ShouldNot(HaveOccurred())

		n := clientChat.MessagesCount()
		Expect(n).Should(Equal(0))
	})

	It("client send the 1st message without errors", func() {
		err := clientChat.SendMessage(ctx, "Hello!")
		Expect(err).ShouldNot(HaveOccurred())

		msg, ok := clientChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(msg.AuthorID.String()).Should(Equal(clientChat.ClientID().String()))
	})

	secondMsgReqID := types.NewRequestID()

	It("client send the 2nd message without errors", func() {
		err := clientChat.SendMessage(ctx, "I have a problem :(", clientchat.WithRequestID(secondMsgReqID))
		Expect(err).ShouldNot(HaveOccurred())

		msg, ok := clientChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(msg.AuthorID.String()).Should(Equal(clientChat.ClientID().String()))
	})

	It("client retries the 2nd message without errors", func() {
		err := clientChat.SendMessage(ctx, "I have a problem :(", clientchat.WithRequestID(secondMsgReqID))
		Expect(err).ShouldNot(HaveOccurred())

		msg, ok := clientChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(msg.AuthorID.String()).Should(Equal(clientChat.ClientID().String()))
	})

	It("we still have two messages in the chat", func() {
		n := clientChat.MessagesCount()
		Expect(n).Should(Equal(2))
	})

	It("we have two messages in the history", func() {
		sentMessages := clientChat.Messages()

		err := clientChat.Refresh(ctx)
		Expect(err).ShouldNot(HaveOccurred())

		for i := 0; i < 3; i++ {
			err := clientChat.GetHistory(ctx)
			Expect(err).ShouldNot(HaveOccurred())
		}
		history := clientChat.Messages()
		Expect(sentMessages, Equal(history))

		for _, m := range history {
			Expect(m.AuthorID.String()).Should(Equal(clientChat.ClientID().String()))
		}
	})

	It("some garbage collection: assign chat to manager", func() {
		err := managerWs.Refresh(ctx)
		Expect(err).ShouldNot(HaveOccurred())

		err = managerWs.ReadyToNewProblems(ctx)
		Expect(err).ShouldNot(HaveOccurred())

		waitForEvent(managerStream) // NewChatEvent.
	})
})
