//go:build e2e

package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	"github.com/lapitskyss/chat-service/internal/types"
	apiclientv1 "github.com/lapitskyss/chat-service/tests/e2e/api/client/v1"
	apimanagerv1 "github.com/lapitskyss/chat-service/tests/e2e/api/manager/v1"
	clientchat "github.com/lapitskyss/chat-service/tests/e2e/client-chat"
	managerworkspace "github.com/lapitskyss/chat-service/tests/e2e/manager-workspace"
	wsstream "github.com/lapitskyss/chat-service/tests/e2e/ws-stream"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

var (
	suiteCtx       context.Context
	suiteCtxCancel context.CancelFunc

	kc *keycloakclient.Client

	apiClientV1Endpoint string
	wsClientEndpoint    string
	wsClientOrigin      string
	wsClientSecProtocol string

	apiManagerV1Endpoint string
	wsManagerEndpoint    string
	wsManagerOrigin      string
	wsManagerSecProtocol string

	clientsPool  *usersPool
	managersPool *usersPool
)

var _ = BeforeSuite(func() {
	suiteCtx, suiteCtxCancel = context.WithCancel(context.Background())
	DeferCleanup(suiteCtxCancel)

	apiClientV1Endpoint = expectEnv("E2E_CLIENT_V1_API_ENDPOINT")
	wsClientEndpoint = expectEnv("E2E_CLIENT_WS_ENDPOINT")
	wsClientSecProtocol = expectEnv("E2E_CLIENT_WS_SEC_PROTOCOL")
	wsClientOrigin = expectEnv("E2E_CLIENT_WS_ORIGIN")

	apiManagerV1Endpoint = expectEnv("E2E_MANAGER_V1_API_ENDPOINT")
	wsManagerEndpoint = expectEnv("E2E_MANAGER_WS_ENDPOINT")
	wsManagerSecProtocol = expectEnv("E2E_MANAGER_WS_SEC_PROTOCOL")
	wsManagerOrigin = expectEnv("E2E_MANAGER_WS_ORIGIN")

	kcBasePath := expectEnv("E2E_KEYCLOAK_BASE_PATH")
	kcRealm := expectEnv("E2E_KEYCLOAK_REALM")
	kcClientID := expectEnv("E2E_KEYCLOAK_CLIENT_ID")
	kcClientSecret := expectEnv("E2E_KEYCLOAK_CLIENT_SECRET")
	kcClientDebug, _ := strconv.ParseBool(expectEnv("E2E_KEYCLOAK_CLIENT_DEBUG"))
	kcClients := expectEnv("E2E_KEYCLOAK_CLIENTS")   // "client1,client2,client3"
	kcManagers := expectEnv("E2E_KEYCLOAK_MANAGERS") // "manager1,manager2,manager3"

	var err error
	kc, err = keycloakclient.New(keycloakclient.NewOptions(
		kcBasePath,
		kcRealm,
		kcClientID,
		kcClientSecret,
		keycloakclient.WithDebugMode(kcClientDebug),
	))
	Expect(err).ShouldNot(HaveOccurred())

	clients, err := parseUsers(kcClients)
	Expect(err).ShouldNot(HaveOccurred())
	GinkgoWriter.Println("clients:", clients)
	clientsPool = newUsersPool(clients)

	managers, err := parseUsers(kcManagers)
	Expect(err).ShouldNot(HaveOccurred())
	GinkgoWriter.Println("managers:", managers)
	managersPool = newUsersPool(managers)
})

func expectEnv(k string) string {
	v := os.Getenv(k)
	Expect(v).NotTo(BeZero(), fmt.Sprintf("Please make sure %q is set correctly.", k))
	return v
}

func parseUsers(s string) ([]user, error) {
	userNames := strings.Split(s, ",")
	if len(userNames) == 0 {
		return nil, errors.New("no users specified")
	}

	known := make(map[string]struct{}, len(userNames))
	result := make([]user, 0, len(userNames))
	for _, uname := range userNames {
		if _, ok := known[uname]; ok {
			return nil, fmt.Errorf("duplicated user: %v", uname)
		}

		result = append(result, user{
			Name:     uname,
			Password: uname, // NOTE: E2E client username & password must be equal.
		})
		known[uname] = struct{}{}
	}
	return result, nil
}

func newClientChat(ctx context.Context, client user) *clientchat.Chat {
	apiClientV1, token := newClientAPI(ctx, client)

	var cl simpleClaims
	t, _, err := new(jwt.Parser).ParseUnverified(token, &cl)
	Expect(err).ShouldNot(HaveOccurred())

	clientID, err := types.Parse[types.UserID](t.Claims.(*simpleClaims).Subject)
	Expect(err).ShouldNot(HaveOccurred())
	GinkgoWriter.Printf("client %v has token sub %v\n", client.Name, clientID)

	clientChat, err := clientchat.New(clientchat.NewOptions(clientID, token, apiClientV1))
	Expect(err).ShouldNot(HaveOccurred())

	return clientChat
}

func newClientAPI(ctx context.Context, client user) (*apiclientv1.ClientWithResponses, string) {
	token, err := kc.Auth(ctx, client.Name, client.Password)
	Expect(err).ShouldNot(HaveOccurred())

	authorizator := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		return nil
	}
	apiClientV1, err := apiclientv1.NewClientWithResponses(
		apiClientV1Endpoint,
		apiclientv1.WithRequestEditorFn(authorizator),
	)
	Expect(err).ShouldNot(HaveOccurred())

	return apiClientV1, token.AccessToken
}

func newManagerWs(ctx context.Context, manager user) *managerworkspace.Workspace {
	apiManagerV1, token := newManagerAPI(ctx, manager)

	var cl simpleClaims
	t, _, err := new(jwt.Parser).ParseUnverified(token, &cl)
	Expect(err).ShouldNot(HaveOccurred())

	managerID, err := types.Parse[types.UserID](t.Claims.(*simpleClaims).Subject)
	Expect(err).ShouldNot(HaveOccurred())
	GinkgoWriter.Printf("manager %v has token sub %v\n", manager.Name, managerID)

	managerWs, err := managerworkspace.New(managerworkspace.NewOptions(managerID, token, apiManagerV1))
	Expect(err).ShouldNot(HaveOccurred())

	return managerWs
}

func newManagerAPI(ctx context.Context, manager user) (*apimanagerv1.ClientWithResponses, string) {
	token, err := kc.Auth(ctx, manager.Name, manager.Password)
	Expect(err).ShouldNot(HaveOccurred())

	authorizator := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		return nil
	}
	apiManagerV1, err := apimanagerv1.NewClientWithResponses(
		apiManagerV1Endpoint,
		apimanagerv1.WithRequestEditorFn(authorizator),
	)
	Expect(err).ShouldNot(HaveOccurred())

	return apiManagerV1, token.AccessToken
}

type user struct {
	Name     string
	Password string
}

type usersPool struct {
	users []user
	mu    sync.Mutex
}

func newUsersPool(users []user) *usersPool {
	p := &usersPool{
		users: make([]user, len(users)),
		mu:    sync.Mutex{},
	}
	copy(p.users, users)
	return p
}

func (p *usersPool) Get() user {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.users) == 0 {
		AbortSuite("there are no users in the pool - let's add a new one to Keycloak")
	}

	u := p.users[0]
	p.users = p.users[1:]

	GinkgoWriter.Printf("user %s removed from pool\n", u.Name)
	return u
}

type simpleClaims struct {
	Subject string `json:"sub,omitempty"`
}

func (sc simpleClaims) Valid() error {
	return nil
}

func waitForEvent(stream *wsstream.Stream) {
	select {
	case <-stream.EventSignals():
	case <-time.After(7 * time.Second):
		Fail("no expected event in the stream")
	}
}

func waitForOptionalEvent(stream *wsstream.Stream) {
	select {
	case <-stream.EventSignals():
	case <-time.After(7 * time.Second):
	}
}
