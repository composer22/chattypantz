'use strict';
//  This module contains the behavior for controlling the application interaction for the
// client app to the server.

angular.module('chattypantzApp').controller('MainCtrl', function($scope, $route, $timeout, chatService) {
  // Requests sent to the server.
  var REQUEST_TYPE = {
    SET_NICKNAME: 101,
    GET_NICKNAME: 102,
    LIST_ROOMS: 103,
    JOIN: 104,
    LIST_NAMES: 105,
    HIDE: 106,
    UNHIDE: 107,
    MSG: 108,
    LEAVE: 109
  }

  // Responses coming from the server.
  var RESPONSE_TYPE = {
    // Command responses
    SET_NICKNAME: 101,
    GET_NICKNAME: 102,
    LIST_ROOMS: 103,
    JOIN: 104,
    LIST_NAMES: 105,
    HIDE: 106,
    UNHIDE: 107,
    MSG: 108,
    LEAVE: 109,

    // Error Conditions
    ERR_ROOM_MANDATORY: 1001,
    ERR_MAX_ROOMS_REACHED: 1002,
    ERR_ROOM_UNAVAILABLE: 1003,
    ERR_NICKNAME_MANDATORY: 1004,
    ERR_ALREADY_JOINED: 1005,
    ERR_NICKNAME_USED: 1006,
    ERR_HIDDEN_NICKNAME: 1007,
    ERR_UNKNOWN_REQUEST: 1008
  }

  // Chat object contains data about the connection and state
  $scope.chat = {
    socket: null,
    status: {
      error: null,
      starting: false,
      started: false
    },
    data: {
      users: [],
      room: 'Demo',
      nickname: '',
      history: '',
      messageText: ''
    }
  }

  // Checks the nickname and create the socket object
  $scope.init = function() {
    if ($scope.nickname) {
      // Open a socket and register event handlers.
      var socket = chatService.socket();
      socket.onopen = $scope.onOpen;
      socket.onmessage = $scope.onMessage;
      socket.onerror = $scope.onError;
      socket.onclose = $scope.onError;

      // Initialize socket object.
      $scope.chat.socket = socket;
      $scope.chat.status.starting = true;
    }
  }

  //////// SOCKET EVENT HANDLERS ////////

  // Connected to the server.
  $scope.onOpen = function(e) {
    $scope.chat.status.starting = false;
    $scope.chat.status.started = true;
    $scope.chat.status.error = null;
    $scope.chat.data.nickname = $scope.nickname;
    // Set the nickname, join the demo room, and get a list of names.
    $scope.sendRequest('', REQUEST_TYPE.SET_NICKNAME, $scope.chat.data.nickname);
    $scope.sendRequest($scope.chat.data.room, REQUEST_TYPE.JOIN, $scope.chat.data.room);
    $scope.$apply(); // update UI
  };

  // Response arrived from the server.
  $scope.onMessage = function(message) {
    var response = angular.fromJson(message.data);
    switch (response.rspType) {
      case RESPONSE_TYPE.SET_NICKNAME:
        $scope.chat.data.history += "Chattypantz server: " + response.content + '\n';
        $scope.scrollBottom();
        $scope.$apply();
        break;
      case RESPONSE_TYPE.JOIN:
        var users = angular.fromJson(response.list);
        $scope.chat.data.users = users;
        $scope.chat.data.history += "Chattypantz server: " + response.content + '\n';
        $scope.scrollBottom();
        $scope.$apply();
        break;
      case RESPONSE_TYPE.LIST_NAMES:
        var users = angular.fromJson(response.list);
        $scope.chat.data.users = users;
        $scope.$apply();
        break;
      case RESPONSE_TYPE.MSG:
        $scope.chat.data.history += response.content + '\n';
        $scope.scrollBottom();
        $scope.$apply();
        break;
      case RESPONSE_TYPE.LEAVE:
        var users = angular.fromJson(response.list);
        $scope.chat.data.users = users;
        $scope.chat.data.history += "Chattypantz server: " + response.content + '\n';
        $scope.scrollBottom();
        $scope.$apply();
        break;
      case RESPONSE_TYPE.ERR_NICKNAME_USED:
        $scope.chat.data.history += "Chattypantz server: " + response.content + '\n';
        $scope.chat.data.history += "Quitting Chattypantz.\n";
        $scope.quit();
        break;
      default:
        $scope.chat.data.history += response.content + '\n';
        $scope.chat.data.history += "Quitting Chattypantz.\n";
        $scope.quit();
    }
  };

  // Connection error.
  $scope.onError = function(e) {
    // Stop run and show err.
    $scope.chat.status.error = "Server disconnected: " + e.reason;
    $scope.chat.status.started = false;
    console.log('onerror: ', e);
    $scope.$apply();
  };

  //////// UI EVENT HANDLERS ////////

  // Send button clicked; send message
  $scope.sendButtonClicked = function() {
    if ($scope.chat.data.messageText != '') {
      $scope.sendRequest($scope.chat.data.room, REQUEST_TYPE.MSG, $scope.chat.data.messageText);
      $scope.chat.data.messageText = '';
    }
  };

  // Quit button disconnects
  $scope.quit = function() {
    $scope.chat.socket.close();
    $route.reload();
  };

  //////// UTILITY FUNCTIONS ////////

  // Sends a request to the server
  $scope.sendRequest = function(room, type, content) {
    $scope.chat.socket.send(JSON.stringify({
      roomName: room,
      reqtype: type,
      content: content
    }));
  };

  // Scroll to the bottom of chat history
  $scope.scrollBottom = function() {
    angular.element("#chat-history")
      .scrollTop(angular.element("#chat-history")[0]
      .scrollHeight - angular.element("#chat-history")
      .height());
  };
});
