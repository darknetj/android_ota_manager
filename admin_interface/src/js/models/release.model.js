var $        = require('jquery'),
    _               = require('backbone'),
    Backbone        = require('backbone-crossdomain');

module.exports = Backbone.Model.extend({
    // Default attributes for the todo item.
    defaults: function() {
      return {
        title: "empty todo...",
        description: "",
        author: ""
      };
    },

    validate: function (attrs) {
        var errors = {};
        if (!attrs.title) errors.title = "Hey! Give this thing a title.";
        if (!attrs.description) errors.description = "You gotta write a description, duh!";
        if (!attrs.author) errors.author = "Put your name in dumb dumb...";

        if (!_.isEmpty(errors)) {
          return errors;
        }
    }
  });
