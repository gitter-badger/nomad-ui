const websocket = (function() {
    var socket = null;

    const onMessage = (store) => event => {
        store.dispatch(event)
    }

    const connect = (store) => {
        socket = new WebSocket(`${location.protocol}//${location.hostname}:${location.port}/ws`)
        socket.onmessage = onMessage(store);
    }

    return store => next => action => {
        if (socket == null || socket.readyState === socket.CLOSED) {
            connect(store)
        }
    }
})();

export default websocket
