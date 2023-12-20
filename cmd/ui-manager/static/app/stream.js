const clearTypingInterval = 900; //0.9 seconds
let clearTypingTimerId;

const eventHandlers = {
    'NewChatEvent': (event) => {
        if ($(`*[data-message-id="${event.chatId}"]`).length > 0) {
            return;
        }
        App.DisplayNewChat(event);

        if ($(App.problemSelector).length === 1) {
            $(`*[data-chat-id="${event.chatId}"]`).trigger('click');
        }

        App.readyToProblemsBtn.text('Ready to Problems 🙋‍');
        event.canTakeMoreProblems ?
            App.readyToProblemsBtn.removeClass('disabled') :
            App.readyToProblemsBtn.addClass('disabled');
    },

    'NewMessageEvent': (event) => {
        if ($(`*[data-message-id="${event.messageId}"]`).length > 0) {
            return;
        }
        if (App.currentChatID !== event.chatId) {
            // Do not display new messages from other chats.
            return;
        }
        App.DisplayNewMessage(event);
    },

    'ChatClosedEvent': (event) => {
        if (App.currentChatID === event.chatId) {
            App.currentChatID = undefined;
            App.currentClientID = undefined;
            App.chatArea.empty();
        }
        $(`*[data-chat-id="${event.chatId}"]`).remove();

        if (!App.readyToProblemsBtn.hasClass('waiting')) {
            event.canTakeMoreProblems ?
                App.readyToProblemsBtn.removeClass('disabled') :
                App.readyToProblemsBtn.addClass('disabled');
        }
    },


    'TypingEvent': (event) => {
        if (App.currentClientID !== event.clientId) {
            return;
        }
        $('#user-is-typing').html(event.clientId + ' is typing...');

        clearTimeout(clearTypingTimerId);
        clearTypingTimerId = setTimeout(function () {
            //clear user is typing message
            $('#user-is-typing').html('');
        }, clearTypingInterval);
    },
};

function initWsStream(token) {
    const sock = new WebSocket(wsEndpoint, [wsProtocol, token]);

    window.addEventListener('unload', function () {
        if (sock.readyState === WebSocket.OPEN) {
            sock.close();
        }
    });

    sock.onopen = function () {
        console.info('ws: connection established');
    };

    sock.onclose = function (event) {
        if (!event.wasClean) {
            console.error('ws: unexpected connection lost');
            console.error('code: ' + event.code + ', reason: ' + event.reason);
        }
    };

    sock.onerror = function (event) {
        console.error('ws: error: ' + JSON.stringify(event));

        // If error occurred then try to reconnect.
        (async () => {
            let promise = new Promise(resolve => setTimeout(resolve, 2000));

            await promise;

            initWsStream(token);
        })();
    };

    sock.onmessage = function (event) {
        console.info('ws: new event: ' + event.data);

        const payload = JSON.parse(event.data);
        const eventType = payload.eventType;

        if (!(eventType in eventHandlers)) {
            console.error('ws: unknown event: ' + eventType);
            return;
        }

        eventHandlers[eventType](payload);
    };

    return sock
}
