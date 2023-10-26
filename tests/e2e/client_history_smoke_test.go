//go:build e2e

package e2e_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lapitskyss/chat-service/internal/types"
	clientchat "github.com/lapitskyss/chat-service/tests/e2e/client-chat"
)

var _ = Describe("Client History Smoke", Ordered, func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc

		clientChat *clientchat.Chat
	)

	BeforeAll(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Setup client.
		clientChat = newClientChat(ctx, clientsPool.Get())
	})

	AfterAll(func() {
		cancel()
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
		Expect(msg.AuthorID).Should(Equal(clientChat.ClientID()))
	})

	secondMsgReqID := types.NewRequestID()

	It("client send the 2nd message without errors", func() {
		err := clientChat.SendMessage(ctx, "I have a problem :(", clientchat.WithRequestID(secondMsgReqID))
		Expect(err).ShouldNot(HaveOccurred())

		msg, ok := clientChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(msg.AuthorID).Should(Equal(clientChat.ClientID()))
	})

	It("client retries the 2nd message without errors", func() {
		err := clientChat.SendMessage(ctx, "I have a problem :(", clientchat.WithRequestID(secondMsgReqID))
		Expect(err).ShouldNot(HaveOccurred())

		msg, ok := clientChat.LastMessage()
		Expect(ok).Should(BeTrue())
		Expect(msg.AuthorID).Should(Equal(clientChat.ClientID()))
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
			Expect(m.AuthorID).Should(Equal(clientChat.ClientID()))
		}
	})
})
