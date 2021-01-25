import Vue from 'vue'
import { setup } from "../game";

let app = new Vue({
    el: "#app",
    data: {
        title: "GoBA - Go Online Battle Arena",
        tps: "TPS: 0",
        teams: {}
    },

    methods: {
        updateTPS: function (message) {
            if (message !== this.tps) {
                this.tps = message;
            }
        },

        updateTeams: function (teams) {
            this.teams = teams;
        }
    }
});

setup("ws://localhost:5000/ws", "canvas", app);