new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        email: null, // Email address used for grabbing an avatar
        pseudo: null, // Our username
        joined: false // True if email and username have been filled in
    },
    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {

            var msg = JSON.parse(e.data);

            console.log(msg.peer);
            console.log("My pseudo is : " + self.pseudo);
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
                console.log("Sending with pseudo : " + this.pseudo)
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
        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});
