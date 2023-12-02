//go:build e2e

package e2e_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	clientchat "github.com/lapitskyss/chat-service/tests/e2e/client-chat"
	wsstream "github.com/lapitskyss/chat-service/tests/e2e/ws-stream"
)

var _ = Describe("Client Events Smoke", Ordered, func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc

		clientChat        *clientchat.Chat
		clientStream      *wsstream.Stream
		clientStreamErrCh = make(chan error, 1)
	)

	BeforeAll(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Setup client.
		clientChat = newClientChat(ctx, clientsPool.Get())

		var err error
		clientStream, err = wsstream.New(wsstream.NewOptions(
			wsClientEndpoint,
			wsClientOrigin,
			wsClientSecProtocol,
			clientChat.AccessToken(),
			clientChat.HandleEvent,
		))
		Expect(err).ShouldNot(HaveOccurred())
		go func() { clientStreamErrCh <- clientStream.Run(ctx) }()
	})

	AfterAll(func() {
		cancel()
		Expect(<-clientStreamErrCh).ShouldNot(HaveOccurred())
	})

	It("client message was sent to manager", func() {
		err := clientChat.SendMessage(ctx, "Hello, sir!")
		Expect(err).ShouldNot(HaveOccurred())

		waitForEvent(clientStream) // NewMessageEvent.
		waitForEvent(clientStream) // MessageSentEvent.

		msg, ok := clientChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(msg.IsReceived).Should(BeTrue())
		Expect(msg.IsBlocked).Should(BeFalse())
	})

	It("client message was blocked", func() {
		err := clientChat.SendMessage(ctx, "My CVC is 678")
		Expect(err).ShouldNot(HaveOccurred())

		waitForEvent(clientStream) // NewMessageEvent.
		waitForEvent(clientStream) // MessageBlockedEvent.

		msg, ok := clientChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(msg.IsReceived).Should(BeFalse())
		Expect(msg.IsBlocked).Should(BeTrue())
	})
})
