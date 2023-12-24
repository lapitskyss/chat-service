//go:build e2e

package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/lapitskyss/chat-service/internal/types"
	apiclientv1 "github.com/lapitskyss/chat-service/tests/e2e/api/client/v1"
	apimanagerv1 "github.com/lapitskyss/chat-service/tests/e2e/api/manager/v1"
)

var _ = Describe("Error Responses", Ordered, func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc

		apiClientV1  *apiclientv1.ClientWithResponses
		apiManagerV1 *apimanagerv1.ClientWithResponses
	)

	BeforeAll(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		apiClientV1, _ = newClientAPI(ctx, clientsPool.Get())
		apiManagerV1, _ = newManagerAPI(ctx, managersPool.Get())
	})

	AfterAll(func() {
		cancel()
	})

	It("401 unauthorized", func() {
		// Arrange.
		invalidAuthorizator := func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer invalid_token")
			return nil
		}
		apiInvalidAccessToken, err := apiclientv1.NewClientWithResponses(
			apiClientV1Endpoint,
			apiclientv1.WithRequestEditorFn(invalidAuthorizator),
		)
		Expect(err).ShouldNot(HaveOccurred())

		// Action.
		resp, err := apiInvalidAccessToken.PostSendMessageWithResponse(ctx,
			&apiclientv1.PostSendMessageParams{XRequestID: types.NewRequestID()},
			apiclientv1.SendMessageRequest{},
		)

		// Assert.
		Expect(err).ShouldNot(HaveOccurred())
		expectSendClientMsgRespCode(resp, http.StatusUnauthorized)
	})

	It("413 bad request, too long message body", func() {
		// Arrange.
		tooLongMsgBody := strings.Repeat("a", 100_000)

		// Action.
		resp, err := apiClientV1.PostSendMessageWithResponse(ctx,
			&apiclientv1.PostSendMessageParams{XRequestID: types.NewRequestID()},
			apiclientv1.PostSendMessageJSONRequestBody{MessageBody: tooLongMsgBody},
		)

		// Assert.
		Expect(err).ShouldNot(HaveOccurred())
		expectSendClientMsgRespCode(resp, http.StatusRequestEntityTooLarge)
	})

	It("400 bad request, message body is empty", func() {
		// Action.
		resp, err := apiClientV1.PostSendMessageWithResponse(ctx,
			&apiclientv1.PostSendMessageParams{XRequestID: types.NewRequestID()},
			apiclientv1.PostSendMessageJSONRequestBody{MessageBody: ""},
		)
		Expect(err).ShouldNot(HaveOccurred())

		// Assert.
		Expect(err).ShouldNot(HaveOccurred())
		expectSendClientMsgRespCode(resp, http.StatusBadRequest)
	})

	It("5001 code no active problem in chat", func() {
		resp, err := apiManagerV1.PostCloseChatWithResponse(
			ctx,
			&apimanagerv1.PostCloseChatParams{XRequestID: types.NewRequestID()},
			apimanagerv1.PostCloseChatJSONRequestBody{ChatId: types.NewChatID()},
		)
		Expect(err).ShouldNot(HaveOccurred())

		// Assert.
		expectCloseChatRespCode(resp, apimanagerv1.ErrorCodeNoActiveProblemInChat)
	})
})

func expectSendClientMsgRespCode[TCode ~int](resp *apiclientv1.PostSendMessageResponse, code TCode) {
	printResponse(resp.Body)
	Expect(resp).ShouldNot(BeNil())
	Expect(resp.JSON200).ShouldNot(BeNil())
	Expect(resp.JSON200.Error).ShouldNot(BeNil())
	Expect(resp.JSON200.Error.Code).Should(BeNumerically("==", code))
}

func expectCloseChatRespCode[TCode ~int](resp *apimanagerv1.PostCloseChatResponse, code TCode) {
	printResponse(resp.Body)
	Expect(resp).ShouldNot(BeNil())
	Expect(resp.JSON200).ShouldNot(BeNil())
	Expect(resp.JSON200.Error).ShouldNot(BeNil())
	Expect(resp.JSON200.Error.Code).Should(BeNumerically("==", code))
}

// printResponse helps to investigate server response in console.
func printResponse(rawBody []byte) {
	b, err := json.MarshalIndent(json.RawMessage(rawBody), "", "\t")
	Expect(err).ShouldNot(HaveOccurred())
	GinkgoWriter.Println(string(b))
}
