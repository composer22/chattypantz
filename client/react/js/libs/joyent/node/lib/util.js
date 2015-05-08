// Chattypants modified and stripped down for demo purposes.
var util = {
  isFunction: function(arg) {
    return typeof arg === 'function';
  },
  isNumber: function(arg) {
    return typeof arg === 'number';
  },
  isObject: function isObject(arg) {
    return typeof arg === 'object' && arg !== null;
  },
  isUndefined: function(arg) {
    return arg === void 0;
  },
};
