import Vue from 'vue';
import axios from 'axios';
import { setup } from "../game";

const host = "goba-env.eba-hiw6diij.ca-central-1.elasticbeanstalk.com"; // "localhost:5000"; 
console.log(host);

let app = new Vue({
    el: "#app",
    data: {
        title: "GoBA - Go Online Battle Arena",
        tps: "TPS: 0",
        inGame: false,

        scores: [],

        code: "",

        createName: "",
        createError: "",
        createErrorMessage: "",

        joinName: "",
        joinError: false,
        joinErrorMessage: false,
    },

    methods: {
        updateTPS: function (message) {
            if (message !== this.tps) {
                this.tps = message;
            }
        },

        updateScores: function (scores) {
            console.log(scores);
            let formatted = [];
            for (let [name, score] of Object.entries(scores)) {
                score.name = name;
                formatted.push(score);
            }
            this.scores = formatted;
        },

        createGame: function () {
            let name = this.createName;

            axios.get(`http://${host}/create?name=${name}`).then(result => {
                this.code = result.data.code;

                if (result.data.success) {
                    this.createGameJoin(this.code, name);
                } else {
                    this.createError = true;
                    this.createErrorMessage = result.data.error;
                }
            });
        },

        createGameJoin: function (code, name) {
            let url = `ws://${host}/join?code=${code}&name=${name}`;
            let socket = new WebSocket(url);

            socket.onmessage = (message) => {
                let event = JSON.parse(message.data);

                if (event.data.success) {
                    this.inGame = true;
                    setup(socket, app);
                }

                else {
                    this.createError = true;
                    this.createErrorMessage = event.data.error;
                    socket.close();
                }
            }
        },

        joinGame: function () {
            let url = `ws://${host}/join?code=${this.code}&name=${this.joinName}`;
            let socket = new WebSocket(url);

            socket.onmessage = (message) => {
                let event = JSON.parse(message.data);

                if (event.data.success) {
                    this.inGame = true;
                    setup(socket, app);
                }

                else {
                    this.joinError = true;
                    this.joinErrorMessage = event.data.error;
                    socket.close();
                }
            }
        }
    }
});