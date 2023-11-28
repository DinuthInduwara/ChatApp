const { createApp, ref } = Vue

const EventNewMSG = "NEW_MSG_EVENT"
const EventNewUser = "NEW_USER_EVENT"
const EventUserDisconnect = "USER_DISCONNECT_EVENT"
const EventTyping = "USER_TYPING_EVENT"
const EventTypingOut = "USER_TYPING_OUT_EVENT"
const EventAllUsers = "ALL_USERS_EVENT"


class Event {
    constructor(type, payload, name) {
        this.type = type;
        this.name = name;
        this.payload = payload;
    }
}

class NewMessageEvent {
    constructor(message, to, from) {
        this.message = message;
        this.to = to
        this.time = Date.now()
        this.from = from
    }
}

class UserEvent {
    constructor(event, payload) {
        this.type = event
        this.payload = payload
    }
}

class NewUser {
    constructor(name, messages = [], typing = false) {
        this.name = name
        this.messages = messages
        this.typing = typing
    }
}


createApp({
    data() {
        return {
            name: "User🤠" + Math.floor(Math.random() * 10) || prompt("What is your Name", this.name),
            online: {},
            conn: null,
            selected: "",
            message: "",
            writing: false,
        }
    },
    async mounted() {
        this.conn = new WebSocket("ws://" + document.location.host + "/ws?name=" + this.name);
        this.conn.onopen = evt => console.log("Connected to Websocket: true")
        this.conn.onclose = evt => console.log("Connected to Websocket: false", evt)
        this.conn.onmessage = evt => this.handleEvent(evt)
        this.conn.addEventListener('open', () => this.conn.send(JSON.stringify(new UserEvent(EventAllUsers,))));
    },
    methods: {
        handleEvent(evt) {
            const eventData = JSON.parse(evt.data);

            const event = Object.assign(new Event, eventData);
            if (event.type === undefined) {
                alert("no 'type' field in event");
            }
            switch (event.type) {

                case EventNewMSG:
                    const messageEvent = Object.assign(new NewMessageEvent, event.payload);
                    if (eventData.payload.from === this.name) {
                        this.online[messageEvent.to].messages.push(messageEvent)
                        console.log(this.online[messageEvent.to])
                    } else {
                        this.online[messageEvent.from].messages.push(messageEvent)
                        console.log(this.online[messageEvent.from])
                    }
                    break;
                case EventNewUser:
                    if (eventData.name !== this.name) {
                        this.online[eventData.name] = new NewUser(eventData.name)
                    }
                    break;
                case EventUserDisconnect:
                    delete this.online[eventData.name]
                    break;
                case EventTyping:
                    if (this.name === eventData.payload.name) return;
                    this.online[eventData.payload.name].typing = true
                    break
                case EventTypingOut:
                    if (this.name === eventData.payload.name) return;
                    this.online[eventData.payload.name].typing = false
                    break
                case EventAllUsers:
                    if (eventData["users"]) {
                        Object.entries(eventData["users"]).forEach(([user, val]) => {
                            if (user !== this.name) {
                                this.online[user] = new NewUser(name)
                            }
                        })
                        this.selected = Object.keys(this.online)[0]
                    }
                    break

                default:
                    alert("unsupported message type");
                    break;
            }
        },
        selectChat() {
            if (!this.selected) {
                const key = Object.keys(this.online)[0]
                if (key.length > 0) {
                    this.selected = key[0];
                } else {
                    alert("No Chat Selected..!");
                    return;
                }
            }
        },
        sendMessage() {
            if (this.message.length === 0) return
            const msgEvent = new NewMessageEvent(this.message, this.selected, this.name)
            this.online[this.selected].messages.push(msgEvent)
            this.conn.send(JSON.stringify(new Event(EventNewMSG, msgEvent, name = this.name)))
            this.message = ""
        },
    },
    watch: {
        message(newMessage) {
            if (newMessage.length !== 0) {
                if (this.writing === false) {
                    this.conn.send(JSON.stringify(new UserEvent(EventTyping, { name: this.name })))
                    this.writing = true
                }
            } else {
                this.conn.send(JSON.stringify(new UserEvent(EventTypingOut, { name: this.name })))
                this.writing = false
            }
        },
        online() {
            if (Object.keys(this.online).length === 0) this.selected = ""
        },



    },

    computed: {}

}).mount('#app')
