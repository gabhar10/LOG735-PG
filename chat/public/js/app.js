new Vue({
    el: '#app',

    data: {
        ws: null, // websocket to server
        newMsg: '', // New message input field
        chatContent: '', // List of all messages
        pseudo: null, // name of current user
        joined: false, // Indicates user has entered pseudo
        connected: true,
        anchor: ''
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
            self = this;
            console.log(self.connected);
            return $.ajax({
                async: false,
                cache:false,
                type: "GET",
                url: "http://localhost:" + location.port + "/disconnect",
                success: function(data){
                    alert(data);
                    self.connected = false;
                },
                error: function(XMLHttpRequest, textStatus, errorThrown) {
                    alert("Status: " + textStatus); alert("Error: " + errorThrown);
                }
            });

        },
        connect: function(){
            return $.ajax({
                async: false,
                cache:false,
                type: "POST",
                url: "http://localhost:" + location.port + "/connect",
                data: {anchor: this.anchor},
                success: function(data){
                    alert(data);
                    this.connected = true;
                }.bind(this)
            })
        }
    }
});
