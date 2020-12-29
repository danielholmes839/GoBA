/*
 * ATTENTION: The "eval" devtool has been used (maybe by default in mode: "development").
 * This devtool is not neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
/******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ({

/***/ "./src/events.ts":
/*!***********************!*\
  !*** ./src/events.ts ***!
  \***********************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"TickCounter\": () => /* binding */ TickCounter\n/* harmony export */ });\nvar TickCounter = /** @class */ (function () {\r\n    function TickCounter(id) {\r\n        this.ticks = 0;\r\n        this.timestamp = 0;\r\n        this.paragraph = document.getElementById(id);\r\n    }\r\n    TickCounter.prototype.update = function (event) {\r\n        if (event.timestamp != this.timestamp) {\r\n            this.paragraph.innerHTML = \"TPS: \" + this.ticks;\r\n            this.timestamp = event.timestamp;\r\n            this.ticks = 0;\r\n        }\r\n        this.ticks++;\r\n    };\r\n    return TickCounter;\r\n}());\r\n\r\n\n\n//# sourceURL=webpack://GoBA-client/./src/events.ts?");

/***/ }),

/***/ "./src/game/index.ts":
/*!***************************!*\
  !*** ./src/game/index.ts ***!
  \***************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"Champion\": () => /* binding */ Champion,\n/* harmony export */   \"Game\": () => /* binding */ Game\n/* harmony export */ });\nvar Screen = /** @class */ (function () {\r\n    function Screen(canvas) {\r\n        this.canvas = canvas;\r\n        this.ctx = this.canvas.getContext(\"2d\", { alpha: false });\r\n        this.centerX = this.canvas.width / 2;\r\n        this.centerY = this.canvas.height / 2;\r\n        this.zoom = 2;\r\n        this.dx = 100;\r\n        this.dy = 100;\r\n    }\r\n    Screen.prototype.clear = function () {\r\n        this.ctx.fillStyle = \"#ffffff\";\r\n        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);\r\n    };\r\n    Screen.prototype.center = function (x, y) {\r\n        this.dx = x;\r\n        this.dy = y;\r\n    };\r\n    /* Game position to canvas position */\r\n    Screen.prototype.transformGamePosition = function (x, y) {\r\n        x = this.zoom * (x - this.dx) + this.centerX;\r\n        y = this.zoom * (y - this.dy) + this.centerY;\r\n        return [x, y];\r\n    };\r\n    /* Canvas position to game position */\r\n    Screen.prototype.transformCanvasPosition = function (x, y) {\r\n        x = Math.round((x - this.centerX) / this.zoom) + this.dx; //+ this.centerX;\r\n        y = Math.round((y - this.centerY) / this.zoom) + this.dy; //+ this.centerY;\r\n        return [x, y];\r\n    };\r\n    /* Draw rect with game position and scale */\r\n    Screen.prototype.drawRect = function (x, y, w, h, color) {\r\n        var _a;\r\n        _a = this.transformGamePosition(x, y), x = _a[0], y = _a[1];\r\n        w *= this.zoom;\r\n        h *= this.zoom;\r\n        this.ctx.strokeStyle = color;\r\n        this.ctx.fillStyle = color;\r\n        this.ctx.fillRect(x, y, w, h);\r\n    };\r\n    /* Draw circle with game position and scale */\r\n    Screen.prototype.drawCircle = function (x, y, r, color) {\r\n        var _a;\r\n        _a = this.transformGamePosition(x, y), x = _a[0], y = _a[1];\r\n        r *= this.zoom;\r\n        this.ctx.strokeStyle = color;\r\n        this.ctx.fillStyle = color;\r\n        this.ctx.beginPath();\r\n        this.ctx.arc(x, y, r, 0, 2 * Math.PI);\r\n        this.ctx.stroke();\r\n        this.ctx.fill();\r\n    };\r\n    /* Draw text with game position */\r\n    Screen.prototype.drawText = function (x, y, text) {\r\n        var _a;\r\n        _a = this.transformGamePosition(x, y), x = _a[0], y = _a[1];\r\n        var size = Math.round(18 * this.zoom);\r\n        this.ctx.font = size + \"px 'Courier New', monospace\";\r\n        this.ctx.fillStyle = \"#000000\";\r\n        this.ctx.fillText(text, x, y);\r\n    };\r\n    return Screen;\r\n}());\r\nvar Champion = /** @class */ (function () {\r\n    function Champion(id, health, x, y) {\r\n        this.r = 50;\r\n        this.id = id;\r\n        this.x = x;\r\n        this.y = y;\r\n        this.health = health;\r\n        this.maxHealth = 100;\r\n    }\r\n    Champion.prototype.draw = function (screen, isAlly, isClient, teamColor) {\r\n        screen.drawCircle(this.x, this.y, this.r, teamColor);\r\n        var color;\r\n        if (isClient) {\r\n            color = Champion.clientHealth;\r\n        }\r\n        else if (isAlly) {\r\n            color = Champion.allyHealth;\r\n        }\r\n        else {\r\n            color = Champion.enemyHealth;\r\n        }\r\n        screen.drawRect(this.x - this.r, this.y - (this.r + 30), this.r * 2.5, 20, color);\r\n        screen.drawText(this.x - this.r, this.y - (this.r + 40), this.id);\r\n    };\r\n    Champion.allyHealth = \"#00ff00\";\r\n    Champion.enemyHealth = \"#ff0000\";\r\n    Champion.clientHealth = \"#ffff00\";\r\n    return Champion;\r\n}());\r\n\r\nvar Game = /** @class */ (function () {\r\n    function Game(setup, canvas, socket) {\r\n        this.canvas = canvas;\r\n        this.socket = socket;\r\n        this.screen = new Screen(canvas);\r\n        this.client = setup.id;\r\n        this.walls = setup.walls;\r\n        this.champions = [];\r\n        this.teams = {};\r\n        this.clients = {};\r\n        this.setOnClicks();\r\n    }\r\n    Game.prototype.tick = function (t) {\r\n        var _this = this;\r\n        this.champions = t.champions.map(function (c) {\r\n            if (_this.championIsClient(c)) {\r\n                _this.screen.center(c.x, c.y);\r\n            }\r\n            return Object.assign(new Champion(\"\", 0, 0, 0), c);\r\n        });\r\n        this.draw();\r\n    };\r\n    Game.prototype.updateTeams = function (_a) {\r\n        var teams = _a.teams, clients = _a.clients;\r\n        this.teams = teams;\r\n        this.clients = clients;\r\n    };\r\n    Game.prototype.draw = function () {\r\n        var _this = this;\r\n        this.screen.clear();\r\n        this.walls.forEach(function (_a) {\r\n            var x = _a.x, y = _a.y, w = _a.w, h = _a.h;\r\n            _this.screen.drawRect(x, y, w, h, \"#dddddd\");\r\n        });\r\n        this.champions.forEach(function (champion) {\r\n            var isAlly = _this.championIsAlly(champion);\r\n            var isClient = _this.championIsClient(champion);\r\n            var teamColor = _this.championGetTeamColor(champion);\r\n            champion.draw(_this.screen, isAlly, isClient, teamColor);\r\n        });\r\n    };\r\n    Game.prototype.championIsClient = function (champion) {\r\n        return champion.id === this.client;\r\n    };\r\n    Game.prototype.championIsAlly = function (champion) {\r\n        // Matching teams\r\n        return this.teams[this.clients[champion.id]].color === this.teams[this.clients[this.client]].color;\r\n    };\r\n    Game.prototype.championGetTeamColor = function (champion) {\r\n        // Matching teams\r\n        return this.teams[this.clients[champion.id]].color;\r\n    };\r\n    Game.prototype.setOnClicks = function () {\r\n        var _this = this;\r\n        this.canvas.addEventListener('contextmenu', function (e) {\r\n            e.preventDefault();\r\n        }, false);\r\n        this.canvas.addEventListener('mousedown', function (e) {\r\n            var _a;\r\n            // e.preventDefault()\r\n            // e.stopPropagation()\r\n            var rect = _this.canvas.getBoundingClientRect();\r\n            var x = Math.round(e.clientX - rect.left);\r\n            var y = Math.round(e.clientY - rect.top);\r\n            _a = _this.screen.transformCanvasPosition(x, y), x = _a[0], y = _a[1];\r\n            var update = { category: \"game\", event: \"move-event\", timstamp: Date.now(), data: { x: x, y: y } };\r\n            _this.socket.send(JSON.stringify(update));\r\n        });\r\n        this.canvas.addEventListener('wheel', function (e) {\r\n            _this.screen.zoom = Math.min(2, Math.max(0.5, _this.screen.zoom + (e.deltaY / -1250)));\r\n            console.log(e);\r\n        });\r\n    };\r\n    return Game;\r\n}());\r\n\r\n\n\n//# sourceURL=webpack://GoBA-client/./src/game/index.ts?");

/***/ }),

/***/ "./src/index.ts":
/*!**********************!*\
  !*** ./src/index.ts ***!
  \**********************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _events__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./events */ \"./src/events.ts\");\n/* harmony import */ var _game__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./game */ \"./src/game/index.ts\");\n\r\n\r\nvar socket = new WebSocket(\"ws://localhost:8080/ws\");\r\nsocket.onopen = function () {\r\n    console.log(\"Websocket opened!\");\r\n};\r\nsocket.onclose = function () {\r\n    console.log(\"Websocket closed :(\");\r\n    socket.close();\r\n};\r\nsocket.onerror = function (ev) {\r\n    console.log(ev);\r\n    socket.close();\r\n};\r\n// \r\nvar game;\r\nvar canvas = document.getElementById(\"canvas\");\r\nvar ticks = new _events__WEBPACK_IMPORTED_MODULE_0__.TickCounter(\"tps\");\r\nsocket.onmessage = function (message) {\r\n    var event = JSON.parse(message.data);\r\n    switch (event.name) {\r\n        case \"setup\":\r\n            var setup = event.data;\r\n            game = new _game__WEBPACK_IMPORTED_MODULE_1__.Game(setup, canvas, socket);\r\n            break;\r\n        case \"tick\":\r\n            var tick = event.data;\r\n            console.log(tick.champions.length);\r\n            game.tick(tick);\r\n            ticks.update(event);\r\n            break;\r\n        case \"update-teams\":\r\n            var update = event.data;\r\n            game.updateTeams(update);\r\n            break;\r\n        default:\r\n            console.log(\"EVENT NOT PROCESSED\", event);\r\n            break;\r\n    }\r\n};\r\n\n\n//# sourceURL=webpack://GoBA-client/./src/index.ts?");

/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		if(__webpack_module_cache__[moduleId]) {
/******/ 			return __webpack_module_cache__[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/define property getters */
/******/ 	(() => {
/******/ 		// define getter functions for harmony exports
/******/ 		__webpack_require__.d = (exports, definition) => {
/******/ 			for(var key in definition) {
/******/ 				if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/ 					Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 				}
/******/ 			}
/******/ 		};
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	(() => {
/******/ 		__webpack_require__.o = (obj, prop) => Object.prototype.hasOwnProperty.call(obj, prop)
/******/ 	})();
/******/ 	
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	(() => {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = (exports) => {
/******/ 			if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	})();
/******/ 	
/************************************************************************/
/******/ 	// startup
/******/ 	// Load entry module
/******/ 	__webpack_require__("./src/index.ts");
/******/ 	// This entry module used 'exports' so it can't be inlined
/******/ })()
;