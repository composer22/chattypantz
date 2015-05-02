var requestType = require('requestType');

module.exports = requestType({
    SET_NICKNAME: 101,
    GET_NICKNAME: 102,
    LIST_ROOMS: 103,
    JOIN: 104,
    LIST_NAMES: 105,
    HIDE: 106,
    UNHIDE: 107,
    MSG: 108,
    LEAVE: 109
	});
