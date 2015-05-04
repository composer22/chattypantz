'use strict';

var server = "ws://127.0.0.1:6660/v1.0/chat";

angular.module('chattypantzApp').factory('chatService', function($http) {
  return {
    socket: function() {
      return new WebSocket(server);
    }
  }
});
