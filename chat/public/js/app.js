new Vue({
    el: '#app',

    data: {
        ws: null, // websocket to server
        newMsg: '', // New message input field
        chatContent: '', // List of all messages
        pseudo: null, // name of current user
        joined: false // Indicates user has entered pseudo
    },
    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {

            var msg = JSON.parse(e.data);
            if (msg.peer == self.pseudo){
                self.chatContent += '<div class="chip teal" >'
                    + msg.peer
                    + '</div>'
                    + emojione.toImage(msg.message) + '<br/>'; // Parse emojis
            }
            else
            {
                self.chatContent += '<div class="chip">'
                    + msg.peer
                    + '</div>'
                    + emojione.toImage(msg.message) + '<br/>'; // Parse emojis
            }

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom

        });
    },
    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                            message: $('<p>').html(this.newMsg).text(),
                            peer: this.pseudo,// Strip out html
                        }
                    ));
                this.newMsg = ''; // Reset newMsg
            }
        },
        set_pseudo: function () {
            if (!this.pseudo) {
                Materialize.toast('You must enter a name', 2000);
                return
            }
            this.pseudo = $('<p>').html(this.pseudo).text();
            this.joined = true;
        },
        disconnect: function() {
            $.ajax({
                type: "GET",
                url: "http://localhost:8000/disconnect",
                success: function(data){
                    // Notify of disconnection
                }
            })
        }
    }
});
