import Vue from "vue";
import axios from "axios";
import { setup } from "../game";

const host = "goba.holmes-dev.com"; // "localhost:5000";
console.log(host);

let app = new Vue({
  el: "#app",
  data: {
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
    updateTPS: function(message) {
      if (message !== this.tps) {
        this.tps = message;
      }
    },

    updateScores: function(scores) {
      let formatted = [];

      for (let [name, score] of Object.entries(scores)) {
        score.name = name;
        formatted.push(score);
      }
      formatted = formatted.sort((a, b) => {
        return a.kills < b.kills ? 1 : -1;
      });
      this.scores = formatted;
    },

    createGame: function() {
      let name = this.createName;

      axios.get(`https://${host}/create?name=${name}`).then((result) => {
        this.code = result.data.code;

        if (result.data.success) {
          this.createGameJoin(this.code, name);
        } else {
          this.createError = true;
          this.createErrorMessage = result.data.error;
        }
      });
    },

    createGameJoin: function(code, name) {
      let url = `wss://${host}/join?code=${code}&name=${name}`;
      let socket = new WebSocket(url);

      socket.onmessage = (message) => {
        let event = JSON.parse(message.data);

        if (event.data.success) {
          this.inGame = true;
          setup(socket, app);
        } else {
          this.createError = true;
          this.createErrorMessage = event.data.error;
          socket.close();
        }
      };
    },

    joinGame: function() {
      let url = `wss://${host}/join?code=${this.code}&name=${this.joinName}`;
      let socket = new WebSocket(url);

      socket.onmessage = (message) => {
        let event = JSON.parse(message.data);

        if (event.data.success) {
          this.inGame = true;
          setup(socket, app);
        } else {
          this.joinError = true;
          this.joinErrorMessage = event.data.error;
          socket.close();
        }
      };
    },

    refresh: function() {
      location.reload();
    },
  },
});
